//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Stomp
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Saurian
//	Type:   Tactic
//	Rarity: Uncommon
//	Æmber:  1
//
//	Play: Deal 5D to a creature. If this damage destroys that creature, exalt a friendly creature.
var Stomp = card.New(
	"Stomp",
	card.House.Saurian,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.WC, 210),
	card.WithAemberBonus(1),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
