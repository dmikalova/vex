//go:build todo

package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// FaustTheGreat
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Saurian
//	Type:   Creature
//	Rarity: Common
//	Power:  4
//	Traits: Dinosaur
//
//	Your opponent's keys cost +1A for each friendly creature with A on it.
//	Play: You may exalt a friendly creature.
var FaustTheGreat = set.New(
	"Faust the Great",
	card.House.Saurian,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.MM, "192"),
	card.WithPower(4),
	card.WithTraits(card.Traits.Dinosaur),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
