// Package cards is the card database aggregator. It blank-imports every released
// set package so their cards self-register (each card calls card.Register from
// its own init), then exposes the assembled database to the rest of the program.
//
// It lives outside package engine on purpose: set packages import the engine (via
// the card facade), and the engine must not import them back. Adding a card is
// just adding a self-registering file to a set; adding a set is one blank import
// below — neither requires listing individual cards here.
package cards

import (
	"github.com/dmikalova/vactrol/internal/card"
	// Blank-imported so each set's cards self-register through its package init.
	_ "github.com/dmikalova/vactrol/internal/cards/sets/callofthearchons"
	"github.com/dmikalova/vactrol/internal/deckgen"
)

// All returns every registered card across every imported set.
func All() []card.Definition {
	return card.Registered()
}

// DeckSet assembles the deck-generation Set from the registered cards, carrying
// each card's generation profile and template materializer. With a single
// released set every card belongs to it; when a second set is added this will
// build one deckgen.Set per set package.
func DeckSet() deckgen.Set {
	regs := card.Cards()
	cs := make([]deckgen.Card, 0, len(regs))
	for _, rc := range regs {
		cs = append(
			cs,
			deckgen.Card{Def: rc.Def, Profile: rc.Profile, Materializer: rc.Materializer},
		)
	}
	return deckgen.NewSet("Call of the Archons", cs, deckgen.DefaultTuning())
}

// Set is a named group of cards, used for reporting such as the Card statistics
// view.
type Set struct {
	Name  string
	Cards []card.Definition
}

// Sets returns the card database grouped by set. Cards do not yet carry a set
// tag, so every registered card is reported under the single released set; when
// a second set is added, this must switch to grouping by a per-card set field.
func Sets() []Set {
	return []Set{{Name: "Call of the Archons", Cards: All()}}
}
