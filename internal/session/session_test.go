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

func newSession() *Session { return New(7, [2]string{}, setup, drive) }

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

	loaded, err := Load(rec, setup, drive)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if loaded.Game().State != s.Game().State {
		t.Fatal("a session loaded from the record replayed to a different state")
	}

	rec.Version = Version + 1
	if _, err := Load(rec, setup, drive); !errors.Is(err, ErrVersion) {
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
