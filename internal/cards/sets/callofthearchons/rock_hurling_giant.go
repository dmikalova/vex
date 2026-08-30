//go:build todo

package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// RockHurlingGiant
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Brobnar
//	Type:   Creature
//	Rarity: Rare
//	Power:  6
//	Traits: Giant
//
//	During your turn, each time you discard a Brobnar card from your hand, you may deal 4 Damage to a creature.
var RockHurlingGiant = card.New(
	"Rock-Hurling Giant",
	card.House.Brobnar,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 44),
	card.WithPower(6),
	card.WithTraits("Giant"),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
