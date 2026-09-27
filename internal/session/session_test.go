package session

import (
	"errors"
	"testing"

	"github.com/dmikalova/vex/internal/engine"
)

// setup deals a bare game for the session tests — no decks, just the engine
// harness the driving action pulls choices from.
func setup(seed int64, _ [2]string) *engine.Game {
	return engine.NewGame("A", "B", seed)
}

// drive is a scripted action: pick one of two cards, draw from the PRNG (an
// information barrier), then choose an option. It folds each answer into the Æmber
// pools so a test can compare states by value after a replay.
func drive(g *engine.Game) {
	id, _ := g.ChooseCreature(0, 0, "pick one", []engine.LocalID{1, 2})
	g.State.Aember[0] = int16(id)
	g.State.Aember[1] = int16(g.State.PRNG.Intn(100))
	i := g.ChooseOption(0, 0, "choose", []string{"x", "y"})
	g.State.Aember[0] += int16(i * 10)
}

// goblin is the one card the test pool holds — a Brobnar creature, so playing it is
// legal exactly when Brobnar is the active house.
var goblin = engine.CardDefinition{
	Name:  "Test Goblin",
	House: engine.Brobnar,
	Type:  engine.Creature,
	Power: 3,
}

// resolveCard is the card pool a session's manual add-card edit looks names up in.
func resolveCard(name string) (engine.CardDefinition, bool) {
	if name == goblin.Name {
		return goblin, true
	}
	return engine.CardDefinition{}, false
}

func newSession() *Session { return New(7, [2]string{}, setup, drive, resolveCard) }

// mustApply applies cmd and fails the test if it is rejected.
func mustApply(t *testing.T, s *Session, cmd engine.Command) {
	t.Helper()
	if err := s.Apply(cmd); err != nil {
		t.Fatalf("apply %+v: %v", cmd, err)
	}
}

// A session records each answer, drives the action to completion, and flags the
// step that drew from the PRNG as crossing an information barrier.
func TestSessionApplyRecordsAndAdvances(t *testing.T) {
	s := newSession()

	req, done := s.Pending()
	if done || req.Kind != engine.RequestPickCard {
		t.Fatalf("first request = %+v done=%v, want a card pick", req, done)
	}

	if err := s.Apply(engine.Command{
		Kind: engine.CommandPickCard,
		Card: 2,
	}); err != nil {
		t.Fatalf("apply card pick: %v", err)
	}
	req, done = s.Pending()
	if done || req.Kind != engine.RequestOption {
		t.Fatalf("after card pick got %+v done=%v, want an option request", req, done)
	}

	if err := s.Apply(engine.Command{
		Kind:  engine.CommandOption,
		Index: 1,
	}); err != nil {
		t.Fatalf("apply option: %v", err)
	}
	if _, done := s.Pending(); !done {
		t.Fatal("the action did not finish after both answers")
	}

	if s.Len() != 2 {
		t.Fatalf("recorded %d commands, want 2", s.Len())
	}
	if g := s.Game(); g.State.Aember[0] != 12 {
		t.Fatalf("Aember[0] = %d, want 12 (card 2 + option 1*10)", g.State.Aember[0])
	}
	// The PRNG step happened while resolving the first command, not the second.
	if !s.CrossesBarrier(0) {
		t.Error("undo past the first command should cross a barrier")
	}
	if s.CrossesBarrier(1) {
		t.Error("undo past only the second command should not cross a barrier")
	}
}

// Apply rejects a command that does not answer the pending request, and refuses
// any command once the match is finished.
func TestSessionApplyRejects(t *testing.T) {
	s := newSession()
	// The pending request is a card pick; an option command does not answer it.
	if err := s.Apply(
		engine.Command{
			Kind:  engine.CommandOption,
			Index: 0,
		},
	); !errors.Is(
		err,
		ErrIllegal,
	) {
		t.Fatalf("apply wrong-kind command: got %v, want ErrIllegal", err)
	}
	// A card not among the candidates is also illegal.
	if err := s.Apply(
		engine.Command{
			Kind: engine.CommandPickCard,
			Card: 9,
		},
	); !errors.Is(
		err,
		ErrIllegal,
	) {
		t.Fatalf("apply out-of-set card: got %v, want ErrIllegal", err)
	}

	mustApply(t, s, engine.Command{
		Kind: engine.CommandPickCard,
		Card: 1,
	})
	mustApply(t, s, engine.Command{
		Kind:  engine.CommandOption,
		Index: 0,
	})
	if err := s.Apply(engine.Command{Kind: engine.CommandOption}); !errors.Is(err, ErrFinished) {
		t.Fatalf("apply after done: got %v, want ErrFinished", err)
	}
}

// Undo rewinds by replaying the kept prefix from a fresh deal, so the state after
// Undo(n) equals the state a fresh session reaches applying the same n commands.
func TestSessionUndoIsReplay(t *testing.T) {
	full := newSession()
	mustApply(t, full, engine.Command{
		Kind: engine.CommandPickCard,
		Card: 2,
	})
	mustApply(t, full, engine.Command{
		Kind:  engine.CommandOption,
		Index: 1,
	})

	if err := full.Undo(1); err != nil {
		t.Fatalf("undo: %v", err)
	}
	if full.Len() != 1 {
		t.Fatalf("after undo Len = %d, want 1", full.Len())
	}

	ref := newSession()
	mustApply(t, ref, engine.Command{
		Kind: engine.CommandPickCard,
		Card: 2,
	})
	if full.Game().State != ref.Game().State {
		t.Fatal("undo did not reproduce the state of the same prefix replayed fresh")
	}

	if err := full.Undo(-1); !errors.Is(err, ErrRange) {
		t.Fatalf("undo(-1): got %v, want ErrRange", err)
	}
	if err := full.Undo(99); !errors.Is(err, ErrRange) {
		t.Fatalf("undo(99): got %v, want ErrRange", err)
	}
}

// The record round-trips: a session loaded from another's record replays to the
// exact same state, and a mismatched version is refused.
func TestSessionRecordRoundTrip(t *testing.T) {
	s := newSession()
	mustApply(t, s, engine.Command{
		Kind: engine.CommandPickCard,
		Card: 2,
	})
	mustApply(t, s, engine.Command{
		Kind:  engine.CommandOption,
		Index: 1,
	})

	rec := s.Record()
	if rec.Version != Version || rec.Seed != 7 || len(rec.Commands) != 2 {
		t.Fatalf("record = %+v, want version %d seed 7 with 2 commands", rec, Version)
	}

	loaded, err := Load(rec, setup, drive, resolveCard)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if loaded.Game().State != s.Game().State {
		t.Fatal("a session loaded from the record replayed to a different state")
	}

	rec.Version = Version + 1
	if _, err := Load(rec, setup, drive, resolveCard); !errors.Is(err, ErrVersion) {
		t.Fatalf("load with wrong version: got %v, want ErrVersion", err)
	}
}

// View projects the current state for a viewer through the engine seam (identity
// today).
func TestSessionView(t *testing.T) {
	s := newSession()
	mustApply(t, s, engine.Command{
		Kind: engine.CommandPickCard,
		Card: 2,
	})
	v := s.View(1)
	if v.Viewer != 1 || v.State != s.Game().State {
		t.Fatalf("View(1) = %+v, want the identity projection for viewer 1", v)
	}
}

// match drives a bare game through the engine's canonical turn loop, so the session
// tests that care about root actions see the same RequestAction a real match yields.
// The game has no decks, which is the point: an empty hand makes the legal set small
// and predictable, so a manual edit that widens it is unmistakable.
func match(g *engine.Game) { g.RunMatch() }

func newMatchSession() *Session { return New(7, [2]string{}, setup, match, resolveCard) }

// startTurn answers the match's setup decisions — who goes first, then each
// player's mulligan prompt — and names Brobnar as the active house, leaving the
// session in the play phase with a pending action request.
func startTurn(t *testing.T, s *Session) {
	t.Helper()
	mustApply(t, s, engine.Command{
		Kind:   engine.CommandSetFirstPlayer,
		Player: 0,
		Index:  engine.RolledFirstPlayer,
	})
	for {
		req, done := s.Pending()
		if done {
			t.Fatal("the match finished before reaching the first action request")
		}
		if req.Kind == engine.RequestAction {
			break
		}
		// Every setup prompt before the first action is a mulligan; keep the hand.
		mustApply(t, s, engine.Command{Kind: engine.CommandOption})
	}
	mustApply(t, s, engine.Command{
		Kind:  engine.CommandChooseHouse,
		House: engine.Brobnar,
	})
}

// addGoblin is the manual force-edit the tests below record: it drops a Brobnar
// creature into player 0's hand, which no rule would have put there.
var addGoblin = engine.Command{
	Kind:   engine.CommandManualAddCard,
	Player: 0,
	Name:   goblin.Name,
}

// playGoblin plays the manually added creature from hand. The battleline is empty,
// so the two flanks coincide and the play names neither.
var playGoblin = engine.Command{
	Kind: engine.CommandPlayCreature,
	Hand: 0,
}

// A manual force-edit is applied to the live game without consuming the pending
// request, is recorded in the log, and is replayed by Load — so a match holding a
// force-edit still rebuilds from its command log alone.
func TestSessionApplyManualIsRecordedAndReplayed(t *testing.T) {
	s := newMatchSession()
	startTurn(t, s)
	setupSteps := s.Len()

	before, _ := s.Pending()
	if err := s.ApplyManual(addGoblin); err != nil {
		t.Fatalf("apply manual add: %v", err)
	}
	if after, done := s.Pending(); done || after.Kind != before.Kind {
		t.Fatalf("manual edit changed the pending request to %+v done=%v", after, done)
	}
	if got := len(s.Game().Hand(0)); got != 1 {
		t.Fatalf("hand holds %d cards after the manual add, want 1", got)
	}
	if s.Len() != setupSteps+1 {
		t.Fatalf("recorded %d commands, want %d (the manual edit is recorded)",
			s.Len(), setupSteps+1)
	}
	if s.CrossesBarrier(setupSteps) {
		t.Error("a manual edit crosses no information barrier")
	}

	loaded, err := Load(s.Record(), setup, match, resolveCard)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if loaded.Game().State != s.Game().State {
		t.Fatal("a session loaded from a record holding a manual edit replayed to a " +
			"different state")
	}

	// A record naming a card the pool does not know fails to load, so a diverged
	// replay stops where it diverged instead of silently dropping the edit.
	rec := s.Record()
	rec.Commands[setupSteps].Name = "Not In The Pool"
	if _, err := Load(rec, setup, match, resolveCard); err == nil {
		t.Fatal("loading a record naming an unknown card should fail")
	}
}

// Undo rewinds across a manual edit by replaying the kept prefix, so a prefix
// holding a force-edit rebuilds the edited board and a prefix stopping short of one
// rebuilds the board without it.
func TestSessionUndoAcrossAManualEdit(t *testing.T) {
	s := newMatchSession()
	startTurn(t, s)
	setupSteps := s.Len()
	if err := s.ApplyManual(addGoblin); err != nil {
		t.Fatalf("apply manual add: %v", err)
	}
	mustApply(t, s, playGoblin)

	// Keeping the edit but dropping the play leaves the goblin in hand.
	if err := s.Undo(setupSteps + 1); err != nil {
		t.Fatalf("undo to just after the edit: %v", err)
	}
	if got := len(s.Game().Hand(0)); got != 1 {
		t.Fatalf("hand holds %d cards after undoing to just after the edit, want 1", got)
	}

	// Rewinding past the edit undoes it too: nothing ever put that card in the deck.
	if err := s.Undo(setupSteps); err != nil {
		t.Fatalf("undo past the edit: %v", err)
	}
	if got := len(s.Game().Hand(0)); got != 0 {
		t.Fatalf("hand holds %d cards after undoing past the edit, want 0", got)
	}
}

// Apply validates a root action against the LIVE legal set, not the one the request
// carried: playing the manually added creature is legal now even though it was not
// when the engine gathered the actions for the pending request.
func TestSessionApplyChecksRootActionsLive(t *testing.T) {
	s := newMatchSession()
	startTurn(t, s)

	req, _ := s.Pending()
	if req.IsLegal(playGoblin) {
		t.Fatal("playing from an empty hand should not be in the cached action set")
	}
	if err := s.Apply(playGoblin); !errors.Is(err, ErrIllegal) {
		t.Fatalf("play before the manual add: got %v, want ErrIllegal", err)
	}

	if err := s.ApplyManual(addGoblin); err != nil {
		t.Fatalf("apply manual add: %v", err)
	}
	// The cached request is now stale — it still offers only the actions an empty
	// hand allowed — so only a live check can accept the play.
	if stale, _ := s.Pending(); stale.IsLegal(playGoblin) {
		t.Fatal("the cached request should still be stale after the manual add")
	}
	mustApply(t, s, playGoblin)
	if got := len(s.Game().Battleline(0)); got != 1 {
		t.Fatalf("battleline holds %d creatures, want 1", got)
	}
}

// boom answers one option request and then panics, standing in for a card whose
// resolution breaks mid-action.
func boom(g *engine.Game) {
	g.ChooseOption(0, 0, "boom", []string{"x", "y"})
	panic("boom")
}

// A panicking action is reported to the caller as an error rather than hanging the
// session, and the command that triggered it stays recorded so the caller can roll
// the broken step back with Undo.
func TestSessionApplyReportsAPanickingAction(t *testing.T) {
	s := New(7, [2]string{}, setup, boom, resolveCard)

	err := s.Apply(engine.Command{Kind: engine.CommandOption})
	if err == nil {
		t.Fatal("applying the command that panics should return an error")
	}
	var panicErr *engine.PanicError
	if !errors.As(err, &panicErr) {
		t.Fatalf("apply = %v, want an engine.PanicError", err)
	}
	if len(panicErr.Stack) == 0 {
		t.Error("the panic error should carry the stack captured at recover time")
	}
	if _, done := s.Pending(); !done {
		t.Error("a panicking action leaves the session done, not awaiting a request")
	}
	if s.Len() != 1 {
		t.Fatalf("recorded %d commands, want the broken step recorded for Undo", s.Len())
	}
	if !errors.As(s.Err(), &panicErr) {
		t.Fatalf("Err = %v, want the same panic error", s.Err())
	}

	if err := s.Undo(0); err != nil {
		t.Fatalf("undo past the broken step: %v", err)
	}
	if _, done := s.Pending(); done {
		t.Error("after undoing the broken step the session should await the request again")
	}
}
