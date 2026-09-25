package match

import (
	"errors"
	"testing"

	"github.com/dmikalova/vex/internal/engine"
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
		// New leaves the whole deck in the deck zone; engine.StartGame deals the
		// opening hands.
		if got := len(g1.Hand(p)); got != 0 {
			t.Errorf("player %d hand = %d, want 0 before StartGame", p, got)
		}
		if got := int(g1.State.Deck[p].Count); got != DeckSize {
			t.Errorf("player %d deck = %d, want %d", p, got, DeckSize)
		}
	}
}

func TestStartGameDealsOpeningHands(t *testing.T) {
	g, _ := New("Alice", "Bob", 42)
	g.StartGame(0)
	// The first player draws one more than the second (7 vs a full HandSize of 6).
	if got := len(g.Hand(0)); got != engine.HandSize+engine.FirstPlayerBonusCards {
		t.Errorf(
			"first player hand = %d, want %d",
			got,
			engine.HandSize+engine.FirstPlayerBonusCards,
		)
	}
	if got := len(g.Hand(1)); got != engine.HandSize {
		t.Errorf("second player hand = %d, want %d", got, engine.HandSize)
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

func TestUnknownSetNameIsAnError(t *testing.T) {
	// An empty name is the caller declining to choose, not a mistake.
	if _, _, _, _, _, err := NewWithSets("Alice", "Bob", 7, [2]string{}); err != nil {
		t.Fatalf("an unnamed set should deal the default set: %v", err)
	}
	// A name no set answers to would otherwise deal the default set, and the game
	// would look fine while ignoring the choice that was made.
	_, _, _, _, _, err := NewWithSets("Alice", "Bob", 7, [2]string{"", "no such set"})
	if !errors.Is(err, ErrUnknownSet) {
		t.Errorf("unknown set name gave %v, want ErrUnknownSet", err)
	}
}
