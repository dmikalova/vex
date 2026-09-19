package engine

import (
	"errors"
	"testing"
)

func TestSearchForName(t *testing.T) {
	e := Search{
		Sources: []Zone{Deck, Discard},
		Filter:  CardFilter{Name: "Timetraveller"},
		Reveal:  true,
		Dest:    ToHand,
	}
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
	if got := (Shuffle{Zones: []Zone{Discard}}).Text(); got != "shuffle your discard pile into your deck" {
		t.Errorf("discard text = %q", got)
	}
	if got := (Shuffle{Zones: []Zone{Hand, Discard}}).Text(); got != "shuffle your hand and discard pile into your deck" {
		t.Errorf("hand+discard text = %q", got)
	}
	if got := (Shuffle{Zones: []Zone{Archives, Discard}}).Text(); got != "shuffle your archives and discard pile into your deck" {
		t.Errorf("archives+discard text = %q", got)
	}
	if (Shuffle{}).validate() != nil {
		t.Error("the bare deck shuffle should be valid")
	}
	if (Shuffle{Zones: []Zone{zoneUnset}}).validate() == nil {
		t.Error("an unshuffleable zone should be invalid")
	}
	if (Shuffle{Zones: []Zone{Hand, Discard}}).validate() != nil {
		t.Error("hand and discard should be valid")
	}

	// Resolve: discard only.
	g := NewGame("A", "B", 1)
	src := g.AddToBattleline(testCreature("helper", 1), 0)
	g.State.Discard[0].add(g.Register(testCreature("a", 1), 0))
	g.State.Discard[0].add(g.Register(testCreature("b", 1), 0))
	ctx := &EffectContext{Resolver: g, Source: src, Controller: 0}
	Shuffle{Zones: []Zone{Discard}}.Resolve(ctx)
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
	Shuffle{Zones: []Zone{Hand, Archives, Discard}}.Resolve(ctx2)
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
	e := Search{
		Sources: []Zone{Deck, Discard},
		Filter:  CardFilter{Name: "Ancient Bear"},
		Any:     true,
		Reveal:  true,
		Dest:    ToHand,
	}
	want := "search your deck and discard pile for any number of Ancient Bears, reveal them, and put them into your hand"
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

func TestSearchDeck(t *testing.T) {
	if (Search{}).validate() == nil {
		t.Error("a search with no source zone should be invalid")
	}
	if (Search{Sources: []Zone{Deck}}).validate() == nil {
		t.Error("a search with no destination should be invalid")
	}
	if (Search{Sources: []Zone{Deck}, Dest: ToDeckShuffled}).validate() == nil {
		t.Error("a search to an unsupported destination should be invalid")
	}
	if (Search{Sources: []Zone{Deck}, Dest: ToHand}).validate() != nil {
		t.Error("a search naming its zone and destination should be valid")
	}
	if got := (Search{Sources: []Zone{Deck}, Dest: ToHand}).Text(); got !=
		"search your deck for a card and put it into your hand" {
		t.Errorf("unrestricted text = %q", got)
	}
	if got := (Search{Sources: []Zone{Deck}, House: namedHouse(Saurian), Reveal: true, Dest: ToHand}).Text(); got !=
		"search your deck for a Saurian card, reveal it, and put it into your hand" {
		t.Errorf("house text = %q", got)
	}
	if got := (Search{Sources: []Zone{Deck}, Filter: CardFilter{Type: Upgrade}, Reveal: true, Dest: ToHand}).Text(); got !=
		"search your deck for an upgrade, reveal it, and put it into your hand" {
		t.Errorf("filter text = %q", got)
	}

	// House-restricted: only the Saurian card is eligible; it is put into hand.
	g := NewGame("A", "B", 1)
	src := g.AddToBattleline(testCreature("rex", 6), 0)
	want := g.Register(NewCard("ally", Saurian, Creature, Common, WithPower(2)), 0)
	other := g.Register(NewCard("outsider", Logos, Creature, Common, WithPower(2)), 0)
	g.State.Deck[0].add(want)
	g.State.Deck[0].add(other)
	Search{
		Sources: []Zone{Deck},
		House:   namedHouse(Saurian),
		Dest:    ToHand,
	}.Resolve(
		&EffectContext{Resolver: g, Source: src, Controller: 0},
	)
	if !g.State.Hand[0].contains(want) {
		t.Error("the Saurian card should be in hand")
	}
	if g.State.Hand[0].contains(other) || !g.State.Deck[0].contains(other) {
		t.Error("the non-Saurian card should stay in the deck")
	}

	// Filter-restricted: only the upgrade is eligible; it is taken.
	gf := NewGame("A", "B", 1)
	sf := gf.AddToBattleline(testCreature("host", 4), 0)
	upgrade := gf.Register(NewCard("gizmo", Logos, Upgrade, Common), 0)
	creature := gf.Register(NewCard("body", Logos, Creature, Common, WithPower(2)), 0)
	gf.State.Deck[0].add(upgrade)
	gf.State.Deck[0].add(creature)
	Search{
		Sources: []Zone{Deck},
		Filter:  CardFilter{Type: Upgrade},
		Dest:    ToHand,
	}.Resolve(
		&EffectContext{Resolver: gf, Source: sf, Controller: 0},
	)
	if !gf.State.Hand[0].contains(upgrade) {
		t.Error("the upgrade should be in hand")
	}
	if gf.State.Hand[0].contains(creature) || !gf.State.Deck[0].contains(creature) {
		t.Error("the non-upgrade card should stay in the deck")
	}

	// Unrestricted: the sole deck card is taken.
	g2 := NewGame("A", "B", 1)
	s2 := g2.AddToBattleline(testCreature("orb-holder", 1), 0)
	only := g2.Register(NewCard("whatever", Logos, Tactic, Common), 0)
	g2.State.Deck[0].add(only)
	Search{
		Sources: []Zone{Deck},
		Dest:    ToHand,
	}.Resolve(
		&EffectContext{Resolver: g2, Source: s2, Controller: 0},
	)
	if !g2.State.Hand[0].contains(only) {
		t.Error("the sole deck card should be put into hand")
	}

	// No matching card: nothing is taken.
	g3 := NewGame("A", "B", 1)
	s3 := g3.AddToBattleline(testCreature("lonely", 1), 0)
	g3.State.Deck[0].add(g3.Register(NewCard("logos", Logos, Creature, Common, WithPower(1)), 0))
	before := len(g3.Hand(0))
	Search{
		Sources: []Zone{Deck},
		House:   namedHouse(Saurian),
		Dest:    ToHand,
	}.Resolve(
		&EffectContext{Resolver: g3, Source: s3, Controller: 0},
	)
	if len(g3.Hand(0)) != before {
		t.Error("a search that finds no match should put nothing into hand")
	}
}

// TestSearchToArchives covers an unrestricted two-zone search that reveals what it
// takes and puts it into the controller's archives (Horizon Saber): the Reveal
// flag forces the reveal an unrestricted search would otherwise skip, and the
// destination is the controller's own archives.
func TestSearchToArchives(t *testing.T) {
	e := Search{
		Sources: []Zone{Deck, Discard},
		Reveal:  true,
		Dest:    ToArchives,
	}
	want := "search your deck and discard pile for a card, reveal it, and put it into your archives"
	if e.Text() != want {
		t.Errorf("text = %q, want %q", e.Text(), want)
	}

	// A discard-pile card is revealed and moved to archives.
	g := NewGame("A", "B", 1)
	src := g.AddToBattleline(testCreature("saber", 11), 0)
	found := g.Register(NewCard("relic", Logos, Tactic, Common), 0)
	g.State.Discard[0].add(found)
	ctx := &EffectContext{Resolver: g, Source: src, Controller: 0}
	if !e.resolveGate(ctx) {
		t.Error("resolveGate reported taking nothing")
	}
	if !g.State.Archives[0].contains(found) {
		t.Error("the found card should be in the controller's archives")
	}
	if g.State.Discard[0].contains(found) {
		t.Error("the found card should leave the discard pile")
	}
}

// TestSearchToTopOfDeck covers an unrestricted-count search that puts every
// matching gigantic half onto the top of the controller's deck (Digging Up the
// Monster): the destination phrase reads "the top of your deck", and a half found
// in the searched zones is moved.
func TestSearchToTopOfDeck(t *testing.T) {
	e := Search{
		Sources: []Zone{Deck, Discard},
		Filter:  CardFilter{Gigantic: true},
		Any:     true,
		Reveal:  true,
		Dest:    ToTopOfDeck,
	}
	want := "search your deck and discard pile for any number of gigantic creatures, " +
		"reveal them, and put them into the top of your deck"
	if e.Text() != want {
		t.Errorf("text = %q, want %q", e.Text(), want)
	}

	// A gigantic half sitting in the discard pile is moved onto the deck.
	g := NewGame("A", "B", 1)
	src := g.AddToBattleline(testCreature("digger", 3), 0)
	half := g.Register(
		NewCard("titan", Logos, Creature, Common, WithPower(9), WithGiganticRole(GiganticBase)),
		0,
	)
	g.State.Discard[0].add(half)
	e.Resolve(&EffectContext{Resolver: g, Source: src, Controller: 0})
	if !g.State.Deck[0].contains(half) {
		t.Error("the gigantic half should be on the deck")
	}
	if g.State.Discard[0].contains(half) {
		t.Error("the gigantic half should leave the discard pile")
	}
}

func TestShuffleDeck(t *testing.T) {
	if got := (Shuffle{}).Text(); got != "shuffle your deck" {
		t.Errorf("text = %q", got)
	}

	g := NewGame("A", "B", 1)
	src := g.AddToBattleline(testCreature("shuffler", 1), 0)
	for range 3 {
		g.State.Deck[0].add(g.Register(NewCard("c", Logos, Creature, Common, WithPower(1)), 0))
	}
	before := g.State.Deck[0].Count
	Shuffle{}.Resolve(&EffectContext{Resolver: g, Source: src, Controller: 0})
	if g.State.Deck[0].Count != before {
		t.Errorf("shuffle changed deck size: %d, want %d", g.State.Deck[0].Count, before)
	}
	var shuffled bool
	for _, rec := range g.Log {
		if _, ok := rec.Entry.(DeckShuffled); ok {
			shuffled = true
		}
	}
	if !shuffled {
		t.Error("Shuffle should record a DeckShuffled log")
	}
}

// TestSearchUpToMaxGiganticHalves covers the gigantic tutors' count-capped
// search: the controller takes up to Max halves of a gigantic creature, one
// optional choice at a time, and the text names "two halves" (Max 2) or "either
// half" (the single-take default).
func TestSearchUpToMaxGiganticHalves(t *testing.T) {
	twoHalves := Search{
		Sources: []Zone{Deck, Discard},
		Filter:  CardFilter{Gigantic: true},
		Max:     2,
		Reveal:  true,
		Dest:    ToArchives,
	}
	if got := twoHalves.Text(); got !=
		"search your deck and discard pile for two halves of a gigantic creature, reveal them, and put them into your archives" {
		t.Errorf("two-halves text = %q", got)
	}
	eitherHalf := Search{
		Sources: []Zone{Deck, Discard},
		Filter:  CardFilter{Gigantic: true},
		Reveal:  true,
	}
	if got := eitherHalf.Text(); got !=
		"search your deck and discard pile for either half of a gigantic creature, reveal it, and put it into your hand" {
		t.Errorf("either-half text = %q", got)
	}

	newHalf := func(g *Game, role GiganticRole) LocalID {
		return g.Register(
			NewCard("Colossus", Logos, Creature, Common, WithPower(9), WithGiganticRole(role)),
			0,
		)
	}

	t.Run("takes up to two halves", func(t *testing.T) {
		g := NewGame("A", "B", 1)
		g.SetChooser(0, optionPicker{idx: 0}) // take the first candidate each time
		src := g.AddToBattleline(testCreature("digger", 3), 0)
		base := newHalf(g, GiganticBase)
		art := newHalf(g, GiganticArt)
		plain := g.Register(NewCard("plain", Logos, Creature, Common, WithPower(1)), 0)
		g.State.Deck[0].add(base)
		g.State.Discard[0].add(art)
		g.State.Deck[0].add(plain)
		ctx := &EffectContext{Resolver: g, Source: src, Controller: 0}

		if !twoHalves.resolveGate(ctx) {
			t.Error("resolveGate reported taking nothing")
		}
		if !g.State.Archives[0].contains(base) || !g.State.Archives[0].contains(art) {
			t.Error("both gigantic halves should be in archives")
		}
		if !g.State.Deck[0].contains(plain) {
			t.Error("the non-gigantic card should stay in the deck")
		}
	})

	t.Run("declining takes nothing", func(t *testing.T) {
		g := NewGame("A", "B", 1)
		src := g.AddToBattleline(testCreature("digger", 3), 0)
		g.State.Deck[0].add(newHalf(g, GiganticBase))
		g.State.Deck[0].add(newHalf(g, GiganticArt))
		g.SetChooser(0, &cardDecliner{decline: true})
		ctx := &EffectContext{Resolver: g, Source: src, Controller: 0}

		if twoHalves.resolveGate(ctx) {
			t.Error("declining every choice should report taking nothing")
		}
		if g.State.Archives[0].Count != 0 {
			t.Errorf("archives should stay empty, count = %d", g.State.Archives[0].Count)
		}
	})

	t.Run("stops when candidates run out before Max", func(t *testing.T) {
		g := NewGame("A", "B", 1)
		g.SetChooser(0, optionPicker{idx: 0})
		src := g.AddToBattleline(testCreature("digger", 3), 0)
		base := newHalf(g, GiganticBase) // the only half available
		g.State.Deck[0].add(base)
		ctx := &EffectContext{Resolver: g, Source: src, Controller: 0}

		if !twoHalves.resolveGate(ctx) {
			t.Error("the sole half should be taken")
		}
		if !g.State.Archives[0].contains(base) {
			t.Error("the sole half should be in archives")
		}
	})
}

// TestSearchShuffleBeforePlacing covers the mid-effect shuffle: the search finds
// and reveals the halves, shuffles the deck, then places them on top — so a half
// that was in the deck lands atop an already-shuffled deck rather than being
// buried by a shuffle after placement (Digging Up the Monster).
func TestSearchShuffleBeforePlacing(t *testing.T) {
	twoHalves := Search{
		Sources:              []Zone{Deck, Discard},
		Filter:               CardFilter{Gigantic: true},
		Max:                  2,
		Reveal:               true,
		ShuffleBeforePlacing: true,
		Dest:                 ToTopOfDeck,
	}
	if got := twoHalves.Text(); got !=
		"search your deck and discard pile for two halves of a gigantic creature, "+
			"reveal them, shuffle your deck, and put them into the top of your deck" {
		t.Errorf("two-halves text = %q", got)
	}

	newHalf := func(g *Game, role GiganticRole) LocalID {
		return g.Register(
			NewCard("Colossus", Logos, Creature, Common, WithPower(9), WithGiganticRole(role)),
			0,
		)
	}
	deckShuffled := func(g *Game) bool {
		for _, rec := range g.Log {
			if _, ok := rec.Entry.(DeckShuffled); ok {
				return true
			}
		}
		return false
	}

	t.Run("places a deck half and a discard half on top after shuffling", func(t *testing.T) {
		g := NewGame("A", "B", 1)
		g.SetChooser(0, optionPicker{idx: 0})
		src := g.AddToBattleline(testCreature("digger", 3), 0)
		base := newHalf(g, GiganticBase)
		art := newHalf(g, GiganticArt)
		plain := g.Register(NewCard("plain", Logos, Creature, Common, WithPower(1)), 0)
		g.State.Deck[0].add(base) // one half starts in the deck
		g.State.Deck[0].add(plain)
		g.State.Discard[0].add(art) // the other in the discard pile
		ctx := &EffectContext{Resolver: g, Source: src, Controller: 0}

		if !twoHalves.resolveGate(ctx) {
			t.Error("resolveGate reported taking nothing")
		}
		if !deckShuffled(g) {
			t.Error("the deck should have been shuffled")
		}
		if !g.State.Deck[0].contains(base) || !g.State.Deck[0].contains(art) {
			t.Error("both halves should be in the deck")
		}
		if g.State.Discard[0].contains(art) {
			t.Error("the discard half should have left the discard pile")
		}
		top := []LocalID{g.State.Deck[0].IDs[0], g.State.Deck[0].IDs[1]}
		if (top[0] != base && top[0] != art) || (top[1] != base && top[1] != art) {
			t.Errorf("both halves should sit on top of the deck, top = %v", top)
		}
	})

	t.Run("a single mandatory take without reveal", func(t *testing.T) {
		single := Search{
			Sources:              []Zone{Deck, Discard},
			Filter:               CardFilter{Gigantic: true},
			ShuffleBeforePlacing: true,
			Dest:                 ToTopOfDeck,
		}
		if got := single.Text(); got !=
			"search your deck and discard pile for either half of a gigantic creature, "+
				"shuffle your deck, and put it into the top of your deck" {
			t.Errorf("single no-reveal text = %q", got)
		}
		g := NewGame("A", "B", 1)
		src := g.AddToBattleline(testCreature("digger", 3), 0)
		base := newHalf(g, GiganticBase) // sole candidate, taken automatically
		g.State.Deck[0].add(base)
		ctx := &EffectContext{Resolver: g, Source: src, Controller: 0}

		if !single.resolveGate(ctx) {
			t.Error("the sole half should be taken")
		}
		if g.State.Deck[0].IDs[0] != base {
			t.Error("the half should sit on top of the deck")
		}
	})

	t.Run("declining still shuffles but places nothing", func(t *testing.T) {
		g := NewGame("A", "B", 1)
		src := g.AddToBattleline(testCreature("digger", 3), 0)
		base := newHalf(g, GiganticBase)
		art := newHalf(g, GiganticArt)
		g.State.Deck[0].add(base)
		g.State.Deck[0].add(art)
		g.SetChooser(0, &cardDecliner{decline: true})
		ctx := &EffectContext{Resolver: g, Source: src, Controller: 0}

		if twoHalves.resolveGate(ctx) {
			t.Error("declining should report taking nothing")
		}
		if !deckShuffled(g) {
			t.Error("the deck shuffles even when nothing is taken")
		}
	})

	t.Run("stops when candidates run out before Max", func(t *testing.T) {
		g := NewGame("A", "B", 1)
		g.SetChooser(0, optionPicker{idx: 0})
		src := g.AddToBattleline(testCreature("digger", 3), 0)
		base := newHalf(g, GiganticBase) // the only half available
		g.State.Deck[0].add(base)
		ctx := &EffectContext{Resolver: g, Source: src, Controller: 0}

		if !twoHalves.resolveGate(ctx) {
			t.Error("the sole half should be taken")
		}
		if !g.State.Deck[0].contains(base) {
			t.Error("the sole half should be back in the deck on top")
		}
	})
}

// TestSearchValidateShuffleBeforePlacing rejects a mid-effect shuffle on a search
// that does not place its finds on top of the deck.
func TestSearchValidateShuffleBeforePlacing(t *testing.T) {
	bad := Search{
		Sources:              []Zone{Deck},
		ShuffleBeforePlacing: true,
		Dest:                 ToHand,
	}
	if err := bad.validate(); !errors.Is(err, errSearchShuffleNotTopOfDeck) {
		t.Errorf("validate = %v, want %v", err, errSearchShuffleNotTopOfDeck)
	}
	ok := Search{
		Sources:              []Zone{Deck},
		ShuffleBeforePlacing: true,
		Dest:                 ToTopOfDeck,
	}
	if err := ok.validate(); err != nil {
		t.Errorf("validate = %v, want nil", err)
	}
}
