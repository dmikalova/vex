//go:build todo

package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// TheSting
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Shadows
//	Type:   Artifact
//	Rarity: Rare
//	Æmber:  1
//	Traits: Vehicle
//
//	Skip your "forge a key" step.
//	You get all Aember spent by your opponent when forging keys.
//	Action: Sacrifice The Sting.
var TheSting = card.New(
	"The Sting",
	card.House.Shadows,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 295),
	card.WithAemberBonus(1),
	card.WithTraits(card.Traits.Vehicle),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
