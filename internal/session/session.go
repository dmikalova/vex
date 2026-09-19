// Package session drives one Vactrol match as an event-sourced replay (ADR 0039,
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

	"github.com/dmikalova/vactrol/internal/engine"
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
	seed   int64
	sets   [2]string
	setup  Setup
	action Action

	commands []engine.Command
	barriers []bool // per applied command: did applying it cross an information barrier

	game    *engine.Game
	stepper *engine.Stepper
	request engine.Request
	done    bool
}

// New deals a fresh match and drives action up to its first decision. setup and
// action are retained so an undo can rebuild the match by replaying.
func New(seed int64, sets [2]string, setup Setup, action Action) *Session {
	s := &Session{seed: seed, sets: sets, setup: setup, action: action}
	s.start()
	return s
}

// Load rebuilds a session from a persisted record and replays its command log. It
// refuses a record whose version does not match this build (no forward compat).
func Load(rec Record, setup Setup, action Action) (*Session, error) {
	if rec.Version != Version {
		return nil, ErrVersion
	}
	s := New(rec.Seed, rec.Sets, setup, action)
	for _, cmd := range rec.Commands {
		if err := s.Apply(cmd); err != nil {
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

// Apply answers the pending request with cmd, records it in the log, and advances
// to the next request. It rejects a command that does not answer the pending
// request, so the log only ever holds legal inputs.
func (s *Session) Apply(cmd engine.Command) error {
	if s.done {
		return ErrFinished
	}
	if !s.request.IsLegal(cmd) {
		return ErrIllegal
	}
	req, done, info := s.stepper.Advance(cmd)
	s.commands = append(s.commands, cmd)
	s.barriers = append(s.barriers, info.CrossedBarrier)
	s.request, s.done = req, done
	return nil
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
		if err := s.Apply(cmd); err != nil {
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
