//go:build todo

package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// TheWarchest
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Brobnar
//	Type:   Artifact
//	Rarity: Uncommon
//	Traits: Item
//
//	Action: Gain 1 Aember for each enemy creature that was destroyed in a fight this turn.
var TheWarchest = card.New(
	"The Warchest",
	card.House.Brobnar,
	card.Type.Artifact,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, 27),
	card.WithTraits("Item"),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
