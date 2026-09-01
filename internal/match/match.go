// Package match sets up Vactrol games in the way every frontend needs — the
// terminal UI, the web client, and the future lobby server all deal the same
// kind of random two-player match. Keeping this here means the frontends share
// one deck-building implementation instead of each copying it.
package match

import (
	"math/rand"

	"github.com/dmikalova/vactrol/internal/cards"
	"github.com/dmikalova/vactrol/internal/deckgen"
	"github.com/dmikalova/vactrol/internal/engine"
)

// DeckSize is the number of cards dealt to each player (opening hand plus draw
// zone).
const DeckSize = deckgen.DeckSize

// DeckHouseCount is how many houses make up a deck — three, as in KeyForge.
const DeckHouseCount = deckgen.PodCount

// New creates a two-player game seeded for deterministic play, deals each player
// a procedurally generated three-house deck, and returns the game together with
// each player's three houses. The caller installs choosers and calls BeginTurn to
// start play, so each frontend can wire in its own interaction model first.
func New(p0Name, p1Name string, seed int64) (*engine.Game, [2][]engine.House) {
	g, houses, _ := NewWithMavericks(p0Name, p1Name, seed)
	return g, houses
}

// NewWithMavericks is New plus, for each player, the LocalID of every Maverick
// card it was dealt — a card played out of its printed house — so a frontend can
// badge those cards. The game and houses are exactly what New returns.
func NewWithMavericks(
	p0Name, p1Name string, seed int64,
) (*engine.Game, [2][]engine.House, [2][]engine.LocalID) {
	g := engine.NewGame(p0Name, p1Name, seed)
	houses, mavericks := SetupDecks(g, seed)
	return g, houses, mavericks
}

// SetupDecks generates each player a deck (see internal/deckgen), shuffles it
// deterministically from the seed so the opening hand is not house-blocked, deals
// an opening hand, and returns each player's three houses (sorted by name)
// together with the LocalID of every Maverick card that player was dealt.
func SetupDecks(g *engine.Game, seed int64) ([2][]engine.House, [2][]engine.LocalID) {
	set := cards.DeckSet()
	var houses [2][]engine.House
	var mavericks [2][]engine.LocalID
	for player := 0; player < 2; player++ {
		deck := deckgen.Generate(set, seed+int64(player)+1)
		houses[player] = deck.Houses()

		// Deal each card with its Maverick flag alongside, so the flag survives the
		// shuffle and can be pinned to the LocalID the engine assigns on deal.
		defs := make([]engine.CardDefinition, 0, deckgen.DeckSize)
		maverick := make([]bool, 0, deckgen.DeckSize)
		for _, pod := range deck.Pods {
			for _, s := range pod.Slots {
				defs = append(defs, s.Card)
				maverick = append(maverick, s.Maverick)
			}
		}

		r := rand.New(rand.NewSource(seed + int64(player) + 100))
		r.Shuffle(len(defs), func(i, j int) {
			defs[i], defs[j] = defs[j], defs[i]
			maverick[i], maverick[j] = maverick[j], maverick[i]
		})
		for i, d := range defs {
			var id engine.LocalID
			if i < engine.HandSize {
				id = g.AddToHand(d, player)
			} else {
				id = g.AddToDeck(d, player)
			}
			if maverick[i] {
				mavericks[player] = append(mavericks[player], id)
			}
		}
	}
	return houses, mavericks
}
