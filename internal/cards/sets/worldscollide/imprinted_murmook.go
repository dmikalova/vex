//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// ImprintedMurmook
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Common
//	Power:  3
//	Traits: Beast
//
//	Elusive. (The first time this creature is attacked each turn, no damage is dealt.)
//	Your keys cost -1A.
var ImprintedMurmook = card.New(
	"Imprinted Murmook",
	card.House.Untamed,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.WC, 358),
	card.WithPower(3),
	card.WithTraits(card.Traits.Beast),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
