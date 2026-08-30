//go:build todo

package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// ScreamingCave
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Dis
//	Type:   Artifact
//	Rarity: Rare
//	Traits: Location
//
//	Action: Shuffle your hand and discard pile into your deck.
var ScreamingCave = card.New(
	"Screaming Cave",
	card.House.Dis,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 79),
	card.WithTraits("Location"),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
