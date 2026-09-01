package match

import (
	"testing"

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

func TestChosenHousesAreDistinct(t *testing.T) {
	_, houses := New("Alice", "Bob", 7)
	for p := range 2 {
		hs := houses[p]
		for i := range hs {
			for j := i + 1; j < len(hs); j++ {
				if hs[i] == hs[j] {
					t.Errorf("player %d houses not distinct: %v", p, hs)
				}
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

func TestDealSpansAllChosenHouses(t *testing.T) {
	g, houses := New("Alice", "Bob", 7)
	for p := range 2 {
		seen := map[engine.House]bool{}
		deck := g.State.Deck[p]
		ids := append(append([]engine.LocalID{}, g.Hand(p)...), deck.IDs[:deck.Count]...)
		for _, id := range ids {
			seen[g.House(id)] = true
		}
		if len(seen) != len(houses[p]) {
			t.Errorf("player %d deck spans %d houses, want all %d (%v)",
				p, len(seen), len(houses[p]), houses[p])
		}
	}
}
