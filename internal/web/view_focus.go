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
// real box reflows it. So only its width is set; the height is whatever its rules
// (a stack of upgrades, say) come out at, with a floor of focusMinGrow and a
// ceiling of what the window has left.

const (
	// focusGrow is how much wider the lifted copy is than the card it was lifted
	// from. This is the dial for how large the lift reads.
	focusGrow = 1.2
	// focusMinGrow is the shortest the copy may be, as a multiple of its card's
	// height: enough to read as lifted even when its rules are a single line.
	focusMinGrow = 1.1
	// focusPad is the margin the copy keeps from the edge of the window.
	focusPad = 8.0
)

// focusCardID is the card the copy is lifted from, if any. Picking a flank keeps
// the lift, since the question is about the card being placed and is asked on it;
// picking a fight target drops it, since the answer is another card on the board
// that a card blown up over it would cover.
func (g *game) focusCardID() (engine.LocalID, bool) {
	if !g.hasSel || (g.phase != phaseMain && g.phase != phaseFlank) {
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
	face := g.cardFace(id)
	// The copy is the card, so a drag has to start from it: it lies over its own
	// neighbours, and a pointer that fell through would grab whichever card the
	// enlarged face happens to cover.
	if g.phase == phaseMain && g.selKind == selHand && g.playableFromHand(id) {
		face.ID = id
		face.Draggable = true
		face.OnDragStart = g.startHandDrag
		face.OnDragEnd = g.endHandDrag
	}
	panel := app.Div().
		OnWheel(g.wheelOverFocus).
		Class(cx("card-focus",
			// The -a/-b pair replays the grow when the lift moves to another card:
			// go-app patches the same element, and a CSS animation only restarts when
			// its animation-name changes.
			ifCls(!g.focusParity, "card-focus--in-a"),
			ifCls(g.focusParity, "card-focus--in-b"),
			ifCls(g.actsUp(), "card-focus--acts-up"),
			// Unmeasured, the copy centres itself on the window: it is better to read
			// the card in the middle of the screen than to lose its buttons with it.
			ifCls(g.hasFocus, "card-focus--placed"))).
		Body(
			face,
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
		// off the card rather than appearing beside it. dy is measured to the same
		// corner y is anchored at, which is the one corner whose position is known
		// without knowing how tall the copy came out.
		Style("--focus-dx", px(g.focusRect.x-x)).
		Style("--focus-dy", px(g.focusDY(y))).
		Style("--focus-from", strconv.FormatFloat(g.focusRect.w/w, 'f', 3, 64))
}

// focusDY is how far the grow animation starts below the copy's anchored corner:
// the offset that lands that corner on the matching corner of the card's own slot.
func (g *game) focusDY(y float64) float64 {
	if g.actsUp() {
		return g.focusRect.y + g.focusRect.h - (g.focusViewH - y)
	}
	return g.focusRect.y - y
}

// focusBox is where the lifted copy goes and how wide it is: centred on the card
// it was lifted from and then pushed back inside the window, so a card on a flank
// or down in the hand still grows evenly in every direction it has the room to.
//
// y is measured from the window edge the card is nearest — the top for a card in
// the opponent's half, the bottom for one in the player's own. Anchoring that way
// is what lets the copy be placed from the card's rect alone: the copy is as tall
// as its text needs, which is not known until it has been laid out, and whatever it
// turns out to be grows away from the near edge rather than through it. minH is
// both the copy's floor and, for want of knowing better, the height it is centred
// as if it had.
func (g *game) focusBox() (x, y, w, minH float64) {
	r := g.focusRect
	w = min(r.w*focusGrow, g.focusViewW-2*focusPad)
	minH = r.h * focusMinGrow
	cy := r.y + r.h/2
	if g.actsUp() {
		cy = g.focusViewH - cy
	}
	x = clampAxis(r.x+r.w/2-w/2, w, g.focusViewW)
	y = clampAxis(cy-minH/2, minH, g.focusViewH)
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

// actsUp is whether the card's verbs are drawn above its face rather than below.
// They go on the side facing the middle of the window — a card down in the player's
// own half puts its buttons up, one in the opponent's half puts them down — so they
// meet the pointer on its way in rather than beyond the edge the card was already
// against.
func (g *game) actsUp() bool {
	return g.hasFocus && g.focusViewH > 0 &&
		g.focusRect.y+g.focusRect.h/2 > g.focusViewH/2
}

// px formats a measured length for a custom property.
func px(v float64) string {
	return strconv.FormatFloat(v, 'f', 1, 64) + "px"
}
