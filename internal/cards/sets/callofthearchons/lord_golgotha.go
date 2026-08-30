//go:build todo

package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// LordGolgotha
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Sanctum
//	Type:   Creature
//	Rarity: Rare
//	Power:  5
//	Armor:  2
//	Traits: Knight • Spirit
//
//	Before Fight: Deal 3 Damage to each neighbor of the creature Lord Golgotha fights.
var LordGolgotha = card.New(
	"Lord Golgotha",
	card.House.Sanctum,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 252),
	card.WithPower(5),
	card.WithArmor(2),
	card.WithTraits("Knight", "Spirit"),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
