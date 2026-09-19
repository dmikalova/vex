package engine

import "testing"

func TestPlayOrUseValidate(t *testing.T) {
	if err := (PlayOrUse{House: exceptHouse(HouseNone)}).validate(); err == nil {
		t.Error("Except with no house should not validate")
	}
	if err := (PlayOrUse{Grant: GrantFight}).validate(); err == nil {
		t.Error("a fight grant is not a play-or-use verb")
	}
	for _, e := range []PlayOrUse{
		{House: exceptHouse(StarAlliance)},
		{House: namedHouse(Mars)},
		{Grant: GrantPlay},
		{Grant: GrantUse},
		{},
	} {
		if err := e.validate(); err != nil {
			t.Errorf("validate(%+v) = %v, want nil", e, err)
		}
	}
}

func TestPlayOrUseText(t *testing.T) {
	cases := []struct {
		e    PlayOrUse
		want string
	}{
		{PlayOrUse{House: exceptHouse(StarAlliance)}, "play or use a non-Star Alliance card"},
		{PlayOrUse{House: namedHouse(Mars)}, "play or use a Mars card"},
		{PlayOrUse{}, "play or use a card"},
		{PlayOrUse{House: namedHouse(Mars), Grant: GrantPlay}, "play a Mars card"},
		{PlayOrUse{House: namedHouse(Mars), Grant: GrantUse}, "use a Mars card"},
	}
	for _, c := range cases {
		if got := c.e.Text(); got != c.want {
			t.Errorf("Text(%+v) = %q, want %q", c.e, got, c.want)
		}
	}
}

// TestPlayOrUsePlaysFromHand covers the play branch: the chosen off-house card in
// hand is played, while the same-house card is not even offered.
func TestPlayOrUsePlaysFromHand(t *testing.T) {
	g := NewGame("A", "B", 1)
	src := g.AddToBattleline(testCreature("taber", 3), 0)
	art := g.AddToHand(NewCard("art", Mars, Artifact, Common), 0)
	g.AddToHand(NewCard("sa", StarAlliance, Artifact, Common), 0)
	g.SetChooser(0, &idQueueChooser{ids: []LocalID{art}})

	PlayOrUse{House: exceptHouse(StarAlliance)}.Resolve(
		&EffectContext{Resolver: g, Controller: 0, Source: src},
	)

	if !g.inPlay(art) {
		t.Error("the off-house artifact should have been played")
	}
}

// TestPlayOrUseUsesInPlay covers the use branch: the chosen off-house creature in
// play reaps, exhausting it and paying its Æmber.
func TestPlayOrUseUsesInPlay(t *testing.T) {
	g := NewGame("A", "B", 1)
	src := g.AddToBattleline(testCreature("taber", 3), 0)
	mars := g.AddToBattleline(NewCard("mars", Mars, Creature, Common, WithPower(3)), 0)
	g.SetChooser(0, &idQueueChooser{ids: []LocalID{mars}})
	before := g.Aember(0)

	PlayOrUse{House: exceptHouse(StarAlliance)}.Resolve(
		&EffectContext{Resolver: g, Controller: 0, Source: src},
	)

	if !g.Exhausted(mars) {
		t.Error("the used creature should be exhausted after reaping")
	}
	if got := g.Aember(0); got != before+1 {
		t.Errorf("Æmber = %d, want %d after a reap", got, before+1)
	}
}

// TestPlayOrUseUsesArtifact covers the artifact branch of the in-play pool: a
// ready off-house artifact with an Action fires it when chosen.
func TestPlayOrUseUsesArtifact(t *testing.T) {
	g := NewGame("A", "B", 1)
	src := g.AddToBattleline(testCreature("taber", 3), 0)
	art := g.AddArtifact(NewCard("gun", Mars, Artifact, Common,
		WithAbility(TriggerAction, GainAember{Player: Controller, Amount: 2})), 0)
	g.SetChooser(0, &idQueueChooser{ids: []LocalID{art}})
	before := g.Aember(0)

	PlayOrUse{House: exceptHouse(StarAlliance)}.Resolve(
		&EffectContext{Resolver: g, Controller: 0, Source: src},
	)

	if got := g.Aember(0); got != before+2 {
		t.Errorf("Æmber = %d, want %d after the artifact's Action", got, before+2)
	}
	if !g.Exhausted(art) {
		t.Error("the artifact should be exhausted after firing its Action")
	}
}

// TestPlayOrUseAdmitsAnyHouse covers the unfiltered form: with no house named,
// every card qualifies.
func TestPlayOrUseAdmitsAnyHouse(t *testing.T) {
	g := NewGame("A", "B", 1)
	src := g.AddToBattleline(testCreature("t", 3), 0)
	art := g.AddToHand(NewCard("a", StarAlliance, Artifact, Common), 0)
	g.SetChooser(0, &idQueueChooser{ids: []LocalID{art}})

	PlayOrUse{}.Resolve(&EffectContext{Resolver: g, Controller: 0, Source: src})

	if !g.inPlay(art) {
		t.Error("with no house filter any card is admitted")
	}
}

// TestPlayOrUseNoCandidates covers the empty pool: nothing qualifies, so the
// controller is never prompted.
func TestPlayOrUseNoCandidates(t *testing.T) {
	g := NewGame("A", "B", 1)
	src := g.AddToBattleline(NewCard("sa", StarAlliance, Creature, Common, WithPower(3)), 0)
	c := &countingChooser{}
	g.SetChooser(0, c)

	PlayOrUse{House: exceptHouse(StarAlliance)}.Resolve(
		&EffectContext{Resolver: g, Controller: 0, Source: src},
	)

	if c.calls != 0 {
		t.Errorf("chooser called %d times, want 0 with no candidates", c.calls)
	}
}

// TestPlayOrUseDeclined covers the declined choice: nobody is played or used.
func TestPlayOrUseDeclined(t *testing.T) {
	g := NewGame("A", "B", 1)
	src := g.AddToBattleline(testCreature("taber", 3), 0)
	mars := g.AddToBattleline(NewCard("mars", Mars, Creature, Common, WithPower(3)), 0)
	g.AddToHand(NewCard("art", Mars, Artifact, Common), 0)
	g.SetChooser(0, orderRejectChooser{})

	PlayOrUse{House: exceptHouse(StarAlliance)}.Resolve(
		&EffectContext{Resolver: g, Controller: 0, Source: src},
	)

	if g.Exhausted(mars) {
		t.Error("declining should use nobody")
	}
}

// TestPlayOrUsePlayGrantSkipsInPlay covers a play-only grant: the in-play creature
// is never offered, only the hand card is played.
func TestPlayOrUsePlayGrantSkipsInPlay(t *testing.T) {
	g := NewGame("A", "B", 1)
	src := g.AddToBattleline(testCreature("taber", 3), 0)
	mars := g.AddToBattleline(NewCard("mars", Mars, Creature, Common, WithPower(3)), 0)
	art := g.AddToHand(NewCard("art", Mars, Artifact, Common), 0)
	g.SetChooser(0, &idQueueChooser{ids: []LocalID{art}})

	PlayOrUse{House: namedHouse(Mars), Grant: GrantPlay}.Resolve(
		&EffectContext{Resolver: g, Controller: 0, Source: src},
	)

	if !g.inPlay(art) {
		t.Error("the hand card should have been played")
	}
	if g.Exhausted(mars) {
		t.Error("a play grant should not offer the in-play creature")
	}
}

// TestPlayOrUseUseGrantSkipsHand covers a use-only grant: the hand card is never
// offered, only the in-play creature is used.
func TestPlayOrUseUseGrantSkipsHand(t *testing.T) {
	g := NewGame("A", "B", 1)
	src := g.AddToBattleline(testCreature("taber", 3), 0)
	mars := g.AddToBattleline(NewCard("mars", Mars, Creature, Common, WithPower(3)), 0)
	art := g.AddToHand(NewCard("art", Mars, Artifact, Common), 0)
	g.SetChooser(0, &idQueueChooser{ids: []LocalID{mars}})

	PlayOrUse{House: namedHouse(Mars), Grant: GrantUse}.Resolve(
		&EffectContext{Resolver: g, Controller: 0, Source: src},
	)

	if !g.Exhausted(mars) {
		t.Error("the in-play creature should have been used")
	}
	if !g.State.Hand[0].contains(art) {
		t.Error("a use grant should not offer the hand card")
	}
}
