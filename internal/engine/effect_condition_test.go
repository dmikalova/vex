package engine

import "testing"

// TestHouseChoicePhrase covers the trailing "of the … house" fragment each choice
// renders, and the empty phrase the no-scope choices return.
func TestHouseChoicePhrase(t *testing.T) {
	cases := map[HouseChoice]string{
		TheChosenHouse:     "of the chosen house",
		TheActiveHouse:     "of the active house",
		TheContextualHouse: "of that card's house",
		AnyHouse:           "",
		houseChoiceUnset:   "",
	}
	for choice, want := range cases {
		if got := choice.phrase(); got != want {
			t.Errorf("phrase(%v) = %q, want %q", choice, got, want)
		}
	}
}

func TestConditionalEffect(t *testing.T) {
	g := NewGame("A", "B", 1)
	src := g.AddToBattleline(testCreature("src", 1), 0)
	ctx := &EffectContext{Resolver: g, Source: src, Controller: 0}

	atLeast := Conditional{
		Cond: PoolAember{Player: Opponent, Is: AtLeast, Amount: 7},
		Then: LoseAember{Player: Opponent, Amount: 4},
	}
	if atLeast.Text() != "if your opponent has 7 Æmber or more, your opponent loses 4 Æmber" {
		t.Errorf("at-least text = %q", atLeast.Text())
	}
	g.State.Aember[1] = 5
	atLeast.Resolve(ctx) // condition not met (5 < 7)
	if g.State.Aember[1] != 5 {
		t.Errorf("unmet condition should be a no-op; opp = %d", g.State.Aember[1])
	}
	g.State.Aember[1] = 8
	atLeast.Resolve(ctx) // met
	if g.State.Aember[1] != 4 {
		t.Errorf("met condition should apply; opp = %d, want 4", g.State.Aember[1])
	}

	exact := Conditional{
		Cond: PoolAember{Player: Opponent, Is: Exactly, Amount: 1},
		Then: StealAember{Amount: 1},
	}
	if exact.Text() != "if your opponent has exactly 1 Æmber, steal 1 Æmber" {
		t.Errorf("exact text = %q", exact.Text())
	}
	g.State.Aember[0], g.State.Aember[1] = 0, 2
	exact.Resolve(ctx) // 2 != 1, no-op
	if g.State.Aember[0] != 0 {
		t.Errorf("exact unmet should be a no-op; you = %d", g.State.Aember[0])
	}
	g.State.Aember[1] = 1
	exact.Resolve(ctx) // met
	if g.State.Aember[0] != 1 || g.State.Aember[1] != 0 {
		t.Errorf("exact met: you=%d opp=%d, want 1/0", g.State.Aember[0], g.State.Aember[1])
	}
}

func TestItIsStunned(t *testing.T) {
	c := ItIsStunned{}
	if c.CondText() != "if that creature was already stunned" {
		t.Errorf("text = %q", c.CondText())
	}

	g := NewGame("A", "B", 1)
	stunned := g.AddToBattleline(testCreature("stunned", 3), 1)
	g.SetStunned(stunned, true)
	ready := g.AddToBattleline(testCreature("ready", 3), 1)

	if !c.Met(&EffectContext{Resolver: g, Controller: 0, It: stunned, HasIt: true}) {
		t.Error("a stunned creature should meet the condition")
	}
	if c.Met(&EffectContext{Resolver: g, Controller: 0, It: ready, HasIt: true}) {
		t.Error("an unstunned creature should not meet the condition")
	}
	if c.Met(&EffectContext{Resolver: g, Controller: 0}) {
		t.Error("no creature in context should not meet the condition")
	}
}

func TestControlsMoreCreatures(t *testing.T) {
	g := NewGame("A", "B", 1)
	ctx := &EffectContext{Resolver: g, Controller: 0}

	c := ControlsMoreCreatures{}
	if c.CondText() != "if you control more creatures than your opponent" {
		t.Errorf("CondText = %q", c.CondText())
	}
	g.AddToBattleline(testCreature("mine", 1), 0)
	if !c.Met(ctx) {
		t.Error("1 vs 0 should be met")
	}
	g.AddToBattleline(testCreature("theirs", 1), 1)
	if c.Met(ctx) {
		t.Error("1 vs 1 should not be met")
	}
}

// TestControlsMoreCreaturesTrait covers the trait-scoped comparison (Pismire),
// which counts only creatures carrying the named trait on each side.
func TestControlsMoreCreaturesTrait(t *testing.T) {
	g := NewGame("A", "B", 1)
	ctx := &EffectContext{Resolver: g, Controller: 0}

	c := ControlsMoreCreatures{Trait: Mutant}
	if got := c.CondText(); got != "if you control more Mutant creatures than your opponent" {
		t.Errorf("CondText = %q", got)
	}

	// A plain creature on each side leaves the Mutant counts level at zero.
	g.AddToBattleline(testCreature("plain", 1), 0)
	if c.Met(ctx) {
		t.Error("0 vs 0 Mutants should not be met")
	}
	// One friendly Mutant tips the Mutant comparison even against more creatures.
	g.AddToBattleline(testCreature("mine", 1, WithTraits(Mutant)), 0)
	g.AddToBattleline(testCreature("theirs", 1), 1)
	if !c.Met(ctx) {
		t.Error("1 vs 0 Mutants should be met")
	}
	// An enemy Mutant levels the Mutant counts back out.
	g.AddToBattleline(testCreature("theirMutant", 1, WithTraits(Mutant)), 1)
	if c.Met(ctx) {
		t.Error("1 vs 1 Mutants should not be met")
	}
}

func TestOpponentHasMoreKeys(t *testing.T) {
	g := NewGame("A", "B", 1)
	ctx := &EffectContext{Resolver: g, Controller: 0}

	c := HasMoreForgedKeys{Player: Opponent}
	if c.CondText() != "if your opponent has more forged keys than you" {
		t.Errorf("CondText = %q", c.CondText())
	}
	if c.Met(ctx) {
		t.Error("0 vs 0 should not be met")
	}
	g.State.Keys[1] = 1
	if !c.Met(ctx) {
		t.Error("opponent 1 vs you 0 should be met")
	}
	g.State.Keys[0] = 1
	if c.Met(ctx) {
		t.Error("1 vs 1 should not be met")
	}

	// The controller-facing sense reads and compares from the other side.
	mine := HasMoreForgedKeys{Player: Controller}
	if mine.CondText() != "if you have more forged keys than your opponent" {
		t.Errorf("CondText = %q", mine.CondText())
	}
	g.State.Keys[0] = 2
	if !mine.Met(ctx) {
		t.Error("you 2 vs opponent 1 should be met")
	}
}

// TestUsedCreatureToReap covers the condition that a creature has reaped this turn
// (Bramble Lynx), plus the "enters play ready" rendering when it gates an entry.
func TestUsedCreatureToReap(t *testing.T) {
	c := UsedCreatureToReap{}
	if got := c.CondText(); got != "if you have used a creature to reap this turn" {
		t.Errorf("CondText = %q", got)
	}

	g := NewGame("A", "B", 1)
	ctx := &EffectContext{Resolver: g, Controller: 0}
	if c.Met(ctx) {
		t.Error("no reap yet, condition should not be met")
	}
	g.State.TurnHistory[0][CreaturesReapedThisTurn] = 1
	if !c.Met(ctx) {
		t.Error("after a reap, condition should be met")
	}

	// It gates a conditional "enters play ready" ability and renders naturally.
	a := Ability{
		Trigger: TriggerEntersPlay,
		Effect: Conditional{
			Cond: UsedCreatureToReap{},
			Then: Ready{Target: Target{Kind: TargetThisCreature}},
		},
	}
	want := "If you have used a creature to reap this turn, " + SelfName + " enters play ready."
	if got := RenderAbility(a); got != want {
		t.Errorf("RenderAbility = %q, want %q", got, want)
	}
}

// TestUsedCreatureToFight covers the condition that a creature has fought this
// turn (Alaka), plus the "enters play ready" rendering when it gates an entry.
func TestUsedCreatureToFight(t *testing.T) {
	c := UsedCreatureToFight{}
	if got := c.CondText(); got != "if you have used a creature to fight this turn" {
		t.Errorf("CondText = %q", got)
	}

	g := NewGame("A", "B", 1)
	ctx := &EffectContext{Resolver: g, Controller: 0}
	if c.Met(ctx) {
		t.Error("no fight yet, condition should not be met")
	}
	g.State.TurnHistory[0][CreaturesFoughtThisTurn] = 1
	if !c.Met(ctx) {
		t.Error("after a fight, condition should be met")
	}

	a := Ability{
		Trigger: TriggerEntersPlay,
		Effect: Conditional{
			Cond: UsedCreatureToFight{},
			Then: Ready{Target: Target{Kind: TargetThisCreature}},
		},
	}
	want := "If you have used a creature to fight this turn, " + SelfName + " enters play ready."
	if got := RenderAbility(a); got != want {
		t.Errorf("RenderAbility = %q, want %q", got, want)
	}
}

// TestUsedNoCreatures covers the inverse condition — met while the controller has
// used no creatures this turn (Sloth).
func TestUsedNoCreatures(t *testing.T) {
	c := UsedNoCreatures{}
	if got := c.CondText(); got != "if you did not use any creatures this turn" {
		t.Errorf("CondText = %q", got)
	}

	g := NewGame("A", "B", 1)
	ctx := &EffectContext{Resolver: g, Controller: 0}
	if !c.Met(ctx) {
		t.Error("no creature used, condition should be met")
	}
	g.State.TurnHistory[0][CreaturesUsedThisTurn] = 1
	if c.Met(ctx) {
		t.Error("after using a creature, condition should not be met")
	}
}

// TestHousesRepresented covers the condition comparing distinct houses among a
// surveyed set (Galactic Census) across every comparison.
func TestHousesRepresented(t *testing.T) {
	g := NewGame("A", "B", 1)
	g.AddToBattleline(NewCard("m", Mars, Creature, Common, WithPower(4)), 0)
	g.AddToBattleline(NewCard("l", Logos, Creature, Common, WithPower(4)), 0)
	g.AddToBattleline(NewCard("s", Sanctum, Creature, Common, WithPower(4)), 1)
	ctx := &EffectContext{Resolver: g, Controller: 0}

	among := HousesAmong{Player: EachPlayer, Type: Creature}
	cases := []struct {
		is   Comparison
		amt  int
		met  bool
		text string
	}{
		{AtLeast, 3, true, "if there are 3 or more houses represented among creatures in play"},
		{AtLeast, 4, false, "if there are 4 or more houses represented among creatures in play"},
		{AtMost, 3, true, "if there are 3 or fewer houses represented among creatures in play"},
		{AtMost, 2, false, "if there are 2 or fewer houses represented among creatures in play"},
		{Exactly, 3, true, "if there are exactly 3 houses represented among creatures in play"},
		{Exactly, 2, false, "if there are exactly 2 houses represented among creatures in play"},
	}
	for _, tc := range cases {
		c := HousesRepresented{Among: among, Is: tc.is, Amount: tc.amt}
		if err := c.validate(); err != nil {
			t.Errorf("validate(%v, %d) = %v", tc.is, tc.amt, err)
		}
		if got := c.Met(ctx); got != tc.met {
			t.Errorf("Met(%v, %d) = %v, want %v", tc.is, tc.amt, got, tc.met)
		}
		if got := c.CondText(); got != tc.text {
			t.Errorf("CondText(%v, %d) = %q", tc.is, tc.amt, got)
		}
	}
	if err := (HousesRepresented{Among: among, Is: comparisonUnset}).validate(); err == nil {
		t.Error("unset comparison should fail validation")
	}
}

func TestNot(t *testing.T) {
	// A negatable inner condition validates and surfaces an invalid inner.
	if err := (Not{Cond: OnFlank{}}).validate(); err != nil {
		t.Errorf("validate: %v", err)
	}
	if err := (Not{Cond: ForgedKey{}}).validate(); err == nil {
		t.Error("Not should surface an invalid inner condition")
	}
	// A condition that renders no negated text cannot be wrapped.
	if err := (Not{Cond: SourceReady{}}).validate(); err == nil {
		t.Error("a non-negatable inner condition should fail validation")
	}
}

func TestSourceReadyCondition(t *testing.T) {
	g := NewGame("A", "B", 1)
	src := g.AddToBattleline(testCreature("src", 2), 0)
	c := SourceReady{}
	if c.CondText() != "if "+SelfName+" is ready" {
		t.Errorf("CondText = %q", c.CondText())
	}
	if !c.Met(&EffectContext{Resolver: g, Source: src}) {
		t.Error("a ready source should be met")
	}
	g.SetExhausted(src, true)
	if c.Met(&EffectContext{Resolver: g, Source: src}) {
		t.Error("an exhausted source should not be met")
	}
}

func TestOnFlankSource(t *testing.T) {
	g := NewGame("A", "B", 1)
	left := g.AddToBattleline(testCreature("left", 2), 0)
	mid := g.AddToBattleline(testCreature("mid", 2), 0)
	g.AddToBattleline(testCreature("right", 2), 0)

	on := OnFlank{}
	off := Not{Cond: OnFlank{}}
	if on.CondText() != "if "+SelfName+" is on a flank" {
		t.Errorf("on CondText = %q", on.CondText())
	}
	if off.CondText() != "if "+SelfName+" is not on a flank" {
		t.Errorf("off CondText = %q", off.CondText())
	}
	if !on.Met(&EffectContext{Resolver: g, Source: left}) {
		t.Error("left flank source should satisfy OnFlank{}")
	}
	if on.Met(&EffectContext{Resolver: g, Source: mid}) {
		t.Error("interior source should not satisfy OnFlank{}")
	}
	if !off.Met(&EffectContext{Resolver: g, Source: mid}) {
		t.Error("interior source should satisfy Not{OnFlank}")
	}
	if off.Met(&EffectContext{Resolver: g, Source: left}) {
		t.Error("flank source should not satisfy Not{OnFlank}")
	}
}

func TestOnFlankNamed(t *testing.T) {
	g := NewGame("A", "B", 1)
	left := g.AddToBattleline(testCreature("left", 2), 1)
	g.AddToBattleline(testCreature("mid", 2), 1)
	right := g.AddToBattleline(testCreature("right", 2), 1)

	leftCond := OnFlank{OfIt: true, Where: LeftFlank}
	rightCond := OnFlank{OfIt: true, Where: RightFlank}
	if leftCond.CondText() != "if it is on the left flank" {
		t.Errorf("left CondText = %q", leftCond.CondText())
	}
	if rightCond.CondText() != "if it is on the right flank" {
		t.Errorf("right CondText = %q", rightCond.CondText())
	}
	if !leftCond.Met(&EffectContext{Resolver: g, It: left, HasIt: true}) {
		t.Error("left flank creature should satisfy the left OnFlank")
	}
	if leftCond.Met(&EffectContext{Resolver: g, It: right, HasIt: true}) {
		t.Error("right flank creature should not satisfy the left OnFlank")
	}
	if !rightCond.Met(&EffectContext{Resolver: g, It: right, HasIt: true}) {
		t.Error("right flank creature should satisfy the right OnFlank")
	}
	if rightCond.Met(&EffectContext{Resolver: g, It: left, HasIt: true}) {
		t.Error("left flank creature should not satisfy the right OnFlank")
	}
	if leftCond.Met(&EffectContext{Resolver: g, HasIt: false}) {
		t.Error("no context creature should not satisfy a named OnFlank")
	}

	// A creature whose controller has no battleline is on no flank.
	g2 := NewGame("A", "B", 1)
	loose := g2.AddToHand(testCreature("loose", 2), 0)
	if leftCond.Met(&EffectContext{Resolver: g2, It: loose, HasIt: true}) {
		t.Error("a creature off the battleline should be on no flank")
	}
	// A non-creature context card is on no flank either.
	art := g2.AddArtifact(NewCard("art", Logos, Artifact, Common), 0)
	if leftCond.Met(&EffectContext{Resolver: g2, It: art, HasIt: true}) {
		t.Error("an artifact should not satisfy a named OnFlank")
	}
}

func TestControlsNamedCondition(t *testing.T) {
	c := ControlsNamed{Name: "Velum"}
	if c.CondText() != "if you control Velum" {
		t.Errorf("CondText = %q", c.CondText())
	}
	g := NewGame("A", "B", 1)
	if c.Met(&EffectContext{Resolver: g, Controller: 0}) {
		t.Error("no named card in play should not be met")
	}
	g.AddToBattleline(NewCard("Velum", Logos, Creature, Common, WithPower(2)), 0)
	if !c.Met(&EffectContext{Resolver: g, Controller: 0}) {
		t.Error("controlling the named card should be met")
	}
}

func TestOnFlankIt(t *testing.T) {
	c := OnFlank{OfIt: true}
	if c.CondText() != "if it is on a flank" {
		t.Errorf("CondText = %q", c.CondText())
	}
	g := NewGame("A", "B", 1)
	left := g.AddToBattleline(testCreature("left", 2), 1)
	mid := g.AddToBattleline(testCreature("mid", 2), 1)
	g.AddToBattleline(testCreature("right", 2), 1)
	if !c.Met(&EffectContext{Resolver: g, It: left, HasIt: true}) {
		t.Error("flank creature should satisfy OnFlank{OfIt}")
	}
	if c.Met(&EffectContext{Resolver: g, It: mid, HasIt: true}) {
		t.Error("interior creature should not satisfy OnFlank{OfIt}")
	}
	if c.Met(&EffectContext{Resolver: g, HasIt: false}) {
		t.Error("no context creature should not satisfy OnFlank{OfIt}")
	}
}

func TestArchivedCreaturesShareHouseCondition(t *testing.T) {
	c := ArchivedCreaturesShareHouse{}
	if c.CondText() != "if those creatures share a house" {
		t.Errorf("CondText = %q", c.CondText())
	}
	g := NewGame("A", "B", 1)
	a := g.AddToBattleline(NewCard("a", Logos, Creature, Common, WithPower(2)), 0)
	b := g.AddToBattleline(NewCard("b", Logos, Creature, Common, WithPower(2)), 0)
	d := g.AddToBattleline(testCreature("d", 2), 0) // Brobnar

	if c.Met(&EffectContext{Resolver: g}) {
		t.Error("no archived creatures cannot share a house")
	}
	ctx := &EffectContext{Resolver: g}
	ctx.Produced.Archived = []LocalID{a, d}
	if c.Met(ctx) {
		t.Error("creatures of different houses should not be met")
	}
	ctx.Produced.Archived = []LocalID{a, b}
	if !c.Met(ctx) {
		t.Error("creatures of the same house should be met")
	}
}

func TestSourceNeighborsAllOfHouseCondition(t *testing.T) {
	g := NewGame("A", "B", 1)
	marsCreature := func(name string) CardDefinition {
		return NewCard(name, Mars, Creature, Common, WithPower(2))
	}
	g.AddToBattleline(marsCreature("left"), 0)
	src := g.AddToBattleline(marsCreature("mid"), 0)
	g.AddToBattleline(testCreature("right", 2), 0) // Brobnar

	c := SourceNeighborsAllOfHouse{House: Mars}
	if c.CondText() != "if it has no non-Mars neighbor" {
		t.Errorf("CondText = %q", c.CondText())
	}
	ctx := &EffectContext{Resolver: g, Source: src}
	if c.Met(ctx) {
		t.Error("Mars creature with a Brobnar neighbor should not be met")
	}

	// A battleline of only Mars creatures satisfies the condition.
	g2 := NewGame("A", "B", 1)
	g2.AddToBattleline(marsCreature("l"), 0)
	src2 := g2.AddToBattleline(marsCreature("m"), 0)
	g2.AddToBattleline(marsCreature("r"), 0)
	if !c.Met(&EffectContext{Resolver: g2, Source: src2}) {
		t.Error("Mars creature with only Mars neighbors should be met")
	}
}

func TestSourceHasNoNeighborOfHouseCondition(t *testing.T) {
	marsCreature := func(name string) CardDefinition {
		return NewCard(name, Mars, Creature, Common, WithPower(2))
	}

	// A neighbor of the named house fails the condition.
	g := NewGame("A", "B", 1)
	g.AddToBattleline(marsCreature("left"), 0)
	src := g.AddToBattleline(testCreature("mid", 2), 0)

	c := SourceHasNoNeighborOfHouse{House: Mars}
	if c.CondText() != "if "+SelfName+" has no Mars neighbor" {
		t.Errorf("CondText = %q", c.CondText())
	}
	if c.Met(&EffectContext{Resolver: g, Source: src}) {
		t.Error("creature with a Mars neighbor should not be met")
	}

	// No neighbor of the named house satisfies the condition.
	g2 := NewGame("A", "B", 1)
	g2.AddToBattleline(testCreature("l", 2), 0)
	src2 := g2.AddToBattleline(testCreature("m", 2), 0)
	g2.AddToBattleline(testCreature("r", 2), 0)
	if !c.Met(&EffectContext{Resolver: g2, Source: src2}) {
		t.Error("creature with no Mars neighbor should be met")
	}
}

func TestControlsCreaturesOfHouses(t *testing.T) {
	g := NewGame("A", "B", 1)
	ctx := &EffectContext{Resolver: g, Controller: 0}

	c := ControlsCreaturesOfHouses{Amount: 3}
	if c.CondText() != "if you control creatures from 3 or more houses" {
		t.Errorf("CondText = %q", c.CondText())
	}
	if err := c.validate(); err != nil {
		t.Errorf("valid Count rejected: %v", err)
	}
	if err := (ControlsCreaturesOfHouses{}).validate(); err == nil {
		t.Error("non-positive Count should be rejected")
	}

	g.AddToBattleline(NewCard("b", Brobnar, Creature, Common, WithPower(2)), 0)
	g.AddToBattleline(NewCard("l", Logos, Creature, Common, WithPower(2)), 0)
	if c.Met(ctx) {
		t.Error("2 houses should not meet a Count of 3")
	}
	g.AddToBattleline(NewCard("s", Sanctum, Creature, Common, WithPower(2)), 0)
	if !c.Met(ctx) {
		t.Error("3 houses should meet a Count of 3")
	}
	// A second creature of an already-counted house does not raise the tally.
	g.AddToBattleline(NewCard("b2", Brobnar, Creature, Common, WithPower(2)), 0)
	if c.Met(ctx) != true {
		t.Error("duplicate house should still leave 3 distinct houses")
	}
	if (ControlsCreaturesOfHouses{Amount: 4}).Met(ctx) {
		t.Error("only 3 distinct houses should not meet a Count of 4")
	}
}

func TestCardsPlayed(t *testing.T) {
	g := NewGame("A", "B", 1)
	ctx := &EffectContext{Resolver: g, Controller: 0}
	cond := CardsPlayed{Player: Controller, House: namedHouse(Sanctum), Amount: 7}
	if cond.CondText() != "if you have played 7 or more Sanctum cards this turn" {
		t.Errorf("cond text = %q", cond.CondText())
	}
	if cond.CountText() != "Sanctum card you have played this turn" {
		t.Errorf("count text = %q", cond.CountText())
	}
	playSanctum := func(n int) {
		g.State.PlayedThisTurn[0].reset()
		for range n {
			g.State.PlayedThisTurn[0].add(
				g.AddToHand(NewCard("sanctum card", Sanctum, Tactic, Common), 0),
			)
		}
	}
	playSanctum(6)
	if cond.Met(ctx) || cond.Value(ctx) != 6 {
		t.Errorf("six Sanctum cards: met=%v value=%d, want false/6", cond.Met(ctx), cond.Value(ctx))
	}
	playSanctum(7)
	if !cond.Met(ctx) {
		t.Error("seven Sanctum cards should satisfy the condition")
	}
	// The default threshold is one played card.
	if (CardsPlayed{Player: Controller, House: namedHouse(Mars)}).Met(ctx) {
		t.Error("no Mars cards played should not meet the default threshold")
	}
}

func TestRepeatMayWhile(t *testing.T) {
	e := Repeat{
		Do:   StealAember{Amount: 1},
		Gate: MayWhile{Cond: PoolAember{Player: Opponent, Is: MoreThanYou}},
	}
	if got := e.Text(); got != "steal 1 Æmber -> if your opponent has more Æmber than you, you may repeat this effect" {
		t.Errorf("text = %q", got)
	}

	// Accepting every offer repeats until the condition fails.
	g := NewGame("A", "B", 1)
	g.State.Aember[0], g.State.Aember[1] = 0, 5 // opponent leads
	g.SetChooser(0, optionPicker{idx: 0})       // always "Yes"
	e.Resolve(&EffectContext{Resolver: g, Controller: 0})
	// 5/0 -> 4/1 -> 3/2 -> 2/3 (opponent no longer leads): Do once plus 2 repeats.
	if g.Aember(0) != 3 || g.Aember(1) != 2 {
		t.Errorf("accepted: you=%d opp=%d, want 3/2", g.Aember(0), g.Aember(1))
	}

	// Declining the first offer runs Do exactly once even though the condition holds.
	g2 := NewGame("A", "B", 1)
	g2.State.Aember[0], g2.State.Aember[1] = 0, 5
	g2.SetChooser(0, optionPicker{idx: 1}) // "No"
	e.Resolve(&EffectContext{Resolver: g2, Controller: 0})
	if g2.Aember(0) != 1 || g2.Aember(1) != 4 {
		t.Errorf("declined: you=%d opp=%d, want 1/4", g2.Aember(0), g2.Aember(1))
	}

	if err := validateEffect(
		Repeat{
			Do:   StealAember{Amount: 1},
			Gate: MayWhile{Cond: InPlay{Player: Controller, Type: Creature}},
		},
	); err != nil {
		t.Errorf("validate = %v", err)
	}
}

// When Do is a single clickable choice, the repeat is offered by letting the
// player keep picking rather than answering Yes/No.
func TestRepeatMayWhileDrivenByChoice(t *testing.T) {
	e := Repeat{
		Gate: MayWhile{Cond: InPlay{Player: Controller, Type: Creature}},
		Do: Sequence{Effects: []Effect{
			Destroy{Target: Target{Kind: TargetChosenEnemyCreature}},
			Destroy{Target: Target{Kind: TargetChosenFriendlyCreature}},
		}},
	}

	// Picking again keeps repeating until no friendly creature remains.
	accepted := NewGame("A", "B", 1)
	accepted.SetChooser(0, &cardDecliner{})
	foeA := accepted.AddToBattleline(testCreature("FoeA", 3), 1)
	foeB := accepted.AddToBattleline(testCreature("FoeB", 3), 1)
	allyA := accepted.AddToBattleline(testCreature("AllyA", 3), 0)
	allyB := accepted.AddToBattleline(testCreature("AllyB", 3), 0)
	e.Resolve(&EffectContext{Resolver: accepted, Controller: 0})
	for _, id := range []LocalID{foeA, foeB, allyA, allyB} {
		if onAnyLine(accepted, id) {
			t.Errorf("creature %d should have been destroyed", id)
		}
	}

	// Passing on the first offer runs Do exactly once even though a friendly remains.
	declined := NewGame("A", "B", 1)
	declined.SetChooser(0, &cardDecliner{decline: true})
	declined.AddToBattleline(testCreature("FoeA", 3), 1)
	survivorFoe := declined.AddToBattleline(testCreature("FoeB", 3), 1)
	declined.AddToBattleline(testCreature("AllyA", 3), 0)
	survivorAlly := declined.AddToBattleline(testCreature("AllyB", 3), 0)
	e.Resolve(&EffectContext{Resolver: declined, Controller: 0})
	if !onAnyLine(declined, survivorFoe) || !onAnyLine(declined, survivorAlly) {
		t.Error("declining the repeat should leave the second pair in play")
	}
}

func TestItIs(t *testing.T) {
	// CondText renders the filtered noun with the right article.
	cases := map[string]ItIs{
		"if it is a Mars creature":          {House: namedHouse(Mars), Type: Creature},
		"if it is an artifact":              {Type: Artifact},
		"if it is a Mars card":              {House: namedHouse(Mars)},
		"if it is a card":                   {},
		"if it is a non-Star Alliance card": {House: exceptHouse(StarAlliance)},
		"if it is of the chosen house":      {House: chosenHouse},
		"if it is of the active house":      {House: activeHouse},
		"if it is another creature":         {Type: Creature, Other: true},
	}
	for want, e := range cases {
		if got := e.CondText(); got != want {
			t.Errorf("CondText() = %q, want %q", got, want)
		}
	}

	// A Not wrapper renders each condition's own negated clause.
	negated := map[string]ItIs{
		"if it is not a Logos card": {House: namedHouse(Logos)},
		"if the discarded card is not a Logos card": {
			House:   namedHouse(Logos),
			Subject: DiscardedCard,
		},
	}
	for want, e := range negated {
		if got := (Not{Cond: e}).CondText(); got != want {
			t.Errorf("Not CondText() = %q, want %q", got, want)
		}
	}

	g := NewGame("Alice", "Bob", 1)
	mars := g.AddToDeck(NewCard("Martian", Mars, Creature, Common, WithPower(1)), 0)
	ctx := &EffectContext{Resolver: g, Controller: 0}

	// No card in context is never met.
	if (ItIs{Type: Creature}).Met(ctx) {
		t.Error("Met with no context card should be false")
	}
	ctx.It, ctx.HasIt = mars, true
	if !(ItIs{House: namedHouse(Mars), Type: Creature}).Met(ctx) {
		t.Error("a Mars creature should match a Mars-creature filter")
	}
	if (ItIs{House: namedHouse(Logos)}).Met(ctx) {
		t.Error("a Mars creature should not match a Logos filter")
	}
	if (ItIs{Type: Artifact}).Met(ctx) {
		t.Error("a creature should not match an artifact filter")
	}

	// Other bars the source card itself: a creature never counts its own play.
	ctx.Source = mars + 1 // any card that is not the context card
	if !(ItIs{Type: Creature, Other: true}).Met(ctx) {
		t.Error("another creature should meet an Other filter")
	}
	ctx.Source = mars
	if (ItIs{Type: Creature, Other: true}).Met(ctx) {
		t.Error("the source card should not meet an Other filter")
	}
	ctx.Source = 0

	// Not inverts the match, so the condition holds for everything that does not fit.
	if !(Not{Cond: ItIs{House: namedHouse(Logos)}}).Met(ctx) {
		t.Error("a Mars creature should meet a not-Logos filter")
	}
	if (Not{Cond: ItIs{House: namedHouse(Mars)}}).Met(ctx) {
		t.Error("a Mars creature should not meet a not-Mars filter")
	}
}

func TestChoseHouse(t *testing.T) {
	if got := (ChoseHouse{House: Sanctum}).CondText(); got != "you choose Sanctum as your active house" {
		t.Errorf("cond text = %q", got)
	}

	g := started(t) // active house Brobnar
	ctx := &EffectContext{Resolver: g, Controller: 0}
	if (ChoseHouse{House: Sanctum}).Met(ctx) {
		t.Error("Sanctum should not be met while Brobnar is active")
	}
	if !(ChoseHouse{House: Brobnar}).Met(ctx) {
		t.Error("Brobnar should be met while Brobnar is active")
	}
}

func TestAfterChooseHouseTrigger(t *testing.T) {
	def := NewCard("Bureaucrat", Sanctum, Creature, Rare, WithPower(3),
		WithAbility(TriggerAfterChooseHouse, Conditional{
			Cond: ChoseHouse{House: Sanctum},
			Then: GainAember{Player: Controller, Amount: 2},
		}))

	t.Run("gains when the watched house is chosen", func(t *testing.T) {
		g := NewGame("A", "B", 1)
		g.StartTurn(0)
		g.AddToBattleline(def, 0)
		if err := g.ChooseHouse(0, Sanctum); err != nil {
			t.Fatalf("ChooseHouse: %v", err)
		}
		if g.Aember(0) != 2 {
			t.Errorf("aember = %d, want 2", g.Aember(0))
		}
	})

	t.Run("does nothing when a different house is chosen", func(t *testing.T) {
		g := NewGame("A", "B", 1)
		g.StartTurn(0)
		g.AddToBattleline(def, 0)
		if err := g.ChooseHouse(0, Logos); err != nil {
			t.Fatalf("ChooseHouse: %v", err)
		}
		if g.Aember(0) != 0 {
			t.Errorf("aember = %d, want 0", g.Aember(0))
		}
	})
}

func TestAfterChooseHouseRendering(t *testing.T) {
	a := Ability{Trigger: TriggerAfterChooseHouse, Effect: Conditional{
		Cond: ChoseHouse{House: Sanctum},
		Then: GainAember{Player: Controller, Amount: 2},
	}}
	if got := RenderAbility(
		a,
	); got != "After you choose Sanctum as your active house, gain 2 Æmber." {
		t.Errorf("render = %q", got)
	}

	// Effect shapes that are not a Conditional{ChoseHouse} do not fold.
	if _, ok := afterChooseHouseText(GainAember{Player: Controller, Amount: 1}); ok {
		t.Error("a non-conditional effect should not fold")
	}
	if _, ok := afterChooseHouseText(
		Conditional{Cond: ControlsMoreCreatures{}, Then: GainAember{Player: Controller, Amount: 1}},
	); ok {
		t.Error("a conditional without ChoseHouse should not fold")
	}
}

func TestAfterAnyPlayerChooseHouseRendering(t *testing.T) {
	a := Ability{Trigger: TriggerAfterAnyPlayerChoosesHouse, Effect: Conditional{
		Cond: ChoseHouse{House: Brobnar},
		Then: GainAember{Player: Controller, Amount: 1},
	}}
	if got := RenderAbility(
		a,
	); got != "After a player chooses Brobnar as their active house, gain 1 Æmber." {
		t.Errorf("render = %q", got)
	}

	// Effect shapes that are not a Conditional{ChoseHouse} do not fold.
	if _, ok := afterAnyPlayerChooseHouseText(GainAember{Player: Controller, Amount: 1}); ok {
		t.Error("a non-conditional effect should not fold")
	}
	if _, ok := afterAnyPlayerChooseHouseText(
		Conditional{Cond: ControlsMoreCreatures{}, Then: GainAember{Player: Controller, Amount: 1}},
	); ok {
		t.Error("a conditional without ChoseHouse should not fold")
	}
}

func TestPoolAember(t *testing.T) {
	// Player must name a pool, and a comparison must be set.
	if (PoolAember{}).validate() == nil {
		t.Error("an unset player should be invalid")
	}
	if (PoolAember{Player: Opponent}).validate() == nil {
		t.Error("an unset comparison should be invalid")
	}
	// The threshold comparisons validate for either side.
	for _, p := range []Player{Opponent, Controller} {
		for _, is := range []Comparison{AtLeast, AtMost, Exactly} {
			if (PoolAember{Player: p, Is: is}).validate() != nil {
				t.Errorf("player %d comparison %d should validate", p, is)
			}
		}
	}
	// Each relative comparison is tied to the side it reads from.
	if (PoolAember{Player: Opponent, Is: MoreThanYou}).validate() != nil {
		t.Error("MoreThanYou with Opponent should validate")
	}
	if (PoolAember{Player: Controller, Is: MoreThanYou}).validate() == nil {
		t.Error("MoreThanYou with Controller should be invalid")
	}
	if (PoolAember{Player: Controller, Is: MoreThanOpponent}).validate() != nil {
		t.Error("MoreThanOpponent with Controller should validate")
	}
	if (PoolAember{Player: Opponent, Is: MoreThanOpponent}).validate() == nil {
		t.Error("MoreThanOpponent with Opponent should be invalid")
	}

	// Opponent-pool wording.
	if got := (PoolAember{Player: Opponent, Is: AtLeast, Amount: 7}).CondText(); got != "if your opponent has 7 Æmber or more" {
		t.Errorf("opp at-least text = %q", got)
	}
	if got := (PoolAember{Player: Opponent, Is: AtMost, Amount: 3}).CondText(); got != "if your opponent has 3 Æmber or fewer" {
		t.Errorf("opp at-most text = %q", got)
	}
	if got := (PoolAember{Player: Opponent, Is: Exactly, Amount: 1}).CondText(); got != "if your opponent has exactly 1 Æmber" {
		t.Errorf("opp exact text = %q", got)
	}
	if got := (PoolAember{Player: Opponent, Is: Exactly, Amount: 0}).CondText(); got != "if your opponent has no Æmber" {
		t.Errorf("opp zero text = %q", got)
	}
	if got := (PoolAember{Player: Opponent, Is: MoreThanYou}).CondText(); got != "if your opponent has more Æmber than you" {
		t.Errorf("more-than-you text = %q", got)
	}
	// Controller-pool wording.
	if got := (PoolAember{Player: Controller, Is: AtLeast, Amount: 3}).CondText(); got != "if you have 3 Æmber or more" {
		t.Errorf("you at-least text = %q", got)
	}
	if got := (PoolAember{Player: Controller, Is: AtMost, Amount: 2}).CondText(); got != "if you have 2 Æmber or fewer" {
		t.Errorf("you at-most text = %q", got)
	}
	if got := (PoolAember{Player: Controller, Is: Exactly, Amount: 1}).CondText(); got != "if you have exactly 1 Æmber" {
		t.Errorf("you exact text = %q", got)
	}
	if got := (PoolAember{Player: Controller, Is: Exactly, Amount: 0}).CondText(); got != "if you have no Æmber" {
		t.Errorf("you zero text = %q", got)
	}
	if got := (PoolAember{Player: Controller, Is: MoreThanOpponent}).CondText(); got != "if you have more Æmber than your opponent" {
		t.Errorf("more-than-opponent text = %q", got)
	}

	// Met: opponent leads 3 to 2.
	g := NewGame("A", "B", 1)
	g.State.Aember[0] = 2 // controller
	g.State.Aember[1] = 3 // opponent
	ctx := &EffectContext{Resolver: g, Controller: 0}
	if !(PoolAember{Player: Opponent, Is: AtLeast, Amount: 3}).Met(ctx) ||
		(PoolAember{Player: Opponent, Is: AtLeast, Amount: 4}).Met(ctx) {
		t.Error("opponent AtLeast wrong")
	}
	if !(PoolAember{Player: Opponent, Is: AtMost, Amount: 3}).Met(ctx) ||
		(PoolAember{Player: Opponent, Is: AtMost, Amount: 2}).Met(ctx) {
		t.Error("opponent AtMost wrong")
	}
	if !(PoolAember{Player: Opponent, Is: Exactly, Amount: 3}).Met(ctx) ||
		(PoolAember{Player: Opponent, Is: Exactly, Amount: 2}).Met(ctx) {
		t.Error("opponent Exactly wrong")
	}
	if !(PoolAember{Player: Controller, Is: AtLeast, Amount: 2}).Met(ctx) ||
		(PoolAember{Player: Controller, Is: AtMost, Amount: 1}).Met(ctx) {
		t.Error("controller threshold wrong")
	}
	if !(PoolAember{Player: Opponent, Is: MoreThanYou}).Met(ctx) {
		t.Error("MoreThanYou should hold when opp 3 > you 2")
	}
	if (PoolAember{Player: Controller, Is: MoreThanOpponent}).Met(ctx) {
		t.Error("MoreThanOpponent should not hold when you 2 < opp 3")
	}
	// Level the pools: neither relative comparison holds.
	g.State.Aember[0] = 3
	if (PoolAember{Player: Opponent, Is: MoreThanYou}).Met(ctx) {
		t.Error("MoreThanYou should not hold when 3 == 3")
	}
	if (PoolAember{Player: Controller, Is: MoreThanOpponent}).Met(ctx) {
		t.Error("MoreThanOpponent should not hold when 3 == 3")
	}

	// Parity comparisons ignore Amount and apply to either player.
	if (PoolAember{Player: Opponent, Is: Even}).validate() != nil ||
		(PoolAember{Player: Controller, Is: Odd}).validate() != nil {
		t.Error("Even/Odd should validate for either player")
	}
	if got := (PoolAember{Player: Opponent, Is: Even}).CondText(); got != "if your opponent has an even amount of Æmber" {
		t.Errorf("opp even text = %q", got)
	}
	if got := (PoolAember{Player: Opponent, Is: Odd}).CondText(); got != "if your opponent has an odd amount of Æmber" {
		t.Errorf("opp odd text = %q", got)
	}
	if got := (PoolAember{Player: Controller, Is: Even}).CondText(); got != "if you have an even amount of Æmber" {
		t.Errorf("you even text = %q", got)
	}
	if got := (PoolAember{Player: Controller, Is: Odd}).CondText(); got != "if you have an odd amount of Æmber" {
		t.Errorf("you odd text = %q", got)
	}
	// Both pools are 3 (odd) here.
	if !(PoolAember{Player: Opponent, Is: Odd}).Met(ctx) ||
		(PoolAember{Player: Opponent, Is: Even}).Met(ctx) {
		t.Error("odd pool parity wrong")
	}
	g.State.Aember[1] = 2 // opponent even
	if !(PoolAember{Player: Opponent, Is: Even}).Met(ctx) ||
		(PoolAember{Player: Opponent, Is: Odd}).Met(ctx) {
		t.Error("even pool parity wrong")
	}

	// Conditional surfaces an unset condition at validation time.
	if (Conditional{Cond: PoolAember{}, Then: GainAember{Player: Controller, Amount: 1}}).validate() == nil {
		t.Error("Conditional should surface an invalid condition")
	}
}

func TestOrCondition(t *testing.T) {
	g := NewGame("A", "B", 1)
	dino := g.AddToBattleline(
		NewCard("dino", Brobnar, Creature, Common, WithPower(3), WithTraits(Dinosaur)), 1)
	rich := g.AddToBattleline(testCreature("rich", 3), 1)
	plain := g.AddToBattleline(testCreature("plain", 3), 1)
	g.State.Cards[rich].Amber = 1
	ctx := &EffectContext{Resolver: g, Controller: 0}

	trait := ItIsOfTrait{Trait: Dinosaur}
	aember := ItHasAember{}
	or := Or{Conditions: []Condition{trait, aember}}

	if got := trait.CondText(); got != "if it is a Dinosaur creature" {
		t.Errorf("trait text = %q", got)
	}
	if got := aember.CondText(); got != "if it has Æmber on it" {
		t.Errorf("aember text = %q", got)
	}
	if got := or.CondText(); got != "if it is a Dinosaur creature or it has Æmber on it" {
		t.Errorf("or text = %q", got)
	}

	// No creature in context: nothing is met.
	if or.Met(ctx) {
		t.Error("Or should not be met with no context creature")
	}
	ctx.It, ctx.HasIt = dino, true
	if !trait.Met(ctx) || !or.Met(ctx) {
		t.Error("a Dinosaur should meet the trait condition and the Or")
	}
	ctx.It = rich
	if trait.Met(ctx) || !aember.Met(ctx) || !or.Met(ctx) {
		t.Error("a creature with Æmber should meet the Æmber condition and the Or")
	}
	ctx.It = plain
	if or.Met(ctx) {
		t.Error("an ordinary creature should meet neither condition")
	}

	// Validation: an Or needs at least two conditions and rejects an invalid one.
	if (Or{Conditions: []Condition{trait}}).validate() == nil {
		t.Error("an Or with fewer than two conditions should be invalid")
	}
	if (Or{Conditions: []Condition{trait, PoolAember{}}}).validate() == nil {
		t.Error("an Or with an invalid sub-condition should be invalid")
	}
	if err := or.validate(); err != nil {
		t.Errorf("a well-formed Or should validate: %v", err)
	}
}

// TestOrCombinesNamedHouses covers the text folding: an Or of ItIs clauses that
// differ only in a single named house renders as one phrase (Ambassador Liu).
func TestOrCombinesNamedHouses(t *testing.T) {
	cards := Or{Conditions: []Condition{
		ItIs{House: HouseMatcher{Kind: MatchNamedHouse, House: Dis}},
		ItIs{House: HouseMatcher{Kind: MatchNamedHouse, House: Shadows}},
	}}
	if got := cards.CondText(); got != "if it is a Dis or Shadows card" {
		t.Errorf("combined house text = %q", got)
	}

	// Sharing a type folds the type into the combined noun.
	creatures := Or{Conditions: []Condition{
		ItIs{House: HouseMatcher{Kind: MatchNamedHouse, House: Untamed}, Type: Creature},
		ItIs{House: HouseMatcher{Kind: MatchNamedHouse, House: Logos}, Type: Creature},
	}}
	if got := creatures.CondText(); got != "if it is an Untamed or Logos creature" {
		t.Errorf("combined creature text = %q", got)
	}

	// A differing type breaks the fold, so each clause renders on its own.
	mixed := Or{Conditions: []Condition{
		ItIs{House: HouseMatcher{Kind: MatchNamedHouse, House: Dis}, Type: Creature},
		ItIs{House: HouseMatcher{Kind: MatchNamedHouse, House: Shadows}, Type: Artifact},
	}}
	if got := mixed.CondText(); got != "if it is a Dis creature or it is a Shadows artifact" {
		t.Errorf("mixed text = %q", got)
	}

	// A non-ItIs clause breaks the fold too.
	nonHouse := Or{Conditions: []Condition{
		ItIs{House: HouseMatcher{Kind: MatchNamedHouse, House: Dis}},
		ItHasAember{},
	}}
	if got := nonHouse.CondText(); got != "if it is a Dis card or it has Æmber on it" {
		t.Errorf("non-house text = %q", got)
	}

	// An ItIs that filters by something other than a single named house
	// (here the active house) breaks the fold, so each clause renders on its own.
	notNamed := Or{Conditions: []Condition{
		ItIs{House: HouseMatcher{Kind: MatchNamedHouse, House: Dis}},
		ItIs{House: activeHouse},
	}}
	if got := notNamed.CondText(); got != "if it is a Dis card or it is of the active house" {
		t.Errorf("not-named text = %q", got)
	}
}
func TestAndCondition(t *testing.T) {
	g := NewGame("A", "B", 1)
	mutant := NewCard("mutant", Sanctum, Creature, Common, WithPower(3), WithTraits(Mutant))
	mine := g.AddToBattleline(mutant, 0)
	theirs := g.AddToBattleline(
		NewCard("theirs", Sanctum, Creature, Common, WithPower(3), WithTraits(Mutant)), 1)
	plain := g.AddToBattleline(testCreature("plain", 3), 0)
	ctx := &EffectContext{Resolver: g, Controller: 0}

	friendly := ItIsFriendly{}
	trait := ItIsOfTrait{Trait: Mutant}
	and := And{Conditions: []Condition{friendly, trait}}

	if got := and.CondText(); got != "if it is a friendly creature and it is a Mutant creature" {
		t.Errorf("and text = %q", got)
	}

	// No creature in context: nothing is met.
	if and.Met(ctx) {
		t.Error("And should not be met with no context creature")
	}
	ctx.It, ctx.HasIt = mine, true
	if !and.Met(ctx) {
		t.Error("a friendly Mutant should meet the And")
	}
	ctx.It = theirs
	if and.Met(ctx) {
		t.Error("an enemy Mutant should not meet the And (not friendly)")
	}
	ctx.It = plain
	if and.Met(ctx) {
		t.Error("a friendly non-Mutant should not meet the And (no trait)")
	}

	// Validation: an And needs at least two conditions and rejects an invalid one.
	if (And{Conditions: []Condition{trait}}).validate() == nil {
		t.Error("an And with fewer than two conditions should be invalid")
	}
	if (And{Conditions: []Condition{trait, PoolAember{}}}).validate() == nil {
		t.Error("an And with an invalid sub-condition should be invalid")
	}
	if err := and.validate(); err != nil {
		t.Errorf("a well-formed And should validate: %v", err)
	}
}

// TestItHasBonusIcon covers the text and the met check: the condition holds only
// for a card in context that prints at least one bonus icon.
func TestItHasBonusIcon(t *testing.T) {
	g := NewGame("A", "B", 1)
	iconed := g.AddToBattleline(
		NewCard("iconed", Logos, Creature, Common, WithPower(3), WithBonus(BonusDamage)), 0)
	plain := g.AddToBattleline(testCreature("plain", 3), 0)
	cond := ItHasBonusIcon{}
	if got := cond.CondText(); got != "if it has a bonus icon" {
		t.Errorf("text = %q", got)
	}
	ctx := &EffectContext{Resolver: g, Controller: 0}
	if cond.Met(ctx) {
		t.Error("with no card in context the condition should not be met")
	}
	ctx.It, ctx.HasIt = iconed, true
	if !cond.Met(ctx) {
		t.Error("a card with a bonus icon should meet the condition")
	}
	ctx.It = plain
	if cond.Met(ctx) {
		t.Error("a card without a bonus icon should not meet the condition")
	}
}

func TestCountIs(t *testing.T) {
	cases := []struct {
		name string
		cond CountIs
		want string
	}{
		{
			name: "at least, plural",
			cond: CountIs{Count: CreaturesUsed{Player: Controller}, Is: AtLeast, Amount: 3},
			want: "if you used 3 or more creatures this turn",
		},
		{
			name: "exactly one, singular",
			cond: CountIs{Count: CardsPlayed{Player: Controller}, Is: Exactly, Amount: 1},
			want: "if you played exactly 1 card this turn",
		},
		{
			name: "house-filtered",
			cond: CountIs{
				Count:  CardsPlayed{Player: Controller, House: namedHouse(Mars)},
				Is:     Exactly,
				Amount: 2,
			},
			want: "if you played exactly 2 Mars cards this turn",
		},
		{
			name: "mass noun ignores plurality",
			cond: CountIs{Count: DamageHealed{}, Is: Exactly, Amount: 3},
			want: "if you healed exactly 3 damage",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.cond.CondText(); got != tc.want {
				t.Errorf("text = %q, want %q", got, tc.want)
			}
			if err := tc.cond.validate(); err != nil {
				t.Errorf("valid condition rejected: %v", err)
			}
		})
	}

	if err := (CountIs{Is: AtLeast, Amount: 1}).validate(); err == nil {
		t.Error("a missing Count should be rejected")
	}
	if err := (CountIs{Count: CardsDestroyed{}, Is: AtLeast}).validate(); err == nil {
		t.Error("a Count with no clause should be rejected")
	}
	if err := (CountIs{Count: DamageHealed{}, Is: MoreThanYou}).validate(); err == nil {
		t.Error("MoreThanYou should be rejected: it compares two Æmber pools")
	}
}

func TestCountIsMet(t *testing.T) {
	g := NewGame("A", "B", 1)
	src := g.AddToBattleline(testCreature("src", 3), 0)
	used := g.AddToBattleline(testCreature("used", 3), 0)
	g.AddToBattleline(testCreature("idle", 3), 0)
	ctx := &EffectContext{Resolver: g, Source: src, Controller: 0}

	atLeastOne := CountIs{Count: CreaturesUsed{Player: Controller}, Is: AtLeast, Amount: 1}
	if atLeastOne.Met(ctx) {
		t.Error("no creature has been used yet")
	}
	g.ReapWith(used)
	if !atLeastOne.Met(ctx) {
		t.Error("the reaped creature should count as used")
	}

	exactlyOne := CountIs{Count: CreaturesUsed{Player: Controller}, Is: Exactly, Amount: 1}
	if !exactlyOne.Met(ctx) {
		t.Error("exactly one creature has been used")
	}
	if got := (CreaturesUsed{Player: Controller}).CountText(); got != "creature you used this turn" {
		t.Errorf("count text = %q", got)
	}
}

func TestCardsPlayedCountsEveryHouseWhenUnset(t *testing.T) {
	g := NewGame("A", "B", 1)
	src := g.AddToBattleline(testCreature("src", 3), 0)
	ctx := &EffectContext{Resolver: g, Source: src, Controller: 0}
	g.State.PlayedThisTurn[0].add(g.AddToHand(testCreature("brobnar", 3), 0))
	g.State.PlayedThisTurn[0].add(
		g.AddToHand(NewCard("logos", Logos, Creature, Common, WithPower(3)), 0),
	)

	if got := (CardsPlayed{Player: Controller}).Value(ctx); got != 2 {
		t.Errorf("unfiltered value = %d, want 2 (every house counts)", got)
	}
	if got := (CardsPlayed{Player: Controller, House: namedHouse(Logos)}).Value(ctx); got != 1 {
		t.Errorf("house-filtered value = %d, want 1", got)
	}
}

func TestFirstCreaturePlayedThisTurn(t *testing.T) {
	cond := FirstCreaturePlayedThisTurn{}
	if got := cond.CondText(); got != "if it is the first creature played this turn" {
		t.Errorf("text = %q", got)
	}

	g := started(t)
	ctx := &EffectContext{Resolver: g, Controller: 0}
	if cond.Met(ctx) {
		t.Error("with no card in context the condition should not be met")
	}

	first := g.AddToHand(NewCard("First", Brobnar, Creature, Common, WithPower(2)), 0)
	second := g.AddToHand(NewCard("Second", Brobnar, Creature, Common, WithPower(2)), 0)
	tactic := g.AddToHand(NewCard("Warm Up", Brobnar, Tactic, Common), 0)

	if err := g.PlayAction(0, handIdxByID(g, 0, tactic)); err != nil {
		t.Fatalf("play tactic: %v", err)
	}
	if cond.Met(&EffectContext{Resolver: g, Controller: 0, It: first, HasIt: true}) {
		t.Error("a creature that was not played should not meet the condition")
	}

	if _, err := g.PlayCreature(0, handIdxByID(g, 0, first), false); err != nil {
		t.Fatalf("play first: %v", err)
	}
	if !cond.Met(&EffectContext{Resolver: g, Controller: 0, It: first, HasIt: true}) {
		t.Error("the first creature played should meet the condition")
	}
	if cond.Met(&EffectContext{Resolver: g, Controller: 0, It: second, HasIt: true}) {
		t.Error("a later creature should not meet the condition")
	}
}

func TestNoCreaturesPlayedThisTurn(t *testing.T) {
	cond := NoCreaturesPlayedThisTurn{}
	if got := cond.CondText(); got != "if you did not play any creatures this turn" {
		t.Errorf("text = %q", got)
	}

	g := started(t)
	ctx := &EffectContext{Resolver: g, Controller: 0}
	if !cond.Met(ctx) {
		t.Error("having played nothing, the condition should be met")
	}

	tactic := g.AddToHand(NewCard("Warm Up", Brobnar, Tactic, Common), 0)
	if err := g.PlayAction(0, handIdxByID(g, 0, tactic)); err != nil {
		t.Fatalf("play tactic: %v", err)
	}
	if !cond.Met(ctx) {
		t.Error("playing only a tactic should keep the condition met")
	}

	creature := g.AddToHand(NewCard("Beef", Brobnar, Creature, Common, WithPower(2)), 0)
	if _, err := g.PlayCreature(0, handIdxByID(g, 0, creature), false); err != nil {
		t.Fatalf("play creature: %v", err)
	}
	if cond.Met(ctx) {
		t.Error("after playing a creature the condition should not be met")
	}
}

// TestItIsYourTurn is met when the controller is the active player and not
// otherwise.
func TestItIsYourTurn(t *testing.T) {
	g := started(t)
	g.State.ActivePlayer = 0
	if got := (ItIsYourTurn{}).CondText(); got != "if it is your turn" {
		t.Errorf("text = %q, want %q", got, "if it is your turn")
	}
	if !(ItIsYourTurn{}).Met(&EffectContext{Resolver: g, Controller: 0}) {
		t.Error("should be met for the active player")
	}
	if (ItIsYourTurn{}).Met(&EffectContext{Resolver: g, Controller: 1}) {
		t.Error("should not be met for the inactive player")
	}
}

func TestAemberOnThisAtLeast(t *testing.T) {
	c := CountIs{Count: AemberOnThis{}, Is: AtLeast, Amount: 4}
	if got := c.CondText(); got != "if there are 4 or more Æmber on it" {
		t.Errorf("text = %q", got)
	}
	g := NewGame("A", "B", 1)
	vault := NewCard("Safe Place", Shadows, Artifact, Rare)
	id := g.AddArtifact(vault, 0)
	ctx := &EffectContext{Resolver: g, Source: id, Controller: 0}
	if c.Met(ctx) {
		t.Error("should not be met below the threshold")
	}
	g.AddAmberOn(id, 4)
	if !c.Met(ctx) {
		t.Error("should be met at the threshold")
	}

	// Not flips the sense to "fewer than", met below the threshold and not above.
	fewer := Not{Cond: CountIs{Count: AemberOnThis{}, Is: AtLeast, Amount: 10}}
	if got := fewer.CondText(); got != "if there are fewer than 10 Æmber on it" {
		t.Errorf("negated text = %q", got)
	}
	if !fewer.Met(ctx) {
		t.Error("fewer-than should be met below the threshold")
	}
	g.AddAmberOn(id, 6) // now 10 on the card
	if fewer.Met(ctx) {
		t.Error("fewer-than should not be met at the threshold")
	}
}

// TestNamedCardPurged covers Igon the Terrible's gate: whether a card of a name
// sits in the controller's purge pile, both senses.
func TestNamedCardPurged(t *testing.T) {
	present := NamedCardPurged{Name: "Igon the Green"}
	absent := Not{Cond: NamedCardPurged{Name: "Igon the Green"}}
	if got := present.CondText(); got != "if Igon the Green has been purged" {
		t.Errorf("present text = %q", got)
	}
	if got := absent.CondText(); got != "if Igon the Green has not been purged" {
		t.Errorf("absent text = %q", got)
	}

	g := NewGame("A", "B", 1)
	src := g.AddToBattleline(testCreature("Igon the Terrible", 8), 0)
	ctx := &EffectContext{Resolver: g, Source: src, Controller: 0}

	if present.Met(ctx) {
		t.Error("nothing purged yet, present should not be met")
	}
	if !absent.Met(ctx) {
		t.Error("nothing purged yet, absent should be met")
	}

	green := g.Register(NewCard("Igon the Green", Brobnar, Creature, Rare, WithPower(4)), 0)
	g.State.Purge[0].add(green)

	if !present.Met(ctx) {
		t.Error("the green is purged, present should be met")
	}
	if absent.Met(ctx) {
		t.Error("the green is purged, absent should not be met")
	}
}

// TestFirstReapOfTurn covers Aember Conduction Unit's gate: met only while the
// reap in context is the first creature to reap this turn.
func TestFirstReapOfTurn(t *testing.T) {
	if got := (FirstReapOfTurn{}).CondText(); got != "if it is the first time a creature has reaped this turn" {
		t.Errorf("CondText = %q", got)
	}
	if (FirstReapOfTurn{}).Met(&EffectContext{Resolver: started(t)}) {
		t.Error("with no creature in context the condition should not be met")
	}

	g := started(t)
	first := g.AddToBattleline(NewCard("first", Brobnar, Creature, Common, WithPower(3)), 0)
	second := g.AddToBattleline(NewCard("second", Brobnar, Creature, Common, WithPower(3)), 0)

	g.reapWith(first)
	if !(FirstReapOfTurn{}).Met(&EffectContext{Resolver: g, It: first, HasIt: true}) {
		t.Error("the first reap of the turn should meet the condition")
	}

	g.reapWith(second)
	if (FirstReapOfTurn{}).Met(&EffectContext{Resolver: g, It: second, HasIt: true}) {
		t.Error("a later reap should not meet the condition")
	}
}

// TestSourceFirstUseThisTurn covers Gladiodontus's gate: met only while the
// source creature's current use is its first this turn.
func TestSourceFirstUseThisTurn(t *testing.T) {
	want := "if this is the first time " + SelfName + " has been used this turn"
	if got := (SourceFirstUseThisTurn{}).CondText(); got != want {
		t.Errorf("CondText = %q", got)
	}

	g := started(t)
	c := g.AddToBattleline(NewCard("glad", Saurian, Creature, Common, WithPower(5)), 0)
	if (SourceFirstUseThisTurn{}).Met(&EffectContext{Resolver: g, Source: c}) {
		t.Error("an unused creature should not meet the condition")
	}

	g.reapWith(c)
	if !(SourceFirstUseThisTurn{}).Met(&EffectContext{Resolver: g, Source: c}) {
		t.Error("the first use of the turn should meet the condition")
	}

	g.SetExhausted(c, false)
	g.reapWith(c)
	if (SourceFirstUseThisTurn{}).Met(&EffectContext{Resolver: g, Source: c}) {
		t.Error("a second use should not meet the condition")
	}
}

// TestCardsInDeckAtMost covers the deck-size threshold Manchego reads before it
// steals.
func TestCardsInDeckAtMost(t *testing.T) {
	e := CountIs{Count: CardsInZone{Zone: Deck, Player: Controller}, Is: AtMost, Amount: 5}
	if got := e.CondText(); got != "if you have 5 or fewer cards in your deck" {
		t.Errorf("CondText = %q", got)
	}

	g := started(t)
	if !e.Met(&EffectContext{Resolver: g, Controller: 0}) {
		t.Error("a small deck should meet the condition")
	}

	for i := 0; i < 6; i++ {
		g.State.Deck[0].add(
			g.Register(NewCard("filler", Brobnar, Creature, Common, WithPower(3)), 0),
		)
	}
	if e.Met(&EffectContext{Resolver: g, Controller: 0}) {
		t.Error("a large deck should not meet the condition")
	}
}
