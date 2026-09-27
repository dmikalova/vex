// Package session drives one Vex match as an event-sourced replay (ADR 0039,
// 0040). The authoritative record is not a game state but the ordered log of
// commands the players entered, together with the seed and sets the match was
// dealt from. State, the typed game log, and undo are all DERIVED by replaying
// that command log from a fresh game through the engine's suspendable step
// function.
//
// The session owns {version, seed, sets, []Command} and the undo cursor, wraps the
// engine, and exposes apply, undo, and view. It knows nothing about how a match is
// dealt (that is internal/match) or how the turn loop is shaped (that is the
// caller's driving action) — it only records the answers to the requests the
// engine yields and can rebuild the exact game by feeding them back.
package session

import (
	"errors"
	"fmt"
	"slices"

	"github.com/dmikalova/vex/internal/engine"
)

// Version tags the on-disk record. A record from a different version is refused
// rather than misread — there is no forward compatibility (ADR 0039).
const Version = 1

// Record is the persisted match: the version that wrote it, the deal seed, the
// per-player set names, and the ordered command log. Replaying the commands from a
// game dealt with seed and sets reproduces the exact match — nothing else needs to
// be stored, because state and log are projections of this.
type Record struct {
	Version  int
	Seed     int64
	Sets     [2]string
	Commands []engine.Command
}

// Setup deals a fresh game for a seed and set pair. The caller supplies it
// (wrapping internal/match) so the session does not depend on deck generation and
// stays a pure driver.
type Setup func(seed int64, sets [2]string) *engine.Game

// Action is the top-level play the session drives through the engine's suspendable
// step function — the whole-game turn loop for a real match, or a scripted
// sequence in a test. It must be deterministic given the game and the answers it
// is fed, so replaying the same commands reproduces the same state.
type Action func(*engine.Game)

// ResolveCard maps a card name to its definition. The session takes one at
// construction because the one manual force-edit that adds a card
// (engine.CommandManualAddCard) needs the card pool, which the engine deliberately
// does not hold (cards import engine, not the reverse), so the caller that owns the
// pool injects the lookup. It must be stable across a replay: a record whose log
// names a card the pool no longer knows fails to load rather than diverging.
type ResolveCard func(name string) (engine.CardDefinition, bool)

var (
	// ErrFinished is returned by Apply after the driving action has finished.
	ErrFinished = errors.New("session: the match has finished")
	// ErrIllegal is returned when a command does not answer the pending request.
	ErrIllegal = errors.New("session: command does not answer the pending request")
	// ErrRange is returned when an undo cursor is out of the applied-command range.
	ErrRange = errors.New("session: undo cursor out of range")
	// ErrVersion is returned when a record's version does not match this build.
	ErrVersion = errors.New("session: record version mismatch")
)

// Session drives one match. It holds the command log and the undo cursor and
// derives everything else by replay.
type Session struct {
	seed        int64
	sets        [2]string
	setup       Setup
	action      Action
	resolveCard ResolveCard

	commands []engine.Command
	barriers []bool // per applied command: did applying it cross an information barrier

	game    *engine.Game
	stepper *engine.Stepper
	request engine.Request
	done    bool
}

// New deals a fresh match and drives action up to its first decision. setup,
// action, and resolveCard are retained so an undo can rebuild the match by
// replaying.
func New(seed int64, sets [2]string, setup Setup, action Action, resolveCard ResolveCard) *Session {
	s := &Session{
		seed:        seed,
		sets:        sets,
		setup:       setup,
		action:      action,
		resolveCard: resolveCard,
	}
	s.start()
	return s
}

// Load rebuilds a session from a persisted record and replays its command log. It
// refuses a record whose version does not match this build (no forward compat).
func Load(rec Record, setup Setup, action Action, resolveCard ResolveCard) (*Session, error) {
	if rec.Version != Version {
		return nil, ErrVersion
	}
	s := New(rec.Seed, rec.Sets, setup, action, resolveCard)
	for _, cmd := range rec.Commands {
		if err := s.replay(cmd); err != nil {
			return nil, err
		}
	}
	return s, nil
}

// start deals a fresh game and drives the action to its first request. It is used
// both on New and on every replay (undo rebuilds from scratch); it closes the
// previous stepper first so a replaced mid-action game leaves no parked goroutine.
func (s *Session) start() {
	if s.stepper != nil {
		s.stepper.Close()
	}
	s.commands = nil
	s.barriers = nil
	s.game = s.setup(s.seed, s.sets)
	s.stepper = engine.NewStepper(s.game, s.action)
	s.request, s.done = s.stepper.Start()
}

// Pending returns the request currently awaiting a command and whether the match
// has finished. When done is true the request is the zero value.
func (s *Session) Pending() (engine.Request, bool) { return s.request, s.done }

// Err reports why the driving action stopped: nil when it ran to completion or is
// still mid-action, and the engine's *PanicError when it broke. Apply already
// surfaces a panic it triggered; this answers the same question for a panic that
// happened before the first request, where there was no Apply to return it.
func (s *Session) Err() error { return s.stepper.Err() }

// Apply answers the pending request with cmd, records it in the log, and advances
// to the next request. It rejects a command that does not answer the pending
// request, so the log only ever holds legal inputs.
func (s *Session) Apply(cmd engine.Command) error {
	if s.done {
		return ErrFinished
	}
	if !s.legal(cmd) {
		return ErrIllegal
	}
	req, done, info := s.stepper.Advance(cmd)
	s.commands = append(s.commands, cmd)
	s.barriers = append(s.barriers, info.CrossedBarrier)
	s.request, s.done = req, done
	if done {
		if err := s.stepper.Err(); err != nil {
			// The command is still recorded, so the caller rolls the broken step back
			// with Undo(Len()-1) and the log stays the reproduction of what was entered.
			return fmt.Errorf("session: the action failed: %w", err)
		}
	}
	return nil
}

// legal reports whether cmd answers the pending request. A RequestAction is checked
// against the LIVE legal set rather than the Actions the request carried, because a
// manual force-edit applied out of band (ApplyManual) since the request was yielded
// makes that cached slice stale — adding a card to hand makes a play legal that was
// not legal when the engine gathered the set. Checking live is safe because nothing
// downstream trusts the cached slice either: suspendChooser.ChooseAction returns
// whatever command it is fed without validating, and ApplyAction re-checks the play
// against live state. Every other request kind is self-contained, so its own
// candidate context is still the authority.
func (s *Session) legal(cmd engine.Command) bool {
	if s.request.Kind == engine.RequestAction {
		return slices.Contains(s.game.LegalActions(s.request.Player), cmd)
	}
	return s.request.IsLegal(cmd)
}

// ApplyManual records a manual force-edit and applies it DIRECTLY to the live game,
// bypassing the stepper — a manual edit is not an answer to the pending request, so
// it must not consume one. Recording it is the point: a playtester who force-edits
// the board and then hits a bug has produced exactly the reproduction worth keeping,
// and a replay rebuilds the edited board by routing the command back here.
//
// Touching s.game directly is safe because the session only ever holds the game
// while the action goroutine is parked: with a request pending, the stepper is
// blocked in suspendChooser.yield waiting for a Command, and once the action is done
// the goroutine is gone. Either way nothing else is reading or writing the state.
//
// The edit crosses no information barrier: it steps no PRNG and reveals no hidden
// zone, so undoing past it is always free.
func (s *Session) ApplyManual(cmd engine.Command) error {
	if err := s.game.ApplyManual(cmd, s.resolveCard); err != nil {
		return err
	}
	s.commands = append(s.commands, cmd)
	s.barriers = append(s.barriers, false)
	return nil
}

// replay reapplies one recorded command by KIND — a manual force-edit straight to
// the live game, everything else as an answer the stepper advances on. It is the
// single routing rule every rebuild goes through (Undo and Load), so the two cannot
// drift into disagreeing about what a recorded command means.
func (s *Session) replay(cmd engine.Command) error {
	if cmd.Kind.IsManual() {
		return s.ApplyManual(cmd)
	}
	return s.Apply(cmd)
}

// Undo rewinds the match to just after its first n commands (n in [0, applied]),
// rebuilding the game by replaying that prefix from a fresh deal. Because state is
// a projection of the log, undo is replay, not an inverse operation.
func (s *Session) Undo(n int) error {
	if n < 0 || n > len(s.commands) {
		return ErrRange
	}
	keep := append([]engine.Command(nil), s.commands[:n]...)
	s.start()
	for _, cmd := range keep {
		if err := s.replay(cmd); err != nil {
			return err
		}
	}
	return nil
}

// Len is the number of commands applied so far — the undo cursor's upper bound.
func (s *Session) Len() int { return len(s.commands) }

// CrossesBarrier reports whether undoing to n would rewind past an information
// barrier: a command whose resolution stepped the PRNG or revealed a hidden zone.
// In hotseat any undo is free; in networked play crossing a barrier needs the
// opponent's consent (ADR 0039). n out of range reports false.
func (s *Session) CrossesBarrier(n int) bool {
	if n < 0 || n > len(s.barriers) {
		return false
	}
	for i := n; i < len(s.barriers); i++ {
		if s.barriers[i] {
			return true
		}
	}
	return false
}

// Record returns the persisted form of the match — the seed, sets, and command
// log — tagged with this build's version. It is the only thing that needs saving.
func (s *Session) Record() Record {
	return Record{
		Version:  Version,
		Seed:     s.seed,
		Sets:     s.sets,
		Commands: append([]engine.Command(nil), s.commands...),
	}
}

// View returns what viewer is allowed to see of the current state, through the
// engine's projection seam (identity today; redactable later).
func (s *Session) View(viewer int) engine.View {
	return engine.Project(s.game.State, viewer)
}

// Game exposes the live game the session is driving, for callers that still read
// the engine directly (the web client, until it renders purely from View).
func (s *Session) Game() *engine.Game { return s.game }
