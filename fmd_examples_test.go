// frostedmd/fmd_examples_test.go -- examples for Frosted Markdown.
// ------------------------------

package frostedmd_test

import (
	"fmt"

	"github.com/biztos/frostedmd"
)

func Example() {

	// NOTE: due to a godoc bug, the following text may not appear correctly.
	// The "# Meta" line and the "Tags" line are both indented (four spaces).
	// https://github.com/golang/go/issues/18446

	// The easiest way to get things done:
	input := `# My Markdown

    # Meta:
    Tags: ["fee","fi","foe"]

Obscurantism threatens clean data.
`

	res, err := frostedmd.MarkdownCommon([]byte(input))
	if err != nil {
		panic(err)
	}
	mm := res.Meta
	fmt.Println("Title:", mm["Title"])
	fmt.Println("Tags:", mm["Tags"])
	fmt.Println("HTML:", string(res.Content))

	// Output:
	// Title: My Markdown
	// Tags: [fee fi foe]
	// HTML: <h1>My Markdown</h1>
	//
	// <p>Obscurantism threatens clean data.</p>
}

func ExampleNew() {

	input := `# Lots of Data

Here we have a full dataset, which (for instance) a template engine
will turn into something cool and dynamic.  Thus we put it at the end
so we can read our nice summary using the *head* command.

    {
        "datasets": {
             "numbers": [11,22,33,44,55,66],
             "letters": ["a","B","ß","í"]
        }
    }

`

	parser := frostedmd.New()
	parser.MetaAtEnd = true
	res, err := parser.Parse([]byte(input))
	if err != nil {
		panic(err)
	}
	mm := res.Meta
	fmt.Println("Title:", mm["Title"])
	fmt.Println("HTML:", string(res.Content))

	// Order within a map is random in Go, so let's make it explicit.
	fmt.Println("Data sets:")
	if ds, ok := mm["datasets"].(map[string]interface{}); ok {
		fmt.Println("  numbers:", ds["numbers"])
		fmt.Println("  letters:", ds["letters"])
	} else {
		fmt.Printf("NOT A MAP: %T\n", mm["datasets"])
	}

	// Output:
	// Title: Lots of Data
	// HTML: <h1>Lots of Data</h1>
	//
	// <p>Here we have a full dataset, which (for instance) a template engine
	// will turn into something cool and dynamic.  Thus we put it at the end
	// so we can read our nice summary using the <em>head</em> command.</p>
	//
	// Data sets:
	//   numbers: [11 22 33 44 55 66]
	//   letters: [a B ß í]
}
