//go:build todo

package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// TheVaultkeeper
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Sanctum
//	Type:   Creature
//	Rarity: Rare
//	Power:  4
//	Armor:  1
//	Traits: Knight • Spirit
//
//	Your Aember cannot be stolen.
var TheVaultkeeper = card.New(
	"The Vaultkeeper",
	card.House.Sanctum,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 261),
	card.WithPower(4),
	card.WithArmor(1),
	card.WithTraits("Knight", "Spirit"),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
