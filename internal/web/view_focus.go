package web

import (
	"strconv"

	"github.com/maxence-charriere/go-app/v11/pkg/app"

	"github.com/dmikalova/vactrol/internal/engine"
)

// This file draws the lifted copy of the selected card: the same face, larger, on
// a layer above the board, with the card's own verbs under it.
//
// The real card never moves. A copy is laid out over the slot it was lifted from
// and grown from there, so deciding what to do with a card costs the board no
// reflow and the card stays readable while its buttons are up — which is the point
// of putting them on the card rather than in a dock the player has to look away to.
//
// The copy is grown by resizing, not by scaling the board card as a picture: the
// reason to enlarge a card is to read the text the board was clipping, and only a
// real box reflows it. Its height is left to the content — at least focusGrow times
// its slot, and taller than that when its rules (a stack of upgrades, say) need the
// room, up to what the window has left.

const (
	// focusGrow is how much larger the lifted copy is than the card it was lifted
	// from, before its own text asks for more.
	focusGrow = 1.5
	// focusPad is the margin the copy keeps from the edge of the window.
	focusPad = 8.0
	// focusPanelID names the copy in the page, so its laid-out height can be measured
	// and the copy centred on its card by the size it really is.
	focusPanelID = "card-focus"
)

// focusCardID is the card the copy is lifted from, if any. Mid-action phases are
// left alone: the board is being picked over for a flank or a fight target, and a
// card blown up over it would cover the thing being picked.
func (g *game) focusCardID() (engine.LocalID, bool) {
	if !g.hasSel || g.phase != phaseMain {
		return 0, false
	}
	if g.busy || g.choosing || g.choosingOption || g.pickerOpen || g.forgingKey >= 0 {
		return 0, false
	}
	return g.sel, true
}

// cardFocus is the lifted copy: the card's face, a button for each thing it can
// do, and the note saying why it can do nothing. It swallows clicks that land on
// it rather than letting them through to the board it covers, so the enlargement
// cannot cost the player a misclick on a card they can no longer see.
func (g *game) cardFocus() app.UI {
	id, ok := g.focusCardID()
	if !ok {
		return app.Div()
	}
	acts, note := g.selActions()
	panel := app.Div().
		ID(focusPanelID).
		Class(cx("card-focus",
			// The -a/-b pair replays the grow when the lift moves to another card:
			// go-app patches the same element, and a CSS animation only restarts when
			// its animation-name changes.
			ifCls(!g.focusParity, "card-focus--in-a"),
			ifCls(g.focusParity, "card-focus--in-b"),
			// Unmeasured, the copy centres itself on the window: it is better to read
			// the card in the middle of the screen than to lose its buttons with it.
			ifCls(g.hasFocus, "card-focus--placed"))).
		Body(
			g.cardFace(id),
			app.Div().Class("card-focus-acts").Body(
				app.Range(acts).Slice(func(i int) app.UI {
					return btn(acts[i].Label, acts[i].On, acts[i].Class)
				}),
				app.If(note != "", func() app.UI {
					return app.Div().Class("hint").Text(note)
				}),
			),
		)
	if !g.hasFocus {
		return panel
	}
	x, y, w, minH := g.focusBox()
	// Where this particular card sits is a runtime measurement, so it is the one
	// thing the markup carries besides class names; app.css owns what is done with
	// it, the same way house colours arrive as --nm/--tp.
	return panel.
		Style("--focus-x", px(x)).
		Style("--focus-y", px(y)).
		Style("--focus-w", px(w)).
		Style("--focus-min-h", px(minH)).
		Style("--focus-max-h", px(g.focusViewH-2*focusPad)).
		// Where the copy grows from: its own slot on the board, so it is seen coming
		// off the card rather than appearing beside it.
		Style("--focus-dx", px(g.focusRect.x-x)).
		Style("--focus-dy", px(g.focusRect.y-y)).
		Style("--focus-from", strconv.FormatFloat(g.focusRect.w/w, 'f', 3, 64))
}

// focusBox is where the lifted copy goes and how wide it is: centred on the card it
// was lifted from, then pushed back inside the window, so a card on a flank or down
// in the hand still grows evenly in every direction it has the room to.
func (g *game) focusBox() (x, y, w, minH float64) {
	r := g.focusRect
	w = min(r.w*focusGrow, g.focusViewW-2*focusPad)
	minH = r.h * focusGrow
	// How tall the copy really is depends on the text it has to fit, which is known
	// only once it has been laid out; until then its floor stands in for it.
	h := max(minH, g.focusPanelH)
	x = clampAxis(r.x+r.w/2-w/2, w, g.focusViewW)
	y = clampAxis(r.y+r.h/2-h/2, h, g.focusViewH)
	return x, y, w, minH
}

// clampAxis pins a span of the given length inside the window, keeping focusPad at
// each edge. A span with no room to spare is pinned to the near edge.
func clampAxis(v, length, window float64) float64 {
	hi := window - length - focusPad
	if hi < focusPad {
		return focusPad
	}
	return min(max(v, focusPad), hi)
}

// px formats a measured length for a custom property.
func px(v float64) string {
	return strconv.FormatFloat(v, 'f', 1, 64) + "px"
}
