// cmd/fmd/main.go - Frosted Markdown tool.
//
// TODO:
// * Batch processing (maybe)
// * Make it stop printing "exit code 1" to STDERR on os.Exit(1)!
// * Additional options for the frostedmd parser.
// * Other output formats?  Anything useful? Meta only?
// * Convert to confluence wiki page?  (Really?) If so what about the meta?
// * Pure md->html function, in order to not keep another one around?
// * Option to lowercase all tags in the JSON and/or YAML (how hard is this?)

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
var STDOUT = os.Stdout      // dittto
var STDERR = os.Stderr      // dittto

const (
	OPTIONS_ERROR       = 1
	FILE_ERROR          = 2
	PARSE_ERROR         = 3
	SERIALIZATION_ERROR = 4
)

type Options struct {
	File     string
	Format   string
	Indent   bool
	NoBase64 bool
	Force    bool
	Silent   bool
	Test     bool
}

type YamlRes struct {
	Meta    map[string]interface{}
	Content string
}

func main() {
	opts := getOpts()
	b, err := ioutil.ReadFile(opts.File)
	failOnError(opts, err, FILE_ERROR)

	// NOTE: we should get back a partial result even when we have an error.
	res, err := frostedmd.MarkdownCommon(b)
	if err != nil {
		if !opts.Silent {
			warnErr(opts, err)
		}
		if !opts.Force {
			EXIT_FUNCTION(PARSE_ERROR)
		}
	}

	if opts.Test {
		EXIT_FUNCTION(0)
	}

	if opts.Format == "yaml" {
		// Bit of a hack here, of course, but I think it works:
		yres := YamlRes{Meta: res.Meta, Content: string(res.Content)}
		yaml, err := yaml.Marshal(yres)
		failOnError(opts, err, SERIALIZATION_ERROR)
		fmt.Fprintln(STDOUT, string(yaml))

	} else {
		// JSON has additional options.
		var src interface{}
		if opts.NoBase64 {
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
		if opts.Indent {
			jsonBytes, err = json.MarshalIndent(src, "", "    ")
			failOnError(opts, err, SERIALIZATION_ERROR)
		} else {
			jsonBytes, err = json.Marshal(src)
			failOnError(opts, err, SERIALIZATION_ERROR)
		}
		fmt.Fprintln(STDOUT, string(jsonBytes))
	}
}

func failOnError(opts *Options, err error, code int) {
	if err != nil {
		warnErr(opts, err)
		EXIT_FUNCTION(code)
	}
}

func warnErr(opts *Options, err error) {
	fmt.Fprintf(STDERR, "%s: %s", opts.File, err.Error())
}

func getOpts() *Options {

	usage := docOptUsageText()

	args, _ := docopt.Parse(
		usage,
		nil,            // use default os args
		true,           // enable help option
		"fmd "+VERSION, // the version string
		false,          // do NOT require options first
		true,           // let DocOpt exit for version, help, user error
	)

	// Let's try to be upstanding OSS citizens here, just in principle.
	if license, _ := args["--license"].(bool); license {
		fmt.Fprintln(STDOUT, frostedmd.LicenseFullText())
		EXIT_FUNCTION(0)
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
			fmt.Fprintln(STDERR, "Only one format allowed.")
			EXIT_FUNCTION(OPTIONS_ERROR)
		}
		format = "yaml"
	}

	return &Options{
		File:     file,
		Format:   format,
		Force:    force,
		Silent:   silent,
		Indent:   indent,
		NoBase64: nobase64,
		Test:     test,
	}
}

// And some text, isolated here because bindoc is overkill.
func docOptUsageText() string {
	return `fmd - Frosted Markdown tool.

    *** WARNING: ALPHA SOFTWARE! API MAY CHANGE AT ANY TIME! ***

Converts Frosted Markdown files into structured data containing two top-level
properties: 'meta' and 'content' -- the latter being the parsed HTML.

Note that in JSON output the HTML content is base64-encoded; this actually
saves significant space in the JSON file for any nontrivial amount of content.

Use the --nobase64 option to output a regular string instead.

In YAML output the HTML content is delviered as-is.

Note that XML is not supported for output, because of difficulties/ambiguities
around serialization.  And because XML is not a serialization language. :-)

More information on Frosted Markdown is available here:

https://github.com/biztos/frostedmd

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
  -n, --nobase64    Do not Base64-encode the JSON 'content' property.
  -f, --force       Do not abort on errors (log them to STDERR).
  -s, --silent      Do not print error messages.
  -t, --test        Parse file but do not print any output on success.
  --license         Print the software license.

Exit codes:

  0: OK.
  1: Parse failure.
  2: Filesystem error.
  3: Document-parsing error.
  4: Serialization error (should never happen).

Examples:

  Consider the following file as sample.md:

    > # Simple FMD
    > 
    >     Title: FMD FTW
    >     Description: Simple is as simple does.
    >     Tags: [fmd, golang, nerdery]
    > 
    > Good enough for me.

  (Stripped of the "> " quotes of course.)

  Basic JSON conversion looks like this:

    $ fmd -i sample.md
    {
        "meta": {
            "Description": "Simple is as simple does.",
            "Tags": [
                "fmd",
                "golang",
                "nerdery"
            ],
            "Title": "FMD FTW"
        },
        "content": "< Base64-Encoded String >"
    }
  
  And YAML looks like this:
  
    $ fmd -y sample.md
    meta:
      Description: Simple is as simple does.
      Tags:
      - fmd
      - golang
      - nerdery
      Title: FMD FTW
    content: |
      <h1>Simple FMD</h1>

      <p>Good enough for me.</p>
  
  A lot of fun can be had with your favorite UNIX shell, for instance with
  the "jq" JSON tool available here: https://stedolan.github.io/jq/

    $ fmd -n sample.md | jq -r .meta.Tags[]
    fmd
    golang
    nerdery
  
Acknowledgements:
  This software would have been immensely harder to write without the
  excellent work of other members of the Open Source Community. On the
  shoulders of giants we stand.  Any utility you get from this program
  is mostly thanks to the authors of the following packages:
  
  Black Friday: https://github.com/russross/blackfriday
  YAML for Go: https://github.com/go-yaml/yaml
  Testify: https://github.com/stretchr/testify
  DocOpt for Go: https://github.com/docopt/docopt.go

  ...and of course the Go Programming Language: https://golang.org

(c) 2016 Kevin A. Frost.  All rights reserved.  This is free software.
For the full license text, use the --license option.

`
}
