package web

import "testing"

// The copy is centred on the card it was lifted from and then pushed back inside
// the window, so a card on a flank or down in the hand still grows evenly in every
// direction it has room for rather than off the edge of the screen. Its y is
// measured from whichever window edge its card is nearest.
func TestTheLiftCentresOnItsCardAndStaysOnScreen(t *testing.T) {
	const w, h = 144, 192 // a card's slot; the copy is at least 1.5x that
	for _, tc := range []struct {
		what           string
		x, y           float64
		wantX, wantY   float64
		wantW, wantMin float64
	}{
		// Below the midline, so y is from the bottom: (800 - 496) - 288/2.
		{"in the player's half", 600, 400, 564, 160, 216, 288},
		{"against the top left", 4, 4, focusPad, focusPad, 216, 288},
		{"against the bottom right", 1130, 610, 1280 - 216 - focusPad, focusPad, 216, 288},
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

// A card in the player's own half puts its verbs above its face; one across the
// midline puts them below, so the buttons always face the middle of the screen.
func TestTheLiftPutsItsVerbsOnTheInsideEdge(t *testing.T) {
	for _, tc := range []struct {
		what   string
		y      float64
		wantUp bool
	}{
		{"in the player's hand", 700, true},
		{"in the opposing battleline", 100, false},
	} {
		g := &game{
			hasFocus:   true,
			focusRect:  cardRect{x: 600, y: tc.y, w: 144, h: 192},
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
	for _, tc := range []struct {
		what   string
		y      float64
		wantDY float64
	}{
		// Anchored by its top: 100 - (196 - 144).
		{"in the opponent's half", 100, 48},
		// Anchored by its bottom: (600 + 192) - (800 - 8).
		{"in the player's half", 600, 0},
	} {
		g := &game{
			hasFocus:   true,
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
