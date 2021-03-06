// fmd.go - Frosted Markdown
// ------

// Package frostedmd converts Markdown files to structured data and HTML.
//
// The structured data is extracted from a Meta Block if present: a code block
// at the beginning (or, optionally, the end) of the file.  At the beginning
// of the file, the Meta Block may optionally be preceded by a single heading.
//
//  # Sample Doc
//
//      # Meta:
//      AmIYaml: true
//      Tags: [foo, bar, baz, baloney]
//
//  There you are.
//
// Parsing and rendering are handled by the excellent Blackfriday package:
// https://godoc.org/github.com/russross/blackfriday
//
// YAML processing is handled with the nearly canonical YAML package from
// Canonical: https://godoc.org/gopkg.in/yaml.v2
//
// The Meta Block position can be reversed globally by setting MetaBlockAtEnd
// to true, or at the Parser level.  In reversed order the meta code block
// must be the last element in the Markdown source.
//
// If the Meta contains no Title (nor "title" nor "TITLE") then the first
// heading is used, if and only if that heading was not preceded by any
// other block besides the Meta Block.
//
// Supported languages for the meta block are JSON and YAML (the default);
// additional languages as well as custom parsers are planned for the future.
//
// If an appropriate meta block is found it will be excluded from the rendered
// HTML content.
package frostedmd

import (
	// Standard Library:
	"encoding/json"
	"errors"

	// Third-Party:
	"gopkg.in/russross/blackfriday.v1"
	"gopkg.in/yaml.v2"
)

// MetaBlockAtEnd defines whether the block of data is expected at the end
// of the Markdown file, or (the default) at the beginning.
var MetaBlockAtEnd = false

// BlackFridayCommonExtensions defines the "Common" set of Blackfriday
// extensions, which are highly recommended for the productive use of
// Markdown.
const BlackFridayCommonExtensions = 0 |
	blackfriday.EXTENSION_NO_INTRA_EMPHASIS |
	blackfriday.EXTENSION_TABLES |
	blackfriday.EXTENSION_FENCED_CODE |
	blackfriday.EXTENSION_AUTOLINK |
	blackfriday.EXTENSION_STRIKETHROUGH |
	blackfriday.EXTENSION_SPACE_HEADERS |
	blackfriday.EXTENSION_HEADER_IDS |
	blackfriday.EXTENSION_BACKSLASH_LINE_BREAK |
	blackfriday.EXTENSION_DEFINITION_LISTS

// BlackFridayCommonHTMLFlags defines the "Common" set of Blackfriday HTML
// flags; also highly recommended.
const BlackFridayCommonHTMLFlags = 0 |
	blackfriday.HTML_USE_XHTML |
	blackfriday.HTML_USE_SMARTYPANTS |
	blackfriday.HTML_SMARTYPANTS_FRACTIONS |
	blackfriday.HTML_SMARTYPANTS_DASHES |
	blackfriday.HTML_SMARTYPANTS_LATEX_DASHES

// Parser defines a parser-renderer used for converting source data to HTML
// and metadata.
type Parser struct {
	MetaAtEnd          bool
	MarkdownExtensions int // uses blackfriday EXTENSION_* constants
	HTMLFlags          int // uses blackfridy HTML_* constants
}

// New returns a new Parser with the common flags and extensions enabled.
func New() *Parser {
	return &Parser{
		MetaAtEnd:          MetaBlockAtEnd,
		MarkdownExtensions: BlackFridayCommonExtensions,
		HTMLFlags:          BlackFridayCommonHTMLFlags,
	}
}

// NewBasic returns a new Parser without the common flags and extensions.
func NewBasic() *Parser {
	return &Parser{}
}

// ParseResult defines the result of a Parse operation.
type ParseResult struct {
	Meta    map[string]interface{} `json:"meta"`
	Content []byte                 `json:"content"`
}

// Parse converts Markdown input into a meta map and HTML content fragment.
// If an error is encountered while parsing the meta block, the rendered
// content is still returned. Thus the caller may choose to handle meta
// errors without interrupting flow.
func (p *Parser) Parse(input []byte) (*ParseResult, error) {

	// cf. renderer.go for the fmdRenderer definition
	renderer := &fmdRenderer{
		bfRenderer: blackfriday.HtmlRenderer(p.HTMLFlags,
			"", // no title
			"", // no css
		),
		metaAtEnd: p.MetaAtEnd,
	}

	htmlBytes := blackfriday.MarkdownOptions(input, renderer,
		blackfriday.Options{Extensions: p.MarkdownExtensions})

	// Partial results are useful sometimes.
	res := &ParseResult{Content: htmlBytes}

	mm, err := p.parseMeta(renderer.metaBytes, renderer.metaLang)
	if err != nil {
		return res, err
	}
	if mm["Title"] == nil && mm["TITLE"] == nil && mm["title"] == nil &&
		renderer.headerTitle != "" {
		mm["Title"] = renderer.headerTitle
	}
	res.Meta = mm
	return res, nil
}

func (p *Parser) parseMeta(input []byte, lang string) (map[string]interface{}, error) {

	mm := map[string]interface{}{}
	if len(input) == 0 {
		return mm, nil
	}

	// Only JSON and YAML are supported for undefined-language code blocks,
	// though we should keep in mind that it's possible the Markdown parser
	// might try at some point to guess.
	if lang == "" {
		// We expect the JSON decoder to bail out fast on bad formats, so:
		if err := json.Unmarshal(input, &mm); err == nil {
			return mm, nil
		}
		lang = "yaml"
	}

	switch lang {
	case "json":
		err := json.Unmarshal(input, &mm)
		if err != nil {
			return mm, err
		}
	case "yaml":
		err := yaml.Unmarshal(input, &mm)
		if err != nil {
			return mm, err
		}
	default:
		return mm, errors.New("Unsupported language for meta block: " + lang)
	}

	return mm, nil

}

// MarkdownBasic converts Markdown input using the same options as
// blackfriday.MarkdownBasic.  This is simply a convenience method for:
//  NewBasic().Parse(input)
func MarkdownBasic(input []byte) (*ParseResult, error) {

	return NewBasic().Parse(input)
}

// MarkdownCommon converts Markdown input using the same options as
// blackfriday.MarkdownCommon.  This is simply a convenience method for:
//  New().Parse(input)
func MarkdownCommon(input []byte) (*ParseResult, error) {

	return New().Parse(input)

}
