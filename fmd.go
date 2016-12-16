// fmd.go - Frosted Markdown
// ------

// Package frostedmd implements Frosted Markdown: standard Markdown to HTML
// conversion with a meta map and a default title. Parsing and rendering are
// handled by the excellent Blackfriday package; the Meta map is extracted
// from the first code block encountered, but only if it is not preceded by
// anything other than an optional header. The order can be reversed globally
// by setting META_AT_END to true, or at the Parser level.  In reversed order
// the meta code block must be the last element in the Markdown source.
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
//
// NOTE: This package will most likely be renamed, and might also be moved out
// of kisipar.  "Greysunday" was pretty tempting but then the sun came out...
package frostedmd

import (
	// Standard Library:
	"encoding/json"
	"errors"

	// Third-Party:
	"github.com/russross/blackfriday"
	"gopkg.in/yaml.v2"
)

// If true, get meta block from end of the Markdown source by default.
var META_AT_END = false

// The "Common" set of Blackfriday extensions; highly recommended for
// productive use of Markdown.
const COMMON_EXTENSIONS = 0 |
	blackfriday.EXTENSION_NO_INTRA_EMPHASIS |
	blackfriday.EXTENSION_TABLES |
	blackfriday.EXTENSION_FENCED_CODE |
	blackfriday.EXTENSION_AUTOLINK |
	blackfriday.EXTENSION_STRIKETHROUGH |
	blackfriday.EXTENSION_SPACE_HEADERS |
	blackfriday.EXTENSION_HEADER_IDS |
	blackfriday.EXTENSION_BACKSLASH_LINE_BREAK |
	blackfriday.EXTENSION_DEFINITION_LISTS

// The "Common" set of Blackfriday HTML flags; also highly recommended.
const COMMON_HTML_FLAGS = 0 |
	blackfriday.HTML_USE_XHTML |
	blackfriday.HTML_USE_SMARTYPANTS |
	blackfriday.HTML_SMARTYPANTS_FRACTIONS |
	blackfriday.HTML_SMARTYPANTS_DASHES |
	blackfriday.HTML_SMARTYPANTS_LATEX_DASHES

// Parser defines a parser-renderer that implements the page.Parser interface.
// It may be also be used to
type Parser struct {
	MetaAtEnd          bool
	MarkdownExtensions int // uses blackfriday EXTENSION_* constants
	HtmlFlags          int // uses blackfridy HTML_* constants
}

// New returns a new Parser with the common flags and extensions enabled.
func New() *Parser {
	return &Parser{
		MetaAtEnd:          META_AT_END,
		MarkdownExtensions: COMMON_EXTENSIONS,
		HtmlFlags:          COMMON_HTML_FLAGS,
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

// Parse converts Markdown input into a meta map and HTML content fragment,
// thus implementing the page.Parser interface. If an error is encountered
// while parsing the meta block, the rendered content is still returned.
// Thus the caller may choose to handle meta errors without interrupting flow.
func (p *Parser) Parse(input []byte) (*ParseResult, error) {

	// cf. renderer.go for the fmdRenderer definition
	renderer := &fmdRenderer{
		bfRenderer: blackfriday.HtmlRenderer(p.HtmlFlags,
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

	// Right now we only support JSON and YAML, so it's pretty easy to choose.
	if lang == "" {
		// We expect the JSON decoder to bail out fast on bad formats, so:
		if err := json.Unmarshal(input, mm); err == nil {
			return mm, nil
		}
		lang = "yaml"
	}

	// Guess the language if possible.
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
