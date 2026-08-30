//go:build todo

package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// SigilOfBrotherhood
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Sanctum
//	Type:   Artifact
//	Rarity: Uncommon
//	Æmber:  1
//	Traits: Power
//
//	Omni: Sacrifice Sigil of Brotherhood. For the remainder of the turn, you may use friendly Sanctum creatures.
var SigilOfBrotherhood = card.New(
	"Sigil of Brotherhood",
	card.House.Sanctum,
	card.Type.Artifact,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, 236),
	card.WithAemberBonus(1),
	card.WithTraits("Power"),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
