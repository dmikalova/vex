package web

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/engine"
)

// These tests cover manual mode: the playtester's controls that move, ready, and
// exhaust cards, adjust the counters, and put any printed card into a hand,
// ignoring the rules the rest of the client obeys.

func TestManualModeTogglesOff(t *testing.T) {
	c := newClient(t)
	c.startTurn()
	c.manual()
	c.do(c.g.toggleManual)
	if c.g.g.Manual() {
		t.Error("manual mode did not turn back off")
	}
}

// Every manual control is inert outside manual mode, so the ordinary match is
// not quietly editable.
func TestManualControlsNeedManualMode(t *testing.T) {
	c := newClient(t)
	c.startTurn()
	me := c.g.active()
	id := c.hand()[0]
	c.g.selectHandID(c.ctx, id)

	amber, chains := c.g.g.State.Aember[me], c.g.g.State.Chains[me]
	keys := c.g.g.Keys(me)
	house := c.g.g.State.ActiveHouse

	c.do(c.g.manualMove(engine.ManualPurge))
	c.do(c.g.manualReady)
	c.do(c.g.manualExhaust)
	c.do(c.g.manualAmberDelta(me, 3))
	c.do(c.g.manualChainsDelta(me, 3))
	c.do(c.g.manualForgeKey(me))
	c.do(c.g.manualUnforgeKey(me))
	c.do(c.g.manualSetHouse(engine.Brobnar))

	if !containsID(c.hand(), id) {
		t.Error("a card was moved out of hand outside manual mode")
	}
	if c.g.g.State.Aember[me] != amber {
		t.Error("Æmber was adjusted outside manual mode")
	}
	if c.g.g.State.Chains[me] != chains {
		t.Error("chains were adjusted outside manual mode")
	}
	if c.g.g.Keys(me) != keys {
		t.Error("keys were adjusted outside manual mode")
	}
	if c.g.forgingKey != -1 {
		t.Error("the forge picker opened outside manual mode")
	}
	if c.g.g.State.ActiveHouse != house {
		t.Error("the active house was changed outside manual mode")
	}
}

// With nothing selected there is no card for the card controls to act on.
func TestManualCardControlsNeedASelection(t *testing.T) {
	c := newClient(t)
	c.manualTurn(testHouse)
	before := len(c.g.rootMarks)
	c.do(c.g.manualMove(engine.ManualPurge))
	c.do(c.g.manualReady)
	c.do(c.g.manualExhaust)
	if len(c.g.rootMarks) != before {
		t.Errorf("a card control with nothing selected recorded %d undo steps",
			len(c.g.rootMarks)-before)
	}
}

func TestManualMovesACard(t *testing.T) {
	c := newClient(t)
	c.manualTurn(testHouse)
	id := c.deal(testCreature)
	c.g.selectHandID(c.ctx, id)
	c.do(c.g.manualMove(engine.ManualArchives))

	if !containsID(c.g.g.Archives(c.g.active()), id) {
		t.Error("the card was not moved to archives")
	}
	if c.g.hasSel {
		t.Error("the moved card was left selected in a zone it is no longer in")
	}
}

func TestManualReadyAndExhaust(t *testing.T) {
	c := newClient(t)
	c.manualTurn(testHouse)
	id := c.deal(testCreature)
	c.playFromHand(id)
	c.g.selectBoardID(c.ctx, id)

	c.do(c.g.manualExhaust)
	if !c.g.g.State.Cards[id].Exhausted {
		t.Error("the creature was not exhausted")
	}
	c.do(c.g.manualReady)
	if c.g.g.State.Cards[id].Exhausted {
		t.Error("the creature was not readied")
	}
}

func TestManualCounters(t *testing.T) {
	c := newClient(t)
	c.manualTurn(testHouse)
	me := c.g.active()

	before := c.g.g.State.Aember[me]
	c.do(c.g.manualAmberDelta(me, 4))
	if got := c.g.g.State.Aember[me]; got != before+4 {
		t.Errorf("Æmber is %d, want %d", got, before+4)
	}

	chains := c.g.g.State.Chains[me]
	c.do(c.g.manualChainsDelta(me, 2))
	if got := c.g.g.State.Chains[me]; got != chains+2 {
		t.Errorf("chains are %d, want %d", got, chains+2)
	}
}

// Forging by hand asks which colour, because a key's colour is not something the
// client can work out for the playtester.
func TestManualForgeAndUnforge(t *testing.T) {
	c := newClient(t)
	c.manualTurn(testHouse)
	me := c.g.active()

	c.do(c.g.manualForgeKey(me))
	if c.g.forgingKey != me {
		t.Fatalf("the forge picker is at %d, want player %d", c.g.forgingKey, me)
	}
	c.do(c.g.pickForgeColor(engine.KeyColorYellow))
	if c.g.forgingKey != -1 {
		t.Error("the forge picker stayed open after a colour was picked")
	}
	if c.g.g.Keys(me) != 1 {
		t.Errorf("the player has %d keys, want 1", c.g.g.Keys(me))
	}

	c.do(c.g.manualUnforgeKey(me))
	if c.g.g.Keys(me) != 0 {
		t.Errorf("the player has %d keys after unforging, want 0", c.g.g.Keys(me))
	}
}

// r, b and y forge directly rather than making the playtester hunt the matching
// button.
func TestKeyColorKeys(t *testing.T) {
	for _, tt := range []struct {
		key  string
		want engine.KeyColor
	}{
		{"r", engine.KeyColorRed},
		{"b", engine.KeyColorBlue},
		{"y", engine.KeyColorYellow},
	} {
		t.Run(tt.key, func(t *testing.T) {
			c := newClient(t)
			c.manualTurn(testHouse)
			me := c.g.active()
			c.do(c.g.manualForgeKey(me))
			c.press(tt.key)
			if c.g.forgingKey != -1 {
				t.Fatalf("%q did not answer the forge picker", tt.key)
			}
			if got := c.g.g.KeyColors(me); len(got) != 1 || got[0] != tt.want {
				t.Errorf("%q forged %v, want %v", tt.key, got, tt.want)
			}
		})
	}
}

// A colour that has already been forged is no longer offered, so its key does
// not answer the picker either.
func TestAKeyColorAlreadyForgedIsNotOffered(t *testing.T) {
	c := newClient(t)
	c.manualTurn(testHouse)
	me := c.g.active()
	c.do(c.g.manualForgeKey(me))
	c.press("r")

	c.do(c.g.manualForgeKey(me))
	c.press("r")
	if c.g.forgingKey != me {
		t.Error("r answered the picker with a colour that was already forged")
	}
	c.do(c.g.cancelForgeKey)
	if c.g.forgingKey != -1 {
		t.Error("cancelling did not close the forge picker")
	}
}

// Picking a colour with no picker open has no player to forge for.
func TestPickForgeColorWithNoPickerOpen(t *testing.T) {
	c := newClient(t)
	c.manualTurn(testHouse)
	c.do(c.g.pickForgeColor(engine.KeyColorRed))
	if c.g.g.Keys(c.g.active()) != 0 {
		t.Error("a colour was forged with no picker open")
	}
}

// Setting the house by hand from the house prompt also gets play under way, the
// same as choosing one normally.
func TestManualSetHouseStartsTheTurn(t *testing.T) {
	c := newClient(t)
	c.manual()
	c.do(c.g.manualSetHouse(engine.Brobnar))
	if c.g.g.State.ActiveHouse != engine.Brobnar {
		t.Errorf("the active house is %v, want Brobnar", c.g.g.State.ActiveHouse)
	}
	if c.g.phase != phaseMain {
		t.Errorf("the phase is %v, want phaseMain", c.g.phase)
	}
}

// A card selected out of the zone viewer becomes the selection so the manual
// controls can act on it, and the viewer gets out of the way.
func TestSelectingFromTheZoneViewer(t *testing.T) {
	c := newClient(t)
	c.manualTurn(testHouse)
	c.g.zonesPlayer = c.g.active()
	id := c.g.g.Deck(c.g.active())[0]

	c.g.selectZoneCard(c.ctx, id)
	if !c.g.hasSel || c.g.sel != id || c.g.selKind != selOther {
		t.Fatalf("selecting %d from the viewer left sel=%d kind=%v",
			id, c.g.sel, c.g.selKind)
	}
	if c.g.zonesPlayer != -1 {
		t.Error("the viewer stayed open over the controls it handed the card to")
	}

	c.do(c.g.manualMove(engine.ManualHand))
	if !containsID(c.hand(), id) {
		t.Error("the card the viewer handed over was not moved to hand")
	}
}

// Put into play arms the Deploy placement picker on a selected hand creature and
// a clicked position drops it there without any play effects.
func TestManualPutIntoPlayByClick(t *testing.T) {
	c := newClient(t)
	c.manualTurn(testHouse)
	c.playFromHand(c.deal(testCreature))
	c.playFromHand(c.deal(testCreature))
	line := c.board()

	hid := c.deal(testCreature)
	c.g.selectHandID(c.ctx, hid)
	c.do(c.g.manualPlay)
	if !c.g.manualPlacing || !c.g.choosingPosition {
		t.Fatal("Put into play did not arm the placement picker")
	}

	c.do(c.g.chooseDeploySide(true)) // right of the clicked creature
	c.g.choosePositionCandidate(c.ctx, line[0])

	nb := c.board()
	if len(nb) != 3 || nb[1] != hid {
		t.Errorf("board = %v, want the hand card at index 1", nb)
	}
	if c.g.manualPlacing || c.g.choosingPosition {
		t.Error("placement state was not cleared after placing")
	}
	if containsID(c.hand(), hid) {
		t.Error("the placed card is still in hand")
	}
}

// With an empty battleline Put into play drops the creature straight in without
// raising the placement picker.
func TestManualPutIntoPlaySkipsPickerOnEmptyLine(t *testing.T) {
	c := newClient(t)
	c.manualTurn(testHouse)

	hid := c.deal(testCreature)
	c.g.selectHandID(c.ctx, hid)
	c.do(c.g.manualPlay)

	if c.g.manualPlacing || c.g.choosingPosition {
		t.Error("an empty line raised the placement picker")
	}
	if !containsID(c.board(), hid) {
		t.Error("the creature was not placed into play")
	}
}

// Cancel backs out of a manual placement without placing, leaving the card in hand.
func TestManualPutIntoPlayCancel(t *testing.T) {
	c := newClient(t)
	c.manualTurn(testHouse)
	c.playFromHand(c.deal(testCreature))

	hid := c.deal(testCreature)
	c.g.selectHandID(c.ctx, hid)
	c.do(c.g.manualPlay)
	if !c.g.manualPlacing {
		t.Fatal("Put into play did not arm the placement picker")
	}

	c.do(c.g.cancelManualPlace)
	if c.g.manualPlacing || c.g.choosingPosition {
		t.Error("Cancel did not clear the placement state")
	}
	if !containsID(c.hand(), hid) {
		t.Error("Cancel should leave the card in hand")
	}
}

func TestThePickerOpensAndCloses(t *testing.T) {
	c := newClient(t)
	c.manualTurn(testHouse)

	c.g.pickerQuery = "left over"
	c.g.pickerFocused = true
	c.do(c.g.openPicker)
	if !c.g.pickerOpen {
		t.Fatal("the picker did not open")
	}
	if c.g.pickerQuery != "" {
		t.Error("the picker opened onto the last search")
	}
	if c.g.pickerFocused {
		t.Error("the picker opened without asking for the caret again")
	}

	c.do(c.g.closePicker)
	if c.g.pickerOpen {
		t.Error("the picker did not close")
	}
}

// A picker row naming a card the pool does not hold adds nothing: the row's card
// is read from the DOM at click time, and off-browser it reads back as nothing
// at all.
func TestAddingACardTheStoreDoesNotKnow(t *testing.T) {
	c := newClient(t)
	c.manualTurn(testHouse)
	c.g.pickerOpen = true
	before := len(c.hand())

	c.do(c.g.addPickedCard)
	if len(c.hand()) != before {
		t.Error("an unnamed picker row added a card")
	}
	if !c.g.pickerOpen {
		t.Error("an unnamed picker row closed the picker")
	}
}

func TestIsInPlay(t *testing.T) {
	c := newClient(t)
	c.manualTurn(testHouse)
	id := c.deal(testCreature)
	if c.g.isInPlay(id) {
		t.Error("a card in hand reads as in play")
	}
	c.playFromHand(id)
	if !c.g.isInPlay(id) {
		t.Error("a creature on the battleline does not read as in play")
	}
}

// Graft threads the selected card face up under an in-play host: the player picks
// Graft, then clicks the host to land it there.
func TestManualGraftsACardUnderAHost(t *testing.T) {
	c := newClient(t)
	c.manualTurn(testHouse)
	host := c.deal(testCreature)
	c.playFromHand(host)
	sub := c.deal(testCreature)
	c.g.selectHandID(c.ctx, sub)

	c.do(c.g.manualGraft)
	if !c.g.hostTargeting || c.g.hostFaceDown {
		t.Fatalf("Graft armed hostTargeting=%v hostFaceDown=%v, want true,false",
			c.g.hostTargeting, c.g.hostFaceDown)
	}
	c.g.attachToHost(c.ctx, host)
	c.settle()

	if !containsID(c.g.g.Under(host), sub) {
		t.Error("the grafted card was not placed under the host")
	}
	if c.g.g.UnderFaceDown(sub) {
		t.Error("a graft placed the card face down, want face up")
	}
	if containsID(c.hand(), sub) {
		t.Error("the grafted card is still in hand")
	}
	if c.g.hostTargeting || c.g.hasSel {
		t.Error("host targeting or the selection lingered after the graft")
	}
}

// Place under threads the selected card face down under a host — the facedown
// counterpart to Graft.
func TestManualPlacesACardUnderAHostFaceDown(t *testing.T) {
	c := newClient(t)
	c.manualTurn(testHouse)
	host := c.deal(testCreature)
	c.playFromHand(host)
	sub := c.deal(testCreature)
	c.g.selectHandID(c.ctx, sub)

	c.do(c.g.manualPlaceUnder)
	if !c.g.hostTargeting || !c.g.hostFaceDown {
		t.Fatalf("Place under armed hostTargeting=%v hostFaceDown=%v, want true,true",
			c.g.hostTargeting, c.g.hostFaceDown)
	}
	c.g.attachToHost(c.ctx, host)
	c.settle()

	if !containsID(c.g.g.Under(host), sub) {
		t.Error("the card was not placed under the host")
	}
	if !c.g.g.UnderFaceDown(sub) {
		t.Error("Place under left the card face up, want face down")
	}
}

// A host pick will not thread a card under itself, so clicking the selected card
// as its own host does nothing.
func TestManualGraftIgnoresSelfAsHost(t *testing.T) {
	c := newClient(t)
	c.manualTurn(testHouse)
	sub := c.deal(testCreature)
	c.playFromHand(sub)
	c.g.selectBoardID(c.ctx, sub)

	c.do(c.g.manualGraft)
	c.g.attachToHost(c.ctx, sub)
	c.settle()

	if containsID(c.g.g.Under(sub), sub) {
		t.Error("a card was grafted under itself")
	}
}

// To hand detaches a selected upgrade from its host and returns it to hand.
func TestManualSendsAnUpgradeToHand(t *testing.T) {
	c := newClient(t)
	c.manualTurn(testHouse)
	host := c.deal(testCreature)
	c.playFromHand(host)
	up := c.g.g.Register(
		engine.NewCard("Test Upgrade", testHouse, engine.Upgrade, engine.Common), c.g.active())
	c.g.g.AttachUpgrade(host, up)

	c.g.selectTab(up)
	if !c.g.isAttached(up) {
		t.Fatal("the attached upgrade does not read as attached")
	}
	c.do(c.g.manualToHand)

	if !containsID(c.hand(), up) {
		t.Error("the upgrade was not sent to hand")
	}
	if containsID(c.g.g.Upgrades(host), up) {
		t.Error("the upgrade is still attached to its host")
	}
}

// To hand also returns a card placed under a host, the other attached state
// isAttached recognises.
func TestManualSendsAnUnderCardToHand(t *testing.T) {
	c := newClient(t)
	c.manualTurn(testHouse)
	host := c.deal(testCreature)
	c.playFromHand(host)
	sub := c.deal(testCreature)
	c.g.selectHandID(c.ctx, sub)
	c.do(c.g.manualPlaceUnder)
	c.g.attachToHost(c.ctx, host)
	c.settle()

	if !c.g.isAttached(sub) {
		t.Fatal("the under-card does not read as attached")
	}
	c.g.selectTab(sub)
	c.do(c.g.manualToHand)

	if !containsID(c.hand(), sub) {
		t.Error("the under-card was not sent to hand")
	}
	if containsID(c.g.g.Under(host), sub) {
		t.Error("the card is still under its host")
	}
}

// The graft, place-under, and to-hand controls are inert outside manual mode, so
// an ordinary match is not quietly editable through them.
func TestManualAttachControlsNeedManualMode(t *testing.T) {
	c := newClient(t)
	c.startTurn()
	id := c.hand()[0]
	c.g.selectHandID(c.ctx, id)

	c.do(c.g.manualGraft)
	c.do(c.g.manualPlaceUnder)
	if c.g.hostTargeting {
		t.Error("host targeting armed outside manual mode")
	}
	before := len(c.g.rootMarks)
	c.do(c.g.manualToHand)
	if len(c.g.rootMarks) != before {
		t.Error("To hand acted outside manual mode")
	}
}
