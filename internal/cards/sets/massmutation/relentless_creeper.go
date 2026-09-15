//go:build todo

package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// RelentlessCreeper
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  2
//	Traits: Imp
//
//	After you choose Dis as your active house, you may return Relentless Creeper from your discard pile to your hand.
var RelentlessCreeper = set.New(
	"Relentless Creeper",
	card.House.Dis,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.MM, "029"),
	card.WithPower(2),
	card.WithTraits(card.Traits.Imp),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
