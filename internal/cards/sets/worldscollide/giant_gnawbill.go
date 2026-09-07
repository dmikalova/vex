//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// GiantGnawbill
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Rare
//	Power:  5
//	Traits: Beast
//
//	After a player chooses an active house, that player destroys an artifact of that house.
var GiantGnawbill = card.New(
	"Giant Gnawbill",
	card.House.Untamed,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.WC, 390),
	card.WithPower(5),
	card.WithTraits(card.Traits.Beast),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
