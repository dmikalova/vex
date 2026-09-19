package deckgen

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/engine"
)

// synthSet builds a deterministic test Set: houses × rarities, several cards
// each, so draws always succeed.
func synthSet() Set {
	return NewSet("Test", synthCards(), DefaultTuning())
}

// synthCards is synthSet's flat card list, shared with tests that build a Set
// with custom tuning.
func synthCards() []Card {
	houses := []engine.House{engine.Brobnar, engine.Dis, engine.Logos, engine.Mars, engine.Sanctum}
	rarities := []engine.Rarity{engine.Common, engine.Uncommon, engine.Rare}
	var cs []Card
	for _, h := range houses {
		for _, rr := range rarities {
			for i := range 6 {
				name := h.String() + "-" + string(rr) + "-" + string(rune('a'+i))
				cs = append(
					cs,
					Card{Def: engine.NewCard(name, h, engine.Creature, rr, engine.WithPower(3))},
				)
			}
		}
	}
	return cs
}

func TestGenerateIsDeterministic(t *testing.T) {
	set := synthSet()
	a := Generate(set, 42)
	b := Generate(set, 42)
	ca, cb := a.Cards(), b.Cards()
	if len(ca) != len(cb) {
		t.Fatalf("card counts differ: %d vs %d", len(ca), len(cb))
	}
	for i := range ca {
		if ca[i].Name != cb[i].Name || ca[i].House != cb[i].House {
			t.Fatalf(
				"slot %d differs: %q/%v vs %q/%v",
				i,
				ca[i].Name,
				ca[i].House,
				cb[i].Name,
				cb[i].House,
			)
		}
	}
}

// legacyPool builds a legacy pool for the same houses as synthSet, with clearly
// distinct names so a legacy draw is recognizable.
func legacyPool() []Card {
	houses := []engine.House{engine.Brobnar, engine.Dis, engine.Logos, engine.Mars, engine.Sanctum}
	var cs []Card
	for _, h := range houses {
		for i := range 4 {
			name := "Legacy-" + h.String() + "-" + string(rune('a'+i))
			cs = append(
				cs,
				Card{
					Def: engine.NewCard(
						name,
						h,
						engine.Creature,
						engine.Common,
						engine.WithPower(2),
					),
				},
			)
		}
	}
	// Entries the legacy pool must skip: a houseless special, a Connected card,
	// and a houseless (HouseNone) card.
	cs = append(
		cs,
		Card{
			Def: engine.NewCard(
				"Legacy-Special",
				engine.Brobnar,
				engine.Creature,
				engine.Common,
			),
			Profile: GenerationProfile{Houseless: true},
		},
		Card{
			Def: engine.NewCard("Legacy-Connected", engine.Dis, engine.Creature, engine.Connected),
		},
		Card{
			Def: engine.NewCard(
				"Legacy-Houseless",
				engine.HouseNone,
				engine.Creature,
				engine.Common,
			),
		},
	)
	return cs
}

func TestLegacyDraws(t *testing.T) {
	tuning := DefaultTuning()
	tuning.LegacyRate = 1 // every non-special slot draws from the legacy pool
	entries := make([]LegacyEntry, 0)
	for _, c := range legacyPool() {
		entries = append(entries, LegacyEntry{
			Card: c,
			Set:  "Other",
		})
	}
	set := NewSet("Test", synthCards(), tuning).WithLegacy(NewLegacy(entries))

	deck := Generate(set, 3)
	legacyCount := 0
	for _, pod := range &deck.Pods {
		for _, s := range &pod.Slots {
			if s.Legacy {
				legacyCount++
				if s.Card.House != pod.House {
					t.Errorf("legacy card %q house = %v, want pod house %v",
						s.Card.Name, s.Card.House, pod.House)
				}
			}
		}
	}
	if legacyCount == 0 {
		t.Fatal("expected some slots to be filled from the legacy pool")
	}
}

// An interloper House pod is drawn entirely from the legacy pool: every slot a
// same-House card from another Set. With InterloperRate at 1 every pod is an
// interloper, so every slot is tagged legacy and keeps its pod's House.
func TestInterloperHousePod(t *testing.T) {
	tuning := DefaultTuning()
	tuning.InterloperRate = 1 // force every pod to draw from the legacy pool
	entries := make([]LegacyEntry, 0)
	for _, c := range legacyPool() {
		entries = append(entries, LegacyEntry{
			Card: c,
			Set:  "Other",
		})
	}
	set := NewSet("Test", synthCards(), tuning).WithLegacy(NewLegacy(entries))

	deck := Generate(set, 3)
	for _, pod := range &deck.Pods {
		for i, s := range &pod.Slots {
			if !s.Legacy {
				t.Errorf("pod %v slot %d not legacy in an interloper pod", pod.House, i)
			}
			if s.Card.House != pod.House {
				t.Errorf("pod %v slot %d house = %v, want pod house",
					pod.House, i, s.Card.House)
			}
		}
	}
	if got := len(deck.Cards()); got != DeckSize {
		t.Fatalf("deck has %d cards, want %d", got, DeckSize)
	}
}

// With no legacy pool the interloper roll never fires, so a single-set build
// draws exactly the deck it did before the overlay existed.
func TestInterloperNeedsLegacyPool(t *testing.T) {
	tuning := DefaultTuning()
	tuning.InterloperRate = 1
	set := NewSet("Test", synthCards(), tuning) // no legacy pool

	deck := Generate(set, 3)
	for _, pod := range &deck.Pods {
		for i, s := range &pod.Slots {
			if s.Legacy {
				t.Errorf("pod %v slot %d tagged legacy without a legacy pool", pod.House, i)
			}
		}
	}
}

// foreignLegacyPool builds a legacy pool of Houses foreign to synthSet — Houses an
// errant pod can bring into a deck.
func foreignLegacyPool() []Card {
	houses := []engine.House{engine.Shadows, engine.Untamed, engine.Saurian}
	var cs []Card
	for _, h := range houses {
		for i := range 4 {
			name := "Foreign-" + h.String() + "-" + string(rune('a'+i))
			cs = append(cs, Card{Def: engine.NewCard(
				name, h, engine.Creature, engine.Common, engine.WithPower(2),
			)})
		}
	}
	return cs
}

// An errant House pod replaces its native House with a foreign one (present in the
// legacy pool but not native to the Set) and draws every slot from the legacy pool
// as that foreign House. With ErrantRate at 1 and three foreign Houses available,
// every pod is an errant foreign pod, all-legacy, keeping its foreign House.
func TestErrantHousePod(t *testing.T) {
	tuning := DefaultTuning()
	tuning.ErrantRate = 1     // force every pod errant
	tuning.InterloperRate = 0 // isolate the errant overlay
	entries := make([]LegacyEntry, 0)
	for _, c := range foreignLegacyPool() {
		entries = append(entries, LegacyEntry{
			Card: c,
			Set:  "Other",
		})
	}
	set := NewSet("Test", synthCards(), tuning).WithLegacy(NewLegacy(entries))

	foreign := map[engine.House]bool{
		engine.Shadows: true, engine.Untamed: true, engine.Saurian: true,
	}
	deck := Generate(set, 3)
	for _, pod := range &deck.Pods {
		if !foreign[pod.House] {
			t.Errorf("pod House %v is not a foreign errant House", pod.House)
		}
		for i, s := range &pod.Slots {
			if !s.Legacy {
				t.Errorf("pod %v slot %d not legacy in an errant pod", pod.House, i)
			}
			if s.Card.House != pod.House {
				t.Errorf("pod %v slot %d House = %v, want pod House", pod.House, i, s.Card.House)
			}
		}
	}
	if got := len(deck.Cards()); got != DeckSize {
		t.Fatalf("deck has %d cards, want %d", got, DeckSize)
	}
}

// With no foreign House in the legacy pool the errant roll never fires, so every
// pod keeps a native House even at ErrantRate 1.
func TestErrantNeedsForeignHouse(t *testing.T) {
	tuning := DefaultTuning()
	tuning.ErrantRate = 1
	tuning.InterloperRate = 0
	entries := make([]LegacyEntry, 0)
	for _, c := range legacyPool() { // same Houses as synthSet: no foreign House
		entries = append(entries, LegacyEntry{
			Card: c,
			Set:  "Other",
		})
	}
	set := NewSet("Test", synthCards(), tuning).WithLegacy(NewLegacy(entries))

	native := map[engine.House]bool{}
	for _, h := range set.Houses() {
		native[h] = true
	}
	deck := Generate(set, 3)
	for _, pod := range &deck.Pods {
		if !native[pod.House] {
			t.Errorf("pod House %v is foreign, but the legacy pool has none", pod.House)
		}
	}
}

// A deck's Houses stay distinct: when more pods roll errant than there are foreign
// Houses, the surplus rolls find no free foreign House and fall back to a native
// pod. With a single foreign House, exactly one pod becomes errant.
func TestErrantExhaustsForeignHouses(t *testing.T) {
	tuning := DefaultTuning()
	tuning.ErrantRate = 1
	tuning.InterloperRate = 0
	entries := make([]LegacyEntry, 0)
	for i := range 4 {
		entries = append(entries, LegacyEntry{
			Card: mkCard("Sau-"+string(rune('a'+i)), engine.Saurian, engine.Common),
			Set:  "Other",
		})
	}
	set := NewSet("Test", synthCards(), tuning).WithLegacy(NewLegacy(entries))

	deck := Generate(set, 3)
	saurian := 0
	for _, pod := range &deck.Pods {
		if pod.House == engine.Saurian {
			saurian++
		}
	}
	if saurian != 1 {
		t.Fatalf("errant deck has %d Saurian pods, want exactly 1", saurian)
	}
}

// candidates drops the entries belonging to the drawing set and keeps the rest, so
// a Set never draws one of its own cards as a legacy card.
func TestLegacyCandidatesExcludesOwnSet(t *testing.T) {
	l := NewLegacy([]LegacyEntry{
		{Card: mkCard("Own", engine.Brobnar, engine.Common), Set: "Mine"},
		{Card: mkCard("Other", engine.Brobnar, engine.Common), Set: "Yours"},
	})
	got := l.candidates(l.byHouseRarity[engine.Brobnar][engine.Common], "Mine")
	if len(got) != 1 || got[0].Def.Name != "Other" {
		t.Fatalf("candidates = %v, want one card named Other", got)
	}
}

// A card the set prints as a reprint stays in the legacy pool (ADR 0021), but when
// a legacy slot draws it the slot is not tagged legacy: it is one of the set's own
// cards, not a guest from another set.
func TestLegacyDrawOfSetMemberNotTaggedLegacy(t *testing.T) {
	tuning := DefaultTuning()
	tuning.LegacyRate = 1 // every non-special slot draws from the legacy pool
	// The legacy pool holds the set's own cards, tagged as another set — the
	// shape of a reprint: the same card printed in this set and pooled as legacy
	// from its native set.
	entries := make([]LegacyEntry, 0)
	for _, c := range synthCards() {
		entries = append(entries, LegacyEntry{
			Card: c,
			Set:  "Other",
		})
	}
	set := NewSet("Test", synthCards(), tuning).WithLegacy(NewLegacy(entries))

	deck := Generate(set, 3)
	for _, pod := range &deck.Pods {
		for _, s := range &pod.Slots {
			if s.Legacy {
				t.Errorf("slot %q tagged legacy, but the set prints that card", s.Card.Name)
			}
		}
	}
	if got := len(deck.Cards()); got != DeckSize {
		t.Fatalf("deck has %d cards, want %d", got, DeckSize)
	}
}

// A reservoir card is undraftable — its set builds no draw pool of its own (that
// is what makes a set like the Anomaly Expansion "not offered for deck
// generation") — but a housed, non-Connected reservoir card is still a legal
// legacy guest: legacy and legacy-maverick slots in other sets can draw it. This
// is what lets Anomaly Expansion cards appear in decks without the set ever being
// draftable. NewLegacy gates only on Houseless/Connected/HouseNone, never on the
// reservoir flag, so the card is kept.
func TestNewLegacyKeepsReservoirCard(t *testing.T) {
	reservoir := Card{
		Def:     engine.NewCard("Anomaly", engine.Brobnar, engine.Creature, engine.Special),
		Profile: GenerationProfile{Reservoir: true},
	}
	if Draftable(reservoir) {
		t.Fatal("a reservoir card must not be draftable")
	}
	l := NewLegacy([]LegacyEntry{{Card: reservoir, Set: "Anomaly Expansion"}})
	got := l.candidates(l.byHouse[engine.Brobnar], "Some Other Set")
	if len(got) != 1 || got[0].Def.Name != "Anomaly" {
		t.Fatalf("legacy pool = %v, want the reservoir card kept as a legacy guest", got)
	}
}

func TestDeckShape(t *testing.T) {
	deck := Generate(synthSet(), 7)
	if got := len(deck.Cards()); got != DeckSize {
		t.Fatalf("deck has %d cards, want %d", got, DeckSize)
	}
	hs := deck.Houses()
	if len(hs) != PodCount {
		t.Fatalf("deck has %d houses, want %d", len(hs), PodCount)
	}
	for i := range hs {
		for j := i + 1; j < len(hs); j++ {
			if hs[i] == hs[j] {
				t.Fatalf("houses not distinct: %v", hs)
			}
		}
	}
}

func TestEveryCardAdoptsItsPodHouse(t *testing.T) {
	deck := Generate(synthSet(), 99)
	for _, pod := range &deck.Pods {
		for i, s := range &pod.Slots {
			if s.Card.House != pod.House {
				t.Errorf("pod %v slot %d has house %v", pod.House, i, s.Card.House)
			}
		}
	}
}

func TestSeedsProduceDifferentDecks(t *testing.T) {
	set := synthSet()
	same := 0
	a, b := Generate(set, 1).Cards(), Generate(set, 2).Cards()
	for i := range a {
		if a[i].Name == b[i].Name {
			same++
		}
	}
	if same == len(a) {
		t.Fatalf("two seeds produced identical decks")
	}
}

func TestRarityWeightsAreHonored(t *testing.T) {
	// A set whose rarity weight is entirely Rare must produce only Rare slots.
	tuning := DefaultTuning()
	tuning.RarityWeights = map[engine.Rarity]float64{engine.Rare: 1}
	tuning.MaverickRate = 0
	tuning.SpecialRate = 0
	set := synthSet()
	set.Tuning = tuning
	for _, pod := range Generate(set, 3).Pods {
		for i, s := range &pod.Slots {
			if s.Rarity != engine.Rare {
				t.Fatalf("pod %v slot %d rarity = %v, want Rare", pod.House, i, s.Rarity)
			}
		}
	}
}

func TestOneCopyPerDeckIsRespected(t *testing.T) {
	// A single common in a one-house set, flagged one-per-deck, must appear at most
	// once even though the house needs twelve slots.
	unique := Card{
		Def: engine.NewCard(
			"Unique",
			engine.Brobnar,
			engine.Creature,
			engine.Common,
			engine.WithPower(3),
		),
		Profile: GenerationProfile{OneCopyPerDeck: true},
	}
	filler := Card{
		Def: engine.NewCard(
			"Filler",
			engine.Brobnar,
			engine.Creature,
			engine.Common,
			engine.WithPower(3),
		),
	}
	set := NewSet("Solo", []Card{unique, filler}, DefaultTuning())
	deck := Generate(set, 5)
	count := 0
	for _, d := range deck.Cards() {
		if d.Name == "Unique" {
			count++
		}
	}
	if count > 1 {
		t.Fatalf("one-copy-per-deck card appeared %d times", count)
	}
}
