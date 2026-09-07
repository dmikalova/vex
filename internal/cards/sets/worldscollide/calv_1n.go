//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// CALV1N
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Staralliance
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  2
//	Armor:  1
//	Traits: Robot
//
//	Fight/Reap: Draw a card.
//	CALV-1N may be played as an upgrade instead of a creature, with the text: "This creature gains, 'Fight/Reap: Draw a card.'"
var CALV1N = card.New(
	"CALV-1N",
	card.House.Staralliance,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.WC, "308"),
	card.WithPower(2),
	card.WithArmor(1),
	card.WithTraits(card.Traits.Robot),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
