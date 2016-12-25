// cmd_test.go

package frostedmd_test

import (
	// Standard library:
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	// Third-party:
	"github.com/biztos/testig" // well, first-party in a way...
	"github.com/stretchr/testify/assert"

	// Under test:
	"github.com/biztos/frostedmd"
)

var StandardUsage = `TESTING Cmd

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

func Test_NewCmd(t *testing.T) {

	assert := assert.New(t)

	cmd := frostedmd.NewCmd("testing", "1.1.0", "anything")

	assert.Equal("testing", cmd.Name, "Name sticks")
	assert.Equal("1.1.0", cmd.Version, "Version sticks")
	assert.Equal("anything", cmd.Usage, "Usage sticks")

	assert.IsType(os.Exit, cmd.Exit, "Exit as expected")
	assert.Equal(os.Stdout, cmd.Stdout, "Stdout as expected")
	assert.Equal(os.Stderr, cmd.Stderr, "Exit as expected")

}

func Test_Fail_SimpleError(t *testing.T) {

	assert := assert.New(t)

	cmd := frostedmd.NewCmd("testing", "1.1.0", "anything")
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	cmd.Stderr = w
	exited := -1
	cmd.Exit = func(c int) {
		exited = c
	}

	cmd.Fail(errors.New("anything"))
	w.Flush()
	assert.Equal("anything\n", b.String(), "error passed as-is")
	assert.Equal(frostedmd.CMD_OTHER_ERROR, exited,
		"exited with 'other' error")

}

func Test_Fail_CmdError(t *testing.T) {

	assert := assert.New(t)

	cmd := frostedmd.NewCmd("testing", "1.1.0", "anything")
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	cmd.Stderr = w
	exited := -1
	cmd.Exit = func(c int) {
		exited = c
	}

	err := frostedmd.CmdError{
		Err:  errors.New("anything"),
		Code: 123,
	}
	cmd.Fail(err)
	w.Flush()
	assert.Equal("anything\n", b.String(), "error passed as-is w/o file")
	assert.Equal(123, exited,
		"exited with passed error code")

}

func Test_Fail_CmdError_WithFile(t *testing.T) {

	assert := assert.New(t)

	cmd := frostedmd.NewCmd("testing", "1.1.0", "anything")
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	cmd.Stderr = w
	exited := -1
	cmd.Exit = func(c int) {
		exited = c
	}

	err := frostedmd.CmdError{
		File: "somefile",
		Err:  errors.New("anything"),
		Code: 123,
	}
	cmd.Fail(err)
	w.Flush()
	assert.Equal("somefile: anything\n", b.String(),
		"error includes filename")
	assert.Equal(123, exited,
		"exited with passed error code")

}

func Test_Fail_CmdError_Silent(t *testing.T) {

	assert := assert.New(t)

	cmd := frostedmd.NewCmd("testing", "1.1.0", "anything")
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	cmd.Stderr = w
	exited := -1
	cmd.Exit = func(c int) {
		exited = c
	}

	err := frostedmd.CmdError{
		Silent: true,
		Err:    errors.New("anything"),
		Code:   123,
	}
	cmd.Fail(err)
	w.Flush()
	assert.Equal("", b.String(), "no error printed")
	assert.Equal(123, exited, "exited with passed error code")

}

func Test_Fail_CmdError_Force(t *testing.T) {

	assert := assert.New(t)

	cmd := frostedmd.NewCmd("testing", "1.1.0", "anything")
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	cmd.Stderr = w
	exited := -1
	cmd.Exit = func(c int) {
		exited = c
	}

	err := frostedmd.CmdError{
		Force: true,
		Err:   errors.New("anything"),
		Code:  123,
	}
	cmd.Fail(err)
	w.Flush()
	assert.Equal("anything\n", b.String(), "error printed")
	assert.Equal(-1, exited, "did not exit")

}

func Test_Fail_CmdError_ZeroCode(t *testing.T) {

	assert := assert.New(t)

	cmd := frostedmd.NewCmd("testing", "1.1.0", "anything")
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	cmd.Stderr = w
	exited := -1
	cmd.Exit = func(c int) {
		exited = c
	}

	err := frostedmd.CmdError{
		Err: errors.New("anything"),
	}
	cmd.Fail(err)
	w.Flush()
	assert.Equal("anything\n", b.String(), "error printed")
	assert.Equal(frostedmd.CMD_OTHER_ERROR, exited,
		"exited with 'other' code")

}

func Test_SetOptions_Minimalist(t *testing.T) {

	assert := assert.New(t)

	// It should be possible to run SetOptions without any of the options
	// directly supported.
	usage := `testing

Usage:
  testing DIR
`

	os.Args = []string{"testing", "could-be-a-dir"}
	cmd := frostedmd.NewCmd("testing", "1.1.0", usage)
	err := cmd.SetOptions()
	assert.Nil(err, "no error on minimalist SetOptions")

	exp := &frostedmd.CmdOptions{Format: "json"}
	assert.EqualValues(exp, cmd.Options, "Options set to default values")

}

func Test_SetOptions_Maximalist(t *testing.T) {

	assert := assert.New(t)

	// Standard usage, with everything possible set.
	os.Args = []string{
		"testing",
		"--yaml",
		"--indent",
		"--nobase64",
		"--content",
		"--force",
		"--silent",
		"--test",
		"somefile",
	}

	cmd := frostedmd.NewCmd("testing", "1.1.0", StandardUsage)
	err := cmd.SetOptions()
	assert.Nil(err, "no error on minimalist SetOptions")

	exp := &frostedmd.CmdOptions{
		Format:      "yaml",
		Indent:      true,
		NoBase64:    true,
		ContentOnly: true,
		Force:       true,
		Silent:      true,
		Test:        true,
		File:        "somefile",
	}
	assert.EqualValues(exp, cmd.Options, "Options set to override values")

}

func Test_SetOptions_OutputContadiction(t *testing.T) {

	assert := assert.New(t)

	// Standard usage, with everything possible set.
	os.Args = []string{
		"testing",
		"--content",
		"--meta",
		"somefile",
	}

	cmd := frostedmd.NewCmd("testing", "1.1.0", StandardUsage)
	err := cmd.SetOptions()
	if assert.Error(err, "error set") {
		assert.Equal("--meta and --content are mutually exclusive.",
			err.Error(), "error string as expected")
		if assert.IsType(frostedmd.CmdError{}, err) {
			// This is rather awkward considering the above...
			e, _ := err.(frostedmd.CmdError)
			assert.Equal(frostedmd.CMD_OPTIONS_ERROR, e.Code,
				"error code is 'options'")
		}
	}

}

func Test_SetOptions_LicenseOption(t *testing.T) {

	assert := assert.New(t)

	os.Args = []string{"testing", "--license"}

	cmd := frostedmd.NewCmd("testing", "1.1.0", StandardUsage)

	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	cmd.Stdout = w
	exited := -1
	cmd.Exit = func(c int) {
		exited = c
	}
	err := cmd.SetOptions()
	w.Flush()
	assert.Nil(err, "no error from SetOptions with --license")
	assert.Regexp("^SOFTWARE LICENSES", b.String(), "looks license-y")
	assert.Equal(0, exited, "exited with success")
}

func Test_SetOptions_BadUsageForBools(t *testing.T) {

	// One is probably enough for the coverage metric but it's silly to not
	// test them all, as long as we're here.
	//
	// NOTE: currently we have no equivalent for the string args because they
	// are not options (actually it's only the one, FILE) and it seems not
	// possible to miscast that in docopt.
	shouldBool := []string{
		"--json",
		"--yaml",
		"--indent",
		"--nobase64",
		"--force",
		"--silent",
		"--test",
		"--content",
		"--meta",
		"--license",
	}
	for _, key := range shouldBool {
		usage := fmt.Sprintf("t\n\nUsage:\n  t [%s=FOO] FILE\n", key)
		exp := "Bad designation of bool opt in docopt usage: " + key
		os.Args = []string{
			"t",
			key + "=123",
			"somefile",
		}
		cmd := frostedmd.NewCmd("testing", "1.1.0", usage)
		f := func() { cmd.SetOptions() }
		testig.AssertPanicsWith(t, f, exp, "panics as expected for "+key)
	}

}

func Test_ParseFile_FileError(t *testing.T) {

	assert := assert.New(t)

	cmd := frostedmd.NewCmd("testing", "1.1.0", StandardUsage)
	cmd.Options = &frostedmd.CmdOptions{File: "no-such-file-here"}
	err := cmd.ParseFile()
	if assert.Error(err) {
		assert.Equal(
			"no-such-file-here: open no-such-file-here: "+
				"no such file or directory",
			err.Error(), "error as expected")
		if assert.IsType(frostedmd.CmdError{}, err, "error has our type") {
			e, _ := err.(frostedmd.CmdError)
			assert.Equal(frostedmd.CMD_FILE_ERROR, e.Code,
				"error has file error exit code")

		}
	}
}

func Test_ParseFile_ParseError(t *testing.T) {

	assert := assert.New(t)

	file := filepath.Join("test", "broken.md")
	cmd := frostedmd.NewCmd("testing", "1.1.0", StandardUsage)
	cmd.Options = &frostedmd.CmdOptions{File: file}
	err := cmd.ParseFile()
	if assert.Error(err) {
		assert.Regexp("^test.*broken.*yaml", err.Error(), "error as expected")
		if assert.IsType(frostedmd.CmdError{}, err, "error has our type") {
			e, _ := err.(frostedmd.CmdError)
			assert.Equal(frostedmd.CMD_PARSE_ERROR, e.Code,
				"error has parse error exit code")

		}
	}
}

func Test_ParseFile_Success(t *testing.T) {

	assert := assert.New(t)

	file := filepath.Join("test", "simple.md")
	expMeta := map[string]interface{}{
		"Title":       "FMD FTW",
		"Description": "Simple is as simple does.",
		"Tags":        []interface{}{"fmd", "golang", "nerdery"},
	}
	expContent := "<h1>Simple FMD</h1>\n\n<p>Good enough for me.</p>\n"

	cmd := frostedmd.NewCmd("testing", "1.1.0", StandardUsage)
	cmd.Options = &frostedmd.CmdOptions{File: file}
	err := cmd.ParseFile()
	if assert.Nil(err, "no error from ParseFile") {
		if assert.NotNil(cmd.Result, "Result was set") {
			assert.Equal(expMeta, cmd.Result.Meta, "Result Meta as expected")
			assert.Equal(expContent, string(cmd.Result.Content),
				"Result Content as expected")
		}
	}
}

func Test_PrintResult_NoResult(t *testing.T) {

	assert := assert.New(t)

	cmd := frostedmd.NewCmd("testing", "1.1.0", StandardUsage)

	rec := testig.NewOutputRecorder()
	cmd.Stdout, cmd.Stderr = rec.Stdout, rec.Stderr

	err := cmd.PrintResult()
	assert.Nil(err, "no error on PrintResult for nil Result")
	assert.Equal("", rec.StdoutString(),
		"no standard output on PrintResult for nil Result")
	assert.Equal("", rec.StderrString(),
		"no standard error on PrintResult for nil Result")
}

func Test_PrintResult_TestMode(t *testing.T) {

	assert := assert.New(t)

	cmd := frostedmd.NewCmd("testing", "1.1.0", StandardUsage)
	cmd.Options = &frostedmd.CmdOptions{Test: true}
	cmd.Result = &frostedmd.ParseResult{Content: []byte("anything")}
	rec := testig.NewOutputRecorder()
	cmd.Stdout, cmd.Stderr = rec.Stdout, rec.Stderr

	err := cmd.PrintResult()
	assert.Nil(err, "no error on PrintResult for Test option")
	assert.Equal("", rec.StdoutString(),
		"no standard output on PrintResult for Test option")
	assert.Equal("", rec.StderrString(),
		"no standard error on PrintResult Test option")
}

func Test_PrintResult_JSON(t *testing.T) {

	assert := assert.New(t)

	cmd := frostedmd.NewCmd("testing", "1.1.0", StandardUsage)
	cmd.Result = &frostedmd.ParseResult{
		Meta:    map[string]interface{}{"foo": 123},
		Content: []byte("here be content"),
	}
	exp := `{"meta":{"foo":123},"content":"aGVyZSBiZSBjb250ZW50"}
`
	cmd.Options = &frostedmd.CmdOptions{} // json is the default

	rec := testig.NewOutputRecorder()
	cmd.Stdout, cmd.Stderr = rec.Stdout, rec.Stderr

	err := cmd.PrintResult()
	assert.Nil(err, "no error on PrintResult")
	assert.Equal(exp, rec.StdoutString(), "json on stdout")
	assert.Equal("", rec.StderrString(), "no standard error")

}

func Test_PrintResult_JSON_NoBase64(t *testing.T) {

	assert := assert.New(t)

	cmd := frostedmd.NewCmd("testing", "1.1.0", StandardUsage)
	cmd.Result = &frostedmd.ParseResult{
		Meta:    map[string]interface{}{"foo": 123},
		Content: []byte("here be content"),
	}

	// Hmm, is the field order deterministic, and thus alphabetic except for
	// byte slices that are presumed long?  Could be... would be a nice trick,
	// but what if it's just randomness here?
	exp := `{"content":"here be content","meta":{"foo":123}}
`
	cmd.Options = &frostedmd.CmdOptions{NoBase64: true} // json is the default

	rec := testig.NewOutputRecorder()
	cmd.Stdout, cmd.Stderr = rec.Stdout, rec.Stderr

	err := cmd.PrintResult()
	assert.Nil(err, "no error on PrintResult")
	assert.Equal(exp, rec.StdoutString(), "json on stdout")
	assert.Equal("", rec.StderrString(), "no standard error")

}

func Test_PrintResult_JSON_Indent(t *testing.T) {

	assert := assert.New(t)

	cmd := frostedmd.NewCmd("testing", "1.1.0", StandardUsage)
	cmd.Result = &frostedmd.ParseResult{
		Meta:    map[string]interface{}{"foo": 123},
		Content: []byte("here be content"),
	}
	exp := `{
    "meta": {
        "foo": 123
    },
    "content": "aGVyZSBiZSBjb250ZW50"
}
`
	cmd.Options = &frostedmd.CmdOptions{Indent: true} // json is the default

	rec := testig.NewOutputRecorder()
	cmd.Stdout, cmd.Stderr = rec.Stdout, rec.Stderr

	err := cmd.PrintResult()
	assert.Nil(err, "no error on PrintResult")
	assert.Equal(exp, rec.StdoutString(), "json on stdout")
	assert.Equal("", rec.StderrString(), "no standard error")

}

func Test_PrintResult_ContentOnly(t *testing.T) {

	assert := assert.New(t)

	cmd := frostedmd.NewCmd("testing", "1.1.0", StandardUsage)
	cmd.Result = &frostedmd.ParseResult{
		Meta:    map[string]interface{}{"foo": 123},
		Content: []byte("here be content"),
	}

	cmd.Options = &frostedmd.CmdOptions{ContentOnly: true}

	rec := testig.NewOutputRecorder()
	cmd.Stdout, cmd.Stderr = rec.Stdout, rec.Stderr

	err := cmd.PrintResult()
	assert.Nil(err, "no error on PrintResult")
	assert.Equal("here be content\n", rec.StdoutString(), "content on stdout")
	assert.Equal("", rec.StderrString(), "no standard error")

}

func Test_PrintResult_MetaOnly_JSON(t *testing.T) {

	assert := assert.New(t)

	cmd := frostedmd.NewCmd("testing", "1.1.0", StandardUsage)
	cmd.Result = &frostedmd.ParseResult{
		Meta:    map[string]interface{}{"foo": 123},
		Content: []byte("here be content"),
	}

	cmd.Options = &frostedmd.CmdOptions{MetaOnly: true}

	rec := testig.NewOutputRecorder()
	cmd.Stdout, cmd.Stderr = rec.Stdout, rec.Stderr

	err := cmd.PrintResult()
	assert.Nil(err, "no error on PrintResult")
	assert.Equal("{\"foo\":123}\n", rec.StdoutString(), "meta on stdout")
	assert.Equal("", rec.StderrString(), "no standard error")

}

func Test_PrintResult_MetaOnly_YAML(t *testing.T) {

	assert := assert.New(t)

	cmd := frostedmd.NewCmd("testing", "1.1.0", StandardUsage)
	cmd.Result = &frostedmd.ParseResult{
		Meta:    map[string]interface{}{"foo": 123},
		Content: []byte("here be content"),
	}

	cmd.Options = &frostedmd.CmdOptions{MetaOnly: true, Format: "yaml"}

	rec := testig.NewOutputRecorder()
	cmd.Stdout, cmd.Stderr = rec.Stdout, rec.Stderr

	err := cmd.PrintResult()
	assert.Nil(err, "no error on PrintResult")
	assert.Equal("foo: 123\n\n", rec.StdoutString(), "meta on stdout")
	assert.Equal("", rec.StderrString(), "no standard error")

}

func Test_PrintResult_YAML(t *testing.T) {

	assert := assert.New(t)

	cmd := frostedmd.NewCmd("testing", "1.1.0", StandardUsage)
	cmd.Result = &frostedmd.ParseResult{
		Meta:    map[string]interface{}{"foo": 123},
		Content: []byte("here be content"),
	}

	cmd.Options = &frostedmd.CmdOptions{Format: "yaml"}

	rec := testig.NewOutputRecorder()
	cmd.Stdout, cmd.Stderr = rec.Stdout, rec.Stderr

	exp := `content: here be content
meta:
  foo: 123

`
	err := cmd.PrintResult()
	assert.Nil(err, "no error on PrintResult")
	assert.Equal(exp, rec.StdoutString(), "yaml on stdout")
	assert.Equal("", rec.StderrString(), "no standard error")

}

func Test_PrintResult_ErrorSerializingJSON(t *testing.T) {

	assert := assert.New(t)

	cmd := frostedmd.NewCmd("testing", "1.1.0", StandardUsage)
	cmd.Result = &frostedmd.ParseResult{
		Meta:    map[string]interface{}{"foo": func() {}},
		Content: []byte("here be content"),
	}

	cmd.Options = &frostedmd.CmdOptions{}

	rec := testig.NewOutputRecorder()
	cmd.Stdout, cmd.Stderr = rec.Stdout, rec.Stderr

	err := cmd.PrintResult()
	if assert.Error(err, "error on PrintResult") {
		if assert.IsType(frostedmd.CmdError{}, err) {
			assert.Equal("json: unsupported type: func()", err.Error(),
				"error string as expected")
			e, _ := err.(frostedmd.CmdError)
			assert.Equal(frostedmd.CMD_SERIALIZATION_ERROR, e.Code,
				"error code is 'serialization'")
		}
	}
	assert.Equal("", rec.StdoutString(), "no standard output")
	assert.Equal("", rec.StderrString(), "no standard error")

}

func Test_PrintResult_ErrorSerializingYAML(t *testing.T) {

	assert := assert.New(t)

	cmd := frostedmd.NewCmd("testing", "1.1.0", StandardUsage)
	cmd.Result = &frostedmd.ParseResult{
		Meta:    map[string]interface{}{"foo": func() {}},
		Content: []byte("here be content"),
	}

	cmd.Options = &frostedmd.CmdOptions{Format: "yaml"}

	rec := testig.NewOutputRecorder()
	cmd.Stdout, cmd.Stderr = rec.Stdout, rec.Stderr

	err := cmd.PrintResult()
	if assert.Error(err, "error on PrintResult") {
		if assert.IsType(frostedmd.CmdError{}, err) {
			assert.Equal("yaml error: cannot marshal type: func()",
				err.Error(), "error string as expected")
			e, _ := err.(frostedmd.CmdError)
			assert.Equal(frostedmd.CMD_SERIALIZATION_ERROR, e.Code,
				"error code is 'serialization'")
		}
	}
	assert.Equal("", rec.StdoutString(), "no standard output")
	assert.Equal("", rec.StderrString(), "no standard error")

}

func Test_Run_SetOptionsError(t *testing.T) {

	assert := assert.New(t)

	os.Args = []string{"test", "-j", "-y", "FILE"} // contradictory options
	cmd := frostedmd.NewCmd("testing", "1.1.0", StandardUsage)
	err := cmd.Run()
	if assert.Error(err, "error on Run") {

	}
	if assert.IsType(frostedmd.CmdError{}, err) {
		e, _ := err.(frostedmd.CmdError)
		assert.Equal(frostedmd.CMD_OPTIONS_ERROR, e.Code,
			"error code is 'options'")
	}

}

func Test_Run_ParseFileError(t *testing.T) {

	assert := assert.New(t)

	os.Args = []string{"test", "no-such-file-here"} // contradictory options
	cmd := frostedmd.NewCmd("testing", "1.1.0", StandardUsage)
	err := cmd.Run()
	if assert.Error(err, "error on Run") {
		if assert.IsType(frostedmd.CmdError{}, err) {
			e, _ := err.(frostedmd.CmdError)
			assert.Equal(frostedmd.CMD_FILE_ERROR, e.Code,
				"error code is 'file'")
		}
	}

}

func Test_Run_Success(t *testing.T) {

	assert := assert.New(t)

	os.Args = []string{"test", "-i", filepath.Join("test", "simple.md")}
	cmd := frostedmd.NewCmd("testing", "1.1.0", StandardUsage)
	rec := testig.NewOutputRecorder()
	cmd.Stdout, cmd.Stderr = rec.Stdout, rec.Stderr
	exp := `{
    "meta": {
        "Description": "Simple is as simple does.",
        "Tags": [
            "fmd",
            "golang",
            "nerdery"
        ],
        "Title": "FMD FTW"
    },
    "content": "PGgxPlNpbXBsZSBGTUQ8L2gxPgoKPHA+R29vZCBlbm91Z2ggZm9yIG1lLjwvcD4K"
}
`

	err := cmd.Run()
	assert.Nil(err, "no error on Run")
	assert.Equal(exp, rec.StdoutString(), "json on stdout")
	assert.Equal("", rec.StderrString(), "no standard error")

}
