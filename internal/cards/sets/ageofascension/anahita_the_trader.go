//go:build todo

package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// AnahitaTheTrader
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Sanctum
//	Type:   Creature
//	Rarity: Rare
//	Power:  2
//	Traits: Human • Merchant
//
//	Reap: Give control of a friendly artifact to your opponent. If you do, they must give you 2A.
var AnahitaTheTrader = card.New(
	"Anahita the Trader",
	card.House.Sanctum,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.AoA, 248),
	card.WithPower(2),
	card.WithTraits(card.Traits.Human, card.Traits.Merchant),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
