package engine

import (
	"strings"
	"testing"
)

// A creature with DealsNoDamageWhenAttacked deals no retaliation damage to an
// attacker that fights it, but still takes the attacker's fight damage — Lollop
// the Titanic.
func TestDealsNoDamageWhenAttacked(t *testing.T) {
	lollop := NewCard("Lollop the Titanic", Brobnar, Creature, Common,
		WithPower(11), WithNoDamageWhenAttacked())
	if got := RenderCardRules(&lollop); !strings.Contains(got,
		"Lollop the Titanic deals no damage when attacked.") {
		t.Fatalf("Lollop rules = %q", got)
	}

	g := NewGame("A", "B", 1)
	attacker := g.AddToBattleline(testCreature("attacker", 5), 0)
	defender := g.AddToBattleline(lollop, 1)
	g.State.ActivePlayer = 0

	if err := g.Fight(0, attacker, defender); err != nil {
		t.Fatalf("Fight: %v", err)
	}
	if !g.inPlay(attacker) {
		t.Error("attacker should survive: Lollop deals no retaliation damage")
	}
	if got := g.Damage(attacker); got != 0 {
		t.Errorf("attacker damage = %d, want 0", got)
	}
	if got := g.Damage(defender); got != 5 {
		t.Errorf("Lollop damage = %d, want 5 (it still takes fight damage)", got)
	}
}
