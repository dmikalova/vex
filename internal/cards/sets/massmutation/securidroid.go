//go:build todo

package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// SecuriDroid
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Staralliance
//	Type:   Creature
//	Rarity: Common
//	Power:  4
//	Traits: Robot
//
//	Taunt.
//	Securi-Droid may be played as an upgrade instead of a creature, with the text: "This creature gains taunt."
var SecuriDroid = set.New(
	"Securi-Droid",
	card.House.Staralliance,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.MM, "312"),
	card.WithPower(4),
	card.WithTraits(card.Traits.Robot),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
