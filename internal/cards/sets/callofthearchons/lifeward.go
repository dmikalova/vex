//go:build todo

package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Lifeward
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Dis
//	Type:   Artifact
//	Rarity: Uncommon
//	Æmber:  1
//	Traits: Power
//
//	Omni: Sacrifice Lifeward. Your opponent cannot play creatures on their next turn.
var Lifeward = card.New(
	"Lifeward",
	card.House.Dis,
	card.Type.Artifact,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, 77),
	card.WithAemberBonus(1),
	card.WithTraits("Power"),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
