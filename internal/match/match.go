// Package match sets up Vactrol games in the way every frontend needs — the
// web client and the future lobby server all deal the same
// kind of random two-player match. Keeping this here means the frontends share
// one deck-building implementation instead of each copying it.
package match

import (
	"github.com/dmikalova/vactrol/internal/cards"
	"github.com/dmikalova/vactrol/internal/deckgen"
	"github.com/dmikalova/vactrol/internal/engine"
)

// DeckSize is the number of cards dealt to each player (opening hand plus draw
// zone).
const DeckSize = deckgen.DeckSize

// DeckHouseCount is how many houses make up a deck — three, as in KeyForge.
const DeckHouseCount = deckgen.PodCount

// Roster is a player's generated deck as a static, read-only list — its three
// House pods each with their twelve cards, kept from deck generation so a deck
// list can show the exact cards dealt (ADR 0025). It holds no draw order: the
// live draw pile lives in the engine deck zone, not here.
type Roster struct {
	Set    string
	Houses [deckgen.PodCount]HouseRoster
}

// HouseRoster is one House of a Roster with its twelve cards in slot order.
type HouseRoster struct {
	House engine.House
	Cards [deckgen.PodSize]RosterCard
}

// RosterCard is one card in a Roster: its definition and its Maverick/Legacy
// provenance. Rarity is read from the card definition, not the slot it filled, so
// a card shows its own rarity wherever it lands (a legacy card and its home-set
// printing match).
type RosterCard struct {
	Def      engine.CardDefinition
	Maverick bool
	Legacy   bool
}

// Empty reports whether the roster is the zero value (no deck generated), so a
// frontend built without a deal — the style gallery — draws no deck list.
func (r Roster) Empty() bool { return r.Houses[0].House == engine.HouseNone }

// rosterOf projects a generated deck into a static Roster: each pod's House and
// each slot's card and Maverick/Legacy provenance, in slot order.
func rosterOf(deck deckgen.Deck) Roster {
	var r Roster
	r.Set = deck.Set
	for i, pod := range deck.Pods {
		r.Houses[i].House = pod.House
		for j, s := range pod.Slots {
			r.Houses[i].Cards[j] = RosterCard{
				Def:      s.Card,
				Maverick: s.Maverick,
				Legacy:   s.Legacy,
			}
		}
	}
	return r
}

// New creates a two-player game seeded for deterministic play, generates each
// player a procedurally generated three-house deck into their deck zone, and
// returns the game together with each player's three houses. The caller installs
// choosers and calls engine.StartGame to deal opening hands and start play, so
// each frontend can wire in its own interaction model first.
func New(p0Name, p1Name string, seed int64) (*engine.Game, [2][]engine.House) {
	g, houses, _, _ := NewWithMavericks(p0Name, p1Name, seed)
	return g, houses
}

// NewWithMavericks is New plus, for each player, the LocalID of every Maverick
// card it was dealt — a card played out of its printed house — and the LocalID of
// every Legacy card — a card drawn from an earlier set's pool — so a frontend can
// badge those cards. The game and houses are exactly what New returns.
func NewWithMavericks(
	p0Name, p1Name string, seed int64,
) (*engine.Game, [2][]engine.House, [2][]engine.LocalID, [2][]engine.LocalID) {
	g, houses, mavericks, legacies, _ := NewWithSets(p0Name, p1Name, seed, [2]string{})
	return g, houses, mavericks, legacies
}

// NewWithSets is NewWithMavericks but each player's deck is generated from a named
// deck-generation set (as cards.DeckSetNamed resolves it). An empty or unknown
// name falls back to the default set, so a caller with no choice to make passes
// the zero value. It also returns each player's static deck Roster (ADR 0025).
func NewWithSets(
	p0Name, p1Name string, seed int64, setNames [2]string,
) (*engine.Game, [2][]engine.House, [2][]engine.LocalID, [2][]engine.LocalID, [2]Roster) {
	g := engine.NewGame(p0Name, p1Name, seed)
	houses, mavericks, legacies, rosters := SetupDecksFor(g, seed, setNames)
	return g, houses, mavericks, legacies, rosters
}

// SetupDecks generates each player a deck (see internal/deckgen) into their deck
// zone and returns each player's three houses (sorted by name) together with the
// LocalID of every Maverick card and every Legacy card that player was dealt.
// engine.StartGame shuffles and deals the opening hands.
func SetupDecks(
	g *engine.Game, seed int64,
) ([2][]engine.House, [2][]engine.LocalID, [2][]engine.LocalID) {
	houses, mavericks, legacies, _ := SetupDecksFor(g, seed, [2]string{})
	return houses, mavericks, legacies
}

// setFor resolves a chosen set name to its deck-generation Set, defaulting an
// empty or unknown name to the base set.
func setFor(name string) deckgen.Set {
	if name != "" {
		if s, ok := cards.DeckSetNamed(name); ok {
			return s
		}
	}
	return cards.DeckSet()
}

// SetupDecksFor is SetupDecks with each player's set named explicitly, so a
// frontend can let the two players play different sets. It also returns each
// player's static deck Roster — the generated deck kept as a read-only list so a
// deck list shows the exact cards dealt, immune to later pool changes (ADR 0025).
func SetupDecksFor(
	g *engine.Game, seed int64, setNames [2]string,
) ([2][]engine.House, [2][]engine.LocalID, [2][]engine.LocalID, [2]Roster) {
	var houses [2][]engine.House
	var mavericks [2][]engine.LocalID
	var legacies [2][]engine.LocalID
	var rosters [2]Roster
	for player := 0; player < 2; player++ {
		deck := deckgen.Generate(setFor(setNames[player]), seed+int64(player)+1)
		houses[player] = deck.Houses()
		rosters[player] = rosterOf(deck)

		// The whole deck goes into the deck zone; engine.StartGame shuffles it and
		// deals the opening hands. Each card's Maverick and Legacy flags are pinned to
		// the LocalID the engine assigns on add, so the badge survives that shuffle.
		for _, pod := range deck.Pods {
			for _, s := range pod.Slots {
				id := g.AddToDeck(s.Card, player)
				if s.Maverick {
					mavericks[player] = append(mavericks[player], id)
				}
				if s.Legacy {
					legacies[player] = append(legacies[player], id)
				}
			}
		}
	}
	return houses, mavericks, legacies, rosters
}
