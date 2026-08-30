//go:build todo

package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Gatekeeper
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Sanctum
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  5
//	Armor:  1
//	Traits: Knight • Spirit
//
//	Play: If your opponent has 7 or more Aember, capture all but 5 of it.
var Gatekeeper = card.New(
	"Gatekeeper",
	card.House.Sanctum,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, 260),
	card.WithPower(5),
	card.WithArmor(1),
	card.WithTraits("Knight", "Spirit"),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
