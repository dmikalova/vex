package match

import (
	"math/rand"
	"testing"

	"github.com/dmikalova/vactrol/internal/cards"
	"github.com/dmikalova/vactrol/internal/engine"
)

func housesEqual(a, b []engine.House) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestNewDealsDeterministicDecks(t *testing.T) {
	g1, h1 := New("Alice", "Bob", 42)
	_, h2 := New("Alice", "Bob", 42)

	for p := range 2 {
		if !housesEqual(h1[p], h2[p]) {
			t.Errorf("player %d houses not deterministic: %v vs %v", p, h1[p], h2[p])
		}
		if len(h1[p]) != DeckHouseCount {
			t.Errorf("player %d got %d houses, want %d", p, len(h1[p]), DeckHouseCount)
		}
		if got := len(g1.Hand(p)); got != engine.HandSize {
			t.Errorf("player %d hand = %d, want %d", p, got, engine.HandSize)
		}
		if got := int(g1.State.Deck[p].Count); got != DeckSize-engine.HandSize {
			t.Errorf("player %d deck = %d, want %d", p, got, DeckSize-engine.HandSize)
		}
	}
}

func TestChosenHousesAreDistinctAndSorted(t *testing.T) {
	_, houses := New("Alice", "Bob", 7)
	for p := range 2 {
		hs := houses[p]
		for i := 1; i < len(hs); i++ {
			if hs[i-1].String() >= hs[i].String() {
				t.Errorf("player %d houses not distinct/sorted: %v", p, hs)
			}
		}
	}
}

func TestDealtCardsBelongToChosenHouses(t *testing.T) {
	g, houses := New("Alice", "Bob", 7)
	for p := range 2 {
		allowed := map[engine.House]bool{}
		for _, h := range houses[p] {
			allowed[h] = true
		}
		check := func(ids []engine.LocalID, zone string) {
			for _, id := range ids {
				if !allowed[g.House(id)] {
					t.Errorf("player %d %s card of house %v not in deck houses %v",
						p, zone, g.House(id), houses[p])
				}
			}
		}
		deck := g.State.Deck[p]
		check(g.Hand(p), "hand")
		check(deck.IDs[:deck.Count], "deck")
	}
}

func TestPoolHousesDistinctAndSorted(t *testing.T) {
	hs := poolHouses(cards.All())
	if len(hs) == 0 {
		t.Fatal("card pool has no houses")
	}
	for i := 1; i < len(hs); i++ {
		if hs[i-1].String() >= hs[i].String() {
			t.Errorf("poolHouses not distinct/sorted: %v", hs)
		}
	}
}

func TestPickHousesCapsAtAvailable(t *testing.T) {
	available := []engine.House{engine.Brobnar, engine.Logos}
	r := rand.New(rand.NewSource(1))
	got := pickHouses(available, r, DeckHouseCount+5)
	if len(got) != len(available) {
		t.Errorf("pickHouses(%d requested) = %d houses, want %d (capped)",
			DeckHouseCount+5, len(got), len(available))
	}
}

func TestCardsOfHousesFiltersToHouses(t *testing.T) {
	pool := cards.All()
	house := poolHouses(pool)[0]
	filtered := cardsOfHouses(pool, []engine.House{house})
	if len(filtered) == 0 {
		t.Fatalf("no cards for house %v", house)
	}
	for _, d := range filtered {
		if d.House != house {
			t.Errorf("card %q has house %v, want %v", d.Name, d.House, house)
		}
	}
}
