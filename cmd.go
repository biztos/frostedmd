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
	ExitFunction func(int)
	Stdout       io.Writer
	Stderr       io.Writer
}

// NewCmd returns a new Cmd for the given name, version and DocOpt usage
// specification.
func NewCmd(name, version, usage string) *Cmd {
	return &Cmd{
		Name:         name,
		Version:      version,
		Usage:        usage,
		ExitFunction: os.Exit,
		Stdout:       os.Stdout,
		Stderr:       os.Stderr,
	}
}

// Fail fails with a useful message based on err; if err is a CmdError
// and has an exit code set, that is used, otherwise the "other" code is
// used: CMD_OPTIONS_ERROR.
func (c *Cmd) Fail(err error) {

	switch e := err.(type) {
	case CmdError:
		if !e.Silent {
			fmt.Fprintln(c.Stderr, e.Error())
		}
		if !e.Force {
			code := e.Code
			if code == 0 {
				code = CMD_OPTIONS_ERROR
			}
			c.ExitFunction(code)
		}

	default:
		fmt.Fprintln(c.Stderr, err.Error())
		c.ExitFunction(CMD_OTHER_ERROR)
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
			Code:   CMD_FILE_ERROR,
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
	args, err := docopt.Parse(
		c.Usage,
		nil,     // use default os args
		true,    // enable help option
		vString, // the version string
		false,   // do NOT require options first
		true,    // let DocOpt exit on -v/-h/opts-error
	)
	if err != nil {
		return CmdError{
			Err:  err,
			Code: CMD_OPTIONS_ERROR,
		}
	}

	// Let's try to be upstanding OSS citizens here, just in principle.
	if license, _ := args["--license"].(bool); license {
		fmt.Fprintln(c.Stdout, LicenseFullText())
		c.ExitFunction(0)
	}

	// These will only fail if we screwed up the docopt stuff:
	json, _ := args["--json"].(bool)
	yaml, _ := args["--yaml"].(bool)
	indent, _ := args["--indent"].(bool)
	nobase64, _ := args["--nobase64"].(bool)
	force, _ := args["--force"].(bool)
	silent, _ := args["--silent"].(bool)
	test, _ := args["--test"].(bool)
	file, _ := args["FILE"].(string)

	// Only one output format please.
	format := "json"
	if yaml {
		if json {
			return CmdError{
				Err:  errors.New("Only one format allowed."),
				Code: CMD_OPTIONS_ERROR,
			}
		}
		format = "yaml"
	}

	c.Options = &CmdOptions{
		File:     file,
		Format:   format,
		Force:    force,
		Silent:   silent,
		Indent:   indent,
		NoBase64: nobase64,
		Test:     test,
	}

	return nil

}
