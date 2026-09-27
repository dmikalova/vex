package web

import (
	"strconv"

	"github.com/maxence-charriere/go-app/v11/pkg/app"
)

// This file renders the Iconography pass's glyph lines (icon.go) into the DOM as
// the card's Icon strip. The strip is a thin band — about 1.5 text lines tall —
// that sits between the stat line and the trait line. Each ability is one line: a
// trigger label followed by its composed glyphs. A glyph is a reused SVG asset
// tinted and decorated by its Target's shape, or, for a mechanic not yet drawn, a
// small text chip so the strip still reads.

// iconStrip renders a card's glyph lines into the .card-icons band. Abilities flow
// on one line separated by a "|" divider; the inner .card-icons-fit span is scaled
// horizontally by cardFitScript when it is too wide, so a busy card squeezes to
// fit rather than wrapping into a second line or clipping.
func iconStrip(lines []glyphLine) app.UI {
	items := make([]app.UI, 0, len(lines)*2)
	for i := range lines {
		if i > 0 {
			items = append(items, app.Span().Class("card-icon-sep").Text("|"))
		}
		items = append(items, iconLine(lines[i]))
	}
	return app.Div().Class("card-icons").Body(
		app.Span().Class("card-icons-fit").Body(items...),
	)
}

// iconLine renders one ability: its trigger glyph(s) and the glyphs it composes.
// Several triggers sharing one effect (Play/Fight/Reap) render as their glyphs in
// a row with no separator, the icon counterpart of the "Play/Fight/Reap:" line. A
// trailing colon follows the trigger glyph(s), mirroring the printed "Play:".
func iconLine(l glyphLine) app.UI {
	return app.Span().Class("card-icon-line").Body(
		app.Range(l.triggers).Slice(func(i int) app.UI {
			return icon(l.triggers[i], "card-glyph-trigger", "icon-outline")
		}),
		app.If(len(l.triggers) > 0, func() app.UI {
			return app.Span().Class("card-glyph-trigger-colon").Text(":")
		}),
		app.Range(l.glyphs).Slice(func(i int) app.UI {
			return iconGlyph(l.glyphs[i])
		}),
	)
}

// iconGlyph renders a single glyph: a leading result-gate arrow when set, then the
// quantity numeral (before the icon, as the card prints it — "3 damage"), then the
// asset with its decoration classes. Every glyph is an icon — unmapped mechanics
// use the abstract glyph, never words.
func iconGlyph(g glyph) app.UI {
	return app.Span().Class("card-glyph").Body(
		app.If(g.arrow, func() app.UI {
			return app.Span().Class("card-glyph-arrow").Text("→")
		}),
		app.If(g.qty != 0, func() app.UI {
			return app.Span().Class("card-glyph-qty").Text(strconv.Itoa(g.qty))
		}),
		app.If(g.qty == 0 && g.text != "", func() app.UI {
			return app.Span().Class("card-glyph-qty").Text(g.text)
		}),
		app.Span().Class(cx("card-glyph-icon", decorClass(g.decor))).Body(
			icon(g.asset, "card-glyph-img", "icon-outline"),
			app.If(g.sigil != sigilNone, func() app.UI {
				return sigilMark(g.sigil)
			}),
		),
	)
}

// sigilMark renders a glyph's verb as the corner mark composited onto it. A verb
// is either a text mark or a small asset; the choice lives here rather than on the
// glyph, so promoting a mark to art touches no family file (ADR 0047).
//
// The tint is a second cue, never the only one: each verb has its own silhouette,
// following the precedent rarityMark set so the strip stays legible to
// colour-blind players.
func sigilMark(s sigil) app.UI {
	cls := cx("card-glyph-sigil", "card-glyph-sigil--"+sigilName(s))
	if stem := sigilAsset(s); stem != "" {
		return app.Span().Class(cls).Body(icon(stem, "card-glyph-sigil-img", "icon-outline"))
	}
	return app.Span().Class(cls).Text(sigilText(s))
}

// sigilName is a sigil's CSS modifier stem, which is also how a test names one.
func sigilName(s sigil) string {
	switch s {
	case sigilGain:
		return "gain"
	case sigilLose:
		return "lose"
	case sigilSteal:
		return "steal"
	case sigilExalt:
		return "exalt"
	case sigilReveal:
		return "reveal"
	case sigilShuffle:
		return "shuffle"
	case sigilMay:
		return "may"
	case sigilGrafted:
		return "grafted"
	case sigilNone:
		return ""
	}
	return ""
}

// sigilAsset is the asset stem a verb draws, or "" when the verb draws a text
// mark instead (sigilText).
func sigilAsset(s sigil) string {
	switch s {
	case sigilSteal:
		return "sigil-steal"
	case sigilExalt:
		return "sigil-exalt"
	case sigilReveal:
		return "sigil-reveal"
	case sigilShuffle:
		return "glyph-swap"
	case sigilGrafted:
		return "glyph-action"
	case sigilNone, sigilGain, sigilLose, sigilMay:
		return ""
	}
	return ""
}

// sigilText is the typographic mark a verb draws when it has no asset.
func sigilText(s sigil) string {
	switch s {
	case sigilGain:
		return "+"
	case sigilLose:
		return "−"
	case sigilMay:
		return "?"
	case sigilNone, sigilSteal, sigilExalt, sigilReveal, sigilShuffle, sigilGrafted:
		return ""
	}
	return ""
}

// decorClass turns a glyph's decoration flags into the space-joined CSS classes
// that tint and mark it (enemy/friendly tint, each-stack, chosen outline, self
// marker).
func decorClass(d decor) string {
	cls := make([]string, 0, 3)
	if d&decorEnemy != 0 {
		cls = append(cls, "card-glyph--enemy")
	}
	if d&decorFriendly != 0 {
		cls = append(cls, "card-glyph--friendly")
	}
	if d&decorEach != 0 {
		cls = append(cls, "card-glyph--each")
	}
	if d&decorChosen != 0 {
		cls = append(cls, "card-glyph--chosen")
	}
	if d&decorThis != 0 {
		cls = append(cls, "card-glyph--this")
	}
	return cx(cls...)
}
