package engine

import "testing"

func TestReturnNamedToHand(t *testing.T) {
	e := ReturnNamedToHand{Name: "Urchin"}
	if e.Text() != "put an Urchin from play or from your discard pile into your hand" {
		t.Errorf("text = %q", e.Text())
	}

	urchin := func(g *Game, player int) LocalID {
		return g.Register(NewCard("Urchin", Shadows, Creature, Common, WithPower(1), WithTraits("Elf", "Thief")), player)
	}

	t.Run("from play", func(t *testing.T) {
		g := NewGame("A", "B", 1)
		src := g.AddToBattleline(testCreature("faygin", 3), 0)
		urch := g.AddToBattleline(NewCard("Urchin", Shadows, Creature, Common, WithPower(1), WithTraits("Elf", "Thief")), 0)
		g.AddToBattleline(testCreature("other", 4), 0)                              // different name: filtered out
		g.State.Discard[0].add(g.Register(NewCard("junk", Dis, Tactic, Common), 0)) // different name in discard: filtered out
		ctx := &EffectContext{Resolver: g, Source: src, Controller: 0}

		e.Resolve(ctx) // the sole Urchin candidate is auto-chosen
		if g.inPlay(urch) {
			t.Error("the Urchin should leave play")
		}
		if !g.State.Hand[0].contains(urch) {
			t.Error("the Urchin should be in the controller's hand")
		}
	})

	t.Run("from discard", func(t *testing.T) {
		g := NewGame("A", "B", 1)
		src := g.AddToBattleline(testCreature("faygin", 3), 0)
		urch := urchin(g, 0)
		g.State.Discard[0].add(urch)
		ctx := &EffectContext{Resolver: g, Source: src, Controller: 0}

		e.Resolve(ctx)
		if g.State.Discard[0].contains(urch) {
			t.Error("the Urchin should leave the discard pile")
		}
		if !g.State.Hand[0].contains(urch) {
			t.Error("the Urchin should be in the controller's hand")
		}
	})

	t.Run("declined", func(t *testing.T) {
		g := NewGame("A", "B", 1)
		src := g.AddToBattleline(testCreature("faygin", 3), 0)
		u1 := g.AddToBattleline(NewCard("Urchin", Shadows, Creature, Common, WithPower(1), WithTraits("Elf", "Thief")), 0)
		g.AddToBattleline(NewCard("Urchin", Shadows, Creature, Common, WithPower(1), WithTraits("Elf", "Thief")), 0)
		g.SetChooser(0, orderRejectChooser{})
		ctx := &EffectContext{Resolver: g, Source: src, Controller: 0}

		e.Resolve(ctx)
		if !g.inPlay(u1) {
			t.Error("nothing should move when the choice is declined")
		}
	})
}

func TestMoveFromPlayToDeck(t *testing.T) {
	g := NewGame("A", "B", 1)
	src := g.AddToBattleline(testCreature("src", 1), 0)
	myArt := g.AddArtifact(NewCard("myrelic", Brobnar, Artifact, Rare), 0)
	enemyArt := g.AddArtifact(NewCard("enemyrelic", Brobnar, Artifact, Rare), 1)
	g.State.Cards[myArt].Exhausted = true
	ctx := &EffectContext{Resolver: g, Source: src, Controller: 0}

	e := PutFromPlay{Target: Target{Kind: TargetEachArtifact}, Destination: ToTopOfDeck}
	if e.Text() != "put each artifact on top of its owner's deck" {
		t.Errorf("text = %q", e.Text())
	}
	e.Resolve(ctx)

	if g.State.Artifacts[0].Count != 0 || g.State.Artifacts[1].Count != 0 {
		t.Errorf("artifact rows not cleared: %d %d", g.State.Artifacts[0].Count, g.State.Artifacts[1].Count)
	}
	if g.State.Deck[0].Count != 1 || g.State.Deck[0].IDs[0] != myArt {
		t.Errorf("player 0 deck top = %v, want %d", g.State.Deck[0].IDs[0], myArt)
	}
	if g.State.Deck[1].Count != 1 || g.State.Deck[1].IDs[0] != enemyArt {
		t.Errorf("player 1 deck top = %v, want %d", g.State.Deck[1].IDs[0], enemyArt)
	}
	if g.State.Cards[myArt].Exhausted {
		t.Errorf("returned artifact should be readied")
	}
}

func TestMoveFromPlayToHand(t *testing.T) {
	g := NewGame("A", "B", 1)
	src := g.AddToBattleline(testCreature("src", 3), 0)
	g.State.Cards[src].Damage = 2
	ctx := &EffectContext{Resolver: g, Source: src, Controller: 0}

	e := PutFromPlay{Target: Target{Kind: TargetThisCreature}, Destination: ToHand}
	if e.Text() != "put "+SelfName+" into its owner's hand" {
		t.Errorf("text = %q", e.Text())
	}
	e.Resolve(ctx)

	if g.State.Battleline[0].Count != 0 {
		t.Errorf("battleline not cleared: %d", g.State.Battleline[0].Count)
	}
	if g.State.Hand[0].Count != 1 || g.State.Hand[0].IDs[0] != src {
		t.Errorf("hand = count %d id %v, want 1 / %d", g.State.Hand[0].Count, g.State.Hand[0].IDs[0], src)
	}
	// Moving to hand clears the per-match state the card accrued in play.
	if g.State.Cards[src].Damage != 0 {
		t.Errorf("damage after move = %d, want 0", g.State.Cards[src].Damage)
	}
}

func TestMoveFromPlayToArchives(t *testing.T) {
	g := NewGame("A", "B", 1)
	src := g.AddToBattleline(testCreature("src", 3), 0)
	attachUpgrade(g, src, NewCard("plating", Mars, Upgrade, Common))
	g.State.Cards[src].Damage = 1
	ctx := &EffectContext{Resolver: g, Source: src, Controller: 0}

	e := PutFromPlay{Target: Target{Kind: TargetThisCreature}, Destination: ToArchives}
	if e.Text() != "put "+SelfName+" into its owner's archives" {
		t.Errorf("text = %q", e.Text())
	}
	e.Resolve(ctx)

	if g.inPlay(src) {
		t.Error("the creature should have left play")
	}
	if g.State.Archives[0].Count != 1 || g.State.Archives[0].IDs[0] != src {
		t.Errorf("creature should be archived, got %v", g.State.Archives[0].IDs[:g.State.Archives[0].Count])
	}
	if len(g.Discard(0)) != 1 { // the upgrade sheds to the discard pile
		t.Errorf("upgrade should be discarded; discard = %v", g.Discard(0))
	}
	if g.State.Cards[src].Damage != 0 {
		t.Error("archiving should clear the creature's in-play state")
	}
}

func TestMoveFromPlayToDeckShuffled(t *testing.T) {
	g := NewGame("A", "B", 1)
	src := g.AddToBattleline(testCreature("chrono", 2), 0)
	ctx := &EffectContext{Resolver: g, Source: src, Controller: 0}

	e := PutFromPlay{Target: Target{Kind: TargetThisCreature}, Destination: ToDeckShuffled}
	if e.Text() != "shuffle {self} into its owner's deck" {
		t.Errorf("text = %q", e.Text())
	}
	e.Resolve(ctx)
	if g.inPlay(src) {
		t.Error("the creature should leave play")
	}
	if g.State.Deck[0].Count != 1 || !g.State.Deck[0].contains(src) {
		t.Errorf("the creature should be shuffled into its owner's deck")
	}
}

func TestMoveFromPlayValidate(t *testing.T) {
	this := Target{Kind: TargetThisCreature}
	for _, d := range []Destination{ToHand, ToTopOfDeck, ToDeckShuffled, ToArchives} {
		if err := (PutFromPlay{Target: this, Destination: d}).validate(); err != nil {
			t.Errorf("destination %d should be valid, got %v", d, err)
		}
	}
	if err := (PutFromPlay{Target: this}).validate(); err == nil {
		t.Error("an unset destination should be rejected")
	}
	if err := (PutFromPlay{Target: this, Destination: ToBottomOfDeck}).validate(); err == nil {
		t.Error("an unsupported destination should be rejected")
	}
	if err := (PutFromPlay{Destination: ToHand}).validate(); err == nil {
		t.Error("an unset target should be rejected")
	}
}

func TestPutUpTo(t *testing.T) {
	eachArt := Target{Kind: TargetEachArtifact}
	// Text renders "up to N" with the plural noun for each destination.
	cases := map[Destination]string{
		ToHand:         "put up to 3 artifacts into their owners' hands",
		ToTopOfDeck:    "put up to 3 artifacts on top of their owners' decks",
		ToDeckShuffled: "shuffle up to 3 artifacts into their owners' decks",
		ToArchives:     "put up to 3 artifacts into their owners' archives",
	}
	for dest, want := range cases {
		if got := (PutUpTo{Max: 3, Target: eachArt, Destination: dest}).Text(); got != want {
			t.Errorf("text(%d) = %q, want %q", dest, got, want)
		}
	}

	// validate rejects an unset target, a non-positive Max, and a bad destination.
	if err := (PutUpTo{Max: 3, Destination: ToHand}).validate(); err == nil {
		t.Error("unset target should be rejected")
	}
	if err := (PutUpTo{Target: eachArt, Destination: ToHand}).validate(); err == nil {
		t.Error("non-positive Max should be rejected")
	}
	if err := (PutUpTo{Max: 1, Target: eachArt, Destination: ToBottomOfDeck}).validate(); err == nil {
		t.Error("unsupported destination should be rejected")
	}
	if err := (PutUpTo{Max: 3, Target: eachArt, Destination: ToHand}).validate(); err != nil {
		t.Errorf("valid PutUpTo = %v", err)
	}

	g := NewGame("A", "B", 1)
	a1 := g.AddArtifact(exAutocannon(), 0)
	a2 := g.AddArtifact(exAutocannon(), 1)
	ctx := &EffectContext{Resolver: g, Controller: 0}
	// Only two artifacts exist, so the loop stops when none remain (below Max).
	PutUpTo{Max: 3, Target: eachArt, Destination: ToHand}.Resolve(ctx)
	if g.inPlay(a1) || g.inPlay(a2) {
		t.Error("both artifacts should have left play")
	}
	if g.State.Hand[0].Count != 1 || g.State.Hand[1].Count != 1 {
		t.Errorf("hands = %d/%d, want 1/1 (each returned to its owner)", g.State.Hand[0].Count, g.State.Hand[1].Count)
	}

	// Choosing "Done" (the option past the sole artifact) stops early, leaving it.
	g2 := NewGame("A", "B", 1)
	art := g2.AddArtifact(exAutocannon(), 0)
	g2.SetChooser(0, optionPicker{idx: 1}) // index 0 is the artifact, 1 is "Done"
	PutUpTo{Max: 3, Target: eachArt, Destination: ToHand}.Resolve(&EffectContext{Resolver: g2, Controller: 0})
	if !g2.inPlay(art) {
		t.Error("choosing Done should leave the artifact in play")
	}
}
