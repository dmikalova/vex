package web

import (
	"errors"
	"testing"

	"github.com/dmikalova/vex/internal/engine"
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

// lockOutHouses bars the active player from every one of their deck houses, so
// the constraint table (ADR 0035) leaves No House the only legal choice this turn
// — the situation Tezmal creates.
func lockOutHouses(c *client) {
	p := c.g.active()
	for _, h := range c.g.deckHouses[p] {
		c.g.g.CannotChooseHouseNextTurn(p, h, 0)
	}
	// CannotChooseHouseNextTurn arms the next turn; promote the armed table onto the
	// current one so the lockout binds the choice the picker is about to offer.
	c.g.g.State.HouseConstraints[p] = c.g.g.State.HouseConstraintsNext[p]
	c.g.g.State.HouseConstraintCount[p] = c.g.g.State.HouseConstraintCountNext[p]
}

// A player locked out of every house is offered a single "No House" button, in
// manual mode as well as ordinary play — offering a barred house would only be
// rejected by ChooseHouse. Choosing No House advances to the main phase so the
// turn can be ended.
func TestLockedOutOfEveryHouseOffersNoHouse(t *testing.T) {
	for _, manual := range []bool{false, true} {
		name := "ordinary"
		if manual {
			name = "manual"
		}
		t.Run(name, func(t *testing.T) {
			c := newClient(t)
			if manual {
				c.manual()
			}
			lockOutHouses(c)
			buttons := c.g.houseButtons()
			if len(buttons) != 1 || buttons[0] != engine.HouseNone {
				t.Fatalf("the picker offers %v, want [No House]", buttons)
			}
			c.do(c.g.pickHouse(engine.HouseNone))
			if c.g.g.State.ActiveHouse != engine.HouseNone {
				t.Errorf("the active house is %v, want No House", c.g.g.State.ActiveHouse)
			}
			if c.g.phase != phaseMain {
				t.Errorf("after choosing No House the phase is %v, want phaseMain", c.g.phase)
			}
		})
	}
}

// via undo. On the game's first turn there is nothing to step back to, so Back is
// greyed out; once a turn has been taken, Back undoes it.
func TestHousePickerUndo(t *testing.T) {
	c := newClient(t)
	if c.g.phase != phaseHouse {
		t.Fatalf("a fresh deal is at phase %v, want phaseHouse", c.g.phase)
	}
	c.g.rootMarks = nil
	c.wants("the first-turn house picker", `title="Undo"`)
	if c.g.canUndo() {
		t.Error("a first-turn house pick has nothing to step back to")
	}

	// Choose a house and end the turn; the opponent then faces the house picker with
	// a turn to step back to, so Undo is live and undoing it returns to that turn.
	c.do(c.g.pickHouse(c.g.pickableHouses()[0]))
	c.pass()
	c.await("the opponent's house choice", func() bool {
		return !c.g.busy && c.g.phase == phaseHouse
	})
	if !c.g.canUndo() {
		t.Fatal("the opponent's house pick has a turn to step back to")
	}
	c.do(c.g.undoAction)
	if c.g.phase != phaseMain {
		t.Errorf("Undo from the house picker left the phase at %v, want phaseMain", c.g.phase)
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

// A Deploy creature skips the flank prompt and instead raises the click-to-place
// position prompt: the player first chooses a side (Deploy left / Deploy right),
// then clicks a battleline creature to land beside it on that side. The ends of
// the line are the flanks. These drive the answer handlers directly against a
// staged prompt (the state ChoosePosition posts), so no background goroutine is
// needed to test the click-to-position map.
func TestDeployPlacementByClick(t *testing.T) {
	stage := func(t *testing.T) (*client, []engine.LocalID) {
		t.Helper()
		c := newClient(t)
		c.manualTurn(testHouse)
		c.playFromHand(c.deal(testCreature))
		c.playFromHand(c.deal(testCreature))
		line := c.board()
		c.g.choosingPosition = true
		c.g.positionLine = line
		c.g.positionRight = false
		c.g.positionSideChosen = false
		return c, line
	}
	answered := func(t *testing.T, c *client) int {
		t.Helper()
		select {
		case pos := <-c.g.chooser.positionReply:
			return pos
		default:
			t.Fatal("the placement handler sent no position")
			return -1
		}
	}

	t.Run("choosing left then clicking a creature deploys to its left", func(t *testing.T) {
		c, line := stage(t)
		c.do(c.g.chooseDeploySide(false))
		c.g.choosePositionCandidate(c.ctx, line[1]) // before the second creature
		if got := answered(t, c); got != 1 {
			t.Errorf("deploy-left of index 1 = position %d, want 1", got)
		}
	})

	t.Run("choosing right places after the clicked creature", func(t *testing.T) {
		c, line := stage(t)
		c.do(c.g.chooseDeploySide(true))
		c.g.choosePositionCandidate(c.ctx, line[0]) // after the first creature
		if got := answered(t, c); got != 1 {
			t.Errorf("deploy-right of index 0 = position %d, want 1", got)
		}
	})

	t.Run("deploy right of the last creature is the right flank", func(t *testing.T) {
		c, line := stage(t)
		c.do(c.g.chooseDeploySide(true))
		c.g.choosePositionCandidate(c.ctx, line[len(line)-1])
		if got := answered(t, c); got != len(line) {
			t.Errorf("deploy-right of the last creature = position %d, want %d",
				got, len(line))
		}
	})

	t.Run("clicking a creature before a side is chosen is ignored", func(t *testing.T) {
		c, line := stage(t)
		c.g.choosePositionCandidate(c.ctx, line[0])
		select {
		case <-c.g.chooser.positionReply:
			t.Fatal("a click before choosing a side answered the prompt")
		default:
		}
	})

	t.Run("Back returns to the side choice", func(t *testing.T) {
		c, _ := stage(t)
		c.do(c.g.chooseDeploySide(true))
		if !c.g.positionSideChosen {
			t.Fatal("choosing a side did not mark it chosen")
		}
		c.do(c.g.deploySideBack)
		if c.g.positionSideChosen {
			t.Error("Back did not return to the side choice")
		}
	})

	t.Run("clicking a creature not on the line is ignored", func(t *testing.T) {
		c, _ := stage(t)
		c.do(c.g.chooseDeploySide(false))
		offLine := c.deal(testCreature) // in hand, not on the battleline
		c.g.choosePositionCandidate(c.ctx, offLine)
		select {
		case <-c.g.chooser.positionReply:
			t.Fatal("an off-line click answered the prompt")
		default:
		}
	})
}

// With no other friendly creatures in play a Deploy creature has only one
// placement, so it is placed without asking: ChoosePosition answers itself with
// position 0 and never raises the prompt.
func TestDeploySkipsPromptWithNoOtherCreatures(t *testing.T) {
	c := newClient(t)
	c.manualTurn(testHouse)
	if pos := c.g.chooser.ChoosePosition("Ranger", "deploy Ranger", nil); pos != 0 {
		t.Errorf("empty-line deploy = position %d, want 0", pos)
	}
	if c.g.choosingPosition {
		t.Error("an empty line raised the deploy placement prompt")
	}
}

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
					tt.key, board, map[bool]string{
						true:  "left",
						false: "right",
					}[tt.left])
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

// On the first turn a player takes exactly one action. After that opening volition
// every other hand card must read as non-discardable, exactly as it reads as
// non-playable: the discard offer is gated by the engine's CanDiscard, not a
// partial active-house check that ignored the first-turn limit.
func TestFirstTurnDiscardRestrictionMatchesPlay(t *testing.T) {
	c := newClient(t)
	// Choose an active house that has at least two cards in the opening hand, so the
	// test has one card to spend the opening action on and a second to check the
	// first-turn limit bars afterwards. Which house that is follows from the deal
	// seed, so the test reads it off the hand rather than pinning a house name.
	var active engine.House
	var same []engine.LocalID
	for _, h := range c.g.pickableHouses() {
		var inHouse []engine.LocalID
		for _, id := range c.hand() {
			if c.g.g.Def(id).House == h {
				inHouse = append(inHouse, id)
			}
		}
		if len(inHouse) >= 2 {
			active, same = h, inHouse
			break
		}
	}
	if len(same) < 2 {
		t.Fatalf("no pickable house has two cards in the seeded first hand")
	}
	c.do(c.g.pickHouse(active))
	if c.g.phase != phaseMain {
		t.Fatalf("after choosing %v the phase is %v, want phaseMain", active, c.g.phase)
	}
	toDiscard, other := same[0], same[1]

	// Before the opening action an active-house card can be discarded.
	if !c.g.discardableFromHand(other) {
		t.Fatal("an active-house card was not discardable before the first action")
	}

	// Spend the one first-turn action on a discard.
	c.g.selectHandID(c.ctx, toDiscard)
	c.do(c.g.discard)

	// The first-turn limit now bars everything else, discard included.
	if c.g.discardableFromHand(other) {
		t.Error("a hand card is still discardable after the first-turn action")
	}
	if c.g.usableFromHand(other) {
		t.Error("a barred hand card still reads as usable")
	}
	if c.g.hasMoves() {
		t.Error("the turn still reports moves after its one action")
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

// An out-of-house artifact offers no Action at all — the same way an out-of-house
// creature offers no reap or fight — rather than offering it and then rejecting
// the use. It asks the engine (CanUseArtifact carries the house check) instead of
// inferring the answer.
func TestOutOfHouseArtifactOffersNoAction(t *testing.T) {
	c := newClient(t)
	c.manualTurn(testHouse) // Untamed
	id := c.deal(testArtifact)
	c.playFromHand(id) // Safe Place is Shadows, played out of house in manual mode
	c.ownNextTurn(testHouse)
	c.do(c.g.toggleManual) // leave manual mode, so the house rules apply again
	if c.g.g.Manual() {
		t.Fatal("manual mode did not turn off")
	}

	c.g.selectBoardID(c.ctx, id)
	acts, note := c.g.selActions()
	for _, a := range acts {
		if a.Label == "Action" {
			t.Fatalf("an out-of-house artifact still offered its Action (note %q)", note)
		}
	}
	if note == "" {
		t.Error("an out-of-house artifact gave no reason for offering nothing")
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

// asUpgradeCreature is a creature that may be played either as a creature or as
// an upgrade (WithPlayableAsUpgrade), so selecting it with a host in play offers
// the two explicit play buttons.
const asUpgradeCreature = "CALV-1N"

// A creature that may go down as a creature or as an upgrade, with a host in
// play, offers two explicit buttons — Play creature and Play upgrade — in place
// of the single Play, so the choice is made up front rather than by a later
// sidebar prompt.
func TestPlayableAsUpgradeOffersBothPlayButtons(t *testing.T) {
	c := newClient(t)
	c.manualTurn(testHouse)
	c.playFromHand(c.deal(asUpgradeCreature)) // a host in play
	c.g.selectHandID(c.ctx, c.deal(asUpgradeCreature))

	acts, note := c.g.handCardActions()
	if note != "" {
		t.Fatalf("the play offer carried a note: %q", note)
	}
	var haveCreature, haveUpgrade bool
	for _, a := range acts {
		switch a.Label {
		case "Play creature":
			haveCreature = true
		case "Play upgrade":
			haveUpgrade = true
		case "Play":
			t.Error("the single Play button is still offered alongside the split ones")
		}
	}
	if !haveCreature || !haveUpgrade {
		t.Errorf("play offer is missing a split button: creature=%v upgrade=%v",
			haveCreature, haveUpgrade)
	}
}

// With no host in play the creature-or-upgrade choice does not arise, so an
// ordinary single Play button is offered.
func TestPlayableAsUpgradeWithoutAHostOffersPlainPlay(t *testing.T) {
	c := newClient(t)
	c.manualTurn(testHouse)
	c.g.selectHandID(c.ctx, c.deal(asUpgradeCreature))

	if c.g.canPlayAsUpgrade() {
		t.Fatal("a creature-as-upgrade with no host in play offered the split")
	}
	acts, _ := c.g.handCardActions()
	var havePlain bool
	for _, a := range acts {
		if a.Label == "Play" {
			havePlain = true
		}
		if a.Label == "Play creature" || a.Label == "Play upgrade" {
			t.Errorf("a split button %q was offered with no host in play", a.Label)
		}
	}
	if !havePlain {
		t.Error("no plain Play button was offered")
	}
}

// The armed Play creature / Play upgrade choice answers the engine's "as a
// creature or an upgrade?" prompt for it, and only for it: an unarmed choice or a
// different prompt is left for the player.
func TestArmedUpgradeChoiceAnswersOnlyItsPrompt(t *testing.T) {
	c := newClient(t)
	pair := []string{"Creature", "Upgrade"}

	c.g.upgradeChoice = choiceUpgrade
	if i, ok := c.g.chooser.armedUpgradeChoice(pair); !ok || i != 1 {
		t.Errorf("upgrade answered (%d,%v), want (1,true)", i, ok)
	}
	c.g.upgradeChoice = choiceCreature
	if i, ok := c.g.chooser.armedUpgradeChoice(pair); !ok || i != 0 {
		t.Errorf("creature answered (%d,%v), want (0,true)", i, ok)
	}
	c.g.upgradeChoice = choiceNone
	if _, ok := c.g.chooser.armedUpgradeChoice(pair); ok {
		t.Error("an unarmed choice answered the prompt")
	}
	c.g.upgradeChoice = choiceUpgrade
	if _, ok := c.g.chooser.armedUpgradeChoice([]string{"Yes", "No"}); ok {
		t.Error("an armed choice answered an unrelated prompt")
	}
}

// An ordinary play clears any armed creature-as-upgrade choice, so a later prompt
// is never answered by a stale one.
func TestPlayClearsAnArmedUpgradeChoice(t *testing.T) {
	c := newClient(t)
	c.manualTurn(testHouse)
	c.g.upgradeChoice = choiceUpgrade
	c.g.selectHandID(c.ctx, c.deal(testCreature))
	c.do(c.g.play)
	if c.g.upgradeChoice != choiceNone {
		t.Errorf("an ordinary play left the armed choice at %v, want choiceNone",
			c.g.upgradeChoice)
	}
}
