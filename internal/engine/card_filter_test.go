package engine

import "testing"

// TestCardFilter exercises the CardFilter predicate directly: an empty filter
// admits any card, a clause conjoins its Type/Trait/Name predicates, and a pure-Or
// filter (no predicate of its own) delegates entirely to its alternatives.
func TestCardFilter(t *testing.T) {
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
	gigBase := g.Register(
		NewCard("colossus", StarAlliance, Creature, Common,
			WithPower(6), WithGiganticRole(GiganticBase)),
		0,
	)
	gigArt := g.Register(
		NewCard("colossus", StarAlliance, Creature, Common, WithGiganticRole(GiganticArt)),
		0,
	)

	cases := []struct {
		name   string
		filter CardFilter
		id     LocalID
		want   bool
	}{
		{"empty admits any", CardFilter{}, human, true},
		{"type match", CardFilter{Type: Creature}, human, true},
		{"type mismatch", CardFilter{Type: Upgrade}, human, false},
		{"trait match", CardFilter{Trait: Robot}, robot, true},
		{"trait mismatch", CardFilter{Trait: Robot}, human, false},
		{"name match", CardFilter{Name: "droid"}, robot, true},
		{"name mismatch", CardFilter{Name: "droid"}, human, false},
		{"conjunction admits", CardFilter{Type: Creature, Trait: Robot}, robot, true},
		{"conjunction rejects on trait", CardFilter{Type: Creature, Trait: Robot}, human, false},
		{
			"disjunction via Or admits type",
			CardFilter{Type: Upgrade, Or: []CardFilter{{Trait: Robot}}},
			upgrade,
			true,
		},
		{
			"disjunction via Or admits alt",
			CardFilter{Type: Upgrade, Or: []CardFilter{{Trait: Robot}}},
			robot,
			true,
		},
		{
			"disjunction via Or rejects neither",
			CardFilter{Type: Upgrade, Or: []CardFilter{{Trait: Robot}}},
			human,
			false,
		},
		{"pure-Or admits alt", CardFilter{Or: []CardFilter{{Trait: Robot}}}, robot, true},
		{"pure-Or rejects", CardFilter{Or: []CardFilter{{Trait: Robot}}}, human, false},
		{"gigantic admits base", CardFilter{Gigantic: true}, gigBase, true},
		{"gigantic admits art", CardFilter{Gigantic: true}, gigArt, true},
		{"gigantic rejects ordinary", CardFilter{Gigantic: true}, human, false},
	}
	for _, c := range cases {
		if got := c.filter.admits(g, c.id); got != c.want {
			t.Errorf("%s: admits = %v, want %v", c.name, got, c.want)
		}
	}
	if got := (CardFilter{Gigantic: true}).noun(); got != "gigantic creature" {
		t.Errorf("Gigantic noun = %q, want %q", got, "gigantic creature")
	}
}
