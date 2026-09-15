//go:build todo

package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// NovuDynamo
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  8
//	Armor:  2
//	Traits: Robot
//
//	At the start of your turn, you may discard a Logos card from your hand or archives. If you do, gain 1A. Otherwise, destroy Novu Dynamo.
var NovuDynamo = set.New(
	"Novu Dynamo",
	card.House.Logos,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.MM, "093"),
	card.WithPower(8),
	card.WithArmor(2),
	card.WithTraits(card.Traits.Robot),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
