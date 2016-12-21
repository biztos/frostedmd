// cmd_test.go

package frostedmd_test

import (
	// Standard library:
	"bufio"
	"bytes"
	"errors"
	"os"
	"testing"

	// Third-party:
	"github.com/stretchr/testify/assert"

	// Under test:
	"github.com/biztos/frostedmd"
)

var STD_USAGE = `TESTING Cmd

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
  -m, --meta        Exclude the content from the meta block.
  -f, --force       Do not abort on errors (log them to STDERR).
  -s, --silent      Do not print error messages.
  -t, --test        Parse file but do not print any output on success.
  --license         Print the software license.
`

func TestNewCmd(t *testing.T) {

	assert := assert.New(t)

	cmd := frostedmd.NewCmd("testing", "1.1.0", "anything")

	assert.Equal("testing", cmd.Name, "Name sticks")
	assert.Equal("1.1.0", cmd.Version, "Version sticks")
	assert.Equal("anything", cmd.Usage, "Usage sticks")

	assert.IsType(os.Exit, cmd.ExitFunction, "ExitFunction as expected")
	assert.Equal(os.Stdout, cmd.Stdout, "Stdout as expected")
	assert.Equal(os.Stderr, cmd.Stderr, "ExitFunction as expected")

}

func TestFail_SimpleError(t *testing.T) {

	assert := assert.New(t)

	cmd := frostedmd.NewCmd("testing", "1.1.0", "anything")
	var b bytes.Buffer
	cmd.Stderr = bufio.NewWriter(&b)
	exited := -1
	cmd.ExitFunction = func(c int) {
		exited = c
	}

	cmd.Fail(errors.New("anything"))
	assert.Equal("anything", b.String(), "error passed as-is")
	assert.Equal(frostedmd.CMD_OTHER_ERROR, exited,
		"exited with 'other' error")

}

func TestLicenseOption(t *testing.T) {

	assert := assert.New(t)

	os.Args = []string{"testing", "--license"}

	cmd := frostedmd.NewCmd("testing", "1.1.0", STD_USAGE)

	var b bytes.Buffer
	cmd.Stdout = bufio.NewWriter(&b)
	exited := -1
	cmd.ExitFunction = func(c int) {
		exited = c
	}
	err := cmd.SetOptions()
	assert.Nil(err, "no error from SetOptions with --license")
	assert.Regexp("^SOFTWARE LICENSES", b.String(), "looks license-y")
	assert.Equal(0, exited, "exited with success")
}
