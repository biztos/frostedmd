// fmd/fmd_test.go -- general tests for Frosted Markdown
// ---------------
package frostedmd_test

import (
	"github.com/biztos/frostedmd"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Parse_SimpleYAML_NoLang(t *testing.T) {

	assert := assert.New(t)

	input := `
# Ima Title

    # I'm a comment!
    OldSchool: "YAML"

Plus "this."
`
	expMap := map[string]interface{}{
		"Title":     "Ima Title",
		"OldSchool": "YAML",
	}
	expContent := "<h1>Ima Title</h1>\n\n<p>Plus &ldquo;this.&rdquo;</p>\n"

	res, err := frostedmd.New().Parse([]byte(input))

	assert.Nil(err, "no error on basic MarkdownCommon")
	assert.Equal(expMap, res.Meta, "meta map as expected")
	assert.Equal(expContent, string(res.Content), "content as expected")

}

func Test_Parse_SimpleJSON_NoLang(t *testing.T) {

	assert := assert.New(t)

	input := `
# Ima Title

    {
        "OldSchool": "YAML"
    }

Plus "this."`

	expMap := map[string]interface{}{
		"Title":     "Ima Title",
		"OldSchool": "YAML",
	}
	expContent := "<h1>Ima Title</h1>\n\n<p>Plus &ldquo;this.&rdquo;</p>\n"

	res, err := frostedmd.New().Parse([]byte(input))

	assert.Nil(err, "no error on basic MarkdownCommon")
	assert.Equal(expMap, res.Meta, "meta map as expected")
	assert.Equal(expContent, string(res.Content), "content as expected")

}

func Test_Parse_SimpleYAML_WithLang(t *testing.T) {

	assert := assert.New(t)

	input := "# Ima Title\n\n```yaml\n" +
		"# I'm a comment!\nOldSchool: \"YAML\"\n" +
		"```\n\nPlus \"this.\""

	expMap := map[string]interface{}{
		"Title":     "Ima Title",
		"OldSchool": "YAML",
	}
	expContent := "<h1>Ima Title</h1>\n\n<p>Plus &ldquo;this.&rdquo;</p>\n"

	res, err := frostedmd.New().Parse([]byte(input))

	assert.Nil(err, "no error on basic MarkdownCommon")
	assert.Equal(expMap, res.Meta, "meta map as expected")
	assert.Equal(expContent, string(res.Content), "content as expected")

}

func Test_Parse_SimpleJSON_WithLang(t *testing.T) {

	assert := assert.New(t)

	input := "# Ima Title\n\n```json\n" +
		"{\"OldSchool\": \"YAML\"}\n" +
		"```\n\nPlus \"this.\""

	expMap := map[string]interface{}{
		"Title":     "Ima Title",
		"OldSchool": "YAML",
	}
	expContent := "<h1>Ima Title</h1>\n\n<p>Plus &ldquo;this.&rdquo;</p>\n"

	res, err := frostedmd.New().Parse([]byte(input))

	assert.Nil(err, "no error on basic MarkdownCommon")
	assert.Equal(expMap, res.Meta, "meta map as expected")
	assert.Equal(expContent, string(res.Content), "content as expected")

}

func Test_Parse_Error_JSON(t *testing.T) {

	assert := assert.New(t)

	input := "# Here\n\n```json\n{ foo: \"bar }\n```\n\nThere."

	expContent := "<h1>Here</h1>\n\n<p>There.</p>\n"

	res, err := frostedmd.New().Parse([]byte(input))

	if assert.Error(err, "error returned") {
		assert.Regexp("invalid character 'f'", err.Error(), "error useful")
	}
	assert.Nil(res.Meta, "empty meta map")
	assert.Equal(expContent, string(res.Content), "content as expected")
}

func Test_Parse_Error_YAML(t *testing.T) {

	assert := assert.New(t)

	input := "# Here\n\n```yaml\nfoo: [1,true,3\n```\n\nThere."

	expContent := "<h1>Here</h1>\n\n<p>There.</p>\n"

	res, err := frostedmd.New().Parse([]byte(input))

	if assert.Error(err, "error returned") {
		assert.Regexp("yaml.*did not find expected", err.Error(),
			"error useful")
	}
	assert.Nil(res.Meta, "empty meta map")
	assert.Equal(expContent, string(res.Content), "content as expected")
}

func Test_Parse_LateMetaIgnored(t *testing.T) {

	assert := assert.New(t)

	input := "# Here\n\nThere!\n\n```yaml\nfoo: [1,true,3\n```\n\nDone."

	expMap := map[string]interface{}{
		"Title": "Here",
	}
	expContent := `<h1>Here</h1>

<p>There!</p>

<pre><code class="language-yaml">foo: [1,true,3
</code></pre>

<p>Done.</p>
`
	res, err := frostedmd.New().Parse([]byte(input))

	assert.Nil(err, "no error returned")
	assert.Equal(expMap, res.Meta, "meta map as expected")
	assert.Equal(expContent, string(res.Content), "content as expected")

}

func Test_Parse_LateHeaderIgnored(t *testing.T) {

	assert := assert.New(t)

	input := "```yaml\nfoo: Bar\n```\n\nThere.\n\n# Elsewhere."

	expMap := map[string]interface{}{
		"foo": "Bar",
	}
	expContent := `<p>There.</p>

<h1>Elsewhere.</h1>
`
	res, err := frostedmd.New().Parse([]byte(input))

	assert.Nil(err, "no error returned")
	assert.Equal(expMap, res.Meta, "meta map as expected")
	assert.Equal(expContent, string(res.Content), "content as expected")

}

func Test_Parse_PostBlockIgnored(t *testing.T) {

	assert := assert.New(t)

	input := "Here.\n\n# There!\n\n```yaml\nfoo: [1,true,3\n```\n"

	expContent := `<p>Here.</p>

<h1>There!</h1>

<pre><code class="language-yaml">foo: [1,true,3
</code></pre>
`
	res, err := frostedmd.New().Parse([]byte(input))

	assert.Nil(err, "no error returned")
	assert.Equal(map[string]interface{}{}, res.Meta, "empty meta map")
	assert.Equal(expContent, string(res.Content), "content as expected")
}

func Test_Parse_MetaAtEnd(t *testing.T) {

	assert := assert.New(t)

	input := "# Ima Title\n\nSome block here.\n\n```yaml\n" +
		"# I'm a comment!\nOldSchool: \"YAML\"\n" +
		"```\n"

	expMap := map[string]interface{}{
		"Title":     "Ima Title",
		"OldSchool": "YAML",
	}
	expContent := `<h1>Ima Title</h1>

<p>Some block here.</p>
`

	parser := frostedmd.New()
	parser.MetaAtEnd = true

	res, err := parser.Parse([]byte(input))

	assert.Nil(err, "no error on basic MarkdownCommon")
	assert.Equal(expMap, res.Meta, "meta map as expected")
	assert.Equal(expContent, string(res.Content), "content as expected")

}

func Test_Parse_MetaAtEnd_MultiBlock(t *testing.T) {

	assert := assert.New(t)

	input := `# Ima Title

Some block here.

    # I'm an earlier code block!
    Truly: "not the meta"

Another here.

    Ahoy: "Not the Meta Yet"

Then finally:

    # I'm the last block!
    OldSchool: "YAML"

`
	expMap := map[string]interface{}{
		"Title":     "Ima Title",
		"OldSchool": "YAML",
	}
	expContent := `<h1>Ima Title</h1>

<p>Some block here.</p>

<pre><code># I'm an earlier code block!
Truly: &quot;not the meta&quot;
</code></pre>

<p>Another here.</p>

<pre><code>Ahoy: &quot;Not the Meta Yet&quot;
</code></pre>

<p>Then finally:</p>
`

	parser := frostedmd.New()
	parser.MetaAtEnd = true

	res, err := parser.Parse([]byte(input))

	assert.Nil(err, "no error on basic MarkdownCommon")
	assert.Equal(expMap, res.Meta, "meta map as expected")
	assert.Equal(expContent, string(res.Content), "content as expected")

}

func Test_Parse_MetaAtEnd_OtherBlockFollows(t *testing.T) {

	assert := assert.New(t)

	input := `# Ima Title

Some block here.

    # I'm an earlier code block!
    Truly: "not the meta"

Another here.

    Ahoy: "Not the Meta Yet"

Then finally:

    # I'm the last block!
    OldSchool: "YAML"

## Because of this.
`
	expMap := map[string]interface{}{
		"Title": "Ima Title",
	}
	expContent := `<h1>Ima Title</h1>

<p>Some block here.</p>

<pre><code># I'm an earlier code block!
Truly: &quot;not the meta&quot;
</code></pre>

<p>Another here.</p>

<pre><code>Ahoy: &quot;Not the Meta Yet&quot;
</code></pre>

<p>Then finally:</p>

<pre><code># I'm the last block!
OldSchool: &quot;YAML&quot;
</code></pre>

<h2>Because of this.</h2>
`

	parser := frostedmd.New()
	parser.MetaAtEnd = true

	res, err := parser.Parse([]byte(input))

	assert.Nil(err, "no error on basic MarkdownCommon")
	assert.Equal(expMap, res.Meta, "meta map as expected")
	assert.Equal(expContent, string(res.Content), "content as expected")

}
