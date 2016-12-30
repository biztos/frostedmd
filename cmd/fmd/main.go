// The Frosted Markdown Tool.
//
// STATUS
//
// This is alpha software and should be considered unstable.  For more
// information please visit the main repo: https://github.com/biztos/frostedmd
//
// SYNOPSIS
//
// The fmd command converts Markdown files to structured data using the
// frostedmd package.
//
//    $ fmd -i sample.md
//    {
//        "meta": {
//            "Description": "Simple is as simple does.",
//            "Tags": [
//                "fmd",
//                "golang",
//                "nerdery"
//            ],
//            "Title": "FMD FTW"
//        },
//        "content": "< Base64-Encoded String >"
//    }
//
// For more information simply invoke the program's help option:
//
//    fmd --help
//
// INSTALLATION
//
// Follow these steps to build your own:
//
//   go get -u github.com/biztos/frostedmd
//   go build github.com/biztos/frostedmd/cmd/fmd
//   ./fmd --version # should work!
//
// Binaries for a number of platforms will be made available as soon as this
// tool is a little more stable.
//
// ROADMAP
//
// (prioritized as of 2016-12-30)
//
//  * Take STDIN if no file provided (this seems more natural than "-").
//    (Needed for sane TextMate previews.)
//
//  * -d option to produce a full HTML5 document in the Content.
//
//  * --style=X option to use an optional canned stylesheet *or* file
//    (included in full) *or* link to stylesheet (if neither file nor
//    canned style option).  Possibly with different options; TBD.
//
//  *** TAG VERSION: 0.9 -- READY FOR USE, PENDING BUGFIXES ***
//
//  * -e option for meta-at-end
//
//  * A set of e2e tests using Run() for all the arg possiblities...
//    Probably keep expected input/output in files for easier diffing...
//    Consider testing this via testig: dir,infile,expfile,func.
//
//  *** TAG VERSION: 1.0 -- READY FOR PUBLIC USE PROPER-LIKE! ***
//
//  * MAYBE -b option to use MarkdownBasic (few/no extensions)
//    PRO: lets you skip extensions and sanity-check vs. dumber parsers.
//    CON: why bother?  is that even a use case?  vs. setting/unsetting flags?
//
//  * Option to lowercase all tags in the JSON and/or YAML (how hard is this?)
//    Use reflection, create a new interface in which all the string keys are
//    lowercased?  This also seems like a generically useful case, if at all
//    useful, so maybe it goes in utli?
//
//  * Options to fine-tune the parser behavior via the ext/flage e.g.
//    blackfriday.EXTENSION_* and blackfriday.HTML_* -- so with on/off
//    switches like -o/--option=FOOBAR, -O/--nooption=FOOBAR.
//    However: should there be a reset flag for this stuff?  Because you'd
//    mostly want to do add things but might want to specify "basic plus"
//    (Nasty bit: need to then also have --list-options or something.)
//
//  * MAYBE Other output formats?  For instance "go-quoted?"
//    It would be nice to be able to tell FOR SURE what the parsed interface{}
//    is, so e.g. timestamps can be debugged. Any other use-case?
//    Any other formats? Asciidoc? (Probably a lot of work.)
//
//  * MAYBE --template=X option to use html/template and a file... MAYBE.
//    Is this really useful for anything?  Wouldn't you need a bunch of
//    template files? The tool is already at >3M, which is a little
//    crazy, and adding template support would be at least another 1.1M.
//    However, templates are a nice way to implement the -d option so
//    it may be a moot point by the time we get here.
package main

import (
	"github.com/biztos/frostedmd"
)

const VERSION = "0.1.0"

func main() {

	cmd := frostedmd.NewCmd("fmd", VERSION, DocOptUsageText)
	if err := cmd.Run(); err != nil {
		cmd.Fail(err)
	} else {
		cmd.Exit(0) // <-- facilitate testing!
	}

}

var DocOptUsageText = `fmd - Frosted Markdown tool.

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
  fmd --license
  fmd -h | --help

Options:
  -v, --version     Show version.
  -h, --help        Show this screen.
  -j, --json        Write output in JSON format (default).
  -y, --yaml        Write output in YAML format.
  -i, --indent      Indent output if applicable.
  -n, --nobase64    Do not Base64-encode the JSON 'content' property.
  -c, --content     Only print the content (as a string), not the meta.
  -m, --meta        Only print the meta block, not the content.
  -p, --plainmd     Convert as "plain" Markdown (not Frosted Markdown).
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
