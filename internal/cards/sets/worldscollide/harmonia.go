//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Harmonia
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Common
//	Power:  2
//	Traits: Human • Witch
//
//	Elusive. (The first time this creature is attacked each turn, no damage is dealt.)
//	After you play a creature, if there are more enemy creatures than friendly creatures, gain 1A.
var Harmonia = card.New(
	"Harmonia",
	card.House.Untamed,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.WC, 357),
	card.WithPower(2),
	card.WithTraits(card.Traits.Human, card.Traits.Witch),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
