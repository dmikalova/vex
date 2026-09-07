//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// BrambleLynx
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Common
//	Power:  3
//	Traits: Beast
//
//	Skirmish. (When you use this creature to fight, it is dealt no damage in return.)
//	If you have used a creature to reap this turn, Bramble Lynx enters play ready.
var BrambleLynx = card.New(
	"Bramble Lynx",
	card.House.Untamed,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.WC, 353),
	card.WithPower(3),
	card.WithTraits(card.Traits.Beast),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
