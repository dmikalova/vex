package engine

import "testing"

// TestFilterIdentityAxes exercises the Filter predicate directly over the axes a
// card in a pile can be asked about: an empty filter admits any card, the axes
// conjoin by default, and MatchAny disjoins them instead.
func TestFilterIdentityAxes(t *testing.T) {
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
		filter Filter
		id     LocalID
		want   bool
	}{
		{"empty admits any", Filter{}, human, true},
		{"type match", Filter{Type: Creature}, human, true},
		{"type mismatch", Filter{Type: Upgrade}, human, false},
		{"trait match", Filter{Trait: Robot}, robot, true},
		{"trait mismatch", Filter{Trait: Robot}, human, false},
		{"name match", Filter{Name: "droid"}, robot, true},
		{"name mismatch", Filter{Name: "droid"}, human, false},
		{"conjunction admits", Filter{
			Type:  Creature,
			Trait: Robot,
		}, robot, true},
		{"conjunction rejects on trait", Filter{
			Type:  Creature,
			Trait: Robot,
		}, human, false},
		{
			"MatchAny admits on type",
			Filter{
				Type:     Upgrade,
				Trait:    Robot,
				MatchAny: true,
			},
			upgrade,
			true,
		},
		{
			"MatchAny admits on trait",
			Filter{
				Type:     Upgrade,
				Trait:    Robot,
				MatchAny: true,
			},
			robot,
			true,
		},
		{
			"MatchAny rejects neither",
			Filter{
				Type:     Upgrade,
				Trait:    Robot,
				MatchAny: true,
			},
			human,
			false,
		},
		{"gigantic admits base", Filter{Gigantic: true}, gigBase, true},
		{"gigantic admits art", Filter{Gigantic: true}, gigArt, true},
		{"gigantic rejects ordinary", Filter{Gigantic: true}, human, false},
	}
	ctx := &EffectContext{Resolver: g, Controller: 0}
	for _, c := range cases {
		if got := c.filter.matches(ctx, c.id); got != c.want {
			t.Errorf("%s: matches = %v, want %v", c.name, got, c.want)
		}
	}
}

// TestFilterNoun pins the unsplit rendering a pile consumer and the
// cannot-be-dealt-damage-by passive print: the card type supplies the noun, a
// MatchAny whose branches are bare adjectives folds onto one noun, and one whose
// branches are not prints each branch in full.
func TestFilterNoun(t *testing.T) {
	cases := []struct {
		name   string
		filter Filter
		base   string
		want   string
	}{
		{"empty keeps the base noun", Filter{}, "card", "card"},
		{"type supplies the noun", Filter{Type: Upgrade}, "card", "upgrade"},
		{"gigantic names both halves", Filter{Gigantic: true}, "card", "gigantic creature"},
		{"name replaces the noun", Filter{Name: "Angry Mob"}, "card", "Angry Mob"},
		{
			"house and trait conjoin",
			Filter{House: namedHouse(Mars), Trait: Robot},
			"creature",
			"Mars Robot creature",
		},
		{
			"adjective branches fold onto one noun",
			Filter{House: namedHouse(Mars), Trait: Robot, MatchAny: true},
			"creature",
			"Mars or Robot creature",
		},
		{
			"a branch with a noun of its own prints in full",
			Filter{Type: Upgrade, Trait: Robot, MatchAny: true},
			"card",
			"upgrade or Robot card",
		},
		{
			"a branch with a clause prints in full",
			Filter{
				Trait:    Mutant,
				Power:    PowerBound{Kind: BoundAtLeast, Amount: 5},
				MatchAny: true,
			},
			"creatures",
			"Mutant creatures or creatures with power 5 or higher",
		},
		{
			"a suffix house is a branch like any other",
			Filter{House: HouseMatcher{Kind: MatchChosenHouse}, Trait: Robot, MatchAny: true},
			"creature",
			"creature of the chosen house or Robot creature",
		},
		{"MatchAny with no axis set adds nothing", Filter{MatchAny: true}, "card", "card"},
	}
	for _, c := range cases {
		if got := c.filter.noun(c.base); got != c.want {
			t.Errorf("%s: noun = %q, want %q", c.name, got, c.want)
		}
	}
}

// TestFilterValidateIdentityOnly pins the guard a pile consumer puts on its
// Filter: the identity axes are legal wherever a card sits, and an in-play axis is
// a definition error rather than a filter that silently matches nothing, because a
// card on top of a deck has no power, damage, or place in a battleline to read.
func TestFilterValidateIdentityOnly(t *testing.T) {
	ok := Search{
		Sources: []Zone{Deck},
		Dest:    ToHand,
		Filter:  Filter{Type: Creature, House: namedHouse(Mars), Trait: Robot},
	}
	if err := ok.validate(); err != nil {
		t.Errorf("identity axes rejected: %v", err)
	}
	bad := Search{
		Sources: []Zone{Deck},
		Dest:    ToHand,
		Filter:  Filter{Power: PowerBound{Kind: BoundAtMost, Amount: 3}},
	}
	if err := bad.validate(); err == nil {
		t.Error("a power bound on a deck search should be rejected")
	}
	if err := (DiscardUntil{Filter: Filter{Stunned: true}}).validate(); err == nil {
		t.Error("a stun axis on a deck dig should be rejected")
	}
	if err := (DiscardUntil{Filter: Filter{Name: "Angry Mob"}}).validate(); err != nil {
		t.Errorf("a name axis on a deck dig rejected: %v", err)
	}
}

// TestCardRejectsBareFilterRefinement pins the rejection at the place it bites: a
// card whose ability passes a Filter to Refine rather than to With never reaches a
// game. Targets sit at many depths in a definition, so the check walks the whole
// definition rather than trusting each effect to validate its own target.
func TestCardRejectsBareFilterRefinement(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Error("a card with a bare Filter refinement should be rejected")
		}
	}()
	NewCard("Bad", Brobnar, Creature, Common, WithPower(1),
		WithAbility(TriggerAction, Destroy{
			Target: Target{Kind: TargetEachCreature}.Refine(Filter{Trait: Robot}),
		}))
}
