package cardtest

import (
	"testing"
	"time"

	"github.com/dmikalova/vactrol/internal/engine"
)

// These tests reach past the clicking vocabulary to the harness's own internals:
// the lookups that only ever render a message, and the guards that exist so a
// wedged engine reports rather than hangs. Neither is reachable by playing a
// card, which is the whole reason they are asserted here.

// shortTimeout shrinks the wedge backstop for the duration of a test, so a test
// that means to trip it takes milliseconds rather than seconds.
func shortTimeout(t *testing.T) {
	t.Helper()
	prev := actionTimeout
	actionTimeout = 10 * time.Millisecond
	t.Cleanup(func() { actionTimeout = prev })
}

// A bridge asked to choose from nothing answers for itself. The engine already
// short-circuits an empty candidate set, so these guards are the harness's own
// belt and braces against an engine that stops doing so.
func TestBridgeAnswersDegenerateChoices(t *testing.T) {
	var lone Card
	h := Play(t, Setup{P1: Side{InPlay: []Entry{Bind(&lone, Creature())}}})
	b := bridgeChooser{h: h, player: 0}
	if _, ok := b.ChooseCreature("src", "Choose a creature", nil); ok {
		t.Error("chose a creature from an empty set")
	}
	if got, ok := b.ChooseCreature("src", "Choose a creature", []engine.LocalID{lone.ID()}); !ok ||
		got != lone.ID() {
		t.Errorf("sole creature = %v/%v, want %v/true", got, ok, lone.ID())
	}
	if _, ok := b.ChooseCardOrDecline("src", "Choose a card", nil); ok {
		t.Error("chose a card from an empty set")
	}
	if got := b.ChooseOption("src", "Choose one", []string{"the only one"}); got != 0 {
		t.Errorf("sole option = %d, want 0", got)
	}
}

// A sole creature needs no prompt: the harness answers it so an ordinary
// single-target card test does not have to script a click.
func TestSoleCreatureNeedsNoClick(t *testing.T) {
	pick := targeting()
	var lone Card
	h := Play(t, Setup{
		P1: Side{Hand: Cards(pick)},
		P2: Side{InPlay: []Entry{Bind(&lone, Creature(Power(3)))}},
	})
	h.P1.Play(pick)
	h.Expect(lone).Damage(1)
}

// With no enemy creature at all the effect finds nothing to damage and the card
// resolves to no effect rather than stalling on an empty prompt.
func TestNoCandidateNeedsNoClick(t *testing.T) {
	pick := targeting()
	h := Play(t, Setup{P1: Side{Hand: Cards(pick)}})
	h.P1.Play(pick)
	h.P1.ExpectAmber(0)
}

// A card may be clicked by definition when exactly one candidate carries that
// name, which is the everyday form a card test writes.
func TestClickByDefinition(t *testing.T) {
	pick := targeting()
	enemy := Creature(Power(3))
	h := Play(t, Setup{
		P1: Side{Hand: Cards(pick)},
		P2: Side{InPlay: Cards(enemy, Creature())},
	})
	h.P1.Play(pick)
	h.P1.ClickCard(enemy)
	h.Expect(enemy).Damage(1)
}

// A definition that names no candidate is reported as an illegal target, the
// same as a handle that names the wrong card.
func TestClickDefinitionThatIsNotACandidate(t *testing.T) {
	pick := targeting()
	got := expectFail(t, func(tb testing.TB) {
		h := Play(tb, twoEnemies(pick))
		h.P1.Play(pick)
		h.P1.ClickCard(Creature())
	})
	wants(t, got, "not a legal target")
}

// An upgrade sits in neither player's battleline but on its host, which is the
// zone a test asserts with At(ct.Attached).
func TestUpgradeIsAttached(t *testing.T) {
	var host, boon Card
	h := Play(t, Setup{P1: Side{InPlay: []Entry{
		Upgraded(Bind(&host, Creature(Power(3))), Bind(&boon, Upgrade(PowerBonus(2)))),
	}}})
	h.Expect(boon).At(Attached)
	h.Expect(host).Power(5)
}

// An id in no zone at all is Gone. Nothing a card test does produces one, so it
// is asserted directly rather than through a card.
func TestLocationOfNothing(t *testing.T) {
	h := Play(t, Setup{})
	if got := h.location(engine.LocalID(200)); got != Gone {
		t.Errorf("location of an unknown id = %v, want %v", got, Gone)
	}
	if got := h.ownerOf(engine.LocalID(200)); got != 0 {
		t.Errorf("owner of an unknown id = %d, want 0", got)
	}
}

// describe and matchesCard render and match the same def-or-handle vocabulary the
// clicking API takes, and both have to cope with a value that is neither.
func TestDescribeAndMatch(t *testing.T) {
	var troll Card
	def := Creature()
	h := Play(t, Setup{P1: Side{InPlay: []Entry{Bind(&troll, def)}}})
	for _, tc := range []struct {
		target any
		want   string
	}{
		{troll, troll.Name()},
		{Card{h: h}, "unbound card"},
		{def, def.Name},
		{42, "?"},
	} {
		if got := h.describe(tc.target); got != tc.want {
			t.Errorf("describe(%v) = %q, want %q", tc.target, got, tc.want)
		}
	}
	if !h.matchesCard(def, troll.ID()) {
		t.Error("a definition does not match its own card")
	}
	if h.matchesCard(42, troll.ID()) {
		t.Error("a non-card matched a card")
	}
}

// The wedge backstop: an engine that neither finishes nor asks anything leaves
// the harness waiting, and it says so instead of hanging the binary.
func TestWedgedEngineIsReported(t *testing.T) {
	shortTimeout(t)
	got := expectFail(t, func(tb testing.TB) {
		h := Play(tb, Setup{})
		h.label = "Wedge"
		h.advance()
	})
	wants(t, got, "Wedge: the engine did not finish or ask for a choice")
}

// The ready check drains whatever the declined effect asks next, and gives up if
// the action goroutine stops answering rather than blocking the run forever.
func TestReadyCheckDrainsThenGivesUp(t *testing.T) {
	shortTimeout(t)
	got := expectFail(t, func(tb testing.TB) {
		h := Play(tb, Setup{})
		// Two prompts nobody is listening to: declining the first hands back a
		// second, and declining that one leaves the harness waiting on an action
		// goroutine that is not there.
		first := promptReq{text: "first", reply: make(chan int, 1)}
		second := promptReq{text: "second", reply: make(chan int, 1)}
		h.current = &first
		go func() { h.prompt <- second }()
		h.checkReady()
	})
	wants(t, got, "unanswered prompt \"first\"")
}

// A scenario naming a house that does not exist would be one where no card is
// playable, so the setup rejects it instead of quietly playing it out.
func TestSetupRejectsAnImpossibleHouse(t *testing.T) {
	got := expectFail(t, func(tb testing.TB) {
		Play(tb, Setup{P1: Side{House: engine.House(engine.NumHouses)}})
	})
	wants(t, got, "is not one of the")
}
