//go:build todo

package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Masterplan
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Shadows
//	Type:   Artifact
//	Rarity: Rare
//	Æmber:  1
//	Traits: Item
//
//	Play: Put a card from your hand facedown beneath Masterplan.
//	Omni: Play the card beneath Masterplan. Sacrifice Masterplan.
var Masterplan = card.New(
	"Masterplan",
	card.House.Shadows,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 288),
	card.WithAemberBonus(1),
	card.WithTraits("Item"),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
