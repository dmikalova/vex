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
)

// All returns every registered card across every imported set.
func All() []card.Definition {
	return card.Registered()
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
