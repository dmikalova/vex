package engine

import (
	"strings"
	"testing"
)

// TestDamageSourceMatcher exercises the matcher's rendering and its live
// power/trait predicate, since the card test that consumes the passive lives in
// another package and does not cover the engine here.
func TestDamageSourceMatcher(t *testing.T) {
	// Active reflects whether any criterion is set.
	if (DamageSourceMatcher{}).Active() {
		t.Error("the zero matcher is inactive")
	}
	if !(DamageSourceMatcher{Trait: Mutant}).Active() {
		t.Error("a trait matcher is active")
	}
	if !(DamageSourceMatcher{MinPower: 5}).Active() {
		t.Error("a power matcher is active")
	}

	// clause renders each of the three shapes.
	cases := []struct {
		m    DamageSourceMatcher
		want string
	}{
		{DamageSourceMatcher{Trait: Mutant}, "Mutant creatures"},
		{DamageSourceMatcher{MinPower: 5}, "creatures with power 5 or higher"},
		{
			DamageSourceMatcher{Trait: Mutant, MinPower: 5},
			"Mutant creatures or creatures with power 5 or higher",
		},
	}
	for _, c := range cases {
		if got := c.m.clause(); got != c.want {
			t.Errorf("clause = %q, want %q", got, c.want)
		}
	}

	// The passive renders into a card's rules text.
	def := NewCard("Ardent Hero", Sanctum, Creature, Common,
		WithPower(4),
		WithCannotBeDealtDamageBy(DamageSourceMatcher{Trait: Mutant, MinPower: 5}))
	if got := RenderCardRules(&def); !strings.Contains(
		got,
		"Ardent Hero cannot be dealt damage by Mutant creatures or creatures with power 5 or higher.",
	) {
		t.Errorf("rules = %q", got)
	}
}

// TestRefusesDamageFrom checks the live predicate against trait, power, and the
// uncredited-source and inactive cases.
func TestRefusesDamageFrom(t *testing.T) {
	m := DamageSourceMatcher{Trait: Mutant, MinPower: 5}
	g := NewGame("A", "B", 1)
	hero := g.AddToBattleline(testCreature("hero", 4, WithCannotBeDealtDamageBy(m)), 0)
	plain := g.AddToBattleline(testCreature("plain", 4), 0)
	mutant := g.AddToBattleline(testCreature("mutant", 2, WithTraits(Mutant)), 1)
	brute := g.AddToBattleline(testCreature("brute", 6), 1)
	weakling := g.AddToBattleline(testCreature("weakling", 3), 1)

	if !g.refusesDamageFrom(hero, mutant) {
		t.Error("a Mutant source should be refused")
	}
	if !g.refusesDamageFrom(hero, brute) {
		t.Error("a power-5-or-higher source should be refused")
	}
	if g.refusesDamageFrom(hero, weakling) {
		t.Error("a weak non-Mutant source should not be refused")
	}
	if g.refusesDamageFrom(hero, 0) {
		t.Error("an uncredited source should never be refused")
	}
	if g.refusesDamageFrom(plain, mutant) {
		t.Error("a creature without the passive should not refuse anything")
	}

	// The refusal reaches the damage pipeline: a Mutant's blow lands nothing, a
	// weakling's lands in full.
	g.applyRawDamage(DamageTarget{ID: hero, Amount: 3, Source: mutant})
	if g.Damage(hero) != 0 {
		t.Errorf("Mutant-dealt damage landed: hero has %d", g.Damage(hero))
	}
	g.applyRawDamage(DamageTarget{ID: hero, Amount: 3, Source: weakling})
	if g.Damage(hero) != 3 {
		t.Errorf("weak source's damage refused: hero has %d, want 3", g.Damage(hero))
	}
}
