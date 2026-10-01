package web

import (
	"strconv"

	"github.com/maxence-charriere/go-app/v11/pkg/app"

	"github.com/dmikalova/vex/internal/engine"
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

// iconGlyph renders a single glyph: its leading punctuation mark when set, then
// the quantity numeral (before the icon, as the card prints it — "3 damage"), then
// the asset with its decoration classes. Every glyph is an icon — unmapped
// mechanics use the abstract glyph, never words.
func iconGlyph(g glyph) app.UI {
	return app.Span().Class("card-glyph").Body(
		app.If(g.punct != punctNone, func() app.UI {
			return punctMark(g.punct)
		}),
		app.If(hasQtyBadge(g), func() app.UI {
			return qtyBadge(g)
		}),
		app.Span().Class(cx("card-glyph-icon", decorClass(g.decor))).Body(
			icon(g.asset, "card-glyph-img", "icon-outline"),
			app.If(g.sigil != sigilNone, func() app.UI {
				return sigilMark(g.sigil)
			}),
		),
	)
}

// hasQtyBadge reports whether a glyph prints anything in its quantity badge: a
// count or a bare comparison mark. A parity comparison carries no number at all —
// "while the opponent's pool is even" has nothing to compare against — so the
// mark alone is enough to draw the badge.
func hasQtyBadge(g glyph) bool {
	return g.qty != 0 || g.comparison != comparisonNone
}

// qtyBadge renders the numeral chip that precedes a glyph's icon — the "3" in "3
// damage" — with its comparison mark ahead of the number when the number is a
// threshold rather than a count, so "≥3 power" reads as "power 3 or more".
//
// The comparison rides this badge rather than taking a third badge axis beside
// decor and sigil: it qualifies a number the badge is already printing, and a
// strip where every glyph wears one more badge is noise (ADR 0047). One class
// covers every mark, because unlike a verb or a punctuation mark a comparison has
// no per-member styling to hang off a modifier — the marks differ by character,
// which is shape rather than hue.
func qtyBadge(g glyph) app.UI {
	return app.Span().Class("card-glyph-qty").Body(
		app.If(g.comparison != comparisonNone, func() app.UI {
			return app.Span().Class("card-glyph-cmp").Text(comparisonMark(g.comparison))
		}),
		app.Text(qtyBody(g)),
	)
}

// qtyBody is the badge's number, or "" when the badge carries a parity mark
// alone — "even" and "odd" compare against no number.
func qtyBody(g glyph) string {
	if g.qty == 0 {
		return ""
	}
	return strconv.Itoa(g.qty)
}

// comparisonMark is the mark a comparison draws on the quantity badge it
// qualifies. Every member of engine.Comparison gets one —
// TestComparisonMarkCoversEveryComparison walks engine.Comparisons() — so a
// comparison added to the engine cannot reach the strip undrawn.
//
// The three absolute comparisons take their mathematical signs. The two parity
// comparisons ignore Amount and print no number, so their mark is the whole badge
// and says what shape the quantity must have: 2n, or 2n+1. The two relative
// comparisons share ">": both mean "more than the other player's pool", and which
// pool is the subject is already carried by the glyph's enemy/friendly decor, so a
// second mark would encode the subject twice.
func comparisonMark(c engine.Comparison) string {
	switch c {
	case engine.AtLeast:
		return "≥"
	case engine.AtMost:
		return "≤"
	case engine.Exactly:
		return "="
	case engine.MoreThanYou, engine.MoreThanOpponent:
		return ">"
	case engine.Even:
		return "2n"
	case engine.Odd:
		return "2n+1"
	case comparisonNone:
		return ""
	}
	return ""
}

// punctMark renders the punctuation that opens a glyph's clause: the flow arrow,
// or one of the three structure marks (ADR 0047). It is drawn between glyphs
// rather than on one, so it is a sibling of the glyph's icon span, not a child of
// it the way sigilMark is.
//
// Every mark is typographic today. Promoting one to art is this function's
// business alone, exactly as it is sigilMark's for a verb.
func punctMark(p punct) app.UI {
	return app.Span().
		Class(cx("card-glyph-punct", "card-glyph-punct--"+punctName(p))).
		Text(punctText(p))
}

// punctName is a punctuation mark's CSS modifier stem, which is also how a test
// names one.
func punctName(p punct) string {
	switch p {
	case punctArrow:
		return "arrow"
	case punctGate:
		return "gate"
	case punctBranch:
		return "branch"
	case punctDivider:
		return "divider"
	case punctNone:
		return ""
	}
	return ""
}

// punctText is the character a punctuation mark draws.
//
// The gate and the branch are one shape in two fills, because they are the two
// ways a single gate can fall and reading them as a pair is the whole point; the
// divider is a slash, the plainest "or" the strip can draw, and deliberately not
// the "|" that separates two abilities.
func punctText(p punct) string {
	switch p {
	case punctArrow:
		return "→"
	case punctGate:
		return "◇"
	case punctBranch:
		return "◆"
	case punctDivider:
		return "/"
	case punctNone:
		return ""
	}
	return ""
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
