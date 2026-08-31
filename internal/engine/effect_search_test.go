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
				WithTraits("Human", "Scientist"),
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

func TestShuffleDiscard(t *testing.T) {
	g := NewGame("A", "B", 1)
	src := g.AddToBattleline(testCreature("helper", 1), 0)
	g.State.Discard[0].add(g.Register(testCreature("a", 1), 0))
	g.State.Discard[0].add(g.Register(testCreature("b", 1), 0))
	ctx := &EffectContext{Resolver: g, Source: src, Controller: 0}

	e := ShuffleDiscard{}
	if e.Text() != "shuffle your discard pile into your deck" {
		t.Errorf("text = %q", e.Text())
	}
	e.Resolve(ctx)
	if g.State.Discard[0].Count != 0 {
		t.Errorf("discard should be empty, count = %d", g.State.Discard[0].Count)
	}
	if g.State.Deck[0].Count != 2 {
		t.Errorf("deck should hold the 2 shuffled cards, count = %d", g.State.Deck[0].Count)
	}
}

func TestShuffleHandAndDiscard(t *testing.T) {
	g := NewGame("A", "B", 1)
	src := g.AddToBattleline(testCreature("helper", 1), 0)
	g.State.Hand[0].add(g.Register(testCreature("h1", 1), 0))
	g.State.Hand[0].add(g.Register(testCreature("h2", 1), 0))
	g.State.Discard[0].add(g.Register(testCreature("d1", 1), 0))
	ctx := &EffectContext{Resolver: g, Source: src, Controller: 0}

	e := ShuffleHandAndDiscard{}
	if e.Text() != "shuffle your hand and discard pile into your deck" {
		t.Errorf("text = %q", e.Text())
	}
	e.Resolve(ctx)
	if g.State.Hand[0].Count != 0 {
		t.Errorf("hand should be empty, count = %d", g.State.Hand[0].Count)
	}
	if g.State.Discard[0].Count != 0 {
		t.Errorf("discard should be empty, count = %d", g.State.Discard[0].Count)
	}
	if g.State.Deck[0].Count != 3 {
		t.Errorf("deck should hold the 3 shuffled cards, count = %d", g.State.Deck[0].Count)
	}
}
