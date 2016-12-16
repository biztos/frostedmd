// fmd/fmd_mdfunc_test.go -- tests for the Markdown* convenience functions.
// ----------------------
package frostedmd_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/biztos/frostedmd"
)

func Test_MarkdownCommon_Simple(t *testing.T) {

	assert := assert.New(t)

	input := `
# Ima Title

    # I'm a comment!
    OldSchool: "YAML"

Plus "this."`

	expMap := map[string]interface{}{
		"Title":     "Ima Title",
		"OldSchool": "YAML",
	}
	expContent := "<h1>Ima Title</h1>\n\n<p>Plus &ldquo;this.&rdquo;</p>\n"

	res, err := frostedmd.MarkdownCommon([]byte(input))

	assert.Nil(err, "no error on basic MarkdownCommon")
	assert.Equal(expMap, res.Meta, "meta map as expected")
	assert.Equal(expContent, string(res.Content), "content as expected")

}

func Test_MarkdownBasic_Simple(t *testing.T) {

	assert := assert.New(t)

	input := `
# Ima Title

    # I'm a comment!
    OldSchool: "YAML"

Plus "this."`

	expMap := map[string]interface{}{
		"Title":     "Ima Title",
		"OldSchool": "YAML",
	}
	expContent := "<h1>Ima Title</h1>\n\n<p>Plus &quot;this.&quot;</p>\n"

	res, err := frostedmd.MarkdownBasic([]byte(input))

	assert.Nil(err, "no error on basic MarkdownCommon")
	assert.Equal(expMap, res.Meta, "meta map as expected")
	assert.Equal(expContent, string(res.Content), "content as expected")

}
