//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Ragnarok
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Brobnar
//	Type:   Tactic
//	Rarity: Rare
//
//	Alpha.
//	Play: For the remainder of the turn, creatures cannot reap and you gain 1A whenever a friendly creature fights. At the end of the turn, destroy each creature.
var Ragnarok = card.New(
	"Ragnarok",
	card.House.Brobnar,
	card.Type.Tactic,
	card.Rarity.Rare,
	card.Provenance(card.WC, "47"),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
