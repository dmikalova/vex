//go:build todo

package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// SafePlace
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Shadows
//	Type:   Artifact
//	Rarity: Rare
//	Æmber:  1
//	Traits: Location
//
//	You may spend Aember on Safe Place when forging keys.
//	Action: Move 1 Aember from your pool to Safe Place.
var SafePlace = card.New(
	"Safe Place",
	card.House.Shadows,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 289),
	card.WithAemberBonus(1),
	card.WithTraits("Location"),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
