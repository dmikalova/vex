//go:build todo

package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Nexus
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Shadows
//	Type:   Creature
//	Rarity: Common
//	Power:  3
//	Traits: Cyborg • Thief
//
//	Elusive. (The first time this creature is attacked each turn, no damage is dealt.)
//	Reap: Use an opponent's artifact as if it were yours.
var Nexus = card.New(
	"Nexus",
	card.House.Shadows,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.CotA, 305),
	card.WithPower(3),
	card.WithTraits("Cyborg", "Thief"),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
