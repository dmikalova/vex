//go:build todo

package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// FeedingPit
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Mars
//	Type:   Artifact
//	Rarity: Rare
//	Traits: Location
//
//	Action: Discard a creature from your hand. If you do, gain 1 Aember.
var FeedingPit = card.New(
	"Feeding Pit",
	card.House.Mars,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 184),
	card.WithTraits("Location"),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
