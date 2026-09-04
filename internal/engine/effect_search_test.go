package engine

import "testing"

func TestSearchForName(t *testing.T) {
	e := SearchForName{Name: "Timetraveller"}
	if e.Text() != "search your deck and discard pile for a Timetraveller, reveal it, and put it into your hand" {
		t.Errorf("text = %q", e.Text())
	}

	newTT := func(g *Game, player int) LocalID {
		return g.Register(
			NewCard(
				"Timetraveller",
				Logos,
				Creature,
				Common,
				WithPower(2),
				WithTraits(Human, Scientist),
			),
			player,
		)
	}

	t.Run("from deck", func(t *testing.T) {
		g := NewGame("A", "B", 1)
		src := g.AddToBattleline(testCreature("helper", 1), 0)
		tt := newTT(g, 0)
		g.State.Deck[0].add(tt)
		g.State.Deck[0].add(
			g.Register(NewCard("plain", Logos, Creature, Common, WithPower(1)), 0),
		) // non-match in deck
		g.State.Discard[0].add(
			g.Register(NewCard("junk", Dis, Tactic, Common), 0),
		) // non-match in discard
		ctx := &EffectContext{Resolver: g, Source: src, Controller: 0}

		e.Resolve(ctx) // the sole Timetraveller is auto-chosen
		if g.State.Deck[0].contains(tt) {
			t.Error("the Timetraveller should leave the deck")
		}
		if !g.State.Hand[0].contains(tt) {
			t.Error("the Timetraveller should be in the controller's hand")
		}
	})

	t.Run("from discard", func(t *testing.T) {
		g := NewGame("A", "B", 1)
		src := g.AddToBattleline(testCreature("helper", 1), 0)
		tt := newTT(g, 0)
		g.State.Discard[0].add(tt)
		ctx := &EffectContext{Resolver: g, Source: src, Controller: 0}

		e.Resolve(ctx)
		if g.State.Discard[0].contains(tt) {
			t.Error("the Timetraveller should leave the discard pile")
		}
		if !g.State.Hand[0].contains(tt) {
			t.Error("the Timetraveller should be in the controller's hand")
		}
	})

	t.Run("no match", func(t *testing.T) {
		g := NewGame("A", "B", 1)
		src := g.AddToBattleline(testCreature("helper", 1), 0)
		g.State.Deck[0].add(g.Register(NewCard("plain", Logos, Creature, Common, WithPower(1)), 0))
		ctx := &EffectContext{Resolver: g, Source: src, Controller: 0}

		e.Resolve(ctx) // no Timetraveller anywhere, so nothing moves
		if g.State.Hand[0].Count != 0 {
			t.Errorf("hand should stay empty, count = %d", g.State.Hand[0].Count)
		}
	})
}

func TestShuffleIntoDeck(t *testing.T) {
	// Text and validate.
	if got := (ShuffleIntoDeck{Zones: []Zone{Discard}}).Text(); got != "shuffle your discard pile into your deck" {
		t.Errorf("discard text = %q", got)
	}
	if got := (ShuffleIntoDeck{Zones: []Zone{Hand, Discard}}).Text(); got != "shuffle your hand and discard pile into your deck" {
		t.Errorf("hand+discard text = %q", got)
	}
	if got := (ShuffleIntoDeck{Zones: []Zone{Archives, Discard}}).Text(); got != "shuffle your archives and discard pile into your deck" {
		t.Errorf("archives+discard text = %q", got)
	}
	if (ShuffleIntoDeck{}).validate() == nil {
		t.Error("no zones should be invalid")
	}
	if (ShuffleIntoDeck{Zones: []Zone{zoneUnset}}).validate() == nil {
		t.Error("an unshuffleable zone should be invalid")
	}
	if (ShuffleIntoDeck{Zones: []Zone{Hand, Discard}}).validate() != nil {
		t.Error("hand and discard should be valid")
	}

	// Resolve: discard only.
	g := NewGame("A", "B", 1)
	src := g.AddToBattleline(testCreature("helper", 1), 0)
	g.State.Discard[0].add(g.Register(testCreature("a", 1), 0))
	g.State.Discard[0].add(g.Register(testCreature("b", 1), 0))
	ctx := &EffectContext{Resolver: g, Source: src, Controller: 0}
	ShuffleIntoDeck{Zones: []Zone{Discard}}.Resolve(ctx)
	if g.State.Discard[0].Count != 0 || g.State.Deck[0].Count != 2 {
		t.Errorf(
			"discard shuffle: discard=%d deck=%d, want 0/2",
			g.State.Discard[0].Count,
			g.State.Deck[0].Count,
		)
	}

	// Resolve: hand, discard, and archives all at once.
	g2 := NewGame("A", "B", 1)
	g2.State.Hand[0].add(g2.Register(testCreature("h1", 1), 0))
	g2.State.Hand[0].add(g2.Register(testCreature("h2", 1), 0))
	g2.State.Discard[0].add(g2.Register(testCreature("d1", 1), 0))
	g2.State.Archives[0].add(g2.Register(testCreature("ar1", 1), 0))
	ctx2 := &EffectContext{Resolver: g2, Controller: 0}
	ShuffleIntoDeck{Zones: []Zone{Hand, Archives, Discard}}.Resolve(ctx2)
	if g2.State.Hand[0].Count != 0 || g2.State.Archives[0].Count != 0 ||
		g2.State.Discard[0].Count != 0 {
		t.Errorf("zones should be empty: hand=%d archives=%d discard=%d",
			g2.State.Hand[0].Count, g2.State.Archives[0].Count, g2.State.Discard[0].Count)
	}
	if g2.State.Deck[0].Count != 4 {
		t.Errorf("deck should hold the 4 shuffled cards, count = %d", g2.State.Deck[0].Count)
	}
}

func TestSearchForNameAll(t *testing.T) {
	e := SearchForName{Name: "Ancient Bear", All: true}
	want := "search your deck and discard pile and put each Ancient Bear from them into your hand"
	if e.Text() != want {
		t.Errorf("text = %q, want %q", e.Text(), want)
	}

	newBear := func(g *Game, player int) LocalID {
		return g.Register(NewCard("Ancient Bear", Untamed, Creature, Common, WithPower(6)), player)
	}

	g := NewGame("A", "B", 1)
	src := g.AddToBattleline(testCreature("flute", 1), 0)
	inDeck, inDiscard := newBear(g, 0), newBear(g, 0)
	g.State.Deck[0].add(inDeck)
	g.State.Discard[0].add(inDiscard)
	ctx := &EffectContext{Resolver: g, Source: src, Controller: 0}

	if !e.resolveGate(ctx) {
		t.Error("resolveGate reported finding nothing")
	}
	if !g.State.Hand[0].contains(inDeck) || !g.State.Hand[0].contains(inDiscard) {
		t.Error("both copies should be in hand, from the deck and the discard pile")
	}

	if e.resolveGate(ctx) {
		t.Error("resolveGate reported a find with no copies left")
	}
}
