package engine

import (
	"reflect"
	"testing"
	"time"
)

// LegalCommands enumerates exactly the answers each request kind accepts, and
// IsLegal agrees with that set — the check a holder runs before applying an
// answer it did not itself produce.
func TestRequestLegalCommands(t *testing.T) {
	cases := []struct {
		name string
		req  Request
		want []Command
	}{
		{
			"pick card",
			Request{Kind: RequestPickCard, Cards: []LocalID{1, 2}},
			[]Command{{Kind: CommandPickCard, Card: 1}, {Kind: CommandPickCard, Card: 2}},
		},
		{
			"pick or decline",
			Request{Kind: RequestPickCardOrDecline, Cards: []LocalID{5}},
			[]Command{{Kind: CommandPickCard, Card: 5}, {Kind: CommandDecline}},
		},
		{
			"option",
			Request{Kind: RequestOption, Options: []string{"a", "b", "c"}},
			[]Command{
				{Kind: CommandOption, Index: 0},
				{Kind: CommandOption, Index: 1},
				{Kind: CommandOption, Index: 2},
			},
		},
		{
			"position",
			Request{Kind: RequestPosition, Cards: []LocalID{3}},
			[]Command{{Kind: CommandPosition, Index: 0}, {Kind: CommandPosition, Index: 1}},
		},
		{
			"reaction",
			Request{
				Kind:      RequestReaction,
				Reactions: []OrderableReaction{{Label: "x"}, {Label: "y"}},
			},
			[]Command{{Kind: CommandReaction, Index: 0}, {Kind: CommandReaction, Index: 1}},
		},
		{
			"action",
			Request{Kind: RequestAction, Actions: []Command{
				{Kind: CommandReap, Card: 3},
				{Kind: CommandEndTurn},
			}},
			[]Command{{Kind: CommandReap, Card: 3}, {Kind: CommandEndTurn}},
		},
		{
			"unknown kind",
			Request{Kind: RequestKind(99)},
			nil,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.req.LegalCommands()
			if !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("LegalCommands() = %v, want %v", got, tc.want)
			}
			for _, cmd := range tc.want {
				if !tc.req.IsLegal(cmd) {
					t.Errorf("IsLegal(%v) = false, want true", cmd)
				}
			}
			// A command of an index/id the request never offered is rejected.
			if tc.req.IsLegal(Command{Kind: CommandReaction, Index: 42}) {
				t.Errorf("IsLegal accepted an unrelated command")
			}
		})
	}
}

// A Stepper suspends the action at each decision and resumes it with the answer,
// driving one action through every chooser capability in turn. The action calls
// the choosers the engine would call, and the driver answers each yielded Request.
func TestStepperDrivesEachCapability(t *testing.T) {
	g := NewGame("A", "B", 1)
	var (
		gotCreature LocalID
		gotOption   int
		gotPosition int
		gotDeclined bool
		gotReaction int
	)
	action := func(g *Game) {
		ch := g.chooserFor(0)
		gotCreature, _ = ch.ChooseCreature("src", "pick", []LocalID{1, 2})
		gotOption = ch.(OptionChooser).ChooseOption("src", "opt", []string{"a", "b"})
		gotPosition = ch.(PositionChooser).ChoosePosition("src", "pos", []LocalID{3})
		_, ok := ch.(DeclinableChooser).ChooseCardOrDecline("src", "may", []LocalID{4, 5})
		gotDeclined = !ok
		gotReaction = ch.(ReactionChooser).ChooseReaction(
			"react", []OrderableReaction{{Label: "x"}, {Label: "y"}})
	}
	s := NewStepper(g, action)

	req, done := s.Start()
	if done || req.Kind != RequestPickCard || req.Player != 0 {
		t.Fatalf("first request = %+v done=%v, want a player-0 card pick", req, done)
	}

	req, done, _ = s.Advance(Command{Kind: CommandPickCard, Card: 2})
	if done || req.Kind != RequestOption {
		t.Fatalf("after card pick got %+v done=%v, want an option request", req, done)
	}

	req, done, _ = s.Advance(Command{Kind: CommandOption, Index: 1})
	if done || req.Kind != RequestPosition {
		t.Fatalf("after option got %+v done=%v, want a position request", req, done)
	}

	req, done, _ = s.Advance(Command{Kind: CommandPosition, Index: 1})
	if done || req.Kind != RequestPickCardOrDecline {
		t.Fatalf("after position got %+v done=%v, want a declinable request", req, done)
	}

	req, done, _ = s.Advance(Command{Kind: CommandDecline})
	if done || req.Kind != RequestReaction {
		t.Fatalf("after decline got %+v done=%v, want a reaction request", req, done)
	}

	_, done, _ = s.Advance(Command{Kind: CommandReaction, Index: 0})
	if !done {
		t.Fatal("the action did not finish after the last answer")
	}

	if gotCreature != 2 || gotOption != 1 || gotPosition != 1 || !gotDeclined || gotReaction != 0 {
		t.Fatalf("answers not delivered: creature=%d option=%d position=%d declined=%v reaction=%d",
			gotCreature, gotOption, gotPosition, gotDeclined, gotReaction)
	}

	// Advancing past completion is a no-op that reports done again.
	if _, done, _ := s.Advance(Command{Kind: CommandOption}); !done {
		t.Fatal("Advance after done did not report done")
	}
}

// A Stepper surfaces the turn loop's suspension point: ChooseAction yields a
// RequestAction carrying the whole legal set, and the answered root Command comes
// back to the action to perform.
func TestStepperYieldsActionRequest(t *testing.T) {
	g := NewGame("A", "B", 1)
	actions := []Command{{Kind: CommandReap, Card: 7}, {Kind: CommandEndTurn}}
	var got Command
	action := func(g *Game) {
		got = g.chooserFor(0).(ActionChooser).ChooseAction(actions)
	}
	s := NewStepper(g, action)

	req, done := s.Start()
	if done || req.Kind != RequestAction || req.Player != 0 {
		t.Fatalf("first request = %+v done=%v, want a player-0 action request", req, done)
	}
	if !reflect.DeepEqual(req.Actions, actions) {
		t.Fatalf("request actions = %v, want %v", req.Actions, actions)
	}

	_, done, _ = s.Advance(Command{Kind: CommandEndTurn})
	if !done {
		t.Fatal("the action did not finish after the answer")
	}
	if got != (Command{Kind: CommandEndTurn}) {
		t.Fatalf("ChooseAction returned %+v, want CommandEndTurn", got)
	}
}

// Close releases the action goroutine when a Stepper is abandoned mid-decision —
// an undo or replay deals a fresh game — instead of leaving it parked forever on a
// command that never arrives. The action's deferred cleanup runs either way,
// proving the goroutine unwound rather than leaked.
func TestStepperCloseReleasesGoroutine(t *testing.T) {
	// Parked yielding its first request (nobody ever called Start, so no receiver):
	// Close fires the send-side cancel.
	t.Run("parked yielding a request", func(t *testing.T) {
		g := NewGame("A", "B", 1)
		released := make(chan struct{})
		action := func(g *Game) {
			defer close(released)
			g.chooserFor(0).ChooseCreature("s", "p", []LocalID{1, 2})
		}
		s := NewStepper(g, action)
		s.Close()
		requireReleased(t, released)
	})

	// Parked awaiting the command after Start: Close fires the receive-side cancel,
	// is idempotent, and a later Advance reports done without blocking.
	t.Run("parked awaiting a command", func(t *testing.T) {
		g := NewGame("A", "B", 1)
		released := make(chan struct{})
		action := func(g *Game) {
			defer close(released)
			g.chooserFor(0).ChooseCreature("s", "p", []LocalID{1, 2})
		}
		s := NewStepper(g, action)
		if _, done := s.Start(); done {
			t.Fatal("action finished before its choice; expected a parked request")
		}
		s.Close()
		requireReleased(t, released)
		s.Close() // idempotent: already done, no double close of cancel
		if _, done, _ := s.Advance(Command{Kind: CommandPickCard, Card: 1}); !done {
			t.Fatal("Advance after Close should report done")
		}
	})
}

// requireReleased fails if the action goroutine's deferred cleanup has not run
// shortly after Close, which would mean it leaked rather than unwound.
func requireReleased(t *testing.T, released <-chan struct{}) {
	t.Helper()
	select {
	case <-released:
	case <-time.After(time.Second):
		t.Fatal("Close did not release the parked action goroutine")
	}
}

// An action that declines a plain card pick returns ok=false to the effect, and an
// action that needs no decision reports done from Start.
func TestStepperDeclineAndNoDecision(t *testing.T) {
	g := NewGame("A", "B", 1)
	var pickedOK bool
	declineAction := func(g *Game) {
		_, pickedOK = g.chooserFor(0).ChooseCreature("s", "p", []LocalID{1, 2})
	}
	s := NewStepper(g, declineAction)
	if _, done := s.Start(); done {
		t.Fatal("expected a card pick request, got done")
	}
	if _, done, _ := s.Advance(Command{Kind: CommandDecline}); !done {
		t.Fatal("the action did not finish after declining")
	}
	if pickedOK {
		t.Error("a declined ChooseCreature returned ok=true")
	}

	s = NewStepper(g, func(*Game) {})
	if _, done := s.Start(); !done {
		t.Fatal("an action with no decision should report done from Start")
	}

	// A declinable pick that is taken (not declined) returns the chosen card.
	var (
		pickedID LocalID
		tookIt   bool
	)
	takeAction := func(g *Game) {
		pickedID, tookIt = g.chooserFor(0).(DeclinableChooser).
			ChooseCardOrDecline("s", "may", []LocalID{7, 8})
	}
	s = NewStepper(g, takeAction)
	s.Start()
	if _, done, _ := s.Advance(Command{Kind: CommandPickCard, Card: 8}); !done {
		t.Fatal("the action did not finish after taking the declinable pick")
	}
	if !tookIt || pickedID != 8 {
		t.Fatalf("declinable pick delivered id=%d ok=%v, want id=8 ok=true", pickedID, tookIt)
	}
}

// StepInfo flags a step that stepped the PRNG as crossing an information barrier,
// and a step that did not as not crossing one.
func TestStepperReportsBarrier(t *testing.T) {
	g := NewGame("A", "B", 1)
	action := func(g *Game) {
		ch := g.chooserFor(0)
		ch.ChooseCreature("s", "p", []LocalID{1, 2}) // req 1
		g.State.PRNG.Intn(
			10,
		) // steps the PRNG resolving cmd 1
		ch.(OptionChooser).ChooseOption("s", "o", []string{"a", "b"}) // req 2
		ch.(OptionChooser).ChooseOption("s", "o", []string{"a", "b"}) // req 3
	}
	s := NewStepper(g, action)
	s.Start()
	_, _, info := s.Advance(Command{Kind: CommandPickCard, Card: 1})
	if !info.CrossedBarrier {
		t.Error("a step that drew from the PRNG did not report a barrier crossing")
	}
	_, _, info = s.Advance(Command{Kind: CommandOption})
	if info.CrossedBarrier {
		t.Error("a step that did not touch the PRNG reported a barrier crossing")
	}
	s.Advance(Command{Kind: CommandOption})
}
