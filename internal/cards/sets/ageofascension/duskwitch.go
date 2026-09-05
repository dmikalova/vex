//go:build todo

package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Duskwitch
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Common
//	Power:  1
//	Traits: Human • Witch
//
//	Omega. (After you play this card, end this step.)
//	Elusive. (The first time this creature is attacked each turn, no damage is dealt.)
//	Your creatures enter play ready.
var Duskwitch = card.New(
	"Duskwitch",
	card.House.Untamed,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.AoA, 320),
	card.WithPower(1),
	card.WithTraits(card.Traits.Human, card.Traits.Witch),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
