// Package match sets up Vactrol games in the way every frontend needs — the
// web client and the future lobby server all deal the same
// kind of random two-player match. Keeping this here means the frontends share
// one deck-building implementation instead of each copying it.
package match

import (
	"errors"
	"fmt"

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
	for i := range deck.Pods {
		pod := &deck.Pods[i]
		r.Houses[i].House = pod.House
		for j := range pod.Slots {
			s := &pod.Slots[j]
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
	g, houses, mavericks, legacies, _, _ := NewWithSets(p0Name, p1Name, seed, [2]string{})
	return g, houses, mavericks, legacies
}

// NewWithSets is NewWithMavericks but each player's deck is generated from a named
// deck-generation set (as cards.DeckSetNamed resolves it). An empty name means
// "the default set", so a caller with no choice to make passes the zero value; a
// name no set answers to is a bug in the caller and errors. It also returns each
// player's static deck Roster (ADR 0025).
func NewWithSets(
	p0Name, p1Name string, seed int64, setNames [2]string,
) (*engine.Game, [2][]engine.House, [2][]engine.LocalID, [2][]engine.LocalID, [2]Roster, error) {
	g := engine.NewGame(p0Name, p1Name, seed)
	houses, mavericks, legacies, rosters, err := SetupDecksFor(g, seed, setNames)
	if err != nil {
		return nil, houses, mavericks, legacies, rosters, err
	}
	return g, houses, mavericks, legacies, rosters, nil
}

// SetupDecks generates each player a deck (see internal/deckgen) into their deck
// zone and returns each player's three houses (sorted by name) together with the
// LocalID of every Maverick card and every Legacy card that player was dealt.
// engine.StartGame shuffles and deals the opening hands.
func SetupDecks(
	g *engine.Game, seed int64,
) ([2][]engine.House, [2][]engine.LocalID, [2][]engine.LocalID) {
	houses, mavericks, legacies, _, _ := SetupDecksFor(g, seed, [2]string{})
	return houses, mavericks, legacies
}

// ErrUnknownSet reports a set name no deck-generation set answers to. An empty
// name is not one: it means "the default set".
var ErrUnknownSet = errors.New("unknown deck-generation set")

// setFor resolves a chosen set name to its deck-generation Set. An empty name is
// the caller declining to choose and gets the base set; a name that resolves to
// nothing is a typo or a stale saved game, and silently dealing the default set
// would hide it behind a game that looks fine.
func setFor(name string) (deckgen.Set, error) {
	if name == "" {
		return cards.DeckSet(), nil
	}
	s, ok := cards.DeckSetNamed(name)
	if !ok {
		return deckgen.Set{}, fmt.Errorf("%w: %q", ErrUnknownSet, name)
	}
	return s, nil
}

// SetupDecksFor is SetupDecks with each player's set named explicitly, so a
// frontend can let the two players play different sets. It also returns each
// player's static deck Roster — the generated deck kept as a read-only list so a
// deck list shows the exact cards dealt, immune to later pool changes (ADR 0025).
func SetupDecksFor(
	g *engine.Game, seed int64, setNames [2]string,
) ([2][]engine.House, [2][]engine.LocalID, [2][]engine.LocalID, [2]Roster, error) {
	var houses [2][]engine.House
	var mavericks [2][]engine.LocalID
	var legacies [2][]engine.LocalID
	var rosters [2]Roster
	// A "name a card" choice (Etan's Jar) reaches the whole implemented card
	// database, which the engine cannot read itself (ADR 0003). Offering only the
	// cards in this match would show a player their opponent's deck list.
	all := cards.All()
	names := make([]string, len(all))
	for i := range all {
		names[i] = all[i].Name
	}
	g.SetNameableNames(names)
	for player := range 2 {
		set, err := setFor(setNames[player])
		if err != nil {
			return houses, mavericks, legacies, rosters, err
		}
		deck := deckgen.Generate(set, seed+int64(player)+1)
		houses[player] = deck.Houses()
		rosters[player] = rosterOf(deck)

		// The whole deck goes into the deck zone; engine.StartGame shuffles it and
		// deals the opening hands. Each card's Maverick and Legacy flags are pinned to
		// the LocalID the engine assigns on add, so the badge survives that shuffle.
		for i := range deck.Pods {
			for j := range deck.Pods[i].Slots {
				s := deck.Pods[i].Slots[j]
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
	return houses, mavericks, legacies, rosters, nil
}
