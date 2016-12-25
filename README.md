# Frosted Markdown

**WARNING: this is alpha software and should be considered unstable.**

In other words, a work in progress.  Like most other software.

## Synopsis

```go
package main

import (
    "fmt"
    "github.com/biztos/frostedmd"
)

func main() {

	input := `# My Markdown

    # Meta:
    Tags: ["fee","fi","foe"]

Obscurantism threatens clean data.
`

	res, err := frostedmd.MarkdownCommon([]byte(input))
	if err != nil {
		panic(err)
	}
	mm := res.Meta()
	fmt.Println("Title:", mm["Title"])
	fmt.Println("Tags:", mm["Tags"])
	fmt.Println("HTML:", string(res.Content()))
}
```

Output:

    Title: My Markdown
    Tags: [fee fi foe]
    HTML: <h1>My Markdown</h1>

    <p>Obscurantism threatens clean data.</p>

## Description

Frosted Markdown is Markdown with useful extensions *and* tasty data
frosting on top, written in [Go][go].

[Markdown][wiki-md] is a lightweight plain-text markup language originally
written by [John Gruber][gruber].  It makes it easy to write blog posts and
similar things in plain-ish text, and convert the result to a fragment of
HTML that can be inserted into, say, your blog template.  And as you probably
know, it has become very popular.

The original Markdown feature set being intentionally minimalist, many
extenensions have been implemented by various authors -- and there being no
Markdown standard, it's a bit chaotic. At present Frosted Markdown relies on
the extensions supported by the excellent [BlackFriday][bf], plus of course
the "frosting."

The "frosting" (also known as "[icing][wiki-icing]") that makes Frosted
Markdown special is a code block (or any preformatted text block) either at
the *beginning* of the file, or *following the first heading.*

This block, called the Meta Block, is parsed and (usually) removed from the
rendered HTML fragment.  Thus you end up with a data structure extracted from
the Meta Block and an HTML fragment, instead of just the HTML fragment.

So, why is this useful?

* Data specific to the file is kept in the file, simplifying management.
* In non-Frosted contexts, e.g. editor previews, your data is still at hand.

It's important to remember that *all Markdown is Frosted Markdown* -- but
some of it doesn't have any frosting. If there is no Meta Block, or the
apparent Meta Block can't be parsed, then it is empty -- except for one
thing: Frosted Markdown will set the Title if the Markdown file leads with a
top-level heading.

Likewise, all Frosted Markdown is Markdown (with a caveat for extension
conflicts). Anything `frostedmd` can parse can be parsed just as well by
`blackfriday`.  You'll just expose any Meta Block as a code block.


[go]: https:/golang.org/
[gruber]: http://daringfireball.net/colophon/
[wiki-md]: https://en.wikipedia.org/wiki/Markdown
[wiki-icing]: https://en.wikipedia.org/wiki/Icing_(food)

## FAQ

### Why frosting?

Because it's sweet and goes on top!  And you normally eat the frosting first
and the cake later, right?

### Why not icing?

Because the author is from the United States, and also happens to be named
*Frost.*

### Can I put the Meta Block at the end instead?

Yes, but you have to tell the parser that's what you want.

### Is this production-ready?

Probably not.  Note the version number in the package.  Until it reaches
`1.0` the set of extensions and even the API itself may change.

### Is this a database?

No. If you're building something database-like off of Frosted Markdown --
which isn't as exotic an idea as you might think -- you will have to parse
all the files up front, and then do your databasey things with the set of
meta structures thus harvested.

## Licenses

Frosted Markdown is (c) Copyright 2016 Kevin A. Frost, with humble
acknowledgement of the various authors who did most of the real work (see
below).

Frosted Markdown itself is made available under the *BSD License* -- see the
file `LICENSE` in this repository. This license also applies to most standard
Go packages.

Note however that this package relies heavily on these additional packages,
which have their own licenses:

* [blackfriday][bf] by Russ Ross et al. -- [Simplified BSD License][bf-lic].
* [yaml][yaml] by Canonical et al. -- [Apache License 2.0][yaml-lic].
* [testify][testify] by Mat Ryer and Tyler Bunnell -- [MIT License][testify-lic].
* [docopt][docopt] by Keith Batten et al. -- [MIT License][docopt-lic].

[bf]: https://github.com/russross/blackfriday
[yaml]: https://github.com/go-yaml/yaml
[testify]: https://github.com/stretchr/testify
[docopt]: https://github.com/docopt/docopt.go
[bf-lic]: https://github.com/russross/blackfriday/blob/master/LICENSE.txt
[yaml-lic]: https://github.com/go-yaml/yaml/blob/v2/LICENSE
[testify-lic]: https://github.com/stretchr/testify/blob/master/LICENSE
[docopt-lic]: https://github.com/docopt/docopt.go/blob/master/LICENSE

## TODO

* Add footnotes to the default extension set.
* More documentation for the long-suffering laity.
* Better test coverage (and arguably better tests).  100% at minimum!
* LOTS more edge cases etc.
* At least CONSIDER less-canonical, but faster, parsing.
    * Specifically, just strip off the top of the file w/o MD-processing it.
* Uniform error handling, e.g. for parse failures in unknown meta blocks.
* Also sane handling of potentially innocent errors.
* Disable title guessing if not already possible (as an option).
* Vendor-in deps, once I grok how that works with Travis et al.

