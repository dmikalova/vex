//go:build todo

package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Protectrix
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Sanctum
//	Type:   Creature
//	Rarity: Common
//	Power:  5
//	Traits: Knight • Spirit
//
//	Reap: You may fully heal a creature. If you do, that creature cannot be dealt damage for the remainder of the turn.
var Protectrix = card.New(
	"Protectrix",
	card.House.Sanctum,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.CotA, 254),
	card.WithPower(5),
	card.WithTraits("Knight", "Spirit"),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
