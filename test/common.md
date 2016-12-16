# Complex Example

```yaml
# I am the Meta Block!
OneThing: This!
```

Here we will try to exercise all the functions of the `blackfriday.Renderer`
interface.  That is, *all* the Markdown we know how to convert to HTML.

## Standard Block-Level Elements

Paragraphs, obviously. And lists:

* This
* That

(Ah, one must have something between lists.  Bug?)

1. First
2. Second

Rule:

___

> Blockquote.
> 
> -- the blockheads

## Inline Elements

We have **bold** and *emphasis* and ~~strikethrough~~ and what else?

Conversions of "quotes" -- and 'quotes' -- and dashes!

No intra_emphasis, because we like C-style variables!

## Verbatim HTML

<span style="color:red">Like this</span> (bad idea, really).

### Links

[Great Art](http://kevinfrost.com/)
for [reference][kfcom]
or for [awesome].

[kfcom]: http://kevinfrost.com/
[awesome]: http://kevinfrost.com/

### Images

![random thing](http://lv8.biztos.com/20040212/fff.gif "Random Thing! LV8!")

## Autolinking

Search here: https://duckduckgo.com/

## Fenced Code Blocks

Language is optional:

```
echo "unknown lang"
```

But useful:

```perl
die 'hoooha';
```

And apparently can be spaced:

``` go
fmt.Println("ok by me")
```

Of course the old style must work as well:

    10 PRINT "I AM HERE "
    20 GOTO 10

## LaTeX-style dash parsing

This is where one differentiates between the "en dash" -- *ndash* -- and
the "em dash" --- *mdash* --- by the number of dashes used, i.e.
`--` vs `---`.

For a visual check:

* - dash
* -- ndash (as in: "1970-2016")
* --- mdash (as in: "here --- or there!")

Oddly, it seems to be built into Blackfriday, and not an option per se.

## Tables

For very small tables, this does of course rock.  Stick that in your
`asciidoc` and smoke it.

Ingredient                      | Cost
--------------------------------|------
South American Filet Mignon     | 7000
Ginger                          |  500

## Definition Lists

Fröccs
: 2:1 wine to soda water unless otherwise specified; cf. Kisfröccs.

## Smart Fractions

1/2 of 1/2 is usually 1/4...

## Not in Common

### Hard Line Breaks

For your poetic
Needs and your poetry
Feeds

### Footnotes

Apparently we are doomed.[^1]

[^1]: Doom being relative.
