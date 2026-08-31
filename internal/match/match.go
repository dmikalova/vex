// Package match sets up Vactrol games in the way every frontend needs — the
// terminal UI, the web client, and the future lobby server all deal the same
// kind of random two-player match. Keeping this here means the frontends share
// one deck-building implementation instead of each copying it.
package match

import (
	"math/rand"
	"sort"

	"github.com/dmikalova/vactrol/internal/cards"
	"github.com/dmikalova/vactrol/internal/engine"
)

// DeckSize is the number of cards dealt to each player (opening hand plus draw
// zone).
const DeckSize = 36

// DeckHouseCount is how many houses make up a deck — three, as in KeyForge.
const DeckHouseCount = 3

// New creates a two-player game seeded for deterministic play, deals each player
// a random three-house deck, and returns the game together with each player's
// three houses. The caller installs choosers and calls BeginTurn to start play,
// so each frontend can wire in its own interaction model first.
func New(p0Name, p1Name string, seed int64) (*engine.Game, [2][]engine.House) {
	g := engine.NewGame(p0Name, p1Name, seed)
	houses := SetupDecks(g, seed)
	return g, houses
}

// SetupDecks builds each player a KeyForge-style deck of DeckHouseCount houses
// drawn from the card pool, shuffles it deterministically from the seed, deals an
// opening hand, and returns each player's chosen houses (sorted by name).
func SetupDecks(g *engine.Game, seed int64) [2][]engine.House {
	pool := cards.All()
	available := poolHouses(pool)
	var houses [2][]engine.House
	for player := 0; player < 2; player++ {
		r := rand.New(rand.NewSource(seed + int64(player) + 1))
		houses[player] = pickHouses(available, r, DeckHouseCount)
		defs := dealDeck(pool, houses[player], r)
		for i, d := range defs {
			if i < engine.HandSize {
				g.AddToHand(d, player)
			} else {
				g.AddToDeck(d, player)
			}
		}
	}
	return houses
}

// dealDeck builds a DeckSize deck split evenly across the chosen houses (KeyForge
// deals a fixed share from each of three houses), drawing each house's share from
// its own cards — shuffled, and repeating a house's pool only if it has fewer
// cards than its share. The assembled deck is then shuffled so the opening hand is
// not house-blocked. Drawing per house is what keeps every chosen house present:
// a single house's cards now exceed DeckSize, so filling from one house-ordered
// list and truncating would yield a single-house deck.
func dealDeck(pool []engine.CardDefinition, hs []engine.House, r *rand.Rand) []engine.CardDefinition {
	defs := make([]engine.CardDefinition, 0, DeckSize)
	n := len(hs)
	if n == 0 {
		return defs
	}
	for i, h := range hs {
		share := DeckSize / n
		if i < DeckSize%n { // spread any remainder across the first houses
			share++
		}
		hCards := cardsOfHouses(pool, []engine.House{h})
		if len(hCards) == 0 {
			continue
		}
		r.Shuffle(len(hCards), func(a, b int) { hCards[a], hCards[b] = hCards[b], hCards[a] })
		for c := 0; c < share; c++ {
			defs = append(defs, hCards[c%len(hCards)])
		}
	}
	r.Shuffle(len(defs), func(i, j int) { defs[i], defs[j] = defs[j], defs[i] })
	return defs
}

// poolHouses returns the distinct houses present in the card pool, sorted by name.
func poolHouses(pool []engine.CardDefinition) []engine.House {
	seen := map[engine.House]bool{}
	var out []engine.House
	for _, d := range pool {
		if !seen[d.House] {
			seen[d.House] = true
			out = append(out, d.House)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].String() < out[j].String() })
	return out
}

// pickHouses chooses n distinct houses from available (capped at its length),
// deterministically from r, and returns them sorted by name.
func pickHouses(available []engine.House, r *rand.Rand, n int) []engine.House {
	if n > len(available) {
		n = len(available)
	}
	perm := r.Perm(len(available))[:n]
	picked := make([]engine.House, n)
	for i, k := range perm {
		picked[i] = available[k]
	}
	sort.Slice(picked, func(i, j int) bool { return picked[i].String() < picked[j].String() })
	return picked
}

// cardsOfHouses returns the pool cards whose house is one of hs.
func cardsOfHouses(pool []engine.CardDefinition, hs []engine.House) []engine.CardDefinition {
	inDeck := map[engine.House]bool{}
	for _, h := range hs {
		inDeck[h] = true
	}
	var out []engine.CardDefinition
	for _, d := range pool {
		if inDeck[d.House] {
			out = append(out, d)
		}
	}
	return out
}
