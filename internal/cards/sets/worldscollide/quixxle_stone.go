//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// QuixxleStone
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Staralliance
//	Type:   Artifact
//	Rarity: Rare
//	Æmber:  1
//	Traits: Item
//
//	If a player has more creatures in play than their opponent, they cannot play creatures.
var QuixxleStone = card.New(
	"Quixxle Stone",
	card.House.Staralliance,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.WC, 338),
	card.WithAemberBonus(1),
	card.WithTraits(card.Traits.Item),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
