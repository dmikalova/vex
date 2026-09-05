package web

import (
	"strings"
	"testing"
)

// The copy is centred on the card it was lifted from and then pushed back inside
// the window, so a card on a flank or down in the hand still grows evenly in every
// direction it has room for rather than off the edge of the screen. Its y is
// measured from whichever window edge its buttons sit against — the bottom for a
// card in hand, the top for a card on the board.
func TestTheLiftCentresOnItsCardAndStaysOnScreen(t *testing.T) {
	const w, h = 144, 192 // a card's slot
	// The copy's own size, and so the box it is centred and clamped as.
	gotW := float64(w) * focusGrow
	minH := float64(h) * focusMinGrow
	for _, tc := range []struct {
		what           string
		kind           selKind
		x, y           float64
		wantX, wantY   float64
		wantW, wantMin float64
	}{
		// A hand card anchors from the bottom, so y is from the bottom edge:
		// (800 - 496) - minH/2.
		{"in hand", selHand, 600, 400, 672 - gotW/2, 304 - minH/2, gotW, minH},
		// A board card anchors from the top, so y is its own centre less half its
		// height: 496 - minH/2.
		{"on the board", selYourCreature, 600, 400, 672 - gotW/2, 496 - minH/2, gotW, minH},
		// A board card against the top edge is pinned there.
		{"against the top left", selYourCreature, 4, 4, focusPad, focusPad, gotW, minH},
		// A hand card near the bottom edge anchors from the bottom and is pinned to
		// its near (bottom-measured) edge.
		{"against the bottom right", selHand, 1130, 610, 1280 - gotW - focusPad, focusPad, gotW, minH},
	} {
		g := &game{
			hasFocus:   true,
			selKind:    tc.kind,
			focusRect:  cardRect{x: tc.x, y: tc.y, w: w, h: h},
			focusViewW: 1280,
			focusViewH: 800,
		}
		x, y, gw, minH := g.focusBox()
		if x != tc.wantX || y != tc.wantY || gw != tc.wantW || minH != tc.wantMin {
			t.Errorf("a card %s lifts to (%v,%v) %vx%v, want (%v,%v) %vx%v",
				tc.what, x, y, gw, minH, tc.wantX, tc.wantY, tc.wantW, tc.wantMin)
		}
	}
}

// A card lifted from hand puts its verbs above its face; a card on the board puts
// them below, so a card being played and a card already in play read their buttons
// in the same place.
func TestTheLiftPutsHandVerbsAboveAndBoardVerbsBelow(t *testing.T) {
	for _, tc := range []struct {
		what   string
		kind   selKind
		wantUp bool
	}{
		{"in the player's hand", selHand, true},
		{"in the player's own battleline", selYourCreature, false},
		{"in the opposing battleline", selOther, false},
	} {
		g := &game{
			hasFocus:   true,
			selKind:    tc.kind,
			focusRect:  cardRect{x: 600, y: 400, w: 144, h: 192},
			focusViewW: 1280,
			focusViewH: 800,
		}
		if up := g.actsUp(); up != tc.wantUp {
			t.Errorf("a card %s draws its verbs above=%v, want %v", tc.what, up, tc.wantUp)
		}
	}
}

// With nothing measured there is no midline to be on a side of, so the verbs stay
// under the face — which is where the unplaced, window-centred copy wants them.
func TestAnUnmeasuredLiftKeepsItsVerbsBelow(t *testing.T) {
	if (&game{}).actsUp() {
		t.Error("an unmeasured lift draws its verbs above its face, want below")
	}
}

// A copy with no room to spare is pinned to the near edge rather than centred half
// off the screen, so the anchored end of its face is the part that survives.
func TestALiftTallerThanTheWindowPinsToTheNearEdge(t *testing.T) {
	g := &game{
		hasFocus:   true,
		focusRect:  cardRect{x: 10, y: 10, w: 144, h: 400},
		focusViewW: 1280,
		focusViewH: 400,
	}
	if _, y, _, _ := g.focusBox(); y != focusPad {
		t.Errorf("an oversized lift sits at y=%v, want %v", y, focusPad)
	}
}

// The grow starts at the card's own slot, so the copy is seen coming off the card
// rather than appearing beside it. Whichever corner the copy is anchored by is the
// one run back to the card's matching corner, since it is the only corner whose
// place is known before the copy's text has decided how tall it is.
func TestTheLiftGrowsOutOfItsCardsSlot(t *testing.T) {
	minH := float64(192) * focusMinGrow
	for _, tc := range []struct {
		what   string
		kind   selKind
		y      float64
		wantDY float64
	}{
		// A board card is anchored by its top: 100 - (196 - minH/2).
		{"in the opponent's half", selOther, 100, 100 - (196 - minH/2)},
		// A hand card is anchored by its bottom: (600 + 192) - (800 - 8).
		{"in the player's hand", selHand, 600, 0},
	} {
		g := &game{
			hasFocus:   true,
			selKind:    tc.kind,
			focusRect:  cardRect{x: 600, y: tc.y, w: 144, h: 192},
			focusViewW: 1280,
			focusViewH: 800,
		}
		_, y, _, _ := g.focusBox()
		if dy := g.focusDY(y); dy != tc.wantDY {
			t.Errorf("a card %s grows from dy=%v, want %v", tc.what, dy, tc.wantDY)
		}
	}
}

// Off-browser there is no element to measure, so the lift is never placed and the
// client renders its centred fallback rather than a bogus rect.
func TestMeasuringTheLiftWithNoPageBehindIt(t *testing.T) {
	c := newClient(t)
	c.manualTurn(testHouse)
	c.g.selectHandID(c.ctx, c.deal(testCreature))
	if c.g.measureFocus() {
		t.Error("measuring with no page behind it reported a move")
	}
	if c.g.hasFocus {
		t.Error("measuring with no page behind it placed the lift anyway")
	}
	c.wants("an unplaced lift", "card-focus")
	c.lacks("an unplaced lift", "card-focus--placed")
}

// A resized window re-places the lift, since where its card sits has moved. Off
// browser there is nothing to measure, so the resize is simply absorbed.
func TestResizingRePlacesTheLift(t *testing.T) {
	c := newClient(t)
	c.manualTurn(testHouse)
	c.g.selectHandID(c.ctx, c.deal(testCreature))
	c.g.OnResize(c.ctx)
	if c.g.hasFocus {
		t.Error("resizing with no page behind it placed the lift anyway")
	}
}

// Scrolling a strip slides the card out from under its own copy, so the copy is
// re-placed the same way a resize re-places it.
func TestScrollingARowRePlacesTheLift(t *testing.T) {
	c := newClient(t)
	c.manualTurn(testHouse)
	c.g.selectHandID(c.ctx, c.deal(testCreature))
	c.g.installScrollTracking()
	c.g.placeFocus()
	if c.g.hasFocus {
		t.Error("scrolling with no page behind it placed the lift anyway")
	}
}

// Dropping the selection drops the lift with it, which is a move: the client has

// to render again to take the copy off the board.
func TestDroppingTheSelectionDropsTheLift(t *testing.T) {
	c := newClient(t)
	c.manualTurn(testHouse)
	c.g.selectHandID(c.ctx, c.deal(testCreature))
	c.g.hasFocus = true
	c.g.clearSelection()
	if !c.g.measureFocus() {
		t.Error("dropping the selection did not report a move")
	}
	if c.g.hasFocus {
		t.Error("the lift outlived the selection")
	}
}

// A selection dropped while its copy was measured plays the copy out — the
// grow-in run backwards — rather than blinking it away: the copy stays up with
// the exit class until the shrink-back finishes, then clears itself.
func TestDeselectingACardPlaysTheLiftOut(t *testing.T) {
	c := newClient(t)
	c.manualTurn(testHouse)
	id := c.deal(testCreature)
	c.g.selectHandID(c.ctx, id)
	// Off-browser nothing measures, so stand in the placement the exit shrinks to.
	c.g.hasFocus = true
	c.g.focusShown = focusSnapshot{id: id, w: 100, minH: 100}

	c.g.clearSelection()
	if !c.g.measureFocus() {
		t.Fatal("dropping the selection did not report a move")
	}
	if !c.g.focusExit {
		t.Fatal("dropping the selection did not start the lift's exit")
	}
	c.wants("the lift shrinking back to its slot", "card-focus--out")

	c.await("the exit to clear itself", func() bool { return !c.g.focusExit })
	c.lacks("the lift gone after its exit", "card-focus--out")
}

// The keybar reports what a card is doing on the table — granted keywords and all
// — so the copy of a card still in hand does not draw one. Reading Elusive off a
// card that is not on a battleline says it is dodging a fight it cannot be in.
func TestTheLiftDrawsAKeybarOnlyInPlay(t *testing.T) {
	c := newClient(t)
	c.manualTurn(testHouse)
	id := c.deal("Dew Faerie")

	c.g.selectHandID(c.ctx, id)
	c.lacks("a card lifted out of hand", "card-keybar")

	c.playFromHand(id)
	c.g.selectBoardID(c.ctx, id)
	c.wants("a card lifted off the battleline", "card-keybar")
}

// A drag has to start from the lifted copy itself: it lies over its neighbours,
// so a pointer falling through it would grab whichever card the enlarged face
// happens to cover. A card already on the board has nowhere to be dragged to.
func TestDraggingTheLiftedCard(t *testing.T) {
	c := newClient(t)
	c.manualTurn(testHouse)
	id := c.deal(testCreature)

	c.g.selectHandID(c.ctx, id)
	if !strings.Contains(liftMarkup(t, c.html()), " draggable>") {
		t.Error("the lifted copy of a playable hand card is not a drag source")
	}

	c.playFromHand(id)
	c.ownNextTurn(testHouse)
	c.g.selectBoardID(c.ctx, id)
	if strings.Contains(liftMarkup(t, c.html()), " draggable>") {
		t.Error("the lifted copy of a creature already in play is a drag source")
	}
}

// liftMarkup is the lifted copy's face, cut out of the drawn client so a count of
// drag sources cannot pick up the hand underneath it.
func liftMarkup(t *testing.T, h string) string {
	t.Helper()
	from := strings.Index(h, `class="card-focus`)
	to := strings.Index(h, `class="card-focus-acts`)
	if from < 0 || to < from {
		t.Fatal("the client drew no lifted card")
	}
	return h[from:to]
}

// Hovering the card that is already lifted does not also pop the preview — that
// would be two enlarged copies of one card. Hovering any other card still does.
func TestHoveringTheLiftedCardDoesNotAlsoPreviewIt(t *testing.T) {
	c := newClient(t)
	c.manualTurn(testHouse)
	first, second := c.deal(testCreature), c.deal(testCreature)
	c.playFromHand(first)
	c.g.selectHandID(c.ctx, second)

	c.g.hoverCard(c.ctx, second)
	if c.g.previewUp() {
		t.Error("hovering the lifted card popped the preview too")
	}
	c.g.hoverCard(c.ctx, first)
	if !c.g.previewUp() {
		t.Error("hovering another card did not pop the preview")
	}
}

func TestFormattingAMeasuredLength(t *testing.T) {
	if got := px(12.345); got != "12.3px" {
		t.Errorf("px(12.345) = %q, want \"12.3px\"", got)
	}
}
