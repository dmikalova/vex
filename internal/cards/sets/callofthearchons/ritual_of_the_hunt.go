//go:build todo

package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// RitualOfTheHunt
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Untamed
//	Type:   Artifact
//	Rarity: Rare
//	Æmber:  1
//	Traits: Power
//
//	Omni: Sacrifice Ritual of the Hunt. For the remainder of the turn, you may use friendly Untamed creatures.
var RitualOfTheHunt = card.New(
	"Ritual of the Hunt",
	card.House.Untamed,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 343),
	card.WithAemberBonus(1),
	card.WithTraits("Power"),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
