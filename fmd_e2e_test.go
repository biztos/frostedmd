// frostedmd/fmd_e2e_test.go - end-to-end tests of Frosted Markdown.
// -------------------------
// NOTE: it's easier, if less pure, to keep this stuff in files.

package frostedmd_test

import (
	// Standard Library:
	"io/ioutil"
	"path/filepath"
	"testing"

	// Third Party:
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"

	// Under Test:
	"github.com/biztos/frostedmd"
)

func readTestFile(name string) []byte {
	path := filepath.Join("test", name)
	b, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return b
}

func Test_FilesEndToEnd_Common(t *testing.T) {

	assert := assert.New(t)

	input := readTestFile("common.md")
	expContent := readTestFile("common.html")
	expYaml := readTestFile("common.yaml")

	res, error := frostedmd.MarkdownCommon(input)
	assert.Nil(error, "no error from MarkdownCommon")
	assert.Equal(string(expContent), string(res.Content),
		"content as expected")
	if assert.NotNil(res.Meta, "Meta not nil") {
		yaml, err := yaml.Marshal(res.Meta)
		if assert.Nil(err, "Meta convertible to YAML") {
			assert.Equal(string(expYaml), string(yaml),
				"converted YAML as expected")
		}
	}
}

// TODO: test with all options on
// TODO: test Basic too
