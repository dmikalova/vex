package engine

import (
	"errors"
	"testing"
)

// nameChooser answers a labeled option prompt by picking the option whose label
// equals want, so a test can name a specific card. It falls back to the first
// option when no label matches.
type nameChooser struct{ want string }

func (nameChooser) ChooseCreature(_, _ string, cands []LocalID) (LocalID, bool) {
	if len(cands) == 0 {
		return 0, false
	}
	return cands[0], true
}

func (c nameChooser) ChooseOption(_, _ string, options []string) int {
	for i, o := range options {
		if o == c.want {
			return i
		}
	}
	return 0
}

// outOfRangeOptionChooser returns an option index past the end of the list, so a
// test can exercise the bounds guard in NameCard.Resolve.
type outOfRangeOptionChooser struct{}

func (outOfRangeOptionChooser) ChooseCreature(_, _ string, cands []LocalID) (LocalID, bool) {
	if len(cands) == 0 {
		return 0, false
	}
	return cands[0], true
}

func (outOfRangeOptionChooser) ChooseOption(_, _ string, options []string) int {
	return len(options)
}

func TestNameCardText(t *testing.T) {
	want := "name a card. Cards with that name cannot be played until " +
		SelfName + " leaves play"
	if got := (NameCard{}).Text(); got != want {
		t.Fatalf("Text() = %q, want %q", got, want)
	}
}

func TestNameableCardsDedupsAndSorts(t *testing.T) {
	g := NewGame("Alice", "Bob", 1)
	g.AddToHand(NewCard("Zebra", Brobnar, Creature, Common, WithPower(1)), 0)
	g.AddToHand(NewCard("Apple", Logos, Creature, Common, WithPower(1)), 0)
	// A second copy of Apple must not add a second option.
	g.AddToDeck(NewCard("Apple", Logos, Creature, Common, WithPower(1)), 1)

	names := g.NameableCards()
	var got []string
	for _, id := range names {
		got = append(got, g.Name(id))
	}
	if len(got) != 2 || got[0] != "Apple" || got[1] != "Zebra" {
		t.Fatalf("NameableCards names = %v, want [Apple Zebra]", got)
	}
}

func TestNameCardResolveRecordsName(t *testing.T) {
	g := started(t)
	jar := g.AddArtifact(NewCard("Etan's Jar", Dis, Artifact, Rare), 0)
	troll := g.AddToHand(NewCard("Troll", Brobnar, Creature, Common, WithPower(8)), 1)
	g.SetChooser(0, nameChooser{want: "Troll"})

	NameCard{}.Resolve(&EffectContext{Resolver: g, Source: jar, Controller: 0})

	if got := g.State.Cards[jar].NamedCardPlus; got != uint8(troll)+1 {
		t.Fatalf("NamedCardPlus = %d, want %d", got, uint8(troll)+1)
	}
}

func TestNameCardResolveOutOfRange(t *testing.T) {
	g := started(t)
	jar := g.AddArtifact(NewCard("Etan's Jar", Dis, Artifact, Rare), 0)
	g.AddToHand(NewCard("Troll", Brobnar, Creature, Common, WithPower(8)), 1)
	g.SetChooser(0, outOfRangeOptionChooser{})

	NameCard{}.Resolve(&EffectContext{Resolver: g, Source: jar, Controller: 0})

	if got := g.State.Cards[jar].NamedCardPlus; got != 0 {
		t.Fatalf("NamedCardPlus = %d, want 0 on out-of-range choice", got)
	}
}

// The names offered come from the injected card database, not the match, so a
// player naming a card is not shown their opponent's deck list. A name no card in
// the match carries bars nothing, so it records nothing.
func TestNameCardOffersTheInjectedDatabase(t *testing.T) {
	g := started(t)
	jar := g.AddArtifact(NewCard("Etan's Jar", Dis, Artifact, Rare), 0)
	troll := g.AddToHand(NewCard("Troll", Brobnar, Creature, Common, WithPower(8)), 1)
	// Deliberately unsorted and duplicated: SetNameableNames settles the order.
	g.SetNameableNames([]string{"Troll", "Ancient Bear", "Troll"})
	if got := g.NameableNames(); len(got) != 2 || got[0] != "Ancient Bear" || got[1] != "Troll" {
		t.Fatalf("NameableNames = %v, want [Ancient Bear Troll]", got)
	}

	g.SetChooser(0, nameChooser{want: "Troll"})
	NameCard{}.Resolve(&EffectContext{Resolver: g, Source: jar, Controller: 0})
	if got := g.State.Cards[jar].NamedCardPlus; got != uint8(troll)+1 {
		t.Fatalf("NamedCardPlus = %d, want %d", got, uint8(troll)+1)
	}

	// Ancient Bear is in the database but not in this match, so nothing is recorded.
	g.SetNamedCard(jar, 0)
	g.State.Cards[jar].NamedCardPlus = 0
	g.SetChooser(0, nameChooser{want: "Ancient Bear"})
	NameCard{}.Resolve(&EffectContext{Resolver: g, Source: jar, Controller: 0})
	if got := g.State.Cards[jar].NamedCardPlus; got != 0 {
		t.Fatalf("NamedCardPlus = %d, want 0 for a name no card in the match carries", got)
	}
}

func TestNameCardResolveNoCandidates(t *testing.T) {
	g := NewGame("Alice", "Bob", 1)
	jar := g.AddArtifact(NewCard("Etan's Jar", Dis, Artifact, Rare), 0)
	// With nothing registered but the Jar itself, its own name is a candidate;
	// remove it from the catalog view by resolving on an empty match instead.
	g.cat.defs = g.cat.defs[:0]

	NameCard{}.Resolve(&EffectContext{Resolver: g, Source: jar, Controller: 0})

	if got := g.State.Cards[jar].NamedCardPlus; got != 0 {
		t.Fatalf("NamedCardPlus = %d, want 0 with no candidates", got)
	}
}

func TestBarredByNamedCardBarsBothPlayers(t *testing.T) {
	g := started(t)
	jar := g.AddArtifact(NewCard("Etan's Jar", Dis, Artifact, Rare), 0)
	troll := NewCard("Troll", Brobnar, Creature, Common, WithPower(8))
	other := NewCard("Gob", Brobnar, Creature, Common, WithPower(1))

	// Name the troll on the jar directly.
	named := g.AddToDeck(troll, 0)
	g.SetNamedCard(jar, named)

	if !g.barredByNamedCard(g.cat.def(named)) {
		t.Fatal("named card should be barred")
	}
	if g.barredByNamedCard(&other) {
		t.Fatal("differently named card should not be barred")
	}

	// The bar lifts when the naming permanent leaves play.
	g.removeFromPlay(jar)
	if g.barredByNamedCard(g.cat.def(named)) {
		t.Fatal("bar should lift when the naming permanent leaves play")
	}
}

func TestCanPlayBarredByName(t *testing.T) {
	g := started(t)
	jar := g.AddArtifact(NewCard("Etan's Jar", Dis, Artifact, Rare), 0)
	troll := g.AddToHand(
		NewCard("Troll", Brobnar, Creature, Common, WithPower(8)), 0,
	)
	g.SetNamedCard(jar, troll)

	if err := g.CanPlay(0, troll); !errors.Is(err, ErrCannotPlayName) {
		t.Fatalf("CanPlay = %v, want ErrCannotPlayName", err)
	}
	idx := handIdxByID(g, 0, troll)
	if _, err := g.PlayCreature(0, idx, false); !errors.Is(err, ErrCannotPlayName) {
		t.Fatalf("PlayCreature = %v, want ErrCannotPlayName", err)
	}
}
