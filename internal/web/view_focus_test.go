package web

import "testing"

// The copy is centred on the card it was lifted from, and then pushed back inside
// the window: a card on a flank or down in the hand still grows evenly in every
// direction it has room for, rather than off the edge of the screen.
func TestTheLiftCentresOnItsCardAndStaysOnScreen(t *testing.T) {
	const w, h = 144, 192 // a card's slot; the copy is 1.5x that
	for _, tc := range []struct {
		what           string
		x, y           float64
		wantX, wantY   float64
		wantW, wantMin float64
	}{
		{"in the middle", 600, 400, 564, 352, 216, 288},
		{"against the top left", 4, 4, focusPad, focusPad, 216, 288},
		{"against the bottom right", 1130, 610, 1280 - 216 - focusPad, 800 - 288 - focusPad, 216, 288},
	} {
		g := &game{
			hasFocus:   true,
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

// A copy taller than the window it is in gets pinned to the top edge rather than
// centred off the top of it, so its face is the part that survives.
func TestALiftTallerThanTheWindowPinsToTheTop(t *testing.T) {
	g := &game{
		hasFocus:    true,
		focusRect:   cardRect{x: 10, y: 10, w: 144, h: 192},
		focusPanelH: 900,
		focusViewW:  1280,
		focusViewH:  400,
	}
	if _, y, _, _ := g.focusBox(); y != focusPad {
		t.Errorf("an oversized lift sits at y=%v, want %v", y, focusPad)
	}
}

// Once laid out the copy knows its real height, and centres on its card by that
// rather than by the floor it started from.
func TestTheLiftCentresByItsLaidOutHeight(t *testing.T) {
	g := &game{
		hasFocus:   true,
		focusRect:  cardRect{x: 600, y: 400, w: 144, h: 192},
		focusViewW: 1280,
		focusViewH: 800,
	}
	_, before, _, _ := g.focusBox()
	g.focusPanelH = 400
	_, after, _, _ := g.focusBox()
	if after != before-56 {
		t.Errorf("a 400px-tall lift sits at y=%v, want %v", after, before-56)
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
