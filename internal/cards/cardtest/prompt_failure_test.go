package cardtest

import (
	"strconv"
	"testing"

	"github.com/dmikalova/vactrol/internal/engine"
)

// These tests cover the harness's prompt bridge: the clicking vocabulary a card
// test drives a choice with, and every way a mis-scripted click is reported.

// targeting builds a tactic whose Play ability damages an enemy creature the
// player picks — the smallest card that stops on a card prompt.
func targeting() engine.CardDefinition {
	return engine.NewCard(
		"Targeting "+strconv.FormatUint(vanillaCount.Add(1), 10),
		DefaultHouse, engine.Tactic, engine.Common,
		engine.WithAbility(engine.TriggerAfterPlay, engine.DealDamage{
			Amount: 1,
			Target: engine.Target{Kind: engine.TargetChosenEnemyCreature},
		}),
	)
}

// choosing builds a tactic whose Play ability offers two Æmber amounts — the
// smallest card that stops on an option prompt.
func choosing() engine.CardDefinition {
	return engine.NewCard(
		"Choosing "+strconv.FormatUint(vanillaCount.Add(1), 10),
		DefaultHouse, engine.Tactic, engine.Common,
		engine.WithAbility(engine.TriggerAfterPlay, engine.ChooseOne{
			Options: []engine.Effect{
				engine.GainAember{Player: engine.Controller, Amount: 1},
				engine.GainAember{Player: engine.Controller, Amount: 2},
			},
		}),
	)
}

// twoEnemies is a scenario where a card prompt has more than one candidate, so
// the bridge cannot answer it on the test's behalf.
func twoEnemies(hand ...any) Setup {
	return Setup{
		P1: Side{Hand: Cards(hand...)},
		P2: Side{InPlay: Cards(Creature(), Creature())},
	}
}

func TestClickMisuse(t *testing.T) {
	pick := targeting()
	var mine Card
	for _, tc := range []struct {
		name  string
		hand  engine.CardDefinition
		click func(h *Harness)
		want  []string
	}{
		{
			name:  "an option answer to a card prompt",
			hand:  pick,
			click: func(h *Harness) { h.P1.ClickOption("Gain") },
			want:  []string{"expects a card", "use ClickCard"},
		},
		{
			name:  "a card answer to an option prompt",
			hand:  choosing(),
			click: func(h *Harness) { h.P1.ClickCard(Creature()) },
			want:  []string{"expects an option", "use ClickOption"},
		},
		{
			name:  "declining a prompt that has no Done",
			hand:  pick,
			click: func(h *Harness) { h.P1.ClickDone() },
			want:  []string{"cannot be declined"},
		},
		{
			name:  "a card that is not a legal target",
			hand:  pick,
			click: func(h *Harness) { h.P1.ClickCard(mine) },
			want:  []string{"not a legal target", "clickable: ["},
		},
		{
			name:  "something that is not a card at all",
			hand:  pick,
			click: func(h *Harness) { h.P1.ClickCard(42) },
			want:  []string{"cannot click a int"},
		},
		{
			name:  "an option label that matches nothing",
			hand:  choosing(),
			click: func(h *Harness) { h.P1.ClickOption("forge a key") },
			want:  []string{"not an option for"},
		},
		{
			name:  "an option label that matches several",
			hand:  choosing(),
			click: func(h *Harness) { h.P1.ClickOption("Æmber") },
			want:  []string{"matches 2 options"},
		},
		{
			name:  "a prompt answered by the wrong player",
			hand:  pick,
			click: func(h *Harness) { h.P2.ClickCard(Creature()) },
			want:  []string{"is for P1, not P2"},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got := expectFail(t, func(tb testing.TB) {
				h := Play(tb, Setup{
					P1: Side{
						InPlay: []Entry{Bind(&mine, Creature())},
						Hand:   Cards(tc.hand),
					},
					P2: Side{InPlay: Cards(Creature(), Creature())},
				})
				h.P1.Play(tc.hand)
				tc.click(h)
			})
			wants(t, got, tc.want...)
		})
	}
}

// Two enemy creatures with the same name are indistinguishable by definition, so
// the harness asks for the handle rather than picking one.
func TestClickAmbiguousDefinition(t *testing.T) {
	twin := Creature()
	pick := targeting()
	got := expectFail(t, func(tb testing.TB) {
		h := Play(tb, Setup{
			P1: Side{Hand: Cards(pick)},
			P2: Side{InPlay: Cards(twin, twin)},
		})
		h.P1.Play(pick)
		h.P1.ClickCard(twin)
	})
	wants(t, got, "2 matching targets", "ct.Bind")
}

// A prompt clicked when nothing asked is a mis-scripted test, not a hang.
func TestClickWithNoPromptPending(t *testing.T) {
	got := expectFail(t, func(tb testing.TB) {
		Play(tb, Setup{}).P1.ClickCard(Creature())
	})
	wants(t, got, "no prompt is pending for P1")
}

// Leaving a prompt unanswered is the mistake the ready check exists for: it names
// the prompt, then drains the parked action so the run does not leak a goroutine.
func TestUnansweredPromptAtEnd(t *testing.T) {
	for _, tc := range []struct {
		name string
		card engine.CardDefinition
		want string
	}{
		{"card prompt", targeting(), "clickable: ["},
		{"option prompt", choosing(), "options: ["},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got := expectFail(t, func(tb testing.TB) {
				h := Play(tb, twoEnemies(tc.card))
				h.P1.Play(tc.card)
			})
			wants(t, got, "ended with an unanswered prompt", tc.want)
		})
	}
}

func TestExpectPromptMismatch(t *testing.T) {
	pick := targeting()
	got := expectFail(t, func(tb testing.TB) {
		h := Play(tb, twoEnemies(pick))
		h.P1.Play(pick)
		p := h.P1.ExpectPrompt("Choose a house")
		p.Source(pick.Name + " (not really)")
	})
	wants(t, got, "want \"Choose a house\"", "source = ")
}

// An action the engine refuses fails the test naming the action, rather than the
// test reading a silently unchanged board.
func TestIllegalActionFails(t *testing.T) {
	var troll Card
	got := expectFail(t, func(tb testing.TB) {
		h := Play(tb, Setup{P1: Side{InPlay: []Entry{Bind(&troll, Creature())}}})
		troll.Exhaust()
		h.P1.Reap(troll)
	})
	wants(t, got, "Reap ", "card is exhausted")
}
