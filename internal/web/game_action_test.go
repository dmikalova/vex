package web

import (
	"testing"

	"github.com/dmikalova/vex/internal/engine"
)

// These tests cover the action plumbing: the undo history every root action
// records, the phase settled after one resolves, and the one-shot animations
// queued by diffing the board before and after.

func TestUndoAndRedo(t *testing.T) {
	c := newClient(t)
	c.manualTurn(testHouse)
	id := c.deal(testCreature)

	if c.g.canRedo() {
		t.Error("a match with nothing undone offers a redo")
	}
	c.playFromHand(id)
	if !containsID(c.board(), id) {
		t.Fatal("the creature is not on the battleline")
	}
	if !c.g.canUndo() {
		t.Fatal("a resolved action left nothing to undo")
	}

	c.do(c.g.undoAction)
	if containsID(c.board(), id) {
		t.Error("undo left the creature in play")
	}
	if !containsID(c.hand(), id) {
		t.Error("undo did not put the card back in hand")
	}
	if !c.g.canRedo() {
		t.Fatal("undo left nothing to redo")
	}

	c.do(c.g.redoAction)
	if !containsID(c.board(), id) {
		t.Error("redo did not put the creature back in play")
	}
}

// Undo is refused while an action is resolving with no prompt up: the state it
// would roll back to is not the one the player is looking at. A prompt is the
// exception — see TestUndoAtAPromptBacksTheActionOut.
func TestUndoIsGuarded(t *testing.T) {
	tests := []struct {
		name string
		arm  func(g *game)
	}{
		{"busy", func(g *game) { g.busy = true }},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := newClient(t)
			c.manualTurn(testHouse)
			c.playFromHand(c.deal(testCreature))
			depth := len(c.g.rootMarks)
			tt.arm(c.g)

			if c.g.canUndo() || c.g.canRedo() {
				t.Error("undo was offered while an action was resolving")
			}
			c.do(c.g.undoAction)
			c.do(c.g.redoAction)
			if len(c.g.rootMarks) != depth {
				t.Errorf("the undo history moved from %d to %d", depth, len(c.g.rootMarks))
			}
		})
	}
}

// Undo is offered at a prompt and rewinds the whole action that raised it: the
// prompt is part of that action, so there is no half-resolved state to stop at.
// Redo stays withheld, since the prompt's answers were never recorded.
func TestUndoAtAPromptBacksTheActionOut(t *testing.T) {
	for _, tt := range []struct {
		name string
		arm  func(g *game)
	}{
		{"choosing", func(g *game) { g.busy, g.choosing = true, true }},
		{"choosingOption", func(g *game) { g.busy, g.choosingOption = true, true }},
		{"choosingPosition", func(g *game) { g.busy, g.choosingPosition = true, true }},
	} {
		t.Run(tt.name, func(t *testing.T) {
			c := newClient(t)
			c.manualTurn(testHouse)
			c.playFromHand(c.deal(testCreature))
			tt.arm(c.g)

			if !c.g.canUndo() {
				t.Error("undo was withheld at a prompt")
			}
			if c.g.canRedo() {
				t.Error("redo was offered at a prompt")
			}
			c.do(c.g.undoAction)
			if !c.g.cancelling {
				t.Error("undo at a prompt did not back the action out")
			}
		})
	}
}

// A new action after an undo throws the redo away: the player has taken a
// different branch, and the one they left is not somewhere to go forward to.
func TestANewActionDropsTheRedo(t *testing.T) {
	c := newClient(t)
	c.manualTurn(testHouse)
	c.playFromHand(c.deal(testCreature))
	c.do(c.g.undoAction)
	if !c.g.canRedo() {
		t.Fatal("undo left nothing to redo")
	}
	c.playFromHand(c.deal(testCreature))
	if c.g.canRedo() {
		t.Error("a new action left the abandoned branch redoable")
	}
}

// The command log is the source of truth (ADR 0039), so it is not capped: a long
// game keeps every action undoable back to the deal rather than dropping the
// oldest steps.
func TestUndoHistoryIsNotCapped(t *testing.T) {
	c := newClient(t)
	c.manualTurn(testHouse)
	const n = 120
	for range n {
		c.g.adjustManualAmber(c.ctx, c.g.active(), 1)
		c.settle()
	}
	if len(c.g.rootMarks) < n {
		t.Errorf("the command log holds %d roots, want at least %d",
			len(c.g.rootMarks), n)
	}
}

func TestSettlePhase(t *testing.T) {
	c := newClient(t)
	c.startTurn()

	c.g.settlePhase()
	if c.g.phase != phaseMain {
		t.Errorf("a turn under way settles at %v, want phaseMain", c.g.phase)
	}

	// Start of turn: no house chosen yet and the engine still waits at the choice,
	// so the picker comes up.
	c.g.g.State.ActiveHouse = engine.HouseNone
	c.g.g.State.Phase = engine.PhaseChooseHouse
	c.g.settlePhase()
	if c.g.phase != phaseHouse {
		t.Errorf("a turn awaiting a house settles at %v, want phaseHouse", c.g.phase)
	}

	// A player locked out of every house chooses No House: ActiveHouse stays None
	// but the engine advances past the choice, so play continues (and the turn can
	// be ended) rather than looping back to the picker.
	c.g.g.State.Phase = engine.PhasePlay
	c.g.settlePhase()
	if c.g.phase != phaseMain {
		t.Errorf("a No-House turn settles at %v, want phaseMain", c.g.phase)
	}

	// A won game outranks the house prompt: the match is over whether or not the
	// turn it ended on ever chose a house.
	c.g.g.State.Winner = 0
	c.g.settlePhase()
	if c.g.phase != phaseOver {
		t.Errorf("a won game settles at %v, want phaseOver", c.g.phase)
	}
}

// A card that changed pulses, and the parity bit flips each time so the CSS
// animation restarts on a repeat rather than sitting finished.
func TestFlashesMarkWhatChanged(t *testing.T) {
	c := newClient(t)
	c.manualTurn(testHouse)
	id := c.deal(testCreature)
	c.playFromHand(id)

	if !c.g.flashes[id].enter {
		t.Error("a card that entered play did not pulse")
	}

	c.pass()
	c.manualTurn(testHouse)
	c.pass()
	c.manualTurn(testHouse)

	c.g.selectBoardID(c.ctx, id)
	c.do(c.g.reap)
	f := c.g.flashes[id]
	if !f.reap {
		t.Error("a creature that reaped did not pulse as reaping")
	}
	if !f.exhaust {
		t.Error("a creature that exhausted did not pulse as exhausting")
	}
	if !c.g.poolFlash[c.g.active()] {
		t.Error("the pool the reap paid into did not pulse")
	}
	odd := f.odd

	c.do(c.g.undoAction)
	if len(c.g.flashes) != 0 {
		t.Error("undo replayed the animation of the action it rolled back")
	}
	c.do(c.g.redoAction)
	c.g.selectBoardID(c.ctx, id)
	c.do(c.g.reap)
	if c.g.flashes[id].odd == odd {
		t.Error("a repeated pulse did not flip its parity, so it would not replay")
	}
}

// A card that leaves the board cannot pulse where it was, so it flies into the
// zone pill it landed in instead.
func TestFlightsAimAtTheLandingZone(t *testing.T) {
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

	if len(c.g.flights) != 2 {
		t.Fatalf("the trade queued %d flights, want both creatures", len(c.g.flights))
	}
	for _, f := range c.g.flights {
		if f.zone != "zone-discard" {
			t.Errorf("card %d flew to %q, want the discard pile", f.id, f.zone)
		}
	}
	if !c.g.discardFlash[c.g.active()] {
		t.Error("the discard pile a creature landed in did not pulse")
	}
}

func TestLandingFindsEachZone(t *testing.T) {
	c := newClient(t)
	c.manual()
	me := c.g.active()
	tests := []struct {
		name string
		dest engine.ManualZone
		want string
	}{
		{"discard", engine.ManualDiscard, "zone-discard"},
		{"purge", engine.ManualPurge, "zone-purge"},
		{"archives", engine.ManualArchives, "zone-archives"},
		{"hand", engine.ManualHand, "zone-hand"},
		{"top of deck", engine.ManualDeckTop, "zone-deck"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id := c.deal(testCreature)
			c.g.g.ManualMove(id, tt.dest)
			player, zone, ok := c.g.landing(id)
			if !ok {
				t.Fatalf("card %d landed nowhere", id)
			}
			if zone != tt.want {
				t.Errorf("card %d landed in %q, want %q", id, zone, tt.want)
			}
			if player != me {
				t.Errorf("card %d landed under player %d, want %d", id, player, me)
			}
		})
	}
}

// A card still in play has not landed anywhere, so nothing flies.
func TestLandingOfACardInPlay(t *testing.T) {
	c := newClient(t)
	c.manualTurn(testHouse)
	id := c.deal(testCreature)
	c.playFromHand(id)
	if _, _, ok := c.g.landing(id); ok {
		t.Error("a card in play reads as having landed out of play")
	}
}

// The very first action has no earlier state to diff against, so there is
// nothing to animate.
func TestNoFlashesWithoutASnapshot(t *testing.T) {
	c := newClient(t)
	c.g.prevValid = false
	c.g.computeFlashes()
	if len(c.g.flashes) != 0 {
		t.Errorf("computeFlashes queued %d pulses with no snapshot behind it",
			len(c.g.flashes))
	}
}

func TestInPlaySet(t *testing.T) {
	c := newClient(t)
	c.manualTurn(testHouse)
	if got := c.g.inPlaySet(); len(got) != 0 {
		t.Errorf("a fresh board holds %d cards in play, want none", len(got))
	}
	id := c.deal(testCreature)
	c.playFromHand(id)
	if !c.g.inPlaySet()[id] {
		t.Error("a creature on the battleline is not in the in-play set")
	}
}

// A status message clears itself, and a newer one is not wiped by an older
// one's timer.
func TestStatusMessages(t *testing.T) {
	c := newClient(t)
	c.g.setStatus("first")
	if c.g.status != "first" {
		t.Fatalf("the status is %q, want %q", c.g.status, "first")
	}
	gen := c.g.statusGen
	c.g.setStatus("second")
	if c.g.statusGen == gen {
		t.Error("a second message did not take over the auto-clear")
	}
	c.g.setStatus("")
	if c.g.status != "" {
		t.Error("clearing the status left a message up")
	}
}

// A corrupt state can panic mid-action. Rather than freeze the UI on
// "resolving…", the action peels its record back off the command log and replays
// what remains, then says what happened.
func TestACrashedActionRollsBack(t *testing.T) {
	c := newClient(t)
	c.manualTurn(testHouse)
	c.playFromHand(c.deal(testCreature))
	before := c.g.g.State
	depth := len(c.g.rootMarks)

	// A real root records its input before running, so mimic that and then crash;
	// the rollback should peel the recorded root back off.
	c.g.record(input{Kind: inEndTurn})
	c.g.runAction(c.ctx, func() error { panic("a corrupt state") })
	c.settle()

	if c.g.busy {
		t.Error("the crashed action left the UI resolving")
	}
	if c.g.status == "" {
		t.Error("the crashed action reported nothing")
	}
	if c.g.g.State != before {
		t.Error("the crashed action did not roll the state back")
	}
	if len(c.g.rootMarks) != depth {
		t.Errorf("the crashed action left the history at %d, want %d",
			len(c.g.rootMarks), depth)
	}
}

// Only one action runs at a time: a second click while one is resolving is
// dropped rather than queued behind it.
func TestASecondActionWhileBusyIsDropped(t *testing.T) {
	c := newClient(t)
	c.manualTurn(testHouse)
	c.g.busy = true
	ran := false
	c.g.runAction(c.ctx, func() error { ran = true; return nil })
	c.settle()
	if ran {
		t.Error("a second action ran while one was already resolving")
	}
}
