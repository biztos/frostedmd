package frostedmd

import (
	"bytes"
	"github.com/russross/blackfriday"
)

// Our special renderer is not exposed.
type fmdRenderer struct {
	blocks      int
	metaAtEnd   bool
	metaBuffer  bytes.Buffer
	haveMeta    bool
	metaBytes   []byte
	metaLang    string
	headerTitle string
	bfRenderer  blackfriday.Renderer // Blackfriday's renderer
}

// This is a bit goofy but it helps us support meta-at-end, which is useful
// (in theory) for cases where you have an larger dataset embedded in your
// page.
func (r *fmdRenderer) incrementBlocks(out *bytes.Buffer) {

	r.blocks++
	if r.metaAtEnd && r.haveMeta {
		// We have the meta, but it's not the last block element, so put it
		// back into the main output buffer.
		if out.Len() > 0 {
			out.WriteByte('\n')
		}
		out.Write(r.metaBuffer.Bytes())
		r.haveMeta = false
		r.metaBytes = []byte{}
		r.metaLang = ""
		r.metaBuffer.Reset()
	}

}

// The renderer implments the blackfriday.Renderer interface, mostly via
// pass-through calls.

// block-level callbacks
func (r *fmdRenderer) BlockCode(out *bytes.Buffer, text []byte, lang string) {

	// If we are looking for the meta block at the end, any block could be it.
	if r.metaAtEnd {
		r.haveMeta = true
		r.metaBytes = text
		r.metaLang = lang
		r.metaBuffer.Reset()
		r.bfRenderer.BlockCode(&r.metaBuffer, text, lang)
		return
	}
	// We may already have one block, but only if it's the header title.
	if r.blocks == 0 || (r.blocks == 1 && r.headerTitle != "") {
		if !r.haveMeta {
			r.haveMeta = true
			r.metaBytes = text
			r.metaLang = lang
			return
		}
	}
	r.incrementBlocks(out)
	r.bfRenderer.BlockCode(out, text, lang)
}
func (r *fmdRenderer) BlockQuote(out *bytes.Buffer, text []byte) {
	r.incrementBlocks(out)
	r.bfRenderer.BlockQuote(out, text)
}
func (r *fmdRenderer) BlockHtml(out *bytes.Buffer, text []byte) {
	r.incrementBlocks(out)
	r.bfRenderer.BlockHtml(out, text)
}
func (r *fmdRenderer) Header(out *bytes.Buffer, text func() bool, level int, id string) {

	if r.blocks == 0 && r.headerTitle == "" {
		marker := out.Len()
		text()
		r.headerTitle = string(out.Bytes()[marker:])
		out.Truncate(marker)
	}

	r.incrementBlocks(out)
	r.bfRenderer.Header(out, text, level, id)
}
func (r *fmdRenderer) HRule(out *bytes.Buffer) {
	r.incrementBlocks(out)
	r.bfRenderer.HRule(out)
}
func (r *fmdRenderer) List(out *bytes.Buffer, text func() bool, flags int) {
	r.incrementBlocks(out)
	r.bfRenderer.List(out, text, flags)
}
func (r *fmdRenderer) ListItem(out *bytes.Buffer, text []byte, flags int) {
	r.incrementBlocks(out)
	r.bfRenderer.ListItem(out, text, flags)
}
func (r *fmdRenderer) Paragraph(out *bytes.Buffer, text func() bool) {
	r.incrementBlocks(out)
	r.bfRenderer.Paragraph(out, text)
}
func (r *fmdRenderer) Table(out *bytes.Buffer, header []byte, body []byte, columnData []int) {
	r.incrementBlocks(out)
	r.bfRenderer.Table(out, header, body, columnData)
}
func (r *fmdRenderer) TableRow(out *bytes.Buffer, text []byte) {
	r.incrementBlocks(out)
	r.bfRenderer.TableRow(out, text)
}
func (r *fmdRenderer) TableHeaderCell(out *bytes.Buffer, text []byte, flags int) {
	r.incrementBlocks(out)
	r.bfRenderer.TableHeaderCell(out, text, flags)
}
func (r *fmdRenderer) TableCell(out *bytes.Buffer, text []byte, flags int) {
	r.incrementBlocks(out)
	r.bfRenderer.TableCell(out, text, flags)
}
func (r *fmdRenderer) Footnotes(out *bytes.Buffer, text func() bool) {
	r.incrementBlocks(out)
	r.bfRenderer.Footnotes(out, text)
}
func (r *fmdRenderer) FootnoteItem(out *bytes.Buffer, name, text []byte, flags int) {
	r.incrementBlocks(out)
	r.bfRenderer.FootnoteItem(out, name, text, flags)
}
func (r *fmdRenderer) TitleBlock(out *bytes.Buffer, text []byte) {
	r.incrementBlocks(out)
	r.bfRenderer.TitleBlock(out, text)
}

// Span-level callbacks
func (r *fmdRenderer) AutoLink(out *bytes.Buffer, link []byte, kind int) {
	r.bfRenderer.AutoLink(out, link, kind)
}
func (r *fmdRenderer) CodeSpan(out *bytes.Buffer, text []byte) {
	r.bfRenderer.CodeSpan(out, text)
}
func (r *fmdRenderer) DoubleEmphasis(out *bytes.Buffer, text []byte) {
	r.bfRenderer.DoubleEmphasis(out, text)
}
func (r *fmdRenderer) Emphasis(out *bytes.Buffer, text []byte) {
	r.bfRenderer.Emphasis(out, text)
}
func (r *fmdRenderer) Image(out *bytes.Buffer, link []byte, title []byte, alt []byte) {
	r.bfRenderer.Image(out, link, title, alt)
}
func (r *fmdRenderer) LineBreak(out *bytes.Buffer) {
	r.bfRenderer.LineBreak(out)
}
func (r *fmdRenderer) Link(out *bytes.Buffer, link []byte, title []byte, content []byte) {
	r.bfRenderer.Link(out, link, title, content)
}
func (r *fmdRenderer) RawHtmlTag(out *bytes.Buffer, tag []byte) {
	r.bfRenderer.RawHtmlTag(out, tag)
}
func (r *fmdRenderer) TripleEmphasis(out *bytes.Buffer, text []byte) {
	r.bfRenderer.TripleEmphasis(out, text)
}
func (r *fmdRenderer) StrikeThrough(out *bytes.Buffer, text []byte) {
	r.bfRenderer.StrikeThrough(out, text)
}
func (r *fmdRenderer) FootnoteRef(out *bytes.Buffer, ref []byte, id int) {
	r.bfRenderer.FootnoteRef(out, ref, id)
}

// Low-level callbacks
func (r *fmdRenderer) Entity(out *bytes.Buffer, entity []byte) {
	r.bfRenderer.Entity(out, entity)
}
func (r *fmdRenderer) NormalText(out *bytes.Buffer, text []byte) {
	r.bfRenderer.NormalText(out, text)
}

// Header and footer
func (r *fmdRenderer) DocumentHeader(out *bytes.Buffer) {
	r.bfRenderer.DocumentHeader(out)
}
func (r *fmdRenderer) DocumentFooter(out *bytes.Buffer) {
	r.bfRenderer.DocumentFooter(out)
}

func (r *fmdRenderer) GetFlags() int {
	return r.bfRenderer.GetFlags()
}
