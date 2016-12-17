// license_test.go

package frostedmd_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/biztos/frostedmd"
)

func TestLicenseFullText(t *testing.T) {

	assert := assert.New(t)

	// Mostly just for coverage but what the hell, let's make sure we have
	// all our linkies here.
	linkies := []string{
		"https://github.com/russross/blackfriday",
		"https://github.com/go-yaml/yaml",
		"https://github.com/stretchr/testify",
		"https://github.com/docopt/docopt.go",
		"https://golang.org",
	}
	res := frostedmd.LicenseFullText()
	for _, link := range linkies {
		assert.Contains(res, link, "apparently have license for %s", link)
	}
}
