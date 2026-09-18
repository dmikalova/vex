package engine

import "testing"

// TestShuffleFromDiscardText covers the three selection modes' printed text.
func TestShuffleFromDiscardText(t *testing.T) {
	each := ShuffleIntoDeck{
		Player:    Controller,
		From:      []Zone{Discard},
		Selection: Each{House: namedHouse(Untamed), Type: Creature},
	}
	if got := each.Text(); got !=
		"shuffle each Untamed creature from your discard pile into your deck" {
		t.Errorf("Each text = %q", got)
	}

	anyNum := ShuffleIntoDeck{Player: Controller, From: []Zone{Discard},
		Selection: Chosen{Type: Creature, Optional: true},
		AnyNumber: true,
	}
	if got := anyNum.Text(); got !=
		"shuffle any number of creatures from your discard pile into your deck" {
		t.Errorf("AnyNumber text = %q", got)
	}
	anyHouse := ShuffleIntoDeck{Player: Controller, From: []Zone{Discard},
		Selection: Chosen{House: namedHouse(Untamed), Type: Creature, Optional: true},
		AnyNumber: true,
	}
	if got := anyHouse.Text(); got !=
		"shuffle any number of Untamed creatures from your discard pile into your deck" {
		t.Errorf("AnyNumber house text = %q", got)
	}

	counted := ShuffleIntoDeck{Player: Controller, From: []Zone{Discard},
		Selection: Chosen{},
		Count:     CardsInPlay{Player: Controller, Trait: Shard},
	}
	if got := counted.Text(); got !=
		"for each friendly Shard, shuffle a card from your discard pile into your deck" {
		t.Errorf("Count text = %q", got)
	}
}

// TestShuffleFromDiscardValidate covers the node's validation rules.
func TestShuffleFromDiscardValidate(t *testing.T) {
	if (ShuffleIntoDeck{Player: Controller, From: []Zone{Discard}}).validate() == nil {
		t.Error("a missing Selection should fail validation")
	}
	if (ShuffleIntoDeck{Player: Controller, Selection: Chosen{}}).validate() == nil {
		t.Error("an empty From should fail validation")
	}
	if (ShuffleIntoDeck{Player: Controller, From: []Zone{Archives}, Selection: Chosen{}}).validate() == nil {
		t.Error("a zone no shuffle draws from should fail validation")
	}
	both := ShuffleIntoDeck{Player: Controller, From: []Zone{Discard},
		Selection: Chosen{},
		AnyNumber: true,
		Count:     CardsInPlay{Player: Controller, Trait: Shard},
	}
	if both.validate() == nil {
		t.Error("pairing AnyNumber with a Count should fail validation")
	}
	if err := (ShuffleIntoDeck{Player: Controller, From: []Zone{Discard}, Selection: Chosen{}}).validate(); err != nil {
		t.Errorf("a set Selection should validate, got %v", err)
	}
}

// TestShuffleFromDiscardEach covers Low Dawn: every Untamed creature leaves the
// discard for the deck as one grouped line while other cards stay put.
func TestShuffleFromDiscardEach(t *testing.T) {
	sel := Each{House: namedHouse(Untamed), Type: Creature}

	g := NewGame("A", "B", 1)
	src := g.AddToDiscard(NewCard("src", Untamed, Tactic, Common), 0)
	u1 := g.AddToDiscard(NewCard("u1", Untamed, Creature, Common, WithPower(3)), 0)
	u2 := g.AddToDiscard(NewCard("u2", Untamed, Creature, Common, WithPower(3)), 0)
	tactic := g.AddToDiscard(NewCard("ut", Untamed, Tactic, Common), 0)
	mars := g.AddToDiscard(NewCard("mc", Mars, Creature, Common, WithPower(3)), 0)

	ShuffleIntoDeck{Player: Controller, From: []Zone{Discard}, Selection: sel}.
		Resolve(&EffectContext{Resolver: g, Controller: 0, Source: src})

	if g.State.Deck[0].Count != 2 {
		t.Errorf("deck count = %d, want 2", g.State.Deck[0].Count)
	}
	if containsID(g.Discard(0), u1) || containsID(g.Discard(0), u2) {
		t.Error("Untamed creatures should have left the discard pile")
	}
	if !containsID(g.Discard(0), tactic) {
		t.Error("the Untamed tactic should stay in the discard pile")
	}
	if !containsID(g.Discard(0), mars) {
		t.Error("the Mars creature should stay in the discard pile")
	}

	// The batch narrates one grouped line.
	var grouped int
	for _, rec := range g.Log {
		if _, ok := rec.Entry.(CardsShuffledIntoDeckBy); ok {
			grouped++
		}
		if _, ok := rec.Entry.(CardShuffledIntoDeck); ok {
			t.Error("a batched shuffle should not narrate a passive per-card line")
		}
	}
	if grouped != 1 {
		t.Errorf("grouped shuffle lines = %d, want 1", grouped)
	}

	// With nothing matching, the effect shuffles nothing and narrates nothing.
	g2 := NewGame("A", "B", 1)
	only := g2.AddToDiscard(NewCard("mc", Mars, Creature, Common, WithPower(3)), 0)
	ShuffleIntoDeck{Player: Controller, From: []Zone{Discard}, Selection: sel}.
		Resolve(&EffectContext{Resolver: g2, Controller: 0, Source: only})
	if g2.State.Deck[0].Count != 0 || !containsID(g2.Discard(0), only) {
		t.Error("a non-matching discard pile should be left untouched")
	}
	for _, rec := range g2.Log {
		if _, ok := rec.Entry.(CardsShuffledIntoDeckBy); ok {
			t.Error("an empty batch should narrate nothing")
		}
	}
}

// TestShuffleFromDiscardAnyNumber covers Not Finished with You: the controller
// shuffles any number of chosen creatures out of their discard pile.
func TestShuffleFromDiscardAnyNumber(t *testing.T) {
	sel := Chosen{Type: Creature, Optional: true}

	// The default chooser takes every eligible creature, skipping a non-creature.
	g := NewGame("A", "B", 1)
	a := g.AddToDiscard(testCreature("a", 3), 0)
	b := g.AddToDiscard(testCreature("b", 3), 0)
	tactic := g.AddToDiscard(NewCard("t", Dis, Tactic, Common), 0)
	ShuffleIntoDeck{
		Player:    Controller,
		From:      []Zone{Discard},
		Selection: sel,
		AnyNumber: true,
	}.Resolve(
		&EffectContext{Resolver: g, Controller: 0, Source: a},
	)
	if g.State.Deck[0].Count != 2 {
		t.Errorf("deck count = %d, want 2", g.State.Deck[0].Count)
	}
	if !containsID(g.Discard(0), tactic) {
		t.Error("the non-creature should stay in the discard pile")
	}
	if containsID(g.Discard(0), a) || containsID(g.Discard(0), b) {
		t.Error("chosen creatures should have left the discard pile")
	}

	// The house filter narrows the eligible creatures.
	houseSel := Chosen{House: namedHouse(Untamed), Type: Creature, Optional: true}
	g2 := NewGame("A", "B", 1)
	untamed := g2.AddToDiscard(NewCard("u", Untamed, Creature, Common, WithPower(1)), 0)
	g2.AddToDiscard(NewCard("d", Dis, Creature, Common, WithPower(1)), 0)
	ShuffleIntoDeck{
		Player:    Controller,
		From:      []Zone{Discard},
		Selection: houseSel,
		AnyNumber: true,
	}.
		Resolve(
			&EffectContext{Resolver: g2, Controller: 0, Source: untamed},
		)
	if g2.State.Deck[0].Count != 1 || containsID(g2.Discard(0), untamed) {
		t.Error("only the Untamed creature should have been shuffled away")
	}

	// Declining shuffles nothing.
	g3 := NewGame("A", "B", 1)
	kept := g3.AddToDiscard(testCreature("kept", 3), 0)
	g3.SetChooser(0, &cardDecliner{decline: true})
	ShuffleIntoDeck{Player: Controller, From: []Zone{Discard}, Selection: sel, AnyNumber: true}.
		Resolve(&EffectContext{Resolver: g3, Controller: 0, Source: kept})
	if g3.State.Deck[0].Count != 0 || !containsID(g3.Discard(0), kept) {
		t.Error("declining should shuffle nothing")
	}

	// Outside a shuffle batch the move narrates itself directly.
	g4 := NewGame("A", "B", 1)
	lone := g4.AddToDiscard(testCreature("lone", 3), 0)
	g4.ShuffleFromDiscardIntoDeck(lone)
	if g4.State.Deck[0].Count != 1 || containsID(g4.Discard(0), lone) {
		t.Error("a lone shuffle should move the card into the deck")
	}
}

// TestShuffleFromDiscardCount covers Shard of Life: the controller shuffles one
// card from their discard pile back into their deck for each friendly Shard.
func TestShuffleFromDiscardCount(t *testing.T) {
	sel := ShuffleIntoDeck{Player: Controller, From: []Zone{Discard},
		Selection: Chosen{},
		Count:     CardsInPlay{Player: Controller, Trait: Shard},
	}

	// Two friendly Shards in play shuffle two discard cards back.
	g := NewGame("A", "B", 1)
	g.AddToBattleline(testCreature("s1", 3, WithTraits(Shard)), 0)
	g.AddToBattleline(testCreature("s2", 3, WithTraits(Shard)), 0)
	g.AddToDiscard(NewCard("d1", Dis, Creature, Common, WithPower(1)), 0)
	g.AddToDiscard(NewCard("d2", Dis, Creature, Common, WithPower(1)), 0)
	g.AddToDiscard(NewCard("d3", Dis, Creature, Common, WithPower(1)), 0)
	src := g.AddArtifact(NewCard("Shard of Life", Untamed, Artifact, Rare), 0)
	sel.Resolve(&EffectContext{Resolver: g, Controller: 0, Source: src})
	if g.State.Deck[0].Count != 2 {
		t.Errorf("deck count = %d, want 2", g.State.Deck[0].Count)
	}
	if len(g.Discard(0)) != 1 {
		t.Errorf("discard size = %d, want 1", len(g.Discard(0)))
	}

	// No friendly Shard shuffles nothing.
	g2 := NewGame("A", "B", 1)
	g2.AddToDiscard(NewCard("d", Dis, Creature, Common, WithPower(1)), 0)
	src2 := g2.AddArtifact(NewCard("Shard of Life", Untamed, Artifact, Rare), 0)
	sel.Resolve(&EffectContext{Resolver: g2, Controller: 0, Source: src2})
	if g2.State.Deck[0].Count != 0 || len(g2.Discard(0)) != 1 {
		t.Error("with no friendly Shard nothing should be shuffled")
	}

	// Fewer discard cards than the count shuffles only what is there.
	g3 := NewGame("A", "B", 1)
	g3.AddToBattleline(testCreature("s", 3, WithTraits(Shard)), 0)
	g3.AddToBattleline(testCreature("s2", 3, WithTraits(Shard)), 0)
	only := g3.AddToDiscard(NewCard("only", Dis, Creature, Common, WithPower(1)), 0)
	src3 := g3.AddArtifact(NewCard("Shard of Life", Untamed, Artifact, Rare), 0)
	sel.Resolve(&EffectContext{Resolver: g3, Controller: 0, Source: src3})
	if g3.State.Deck[0].Count != 1 || containsID(g3.Discard(0), only) {
		t.Error("the single discard card should have been shuffled away")
	}
}

// TestShuffleFromDiscardNamed covers Chain Gang: one card of the named
// card leaves the discard for the deck while a differently-named card stays put.
func TestShuffleFromDiscardNamed(t *testing.T) {
	if got := (ShuffleIntoDeck{Player: Controller, From: []Zone{Discard}, Selection: Named{Name: "Subtle Chain"}}).Text(); got !=
		"shuffle Subtle Chain from your discard pile into your deck" {
		t.Errorf("text = %q", got)
	}

	g := NewGame("A", "B", 1)
	chain := g.AddToDiscard(NewCard("Subtle Chain", Dis, Tactic, Common), 0)
	other := g.AddToDiscard(NewCard("Mind Barb", Dis, Tactic, Common), 0)

	ShuffleIntoDeck{
		Player:    Controller,
		From:      []Zone{Discard},
		Selection: Named{Name: "Subtle Chain"},
	}.
		Resolve(
			&EffectContext{Resolver: g, Controller: 0, Source: other},
		)

	if containsID(g.Discard(0), chain) {
		t.Error("the named card should have left the discard pile")
	}
	if !containsID(g.Discard(0), other) {
		t.Error("a differently-named card should stay in the discard pile")
	}
	if g.State.Deck[0].Count != 1 {
		t.Errorf("deck count = %d, want 1", g.State.Deck[0].Count)
	}

	// With no card of that name in the discard pile, the effect shuffles nothing.
	g2 := NewGame("A", "B", 1)
	only := g2.AddToDiscard(NewCard("Mind Barb", Dis, Tactic, Common), 0)
	ShuffleIntoDeck{
		Player:    Controller,
		From:      []Zone{Discard},
		Selection: Named{Name: "Subtle Chain"},
	}.
		Resolve(
			&EffectContext{Resolver: g2, Controller: 0, Source: only},
		)
	if g2.State.Deck[0].Count != 0 || !containsID(g2.Discard(0), only) {
		t.Error("a discard pile without the named card should be left untouched")
	}
}
