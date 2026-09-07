//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Snag
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  5
//	Traits: Demon
//
//	Fight: Your opponent must choose the house of the creature Snag fights as their active house on their next turn.
var Snag = card.New(
	"Snag",
	card.House.Dis,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.WC, 96),
	card.WithPower(5),
	card.WithTraits(card.Traits.Demon),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
