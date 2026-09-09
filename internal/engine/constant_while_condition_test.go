package engine

import (
	"strings"
	"testing"
)

// KeyColorForged is met only while the named player has forged a key of the given
// colour, covering The Red Baron.
func TestKeyColorForgedCondition(t *testing.T) {
	g := NewGame("A", "B", 1)
	ctx := &EffectContext{Resolver: g, Controller: 0}

	own := KeyColorForged{Player: Controller, Color: KeyColorRed}
	opp := KeyColorForged{Player: Opponent, Color: KeyColorRed}

	if own.CondText() != "if your red key is forged" {
		t.Errorf("own CondText = %q", own.CondText())
	}
	if opp.CondText() != "if your opponent's red key is forged" {
		t.Errorf("opp CondText = %q", opp.CondText())
	}

	if own.Met(ctx) || opp.Met(ctx) {
		t.Error("no key forged: neither condition should be met")
	}

	g.State.Keys[0] = 1
	g.State.KeyColors[0][0] = KeyColorRed
	if !own.Met(ctx) {
		t.Error("your red key forged: own condition should be met")
	}
	if opp.Met(ctx) {
		t.Error("your red key forged: opponent condition should not be met")
	}

	g.State.Keys[1] = 1
	g.State.KeyColors[1][0] = KeyColorBlue
	if opp.Met(ctx) {
		t.Error("opponent's blue key forged: red condition should not be met")
	}
	g.State.KeyColors[1][0] = KeyColorRed
	if !opp.Met(ctx) {
		t.Error("opponent's red key forged: opponent condition should be met")
	}
}

// A WhileCondition constant self-grant applies its granted ability only while the
// condition holds, and renders the "While ..." line, covering The Red Baron's reap.
func TestConstantAbilityWhileConditionGranted(t *testing.T) {
	g := NewGame("A", "B", 1)
	def := NewCard("Baron", Brobnar, Creature, Special, WithPower(4),
		WithConstantAbility(ConstantAbility{
			Target:         Target{Kind: TargetThisCreature},
			WhileCondition: KeyColorForged{Player: Controller, Color: KeyColorRed},
			Granted: []Ability{{
				Trigger: TriggerAfterReap,
				Effect:  StealAember{Amount: 1},
			}},
		}))

	if !strings.Contains(RenderCardRules(&def),
		`While your red key is forged, Baron gains, "Reap: Steal 1 Æmber."`) {
		t.Errorf("card rules should render the while-condition granted line, got %q",
			RenderCardRules(&def))
	}

	baron := g.AddToBattleline(def, 0)

	if g.HasTrigger(baron, TriggerAfterReap) {
		t.Error("no key forged: Baron should not have the granted reap")
	}

	g.State.Keys[0] = 1
	g.State.KeyColors[0][0] = KeyColorRed
	if !g.HasTrigger(baron, TriggerAfterReap) {
		t.Error("your red key forged: Baron should have the granted reap")
	}
}

// A WhileCondition constant keyword grant applies only while the condition holds,
// and renders the "While ..." line, covering The Red Baron's elusive.
func TestConstantAbilityWhileConditionKeyword(t *testing.T) {
	g := NewGame("A", "B", 1)
	def := NewCard("Baron", Brobnar, Creature, Special, WithPower(4),
		WithConstantAbility(ConstantAbility{
			Target:         Target{Kind: TargetThisCreature},
			WhileCondition: KeyColorForged{Player: Opponent, Color: KeyColorRed},
			Keywords:       []Keyword{Elusive},
		}))

	if !strings.Contains(RenderCardRules(&def),
		"While your opponent's red key is forged, Baron gains elusive.") {
		t.Errorf("card rules should render the while-condition keyword line, got %q",
			RenderCardRules(&def))
	}

	baron := g.AddToBattleline(def, 0)

	if g.hasKeyword(baron, Elusive) {
		t.Error("no key forged: Baron should not have elusive")
	}

	g.State.Keys[1] = 1
	g.State.KeyColors[1][0] = KeyColorRed
	if !g.hasKeyword(baron, Elusive) {
		t.Error("opponent's red key forged: Baron should have elusive")
	}
}
