// cmd/fmd/main.go - Frosted Markdown tool.
//
// TODO:
// * XML output
// * Batch processing (maybe)

package main

import (
	// Standard library:
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	// Third-Party:
	"github.com/docopt/docopt-go"
	"gopkg.in/yaml.v2"

	// Our package:
	"github.com/biztos/frostedmd"
)

const VERSION = "0.1.0"

var EXIT_FUNCTION = os.Exit // testing-friendly abstraction
var STDERR = os.Stderr      // testing-friendly abstraction

type Options struct {
	File   string
	Format string
	Indent bool
	Force  bool
}

type YamlRes struct {
	Meta    map[string]interface{}
	Content string
}

func main() {
	opts := getOpts()
	b, err := ioutil.ReadFile(opts.File)
	if err != nil {
		fmt.Fprintln(STDERR, err)
		EXIT_FUNCTION(1)
	}

	// NOTE: we should get back a partial result even when we have an error.
	res, err := frostedmd.MarkdownCommon(b)
	if err != nil {
		fmt.Fprintln(STDERR, err)
		if !opts.Force {
			EXIT_FUNCTION(1)
		}
	}

	if opts.Format == "yaml" {
		// Bit of a hack here, of course, but I think it works:
		yres := YamlRes{Meta: res.Meta, Content: string(res.Content)}
		yaml, err := yaml.Marshal(yres)
		if err != nil {
			fmt.Fprintln(STDERR, err)
			EXIT_FUNCTION(1)
		}
		fmt.Println(string(yaml))

	} else if opts.Indent {
		jsonBytes, err := json.MarshalIndent(res, "", "    ")
		if err != nil {
			fmt.Fprintln(STDERR, err)
			EXIT_FUNCTION(1)
		}
		fmt.Println(string(jsonBytes))
	} else {
		jsonBytes, err := json.Marshal(res)
		if err != nil {
			fmt.Fprintln(STDERR, err)
			EXIT_FUNCTION(1)
		}
		fmt.Println(string(jsonBytes))
	}
}

func getOpts() *Options {

	usage := `fmd - Frosted Markdown tool.

    *** WARNING: ALPHA SOFTWARE! API MAY CHANGE AT ANY TIME! ***

Converts Frosted Markdown files into structured data.  Note that in JSON
output the HTML content is base64-encoded; this actually saves significant
space in the JSON file for any nontrivial amount of content.

In YAML output the HTML content is delviered as-is.

Note that XML is not supported for output, because of difficulties/ambiguities
around serialization.  And because XML is not a serialization language. :-)

For more information on Frosted Markdown see:

https://github.com/biztos/kisipar/frostedmd

Usage:
  fmd [options] FILE
  fmd --version
  fmd -h | --help

Options:
  -v, --version     Show version.
  -h, --help        Show this screen.
  -j, --json        Write output in JSON format (default).
  -y, --yaml        Write output in YAML format.
  -i, --indent      Indent output if applicable.
  -f, --force       Do not abort on errors (log them to STDERR).
`

	args, _ := docopt.Parse(
		usage,
		nil,            // use default os args
		true,           // enable help option
		"fmd "+VERSION, // the version string
		false,          // do NOT require options first
		true,           // let DocOpt exit for version, help, user error
	)

	// These will only fail if we screwed up the docopt stuff:
	json, _ := args["--json"].(bool)
	yaml, _ := args["--yaml"].(bool)
	file, _ := args["FILE"].(string)
	force, _ := args["--force"].(bool)
	indent, _ := args["--indent"].(bool)

	// Only one output format please.
	format := "json"
	if yaml {
		if json {
			fmt.Fprintln(STDERR, "Only one format allowed.")
			EXIT_FUNCTION(1)
		}
		format = "yaml"
	}

	return &Options{
		File:   file,
		Format: format,
		Force:  force,
		Indent: indent,
	}
}
