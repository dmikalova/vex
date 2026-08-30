//go:build todo

package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Tolas
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Rare
//	Power:  1
//	Traits: Imp
//
//	Elusive. (The first time this creature is attacked each turn, no damage is dealt.)
//	Each time a creature is destroyed, its opponent gains 1 Aember.
var Tolas = card.New(
	"Tolas",
	card.House.Dis,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 103),
	card.WithPower(1),
	card.WithTraits("Imp"),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
