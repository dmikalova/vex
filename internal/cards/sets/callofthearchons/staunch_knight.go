//go:build todo

package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// StaunchKnight
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Sanctum
//	Type:   Creature
//	Rarity: Common
//	Power:  4
//	Armor:  2
//	Traits: Human • Knight
//
//	Staunch Knight gets +2 power while it is on a flank.
var StaunchKnight = card.New(
	"Staunch Knight",
	card.House.Sanctum,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.CotA, 259),
	card.WithPower(4),
	card.WithArmor(2),
	card.WithTraits("Human", "Knight"),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
