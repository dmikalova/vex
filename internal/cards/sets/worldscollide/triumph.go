//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Triumph
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Saurian
//	Type:   Tactic
//	Rarity: Rare
//	Æmber:  1
//
//	Play: If there are no enemy creatures, exalt each friendly creature. If you do and there are 6 or more friendly creatures, forge a key at no cost.
var Triumph = card.New(
	"Triumph",
	card.House.Saurian,
	card.Type.Tactic,
	card.Rarity.Rare,
	card.Provenance(card.WC, 234),
	card.WithAemberBonus(1),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
