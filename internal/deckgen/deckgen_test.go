package deckgen

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/engine"
)

// synthSet builds a deterministic test Set: houses × rarities, several cards
// each, so draws always succeed.
func synthSet() Set {
	houses := []engine.House{engine.Brobnar, engine.Dis, engine.Logos, engine.Mars, engine.Sanctum}
	rarities := []engine.Rarity{engine.Common, engine.Uncommon, engine.Rare}
	var cs []Card
	for _, h := range houses {
		for _, rr := range rarities {
			for i := 0; i < 6; i++ {
				name := h.String() + "-" + string(rr) + "-" + string(rune('a'+i))
				cs = append(
					cs,
					Card{Def: engine.NewCard(name, h, engine.Creature, rr, engine.WithPower(3))},
				)
			}
		}
	}
	return NewSet("Test", cs, DefaultTuning())
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
	for _, pod := range deck.Pods {
		for i, s := range pod.Slots {
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
		for i, s := range pod.Slots {
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
