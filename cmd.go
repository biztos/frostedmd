// cmd.go - things useful for the "fmd" command (or anything build like it).

package frostedmd

import (

	// Standard library:
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	// Third-Party:
	"github.com/docopt/docopt-go"
	"github.com/russross/blackfriday"
	"gopkg.in/yaml.v2"
)

const (
	CMD_OPTIONS_ERROR       = 1
	CMD_FILE_ERROR          = 2
	CMD_PARSE_ERROR         = 3
	CMD_SERIALIZATION_ERROR = 4
	CMD_OTHER_ERROR         = 99
)

var DefaultExitFunction = os.Exit // override for testing main()

// CmdOptions describes the options available to the command.  The standard
// fmd command exposes all of them.
type CmdOptions struct {
	File          string
	Format        string
	Indent        bool
	NoBase64      bool
	ContentOnly   bool
	MetaOnly      bool
	PlainMarkdown bool
	Force         bool
	Silent        bool
	Test          bool
}

// CmdError defines an error in the command-running context.
type CmdError struct {
	Err    error  // The source error.
	File   string // The file to include in the error message.
	Code   int    // The exit code, if Force is not set.
	Silent bool   // If true, do not print the error.
	Force  bool   // If true, do not exit.
}

// Error stringifies the error per the error interface.
func (e CmdError) Error() string {

	if e.File == "" {
		return e.Err.Error()
	}
	return fmt.Sprintf("%s: %s", e.File, e.Err.Error())

}

// Cmd defines a command-line program or its equivalent.
type Cmd struct {
	Name    string
	Version string
	Usage   string
	Options *CmdOptions
	Result  *ParseResult

	// In order to make testing realistically possible in the command context
	// we make these standard things overrideable:
	Exit   func(int)
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
}

// NewCmd returns a new Cmd for the given name, version and DocOpt usage
// specification.  The Exit property is set to the DefaultExitFunction,
// allowing a global override for testing.
func NewCmd(name, version, usage string) *Cmd {
	return &Cmd{
		Name:    name,
		Version: version,
		Usage:   usage,
		Exit:    DefaultExitFunction,
		Stdin:   os.Stdin,
		Stdout:  os.Stdout,
		Stderr:  os.Stderr,
	}
}

// Run calls the methods used for a standard command run in order: SetOptions,
// ParseFile, and finally PrintResult.  The first error encountered is
// returned, to (normally) be passed to Fail.  Note that in the interest
// of simplicity, docopt is allowed to exit directly from within SetOptions.
func (c *Cmd) Run() error {
	if err := c.SetOptions(); err != nil {
		return err
	}
	if err := c.ParseFile(); err != nil {
		return err
	}
	return c.PrintResult()
}

// Fail fails with a useful message based on err; if err is a CmdError
// and has an exit code set, that is used, otherwise the "other" code is
// used: CMD_OTHER_ERROR.
func (c *Cmd) Fail(err error) {

	switch e := err.(type) {
	case CmdError:
		if !e.Silent {
			fmt.Fprintln(c.Stderr, e.Error())
		}
		if !e.Force {
			code := e.Code
			if code == 0 {
				code = CMD_OTHER_ERROR
			}
			c.Exit(code)
		}

	default:
		fmt.Fprintln(c.Stderr, err.Error())
		c.Exit(CMD_OTHER_ERROR)
	}

}

// CmdYamlRes is a shim to handle output serialization to YAML.
type CmdYamlRes struct {
	Meta    map[string]interface{}
	Content string
}

// ParseFile parses a single file contained in the Options, according to the
// other options therein, and returning any error.  Note that a useful result
// *may* be returned together with an error.  If the Options.File is the
// empty string, input is read from the command's Stdin (os.Stdin by default).
func (c *Cmd) ParseFile() error {

	var input []byte
	var err error
	if c.Options.File == "" {
		input, err = ioutil.ReadAll(c.Stdin)
	} else {
		input, err = ioutil.ReadFile(c.Options.File)
	}
	if err != nil {
		return CmdError{
			Code: CMD_FILE_ERROR,
			Err:  err,
			File: c.Options.File,
		}
	}

	// In some cases we skip the whole Frosted Markdown bag of tricks, thus
	// allowing the use of the fmd tool as a generic converter (strongly
	// favoring Blackfriday extensions of course).
	if c.Options.PlainMarkdown {
		c.Result = &ParseResult{Content: blackfriday.MarkdownCommon(input)}
		return nil
	}

	// NOTE: we should get back a partial result even when we have an error.
	res, err := MarkdownCommon(input)
	c.Result = res // cf. the Force option
	if err != nil {
		return CmdError{
			Code:   CMD_PARSE_ERROR,
			Err:    err,
			File:   c.Options.File,
			Silent: c.Options.Silent,
			Force:  c.Options.Force,
		}
	}

	return nil

}

// PrintResult prints the Result according to the Options, with output
// going to c.Stdout.  Any error returned should be considered fatal.
// If Result is nil, nothing is printed; this is normal if the Test option
// is set.
func (c *Cmd) PrintResult() error {

	res := c.Result
	if res == nil || c.Options.Test {
		return nil
	}

	// If we only want the content, life is very simple.
	if c.Options.ContentOnly || c.Options.PlainMarkdown {
		fmt.Fprintln(c.Stdout, string(res.Content))
		return nil
	}

	// With or without Content, the nature of the Meta means the encoder will
	// need to use introspection (reflect).
	var src interface{}
	if c.Options.Format == "yaml" {

		if c.Options.MetaOnly {
			src = res.Meta
		} else {
			src = map[string]interface{}{
				"meta":    res.Meta,
				"content": string(res.Content),
			}
		}
		yaml, err := safeMarshalYaml(src)
		if err != nil {
			return CmdError{
				Code: CMD_SERIALIZATION_ERROR,
				Err:  err,
				File: c.Options.File,
			}
		}
		fmt.Fprintln(c.Stdout, string(yaml))
		return nil
	}

	// JSON, the default,  has additional options.
	if c.Options.MetaOnly {
		src = res.Meta
	} else if c.Options.NoBase64 {
		// Only []byte values are Base64-encoded, strings are not.
		src = map[string]interface{}{
			"meta":    res.Meta,
			"content": string(res.Content),
		}
	} else {
		src = res
	}
	var jsonBytes []byte
	var err error
	if c.Options.Indent {
		jsonBytes, err = json.MarshalIndent(src, "", "    ")
	} else {
		jsonBytes, err = json.Marshal(src)

	}
	if err != nil {
		return CmdError{
			Code: CMD_SERIALIZATION_ERROR,
			Err:  err,
			File: c.Options.File,
		}
	}

	fmt.Fprintln(c.Stdout, string(jsonBytes))

	return nil
}

// NOTES:
// 1. This appears to be the idiomatic way to trap a panic in an external
//    package while still returning sane values.  If not, please let me know!
// 2. This is only necessary because the yaml package panics where it should
//    return an error.  TODO: fix there, and send a pull request.
func safeMarshalYaml(input interface{}) ([]byte, error) {

	var output []byte
	var err error
	func() {
		defer func() {
			if r := recover(); r != nil {
				err = fmt.Errorf("yaml error: %s", r)
			}
		}()
		output, err = yaml.Marshal(input)
	}()

	return output, err

}

// SetOptions sets the Options struct according to the Usage. If a --license
// boolean option is found, the LicenseFullText is printed and the program
// exits with success.  The --help and --version (-h and -v) options are
// handled similarly by docopt, and it will also fail directly to os.Exit
// on any option errors.
func (c *Cmd) SetOptions() error {

	vString := c.Name + " " + c.Version
	args, _ := docopt.Parse(
		c.Usage,
		nil,     // use default os args
		true,    // enable help option
		vString, // the version string
		false,   // do NOT require options first
		true,    // let DocOpt exit on -v/-h/opts-error (thus ignore err)
	)

	// These must be boolean or undefined:
	// TODO: come up with a more generic wrapper for DocOpt that does all this
	// stuff with introspection on the options type.
	mustBool := []string{
		"--json",
		"--yaml",
		"--indent",
		"--nobase64",
		"--force",
		"--silent",
		"--test",
		"--content",
		"--meta",
		"--plainmd",
		"--license",
	}
	have := map[string]bool{}
	for _, key := range mustBool {
		if val, ok := args[key]; ok {
			if v, ok := val.(bool); ok {
				have[key] = v
			} else {
				// Here we don't want to return as it's programmer error.
				panic("Bad designation of bool opt in docopt usage: " + key)
			}
		}

	}

	// This is the only string arg; again the caller might leave it out in
	// favor of some other strategy.  However as of now it's not clear one
	// even *could* set this to any type other than string in docopt, so we
	// will not leave a hole in the test coverage for that.
	file, _ := args["FILE"].(string)

	// Let's try to be upstanding OSS citizens here, just in principle.
	if have["--license"] {
		fmt.Fprintln(c.Stdout, LicenseFullText())
		c.Exit(0)
	}

	// Only one output format please.
	format := "json"
	if have["--yaml"] {
		if have["--json"] {
			return CmdError{
				Err:  errors.New("Only one format allowed."),
				Code: CMD_OPTIONS_ERROR,
			}
		}
		format = "yaml"
	}

	// Catch any contradictory options.
	if have["--meta"] && have["--content"] {
		// Obviously can't output MetaOnly and ContentOnly.
		return CmdError{
			Err:  errors.New("--meta and --content are mutually exclusive."),
			Code: CMD_OPTIONS_ERROR,
		}
	}
	if have["--plainmd"] {
		// PlainMarkdown overrides all other options at the moment.
		// TODO: allow the "basic" option when we implement it.
		format = ""
		for k, v := range have {
			if k != "--plainmd" && v == true {
				return CmdError{
					Err:  errors.New("--plainmd excludes other options."),
					Code: CMD_OPTIONS_ERROR,
				}
			}
		}
	}

	c.Options = &CmdOptions{
		File:          file,
		Format:        format,
		Force:         have["--force"],
		Silent:        have["--silent"],
		Indent:        have["--indent"],
		NoBase64:      have["--nobase64"],
		Test:          have["--test"],
		ContentOnly:   have["--content"],
		MetaOnly:      have["--meta"],
		PlainMarkdown: have["--plainmd"],
	}

	return nil

}
