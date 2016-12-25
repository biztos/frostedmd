// renderer_test.go

package frostedmd_test

import (
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/russross/blackfriday"
	"github.com/stretchr/testify/assert"

	"github.com/biztos/frostedmd"
)

func Test_Renderer(t *testing.T) {

	// Try to exercise all of the renderer in one place; arguably a bit sloppy
	// but should cover the basics.
	assert := assert.New(t)

	// Ouch: we have two required methods of the blackfriday.Renderer
	// interface that are very hard (impossible?) to reach.  Hence this:
	assert.Equal(0, frostedmd.RendererTestCoverageShim(), "shim shummed")

	path := filepath.Join("test", "renderer.md")
	input, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	//	expContent := "<h1>Here</h1>\n\n<p>There.</p>\n"

	parser := frostedmd.New()
	t.Log(parser.MarkdownExtensions)
	parser.MarkdownExtensions = parser.MarkdownExtensions |
		blackfriday.EXTENSION_FOOTNOTES | // consider adding to defaults!
		blackfriday.EXTENSION_AUTO_HEADER_IDS | // ditto
		blackfriday.EXTENSION_LAX_HTML_BLOCKS |
		blackfriday.EXTENSION_HARD_LINE_BREAK |
		blackfriday.EXTENSION_TITLEBLOCK

	res, err := parser.Parse(input)
	assert.Nil(err)
	assert.NotNil(res)

	t.Log(string(res.Content))
}
