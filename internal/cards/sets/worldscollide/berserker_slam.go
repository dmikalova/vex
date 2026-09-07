//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// BerserkerSlam
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Brobnar
//	Type:   Tactic
//	Rarity: Common
//	Æmber:  1
//
//	Play: Deal 4D to a flank creature. If this damage destroys that creature, its controller loses 1A.
var BerserkerSlam = card.New(
	"Berserker Slam",
	card.House.Brobnar,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.WC, 5),
	card.WithAemberBonus(1),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
