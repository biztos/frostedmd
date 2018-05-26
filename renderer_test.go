// renderer_test.go

package frostedmd_test

import (
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/russross/blackfriday.v1"

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

	exp := `<h1 class="title">pandoc-style title
</h1>
<p>Here we try to completely exercise the renderer.</p>

<h1 id="first">first!</h1>

<hr />

<p><strong>Here&rsquo;s a quote:</strong></p>

<blockquote>
<p>I am often misquoted. &ndash; Abraham Lincoln</p>
</blockquote>

<p><em>in which</em></p>

<ol>
<li>There is little truth.<br /></li>
<li>There is much meme.<br />

<ul>
<li>Metameme!<br />
<br /></li>
</ul></li>
</ol>

<p>A table, perhaps?</p>

<table>
<thead>
<tr>
<th>Fee</th>
<th align="center">Fi</th>
</tr>
</thead>

<tbody>
<tr>
<td>$12.20</td>
<td align="center">fish-fry</td>
</tr>
</tbody>
</table>

<p>An autolink: <a href="https://frostopolis.com/">https://frostopolis.com/</a></p>

<p class="coverage">Perhaps a 'graph?</p>

<p>Break?<br />
Possibly.</p>

<p>Thing? No_Thing!</p>

<p>Strike? <del>Struck!</del></p>

<p>Note? Noted.<sup class="footnote-ref" id="fnref:1"><a href="#fn:1">1</a></sup></p>

<pre><code class="language-go">func do() bool {
    fmt.Println(&quot;done&quot;)
    return true
}
</code></pre>

<pre><code># and the old school
echo &quot;yep&quot;
</code></pre>

<p>Other quotes are &ldquo;smart.&rdquo;</p>

<p><strong><em>Definitive:</em></strong></p>

<dl>
<dt>Fee<br /></dt>
<dd>A price to pay<br />
<br /></dd>
</dl>

<p>TODO: figure out how to exercise</p>
<div class="footnotes">

<hr />

<ol>
<li id="fn:1">Duly <em>noted</em> no less.<br />
</li>
</ol>
</div>
`

	parser := frostedmd.New()
	parser.MarkdownExtensions = parser.MarkdownExtensions |
		blackfriday.EXTENSION_FOOTNOTES | // consider adding to defaults!
		blackfriday.EXTENSION_AUTO_HEADER_IDS | // ditto
		blackfriday.EXTENSION_LAX_HTML_BLOCKS |
		blackfriday.EXTENSION_HARD_LINE_BREAK |
		blackfriday.EXTENSION_TITLEBLOCK

	res, err := parser.Parse(input)
	assert.Nil(err)
	if assert.NotNil(res) {
		assert.Equal(exp, string(res.Content), "HTML as expected")
	}

}
