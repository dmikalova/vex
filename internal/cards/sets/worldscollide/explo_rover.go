//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// ExploRover
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Staralliance
//	Type:   Creature
//	Rarity: Common
//	Power:  3
//	Traits: Robot
//
//	Skirmish.
//	Explo-rover may be played as an upgrade instead of a creature, with the text: "This creature gains skirmish."
var ExploRover = card.New(
	"Explo-rover",
	card.House.Staralliance,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.WC, 297),
	card.WithPower(3),
	card.WithTraits(card.Traits.Robot),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
