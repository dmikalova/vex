package web

import (
	"errors"
	"testing"

	"github.com/dmikalova/vactrol/internal/engine"
)

// These tests take turns the way a player does — select, play, reap, fight, end
// — and pin down the guards that keep a click in the wrong context from doing
// anything.

// testCreature is a creature with no prompt in its Play ability, so playing one
// never stops the test on a choice.
const testCreature = "Flaxia"

// testHouse is testCreature's house, which a manual turn is set to so the board
// is laid out under the house that can use it.
const testHouse = engine.Untamed

func TestChooseHouseStartsTheTurn(t *testing.T) {
	c := newClient(t)
	if c.g.phase != phaseHouse {
		t.Fatalf("a fresh deal is at phase %v, want phaseHouse", c.g.phase)
	}
	houses := c.g.pickableHouses()
	if len(houses) != 3 {
		t.Fatalf("the deal offers %d houses, want 3", len(houses))
	}
	c.do(c.g.pickHouse(houses[1]))
	if c.g.g.State.ActiveHouse != houses[1] {
		t.Errorf("the active house is %v, want %v", c.g.g.State.ActiveHouse, houses[1])
	}
	if c.g.phase != phaseMain {
		t.Errorf("the phase is %v, want phaseMain", c.g.phase)
	}
}

func TestSelectingACardInHand(t *testing.T) {
	c := newClient(t)
	c.startTurn()
	id := c.hand()[0]

	c.g.selectHandID(c.ctx, id)
	if !c.g.hasSel || c.g.sel != id || c.g.selKind != selHand {
		t.Fatalf("selecting %d left sel=%d kind=%v has=%v",
			id, c.g.sel, c.g.selKind, c.g.hasSel)
	}
	if c.g.selHand < 0 {
		t.Error("the hand index was not recovered from the id")
	}
	if got := c.g.selHandSlot(); got < 0 {
		t.Errorf("selHandSlot is %d, want the card's place in the drawn hand", got)
	}
}

// A card that is not in the active player's hand is not selectable as one, so a
// click left over from a re-render cannot move the selection somewhere
// impossible.
func TestSelectingACardNotInHandDoesNothing(t *testing.T) {
	c := newClient(t)
	c.startTurn()
	c.g.selectHandID(c.ctx, c.g.g.Deck(c.g.active())[0])
	if c.g.hasSel {
		t.Error("a card outside the hand was selected as a hand card")
	}
	if got := c.g.selHandSlot(); got != -1 {
		t.Errorf("selHandSlot with nothing selected is %d, want -1", got)
	}
}

// Playing an Omega card ends the play phase the moment it resolves, which runs
// the turn out. The end-turn button hands the turn to the opponent afterward;
// this makes sure the Omega path does the same handoff rather than leaving the
// ended turn's player waiting to choose a house again.
func TestPlayingAnOmegaCardHandsTheTurnToTheOpponent(t *testing.T) {
	c := newClient(t)
	c.startTurn()
	c.manual() // lift the house restriction so Swindle (Shadows) is playable
	me := c.g.active()
	c.playFromHand(c.deal("Swindle"))

	if got := c.g.active(); got != 1-me {
		t.Fatalf("after an Omega card the active player is %d, want the opponent %d", got, 1-me)
	}
	if c.g.phase != phaseHouse {
		t.Errorf("the opponent's turn is at phase %v, want phaseHouse", c.g.phase)
	}
}

func TestBoardKindOfEachRow(t *testing.T) {
	c := newClient(t)
	c.manualTurn(testHouse)
	mine := c.deal(testCreature)
	c.playFromHand(mine)

	if got := c.g.boardKindOf(mine); got != selYourCreature {
		t.Errorf("boardKindOf(own creature) = %v, want selYourCreature", got)
	}
	other := c.g.g.Deck(1 - c.g.active())[0]
	if got := c.g.boardKindOf(other); got != selOther {
		t.Errorf("boardKindOf(a card not in the active player's rows) = %v, want selOther", got)
	}
}

// Selecting is refused while a prompt or an action owns the screen; the player
// has to answer it before touching the board again.
func TestSelectingIsGuarded(t *testing.T) {
	tests := []struct {
		name string
		arm  func(g *game)
	}{
		{"busy", func(g *game) { g.busy = true }},
		{"choosing", func(g *game) { g.choosing = true }},
		{"picking a fight target", func(g *game) { g.phase = phaseFightTarget }},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := newClient(t)
			c.startTurn()
			id := c.hand()[0]
			tt.arm(c.g)

			c.g.selectHandID(c.ctx, id)
			c.g.selectBoardID(c.ctx, id)
			if c.g.hasSel {
				t.Error("the selection moved while a prompt was up")
			}
		})
	}
}

// A creature played onto an empty battleline has only one place to go, so it
// skips the flank prompt; with a creature already there the player is asked.
func TestFlankPromptOnlyWhenThereIsAChoice(t *testing.T) {
	c := newClient(t)
	c.manualTurn(testHouse)

	first := c.deal(testCreature)
	c.g.selectHandID(c.ctx, first)
	c.do(c.g.play)
	if c.g.phase != phaseMain {
		t.Fatalf("the first creature asked for a flank; phase is %v", c.g.phase)
	}
	if !containsID(c.board(), first) {
		t.Fatal("the first creature did not reach the battleline")
	}

	second := c.deal(testCreature)
	c.g.selectHandID(c.ctx, second)
	c.do(c.g.play)
	if c.g.phase != phaseFlank {
		t.Fatalf("the second creature did not ask for a flank; phase is %v", c.g.phase)
	}
	c.do(c.g.playFlank(true))
	if c.board()[0] != second {
		t.Errorf("the left flank holds %d, want %d", c.board()[0], second)
	}
}

// Escape backs out of the flank prompt, leaving the card in hand.
func TestCancellingTheFlankPrompt(t *testing.T) {
	c := newClient(t)
	c.manualTurn(testHouse)
	c.playFromHand(c.deal(testCreature))

	second := c.deal(testCreature)
	c.g.selectHandID(c.ctx, second)
	c.do(c.g.play)
	if c.g.phase != phaseFlank {
		t.Fatalf("phase is %v, want phaseFlank", c.g.phase)
	}
	c.press("Escape")
	if c.g.phase != phaseMain {
		t.Errorf("Escape left the phase at %v, want phaseMain", c.g.phase)
	}
	if !containsID(c.hand(), second) {
		t.Error("the cancelled creature left the hand")
	}
}

// While placing a creature, l and r commit a flank on their own rather than
// moving the selection, so a deliberate press needs no second key to confirm it.
func TestFlankKeys(t *testing.T) {
	for _, tt := range []struct {
		key  string
		left bool
	}{{"l", true}, {"r", false}} {
		t.Run(tt.key, func(t *testing.T) {
			c := newClient(t)
			c.manualTurn(testHouse)
			c.playFromHand(c.deal(testCreature))

			second := c.deal(testCreature)
			c.g.selectHandID(c.ctx, second)
			c.do(c.g.play)
			c.press(tt.key)

			board := c.board()
			at := 0
			if !tt.left {
				at = len(board) - 1
			}
			if board[at] != second {
				t.Errorf("%q put the creature at %v, want the %s flank",
					tt.key, board, map[bool]string{true: "left", false: "right"}[tt.left])
			}
		})
	}
}

func TestDiscardingFromHand(t *testing.T) {
	c := newClient(t)
	id := c.hand()[0]
	h := c.g.g.House(id)
	c.do(c.g.pickHouse(h))
	if c.g.phase != phaseMain {
		t.Fatalf("after choosing %v the phase is %v, want phaseMain", h, c.g.phase)
	}
	c.g.selectHandID(c.ctx, id)
	c.do(c.g.discard)
	if containsID(c.hand(), id) {
		t.Error("the discarded card is still in hand")
	}
	if !containsID(c.g.g.Discard(c.g.active()), id) {
		t.Error("the discarded card is not in the discard pile")
	}
}

// Dragging a card out of hand and dropping it on the board plays it; letting go
// anywhere else leaves it where it was.
func TestDragAndDrop(t *testing.T) {
	c := newClient(t)
	c.manualTurn(testHouse)
	id := c.deal(testCreature)

	c.g.startHandDrag(c.ctx, id)
	if !c.g.dragging || c.g.sel != id {
		t.Fatalf("the drag did not start: dragging=%v sel=%d", c.g.dragging, c.g.sel)
	}
	if c.g.hasHover {
		t.Error("the hover preview stayed up during the drag")
	}

	c.g.endHandDrag(c.ctx, id)
	if c.g.dragging {
		t.Error("releasing the pointer did not end the drag")
	}

	c.g.startHandDrag(c.ctx, id)
	c.g.dropOnBoard(c.ctx, nullEvent())
	c.settle()
	if c.g.dragging {
		t.Error("the drop did not end the drag")
	}
	if !containsID(c.board(), id) {
		t.Error("the dropped creature did not reach the battleline")
	}
}

// A drop with no drag behind it is a click on the board, not a play.
func TestDropWithoutADragDoesNothing(t *testing.T) {
	c := newClient(t)
	c.startTurn()
	before := len(c.hand())
	c.g.dropOnBoard(c.ctx, nullEvent())
	c.settle()
	if len(c.hand()) != before {
		t.Error("a stray drop played a card")
	}
}

// Ending a turn with moves left arms a confirmation first, so the turn is not
// thrown away by one click.
func TestEndTurnConfirms(t *testing.T) {
	c := newClient(t)
	c.startTurn()
	if !c.g.hasMoves() {
		t.Skip("the deal left the opening turn with nothing to do")
	}
	was := c.g.active()

	c.do(c.g.endTurn)
	if !c.g.confirmEndTurn {
		t.Fatal("ending a turn with moves left did not arm the confirmation")
	}
	if c.g.active() != was {
		t.Fatal("the turn ended without confirmation")
	}

	c.do(c.g.endTurn)
	if c.g.active() == was {
		t.Error("the confirmed end turn did not pass play on")
	}
	if c.g.phase != phaseHouse {
		t.Errorf("the new turn is at phase %v, want phaseHouse", c.g.phase)
	}
}

// Arming the end-turn confirm drops the current selection, so the red "Confirm
// end turn" is not read as ending the turn with the selected card, and the
// jiggling usable cards are what the eye lands on instead.
func TestArmingTheEndTurnConfirmClearsTheSelection(t *testing.T) {
	c := newClient(t)
	c.startTurn()
	if !c.g.hasMoves() {
		t.Skip("the deal left the opening turn with nothing to do")
	}
	id := c.hand()[0]
	c.g.selectHandID(c.ctx, id)
	if !c.g.hasSel {
		t.Fatal("the card was not selected")
	}

	c.do(c.g.endTurn)
	if !c.g.confirmEndTurn {
		t.Fatal("ending a turn with moves left did not arm the confirmation")
	}
	if c.g.hasSel {
		t.Error("arming the confirm left the card selected")
	}
}

// A fight grant (Brothers in Battle) lets an off-house creature fight, so the
// board makes it actionable and offers Fight alone — reaping stays barred out of
// house. Without the grant an off-house creature is inert.
func TestFightGrantOffersFightOnAnOffHouseCreature(t *testing.T) {
	c := newClient(t)
	c.g.g.State.ActiveHouse = engine.Brobnar
	c.g.phase = phaseMain
	p := c.g.active()
	att := c.g.g.AddToBattleline(
		engine.NewCard(
			"Off Fighter",
			engine.Untamed,
			engine.Creature,
			engine.Common,
			engine.WithPower(4),
		),
		p,
	)
	c.g.g.AddToBattleline(
		engine.NewCard("Foe", engine.Dis, engine.Creature, engine.Common, engine.WithPower(2)),
		1-p)

	if c.g.actionable(att, selYourCreature) {
		t.Fatal("an off-house creature is actionable without a grant")
	}

	c.g.g.State.MayFightHouse[p] = engine.Untamed
	if !c.g.actionable(att, selYourCreature) {
		t.Fatal("a fight-granted off-house creature is not actionable")
	}

	c.g.selectBoardID(c.ctx, att)
	acts, note := c.g.creatureCardActions()
	if note != "" {
		t.Fatalf("granted creature reported %q, want an action instead", note)
	}
	if len(acts) != 1 || acts[0].Label != "Fight" {
		t.Errorf("granted creature actions = %+v, want a single Fight", acts)
	}
}
func TestEndTurnOnlyFromMain(t *testing.T) {
	c := newClient(t)
	was := c.g.active()
	c.do(c.g.endTurn) // still at the house prompt
	if c.g.active() != was || c.g.confirmEndTurn {
		t.Error("the turn ended from the house prompt")
	}
}

// A creature is usable on the turn after the one it was played on, so reaping
// with a freshly played one is refused and says why.
func TestReaping(t *testing.T) {
	c := newClient(t)
	c.manualTurn(testHouse)
	id := c.deal(testCreature)
	c.playFromHand(id)

	c.g.selectBoardID(c.ctx, id)
	c.do(c.g.reap)
	if c.g.status == "" {
		t.Error("reaping with a creature played this turn reported nothing")
	}

	c.pass() // the opponent's turn
	c.manualTurn(testHouse)
	c.pass() // back to the creature's controller
	c.manualTurn(testHouse)

	before := c.g.g.State.Aember[c.g.active()]
	c.g.selectBoardID(c.ctx, id)
	c.do(c.g.reap)
	if c.g.status != "" {
		t.Fatalf("reaping reported %q", c.g.status)
	}
	if c.g.g.State.Aember[c.g.active()] <= before {
		t.Errorf("reaping did not gain Æmber: %d then %d",
			before, c.g.g.State.Aember[c.g.active()])
	}
	if !c.g.g.State.Cards[id].Exhausted {
		t.Error("reaping did not exhaust the creature")
	}
}

// A fight with more than one legal target puts the player into target selection,
// which Escape backs out of; picking a target resolves the fight.
func TestFightTargeting(t *testing.T) {
	c := newClient(t)
	c.manualTurn(testHouse)
	attacker := c.deal(testCreature)
	c.playFromHand(attacker)

	c.pass()
	c.manualTurn(testHouse)
	def1 := c.deal(testCreature)
	def2 := c.deal(testCreature)
	c.playFromHand(def1)
	c.playFromHand(def2)

	c.pass()
	c.manualTurn(testHouse)

	c.g.selectBoardID(c.ctx, attacker)
	c.do(c.g.startFight)
	if c.g.phase != phaseFightTarget {
		t.Fatalf("phase is %v, want phaseFightTarget (status %q)", c.g.phase, c.g.status)
	}
	if c.g.attacker != attacker {
		t.Errorf("the attacker is %d, want %d", c.g.attacker, attacker)
	}

	c.press("Escape")
	if c.g.phase != phaseMain || c.g.attacker != 0 {
		t.Errorf("Escape left phase=%v attacker=%d", c.g.phase, c.g.attacker)
	}

	c.g.selectBoardID(c.ctx, attacker)
	c.do(c.g.startFight)
	c.g.fightTargetID(c.ctx, def2)
	c.settle()
	if c.g.phase != phaseMain {
		t.Errorf("after the fight the phase is %v, want phaseMain", c.g.phase)
	}
	// Two creatures of equal power trade, so the target the player picked is the
	// one that died and the one they left alone is still there.
	opp := 1 - c.g.active()
	if containsID(c.g.g.Battleline(opp), def2) {
		t.Error("the creature the fight was aimed at survived")
	}
	if !containsID(c.g.g.Battleline(opp), def1) {
		t.Error("the creature the fight was not aimed at died")
	}
}

// A fight with exactly one legal target has nothing to choose, so it resolves
// straight away instead of asking for the only possible answer.
func TestFightWithASingleTargetSkipsSelection(t *testing.T) {
	c := newClient(t)
	c.manualTurn(testHouse)
	attacker := c.deal(testCreature)
	c.playFromHand(attacker)

	c.pass()
	c.manualTurn(testHouse)
	defender := c.deal(testCreature)
	c.playFromHand(defender)

	c.pass()
	c.manualTurn(testHouse)

	c.g.selectBoardID(c.ctx, attacker)
	c.do(c.g.startFight)
	if c.g.phase != phaseMain {
		t.Errorf("phase is %v, want phaseMain — the only target should not be asked for",
			c.g.phase)
	}
	if containsID(c.g.g.Battleline(1-c.g.active()), defender) {
		t.Error("the fight did not resolve against the only target there was")
	}
}

// A creature that cannot be used is told so before the player is put to the
// trouble of picking a target for it.
func TestFightRefusedBeforeTargeting(t *testing.T) {
	c := newClient(t)
	c.manualTurn(testHouse)
	id := c.deal(testCreature)
	c.playFromHand(id)

	c.g.selectBoardID(c.ctx, id)
	c.do(c.g.startFight)
	if c.g.phase == phaseFightTarget {
		t.Error("a creature played this turn was put into target selection")
	}
	if c.g.status == "" {
		t.Error("the refused fight reported nothing")
	}
}

func TestPlayTypeError(t *testing.T) {
	tests := []struct {
		name string
		in   error
		kind engine.CardType
		want string
	}{
		{"a barred type says which", engine.ErrCannotPlayType, engine.Tactic,
			"Tactic cards cannot be played"},
		{"any other error passes through", errors.New("no"), engine.Creature, "no"},
		{"no error stays nil", nil, engine.Creature, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := playTypeError(tt.in, tt.kind)
			if tt.want == "" {
				if got != nil {
					t.Fatalf("got %v, want nil", got)
				}
				return
			}
			if got == nil || got.Error() != tt.want {
				t.Errorf("got %v, want %q", got, tt.want)
			}
		})
	}
}

// A click on the board's background drops the selection.
func TestClickAway(t *testing.T) {
	c := newClient(t)
	c.startTurn()
	c.g.selectHandID(c.ctx, c.hand()[0])
	c.g.clickAway(c.ctx, nullEvent())
	c.settle()
	if c.g.hasSel {
		t.Error("a click on the background left the selection up")
	}
}

// With nothing selected, or while a mid-action phase owns the screen, a
// background click has nothing to do.
func TestClickAwayIsGuarded(t *testing.T) {
	c := newClient(t)
	c.startTurn()
	c.g.clickAway(c.ctx, nullEvent()) // nothing selected

	c.g.selectHandID(c.ctx, c.hand()[0])
	c.g.phase = phaseFlank
	c.g.clickAway(c.ctx, nullEvent())
	if !c.g.hasSel {
		t.Error("a background click dropped the selection mid-play")
	}
}

// testArtifact is an artifact whose Action ability needs no prompt, so a test can
// play it and use it without stopping on a choice.
const testArtifact = "Safe Place"

// An artifact goes to its own row rather than the battleline, and there is no
// flank to choose for it.
func TestPlayingAnArtifact(t *testing.T) {
	c := newClient(t)
	c.manualTurn(testHouse)
	id := c.deal(testArtifact)
	c.playFromHand(id)

	if !containsID(c.g.g.Artifacts(c.g.active()), id) {
		t.Fatal("the artifact did not reach the artifact row")
	}
	if c.g.phase != phaseMain {
		t.Errorf("playing an artifact left the phase at %v, want phaseMain", c.g.phase)
	}
}

// An artifact enters play exhausted, so its Action waits for the turn after the
// one it was played on — and then moves Æmber into it.
func TestUsingAnArtifactsAction(t *testing.T) {
	c := newClient(t)
	c.manualTurn(testHouse)
	id := c.deal(testArtifact)
	c.playFromHand(id)
	c.ownNextTurn(testHouse)

	me := c.g.active()
	c.g.g.State.Aember[me] = 2
	c.g.selectBoardID(c.ctx, id)
	c.do(c.g.useAction)

	if c.g.status != "" {
		t.Fatalf("using the artifact reported %q", c.g.status)
	}
	if got := c.g.g.AmberOn(id); got != 1 {
		t.Errorf("the artifact holds %d Æmber, want 1", got)
	}
	if !c.g.g.Exhausted(id) {
		t.Error("using the artifact did not exhaust it")
	}
}

// A card with no ability for the action it was asked to take says so rather than
// doing nothing.
func TestUsingWhatCannotBeUsed(t *testing.T) {
	c := newClient(t)
	c.manualTurn(testHouse)
	id := c.deal(testCreature)
	c.playFromHand(id)
	c.ownNextTurn(testHouse)

	c.g.selectBoardID(c.ctx, id)
	c.do(c.g.useAction)
	if c.g.status == "" {
		t.Error("using a creature with no Action ability reported nothing")
	}
}
