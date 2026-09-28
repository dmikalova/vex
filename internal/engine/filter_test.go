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

func TestTargetSharingTrait(t *testing.T) {
	if got := (Target{Kind: TargetEachCreature}).With(Filter{SharesTrait: true}).
		Text(); got != "each creature that shares a trait with it" {
		t.Errorf("shares-trait text = %q", got)
	}

	g := NewGame("A", "B", 1)
	kin := g.AddToBattleline(testCreature("kin", 3, WithTraits(Beast)), 0)
	prey := g.AddToBattleline(testCreature("prey", 5, WithTraits(Beast)), 1)
	g.AddToBattleline(testCreature("spared", 5, WithTraits(Robot)), 1)
	target := Target{Kind: TargetEachCreature}.With(Filter{SharesTrait: true})

	// Without a context card the filter matches nothing.
	noIt := &EffectContext{
		Resolver:   g,
		Controller: 0,
	}
	if ids := target.Select(noIt); len(ids) != 0 {
		t.Errorf("shares-trait without It = %v, want empty", ids)
	}

	// With the Beast in context, only trait-sharing creatures pass.
	ctx := &EffectContext{
		Resolver:   g,
		Controller: 0,
		It:         kin,
		HasIt:      true,
	}
	ids := target.Select(ctx)
	if len(ids) != 2 || ids[0] != kin || ids[1] != prey {
		t.Errorf("shares-trait select = %v, want [%d %d]", ids, kin, prey)
	}
}

func TestTargetPowerFilters(t *testing.T) {
	g := NewGame("A", "B", 1)
	p2 := g.AddToBattleline(testCreature("p2", 2), 0)
	p4 := g.AddToBattleline(testCreature("p4", 4), 0)
	p6 := g.AddToBattleline(testCreature("p6", 6), 0)
	ctx := &EffectContext{
		Resolver:   g,
		Controller: 0,
	}

	if ids := (Target{Kind: TargetEachCreature}).With(Filter{Power: PowerBound{Kind: BoundAtMost, Amount: 3}}).
		Select(ctx); len(
		ids,
	) != 1 ||
		ids[0] != p2 {
		t.Errorf("PowerAtMost(3) = %v, want [%d]", ids, p2)
	}
	if ids := (Target{Kind: TargetEachCreature}).With(Filter{Power: PowerBound{Kind: BoundAtLeast, Amount: 5}}).
		Select(ctx); len(
		ids,
	) != 1 ||
		ids[0] != p6 {
		t.Errorf("PowerAtLeast(5) = %v, want [%d]", ids, p6)
	}
	if ids := (Target{Kind: TargetEachCreature}).With(Filter{Power: PowerBound{Kind: BoundExactly, Amount: 4}}).
		Select(ctx); len(
		ids,
	) != 1 ||
		ids[0] != p4 {
		t.Errorf("PowerExactly(4) = %v, want [%d]", ids, p4)
	}

	if got := (Target{Kind: TargetChosenCreature}).With(Filter{Power: PowerBound{Kind: BoundAtLeast, Amount: 5}}).
		Text(); got != "a creature with power 5 or higher" {
		t.Errorf("PowerAtLeast text = %q", got)
	}
	if got := (Target{Kind: TargetChosenCreature}).With(Filter{Power: PowerBound{Kind: BoundExactly, Amount: 1}}).
		Text(); got != "a creature with power 1" {
		t.Errorf("PowerExactly text = %q", got)
	}
}

func TestTargetPowerParity(t *testing.T) {
	g := NewGame("A", "B", 1)
	p2 := g.AddToBattleline(testCreature("p2", 2), 0)
	p3 := g.AddToBattleline(testCreature("p3", 3), 0)
	p4 := g.AddToBattleline(testCreature("p4", 4), 0)
	ctx := &EffectContext{
		Resolver:   g,
		Controller: 0,
	}

	if ids := (Target{Kind: TargetEachCreature}).With(Filter{Power: PowerBound{Kind: BoundOdd}}).
		Select(ctx); len(ids) != 1 ||
		ids[0] != p3 {
		t.Errorf("OddPower = %v, want [%d]", ids, p3)
	}
	if ids := (Target{Kind: TargetEachCreature}).With(Filter{Power: PowerBound{Kind: BoundEven}}).
		Select(ctx); len(ids) != 2 ||
		ids[0] != p2 ||
		ids[1] != p4 {
		t.Errorf("EvenPower = %v, want [%d %d]", ids, p2, p4)
	}
	if got := (Target{Kind: TargetEachCreature}).With(Filter{Power: PowerBound{Kind: BoundOdd}}).
		Text(); got != "each creature with odd power" {
		t.Errorf("OddPower text = %q", got)
	}
	if got := (Target{Kind: TargetEachCreature}).With(Filter{Power: PowerBound{Kind: BoundEven}}).
		Text(); got != "each creature with even power" {
		t.Errorf("EvenPower text = %q", got)
	}
}

func TestTargetUndamagedAndOther(t *testing.T) {
	g := NewGame("A", "B", 1)
	src := g.AddToBattleline(testCreature("src", 3), 0)
	hurt := g.AddToBattleline(testCreature("hurt", 3), 0)
	g.State.Cards[hurt].Damage = 1
	ctx := &EffectContext{
		Resolver:   g,
		Source:     src,
		Controller: 0,
	}

	if ids := (Target{Kind: TargetEachCreature}).With(Filter{Damage: DamageNone}).
		Select(ctx); len(ids) != 1 ||
		ids[0] != src {
		t.Errorf("Undamaged filter = %v, want [%d]", ids, src)
	}
	if ids := (Target{Kind: TargetEachCreature}).With(Filter{Except: ExcludeSource}).
		Select(ctx); len(ids) != 1 ||
		ids[0] != hurt {
		t.Errorf("Other filter = %v, want [%d]", ids, hurt)
	}
	if got := (Target{Kind: TargetEachCreature}).With(Filter{Except: ExcludeSource, Damage: DamageNone}).
		Text(); got != "each other undamaged creature" {
		t.Errorf("other+undamaged text = %q", got)
	}
}

func TestTargetReady(t *testing.T) {
	g := NewGame("A", "B", 1)
	src := g.AddToBattleline(testCreature("src", 3), 0)
	spent := g.AddToBattleline(testCreature("spent", 3), 0)
	g.State.Cards[spent].Exhausted = true
	ctx := &EffectContext{
		Resolver:   g,
		Source:     src,
		Controller: 0,
	}

	if ids := (Target{Kind: TargetEachCreature}).With(Filter{Ready: true}).
		Select(ctx); len(ids) != 1 ||
		ids[0] != src {
		t.Errorf("Ready filter = %v, want [%d]", ids, src)
	}
	if got := (Target{Kind: TargetEachFriendlyCreature}).With(Filter{Ready: true}).
		Text(); got != "each friendly ready creature" {
		t.Errorf("ready text = %q", got)
	}
}

func TestTargetWithoutAember(t *testing.T) {
	g := NewGame("A", "B", 1)
	rich := g.AddToBattleline(testCreature("rich", 5), 0)
	g.State.Cards[rich].Amber = 2
	bare := g.AddToBattleline(testCreature("bare", 3), 0)
	ctx := &EffectContext{
		Resolver:   g,
		Controller: 0,
	}

	if ids := (Target{Kind: TargetEachCreature}).With(Filter{Aember: AemberNone}).
		Select(ctx); len(ids) != 1 || ids[0] != bare {
		t.Errorf("WithoutAember = %v, want [%d]", ids, bare)
	}
	if got := (Target{Kind: TargetChosenCreature}).With(Filter{Aember: AemberNone}).
		Text(); got != "a creature with no Æmber on it" {
		t.Errorf("WithoutAember text = %q", got)
	}
}

func TestTargetWithoutBonusIcons(t *testing.T) {
	g := NewGame("A", "B", 1)
	g.AddToBattleline(testCreature("iconed", 5, WithBonus(BonusAember)), 0)
	bare := g.AddToBattleline(testCreature("bare", 3), 0)
	ctx := &EffectContext{
		Resolver:   g,
		Controller: 0,
	}

	if ids := (Target{Kind: TargetEachCreature}).With(Filter{NoBonusIcons: true}).
		Select(ctx); len(ids) != 1 || ids[0] != bare {
		t.Errorf("WithoutBonusIcons = %v, want [%d]", ids, bare)
	}
	if got := (Target{Kind: TargetChosenCreature}).With(Filter{NoBonusIcons: true}).
		Text(); got != "a creature with no bonus icons" {
		t.Errorf("WithoutBonusIcons text = %q", got)
	}
}

// TestHouseWithAtLeast covers the Filter axis that counts a creature's house
// across both battlelines. It reads the whole board but is decided one creature
// at a time, which is what makes it a filter rather than a Refinement.
func TestHouseWithAtLeast(t *testing.T) {
	// Text renders the "belongs to a house" clause with the threshold.
	want := "each creature that belongs to a house that has 3 or more creatures in play"
	if got := (Target{Kind: TargetEachCreature}).With(Filter{HouseWithAtLeast: 3}).
		Text(); got != want {
		t.Errorf("text = %q", got)
	}

	// Mars has three creatures split across both players; Sanctum has one. The
	// refinement keeps the three Mars creatures and drops the lone Sanctum creature,
	// proving both players' creatures count toward one house's total.
	g := NewGame("A", "B", 1)
	m1 := g.AddToBattleline(NewCard("m1", Mars, Creature, Common, WithPower(3)), 0)
	m2 := g.AddToBattleline(NewCard("m2", Mars, Creature, Common, WithPower(3)), 0)
	m3 := g.AddToBattleline(NewCard("m3", Mars, Creature, Common, WithPower(3)), 1)
	g.AddToBattleline(NewCard("s1", Sanctum, Creature, Common, WithPower(3)), 0)
	ctx := &EffectContext{
		Resolver:   g,
		Controller: 0,
	}
	got := (Target{Kind: TargetEachCreature}).With(Filter{HouseWithAtLeast: 3}).Select(ctx)
	if len(got) != 3 || !containsID(got, m1) || !containsID(got, m2) || !containsID(got, m3) {
		t.Errorf("HouseWithAtLeast(3) = %v, want the three Mars creatures", got)
	}

	// Raising the threshold above every house's count keeps nothing.
	none := (Target{Kind: TargetEachCreature}).With(Filter{HouseWithAtLeast: 4}).Select(ctx)
	if len(none) != 0 {
		t.Errorf("HouseWithAtLeast(4) = %v, want nothing", none)
	}
}

// TestWithoutSharedTrait covers the Filter axis that checks one creature against
// its own battleline — a per-candidate test like SharesTrait beside it, not a
// comparison between candidates.
func TestWithoutSharedTrait(t *testing.T) {
	// Text renders the "does not share a trait" clause.
	want := "each creature that does not share a trait with another creature in its controller's battleline"
	if got := (Target{Kind: TargetEachCreature}).With(Filter{WithoutSharedTrait: true}).
		Text(); got != want {
		t.Errorf("text = %q", got)
	}

	// P0 has two Beasts (they share a trait, so neither is a loner) and one
	// Human whose only trait-sharer sits in the ENEMY battleline. P1 has that
	// lone Human. The refinement keeps the two Humans (each a loner in its own
	// battleline) and drops the two Beasts.
	g := NewGame("A", "B", 1)
	g.AddToBattleline(NewCard("b1", Mars, Creature, Common, WithPower(3), WithTraits(Beast)), 0)
	g.AddToBattleline(NewCard("b2", Mars, Creature, Common, WithPower(3), WithTraits(Beast)), 0)
	h0 := g.AddToBattleline(
		NewCard("h0", Mars, Creature, Common, WithPower(3), WithTraits(Human)),
		0,
	)
	h1 := g.AddToBattleline(
		NewCard("h1", Sanctum, Creature, Common, WithPower(3), WithTraits(Human)),
		1,
	)
	ctx := &EffectContext{
		Resolver:   g,
		Controller: 0,
	}

	got := (Target{Kind: TargetEachCreature}).With(Filter{WithoutSharedTrait: true}).Select(ctx)
	if len(got) != 2 || !containsID(got, h0) || !containsID(got, h1) {
		t.Errorf("WithoutSharedTrait = %v, want the two Human loners", got)
	}
}

func TestTargetKeyword(t *testing.T) {
	g := NewGame("A", "B", 1)
	elusive := g.AddToBattleline(
		NewCard("elu", Brobnar, Creature, Common, WithKeywords(Elusive)),
		0,
	)
	g.AddToBattleline(testCreature("plain", 3), 0)
	ctx := &EffectContext{
		Resolver:   g,
		Controller: 0,
	}

	if ids := (Target{Kind: TargetEachCreature}).With(Filter{Keyword: Elusive}).
		Select(ctx); len(ids) != 1 ||
		ids[0] != elusive {
		t.Errorf("Keyword(Elusive) = %v, want [%d]", ids, elusive)
	}
	if got := (Target{Kind: TargetEachCreature}).With(Filter{Keyword: Elusive}).
		Text(); got != "each elusive creature" {
		t.Errorf("keyword text = %q", got)
	}
}

func TestTargetOfHouse(t *testing.T) {
	g := NewGame("A", "B", 1)
	mars := g.AddToBattleline(NewCard("m", Mars, Creature, Common, WithPower(3)), 0)
	g.AddToBattleline(NewCard("s", Sanctum, Creature, Common, WithPower(3)), 0)
	ctx := &EffectContext{
		Resolver:   g,
		Source:     mars,
		Controller: 0,
	}

	ids := (Target{Kind: TargetEachFriendlyCreature}).With(Filter{House: namedHouse(Mars)}).
		Select(ctx)
	if len(ids) != 1 || ids[0] != mars {
		t.Errorf("OfHouse(Mars) = %v, want [%d] (Sanctum creature filtered out)", ids, mars)
	}
	if got := (Target{Kind: TargetEachCreature}).With(Filter{House: namedHouse(Mars)}).
		Text(); got != "each Mars creature" {
		t.Errorf("OfHouse text = %q", got)
	}
}

// TestTargetOfActiveHouse covers Techivore Pulpate's target: only artifacts of
// the player's active house are selected, and the phrase reads "of that house".
func TestTargetOfActiveHouse(t *testing.T) {
	g := NewGame("A", "B", 1)
	g.State.ActiveHouse = Mars
	mars := g.AddArtifact(NewCard("m", Mars, Artifact, Common), 0)
	g.AddArtifact(NewCard("s", Sanctum, Artifact, Common), 1)
	ctx := &EffectContext{
		Resolver:   g,
		Controller: 0,
	}

	ids := (Target{Kind: TargetEachArtifact}).With(Filter{House: activeHouse}).Select(ctx)
	if len(ids) != 1 || ids[0] != mars {
		t.Errorf("OfActiveHouse = %v, want [%d] (Sanctum artifact filtered out)", ids, mars)
	}
	if got := (Target{Kind: TargetEachArtifact}).With(Filter{House: activeHouse}).
		Text(); got != "each artifact of that house" {
		t.Errorf("OfActiveHouse text = %q", got)
	}
}

// TestTargetMatchingAny covers EMP Blast's target: house and trait disjoin into
// one set, so a Mars Robot is one member of it rather than a member of two
// sequenced targets — which matters because the same renderer produces "deal 1
// damage to each Mars or Robot creature", where being hit twice would be wrong.
func TestTargetMatchingAny(t *testing.T) {
	g := NewGame("A", "B", 1)
	martian := g.AddToBattleline(
		NewCard("m", Mars, Creature, Common, WithPower(3), WithTraits(Martian)), 0)
	robot := g.AddToBattleline(
		NewCard("r", Logos, Creature, Common, WithPower(3), WithTraits(Robot)), 0)
	marsRobot := g.AddToBattleline(
		NewCard("mr", Mars, Creature, Common, WithPower(3), WithTraits(Robot)), 0)
	neither := g.AddToBattleline(
		NewCard("n", Logos, Creature, Common, WithPower(3)), 0)
	ctx := &EffectContext{
		Resolver:   g,
		Source:     martian,
		Controller: 0,
	}

	either := (Target{Kind: TargetEachCreature}).With(
		Filter{House: namedHouse(Mars), Trait: Robot, MatchAny: true},
	)
	got := either.Select(ctx)
	if len(got) != 3 || containsID(got, neither) {
		t.Errorf("MatchingAny selected %v, want the three matching creatures once each", got)
	}
	hits := 0
	for _, id := range got {
		if id == marsRobot {
			hits++
		}
	}
	if hits != 1 {
		t.Errorf("a creature matching both halves appears %d times in %v, want 1", hits, got)
	}
	if !containsID(got, martian) || !containsID(got, robot) {
		t.Errorf("MatchingAny selected %v, want both halves included", got)
	}
	if text := either.Text(); text != "each Mars or Robot creature" {
		t.Errorf("MatchingAny text = %q", text)
	}

	both := (Target{Kind: TargetEachCreature}).With(Filter{House: namedHouse(Mars), Trait: Robot})
	if conj := both.Select(ctx); len(conj) != 1 || conj[0] != marsRobot {
		t.Errorf("without MatchingAny the axes still conjoin, got %v", conj)
	}
}

// TestTargetMatchingAnyDegenerateAxes covers the edges of the disjunction: a
// MatchAny that sets one axis has nothing to join and reads as that axis alone,
// and a house that renders after the noun is a branch like any other — it prints
// its own phrase rather than folding into the adjectives.
func TestTargetMatchingAnyDegenerateAxes(t *testing.T) {
	lone := (Target{Kind: TargetEachCreature}).With(Filter{House: namedHouse(Mars), MatchAny: true})
	if got := lone.Text(); got != "each Mars creature" {
		t.Errorf("text = %q, want the single axis alone", got)
	}
	suffix := (Target{Kind: TargetEachCreature}).With(
		Filter{House: HouseMatcher{Kind: MatchChosenHouse}, Trait: Robot, MatchAny: true},
	)
	if got := suffix.Text(); got !=
		"each creature of the chosen house or Robot creature" {
		t.Errorf("suffix-house text = %q", got)
	}
}

func TestTargetExceptTrait(t *testing.T) {
	g := NewGame("A", "B", 1)
	agent := g.AddToBattleline(
		NewCard("a", Mars, Creature, Common, WithPower(3), WithTraits(Agent)),
		0,
	)
	martian := g.AddToBattleline(
		NewCard("m", Mars, Creature, Common, WithPower(3), WithTraits(Martian)),
		0,
	)
	ctx := &EffectContext{
		Resolver:   g,
		Source:     agent,
		Controller: 0,
	}

	ids := (Target{Kind: TargetEachCreature}).With(Filter{House: namedHouse(Mars), ExceptTrait: Agent}).
		Select(ctx)
	if len(ids) != 1 || ids[0] != martian {
		t.Errorf("ExceptTrait(Agent) = %v, want [%d] (Agent filtered out)", ids, martian)
	}
	if got := (Target{Kind: TargetChosenCreature}).With(Filter{House: namedHouse(Mars), ExceptTrait: Agent}).
		Text(); got != "a non-Agent Mars creature" {
		t.Errorf("ExceptTrait text = %q", got)
	}
}

// TestOfHouseWithMostCreatures covers Etaromme's target: it keeps only creatures
// of the house with the most creatures in play, counting both battlelines, and
// keeps every tied house eligible on a tie.
func TestOfHouseWithMostCreatures(t *testing.T) {
	if got := (Target{Kind: TargetChosenCreature}).With(Filter{HouseWithMostCreatures: true}).
		Text(); got != "a creature of the house with the most creatures in play" {
		t.Errorf("text = %q", got)
	}

	t.Run("keeps only the most populous house", func(t *testing.T) {
		g := NewGame("A", "B", 1)
		// Mars leads with three creatures; Brobnar has two, Dis has one.
		for range 3 {
			g.AddToBattleline(NewCard("m", Mars, Creature, Common, WithPower(3)), 0)
		}
		g.AddToBattleline(NewCard("b", Brobnar, Creature, Common, WithPower(3)), 0)
		g.AddToBattleline(NewCard("b2", Brobnar, Creature, Common, WithPower(3)), 1)
		dis := g.AddToBattleline(NewCard("d", Dis, Creature, Common, WithPower(3)), 1)
		ctx := &EffectContext{
			Resolver:   g,
			Source:     dis,
			Controller: 0,
		}

		ids := (Target{Kind: TargetEachCreature}).With(Filter{HouseWithMostCreatures: true}).
			Select(ctx)
		if len(ids) != 3 {
			t.Fatalf("selected %v, want the 3 Mars creatures", ids)
		}
		for _, id := range ids {
			if g.House(id) != Mars {
				t.Errorf("selected %d of house %v, want Mars", id, g.House(id))
			}
		}
	})

	t.Run("keeps every tied house on a tie", func(t *testing.T) {
		g := NewGame("A", "B", 1)
		m1 := g.AddToBattleline(NewCard("m", Mars, Creature, Common, WithPower(3)), 0)
		g.AddToBattleline(NewCard("m2", Mars, Creature, Common, WithPower(3)), 0)
		g.AddToBattleline(NewCard("b", Brobnar, Creature, Common, WithPower(3)), 1)
		g.AddToBattleline(NewCard("b2", Brobnar, Creature, Common, WithPower(3)), 1)
		ctx := &EffectContext{
			Resolver:   g,
			Source:     m1,
			Controller: 0,
		}

		ids := (Target{Kind: TargetEachCreature}).With(Filter{HouseWithMostCreatures: true}).
			Select(ctx)
		if len(ids) != 4 {
			t.Fatalf("selected %v, want all 4 creatures of the two tied houses", ids)
		}
	})
}

func TestTargetWithUpgrade(t *testing.T) {
	g := NewGame("A", "B", 1)
	upgraded := g.AddToBattleline(testCreature("up", 3), 0)
	bare := g.AddToBattleline(testCreature("bare", 3), 0)
	attachUpgrade(g, upgraded, NewCard("plating", Mars, Upgrade, Common))
	ctx := &EffectContext{
		Resolver:   g,
		Controller: 0,
	}

	if ids := (Target{Kind: TargetEachCreature}).With(Filter{Upgrade: true}).
		Select(ctx); len(ids) != 1 || ids[0] != upgraded {
		t.Errorf("WithUpgrade = %v, want [%d] (the bare creature %d filtered out)",
			ids, upgraded, bare)
	}
	if got := (Target{Kind: TargetEachCreature}).With(Filter{Upgrade: true}).
		Text(); got != "each creature with an upgrade" {
		t.Errorf("with-upgrade text = %q", got)
	}
}

func TestTargetSharesHouseWithNeighbors(t *testing.T) {
	g := NewGame("A", "B", 1)
	left := g.AddToBattleline(NewCard("l", Mars, Creature, Common, WithPower(3)), 0)
	mid := g.AddToBattleline(NewCard("m", Mars, Creature, Common, WithPower(3)), 0)
	right := g.AddToBattleline(NewCard("r", Mars, Creature, Common, WithPower(3)), 0)
	ctx := &EffectContext{
		Resolver:   g,
		Controller: 0,
	}

	// Every creature shares its house with at least one neighbor.
	if ids := (Target{Kind: TargetEachCreature}).With(Filter{SharesHouseWithNeighbors: 1}).
		Select(ctx); len(ids) != 3 ||
		ids[0] != left || ids[1] != mid || ids[2] != right {
		t.Errorf("SharesHouseWithNeighbors(1) = %v, want [%d %d %d]", ids, left, mid, right)
	}
	// Only the middle creature has two same-house neighbors; the flanks have one.
	if ids := (Target{Kind: TargetEachCreature}).With(Filter{SharesHouseWithNeighbors: 2}).
		Select(ctx); len(ids) != 1 || ids[0] != mid {
		t.Errorf("SharesHouseWithNeighbors(2) = %v, want [%d]", ids, mid)
	}
	if got := (Target{Kind: TargetEachCreature}).With(Filter{SharesHouseWithNeighbors: 1}).
		Text(); got != "each creature that shares a house with at least 1 of its neighbors" {
		t.Errorf("shares-house(1) text = %q", got)
	}
	if got := (Target{Kind: TargetEachCreature}).With(Filter{SharesHouseWithNeighbors: 2}).
		Text(); got != "each creature that shares a house with 2 of its neighbors" {
		t.Errorf("shares-house(2) text = %q", got)
	}
}

// TestPowerLessThanSource covers the source-relative refinement: it keeps the
// creatures whose power is below the source card's own power (Dreadbone Decimus
// destroys a creature with lower power than itself) and renders "... with lower
// power than <self>".
func TestPowerLessThanSource(t *testing.T) {
	g := NewGame("A", "B", 1)
	source := g.AddToBattleline(testCreature("source", 3), 0)
	weak := g.AddToBattleline(testCreature("weak", 1), 1)
	equal := g.AddToBattleline(testCreature("equal", 3), 1)
	strong := g.AddToBattleline(testCreature("strong", 5), 1)
	ctx := &EffectContext{
		Resolver:   g,
		Controller: 0,
		Source:     source,
	}

	got := Target{
		Kind: TargetEachEnemyCreature,
	}.With(Filter{Power: PowerBound{Kind: BoundLessThanSource}}).
		Select(ctx)
	if len(got) != 1 || got[0] != weak || containsID(got, equal) || containsID(got, strong) {
		t.Errorf("PowerLessThanSource = %v, want [weak]", got)
	}

	if text := (Target{Kind: TargetChosenEnemyCreature}).With(Filter{Power: PowerBound{Kind: BoundLessThanSource}}).
		Text(); text != "an enemy creature with lower power than "+SelfName {
		t.Errorf("PowerLessThanSource text = %q", text)
	}
}

func TestTargetNamed(t *testing.T) {
	g := NewGame("A", "B", 1)
	bear := g.Register(NewCard("Ancient Bear", Untamed, Creature, Common, WithPower(6)), 0)
	g.State.Battleline[0].add(bear)
	other := g.AddToBattleline(testCreature("Chuff Ape", 6), 0)
	ctx := &EffectContext{
		Resolver:   g,
		Controller: 0,
	}

	tgt := Target{Kind: TargetEachFriendlyCreature}.With(Filter{Name: "Ancient Bear"})
	// A named card needs no describing, so the name replaces the noun outright.
	if want := "each friendly Ancient Bear"; tgt.Text() != want {
		t.Errorf("text = %q, want %q", tgt.Text(), want)
	}
	got := tgt.Select(ctx)
	if len(got) != 1 || got[0] != bear {
		t.Errorf("selected %v, want just the bear (not %v)", got, other)
	}
}
