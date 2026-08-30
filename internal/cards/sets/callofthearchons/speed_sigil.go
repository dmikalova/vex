//go:build todo

package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// SpeedSigil
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Shadows
//	Type:   Artifact
//	Rarity: Uncommon
//	Æmber:  1
//	Traits: Power
//
//	The first creature played each turn enters play ready.
var SpeedSigil = card.New(
	"Speed Sigil",
	card.House.Shadows,
	card.Type.Artifact,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, 293),
	card.WithAemberBonus(1),
	card.WithTraits("Power"),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
