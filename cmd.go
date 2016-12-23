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
	File        string
	Format      string
	Indent      bool
	NoBase64    bool
	ContentOnly bool
	MetaOnly    bool
	Force       bool
	Silent      bool
	Test        bool
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
	if err := c.PrintResult(); err != nil {
		return err
	}
	return nil
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
// *may* be returned together with an error.
func (c *Cmd) ParseFile() error {
	b, err := ioutil.ReadFile(c.Options.File)
	if err != nil {
		return CmdError{
			Code: CMD_FILE_ERROR,
			Err:  err,
			File: c.Options.File,
		}
	}

	// NOTE: we should get back a partial result even when we have an error.
	res, err := MarkdownCommon(b)
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

	if c.Options.Test {
		c.Result = nil
	}

	return nil

}

// PrintResult prints the Result according to the Options, with output
// going to c.Stdout.  Any error returned should be considered fatal.
// If Result is nil, nothing is printed; this is normal if the Test option
// is set.
func (c *Cmd) PrintResult() error {

	res := c.Result
	if res == nil {
		return nil
	}

	if c.Options.Format == "yaml" {
		// Bit of a hack here, of course, but I think it works:
		yres := CmdYamlRes{Meta: res.Meta, Content: string(res.Content)}
		yaml, err := yaml.Marshal(yres)
		if err != nil {
			return CmdError{
				Code: CMD_SERIALIZATION_ERROR,
				Err:  err,
				File: c.Options.File,
			}
		}
		fmt.Fprintln(c.Stdout, string(yaml))
	} else {
		// JSON has additional options.
		var src interface{}
		if c.Options.NoBase64 {
			// Only []byte values are Base64-encoded, strings are not.
			src = map[string]interface{}{
				"meta":    res.Meta,
				"content": string(res.Content),
			}
		} else {
			// Nb: this is exactly the kind of shit I'm not supposed to get
			// away with in a "strongly-typed language."
			src = res
		}
		var jsonBytes []byte
		var err error
		if c.Options.Indent {
			jsonBytes, err = json.MarshalIndent(src, "", "    ")
			if err != nil {
				return CmdError{
					Code: CMD_SERIALIZATION_ERROR,
					Err:  err,
					File: c.Options.File,
				}
			}
		} else {
			jsonBytes, err = json.Marshal(src)
			if err != nil {
				return CmdError{
					Code: CMD_SERIALIZATION_ERROR,
					Err:  err,
					File: c.Options.File,
				}
			}
		}
		fmt.Fprintln(c.Stdout, string(jsonBytes))
	}

	return nil
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

	// Don't contradict what you want us to output.
	if have["--meta"] && have["--content"] {
		return CmdError{
			Err:  errors.New("--meta and --content are mutually exclusive."),
			Code: CMD_OPTIONS_ERROR,
		}
	}

	// Let's try to be upstanding OSS citizens here, just in principle.
	if have["--license"] {
		fmt.Fprintln(c.Stdout, LicenseFullText())
		c.Exit(0)
	}

	c.Options = &CmdOptions{
		File:        file,
		Format:      format,
		Force:       have["--force"],
		Silent:      have["--silent"],
		Indent:      have["--indent"],
		NoBase64:    have["--nobase64"],
		Test:        have["--test"],
		ContentOnly: have["--content"],
		MetaOnly:    have["--meta"],
	}

	return nil

}
