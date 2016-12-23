// main_test.go
//
// HIC SUNT DRACONES!  As you would imagine: testing main() is hard unless
// it's supported from the, uh, git-go, and in Go it is not.
//
// All the relevant logic is tested in the frostemd package.

package main // Can't be main_test or we can't call main().

import (
	// Standard library:
	"os"
	"testing"

	// Third-party:
	"github.com/stretchr/testify/assert"

	// Required for our Exit override:
	"github.com/biztos/frostedmd"
)

func TestMainForCoverage_Success(t *testing.T) {

	assert := assert.New(t)
	exited := -1
	frostedmd.DefaultExitFunction = func(c int) {
		exited = c
	}

	// Success, parsing a file in test mode so no output:
	os.Args = []string{"fmd", "-t", "sample.md"}
	main()
	assert.Equal(0, exited, "main() exited with success (code zero)")

}

func TestMainForCoverage_Failure(t *testing.T) {

	assert := assert.New(t)

	exited := -1
	frostedmd.DefaultExitFunction = func(c int) {
		exited = c
	}

	// Failure (broken yaml), in test mode and silent:
	os.Args = []string{"fmd", "-t", "-s", "broken.md"}
	main()
	assert.Equal(frostedmd.CMD_PARSE_ERROR, exited,
		"main() exited with a parse error")

}
