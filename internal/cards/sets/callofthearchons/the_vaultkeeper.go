package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// The Vaultkeeper
//
//	House:  Sanctum
//	Type:   Creature
//	Rarity: Rare
//	Power:  4
//	Traits: Knight • Spirit
//
//	Your Æmber cannot be stolen.
var TheVaultkeeper = card.New(
	"The Vaultkeeper",
	card.House.Sanctum,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 261),
	card.WithPower(4),
	card.WithTraits(card.Traits.Knight, card.Traits.Spirit),
	card.WithAemberTheftImmunity(),
)
