//go:build todo

package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Stealthster
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Staralliance
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  3
//	Traits: Robot
//
//	Elusive.
//	Stealthster may be played as an upgrade instead of a creature, with the text: "This creature gains elusive."
var Stealthster = set.New(
	"Stealthster",
	card.House.Staralliance,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.MM, "329"),
	card.WithPower(3),
	card.WithTraits(card.Traits.Robot),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
