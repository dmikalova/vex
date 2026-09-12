package engine

import "testing"

// TestMatch exercises the Match predicate directly: an empty match admits any
// card, a clause conjoins its Type/Trait/Name predicates, and a pure-Or match
// (no predicate of its own) delegates entirely to its alternatives.
func TestMatch(t *testing.T) {
	g := NewGame("A", "B", 1)
	robot := g.Register(
		NewCard("droid", StarAlliance, Creature, Common, WithPower(3), WithTraits(Robot)),
		0,
	)
	human := g.Register(
		NewCard("pilot", StarAlliance, Creature, Common, WithPower(3), WithTraits(Human)),
		0,
	)
	upgrade := g.Register(NewCard("chip", StarAlliance, Upgrade, Common), 0)

	cases := []struct {
		name  string
		match Match
		id    LocalID
		want  bool
	}{
		{"empty admits any", Match{}, human, true},
		{"type match", Match{Type: Creature}, human, true},
		{"type mismatch", Match{Type: Upgrade}, human, false},
		{"trait match", Match{Trait: Robot}, robot, true},
		{"trait mismatch", Match{Trait: Robot}, human, false},
		{"name match", Match{Name: "droid"}, robot, true},
		{"name mismatch", Match{Name: "droid"}, human, false},
		{"conjunction admits", Match{Type: Creature, Trait: Robot}, robot, true},
		{"conjunction rejects on trait", Match{Type: Creature, Trait: Robot}, human, false},
		{
			"disjunction via Or admits type",
			Match{Type: Upgrade, Or: []Match{{Trait: Robot}}},
			upgrade,
			true,
		},
		{
			"disjunction via Or admits alt",
			Match{Type: Upgrade, Or: []Match{{Trait: Robot}}},
			robot,
			true,
		},
		{
			"disjunction via Or rejects neither",
			Match{Type: Upgrade, Or: []Match{{Trait: Robot}}},
			human,
			false,
		},
		{"pure-Or admits alt", Match{Or: []Match{{Trait: Robot}}}, robot, true},
		{"pure-Or rejects", Match{Or: []Match{{Trait: Robot}}}, human, false},
	}
	for _, c := range cases {
		if got := c.match.admits(g, c.id); got != c.want {
			t.Errorf("%s: admits = %v, want %v", c.name, got, c.want)
		}
	}
}
