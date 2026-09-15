//go:build todo

package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// ArmoryOfficerNel
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Staralliance
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  4
//	Traits: Alien
//
//	Enhance Draw.
//	After an upgrade enters play, draw a card.
var ArmoryOfficerNel = set.New(
	"Armory Officer Nel",
	card.House.Staralliance,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.MM, "319"),
	card.WithPower(4),
	card.WithTraits(card.Traits.Alien),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
