//go:build todo

package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Teliga
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Rare
//	Power:  3
//	Traits: Human • Witch
//
//	Each time your opponent plays a creature, gain 1 Aember.
var Teliga = card.New(
	"Teliga",
	card.House.Untamed,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 366),
	card.WithPower(3),
	card.WithTraits("Human", "Witch"),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
