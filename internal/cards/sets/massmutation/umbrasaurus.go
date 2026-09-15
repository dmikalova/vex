//go:build todo

package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// UmbraSaurus
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Saurian
//	Type:   Creature
//	Rarity: Special
//	Power:  4
//	Traits: Mutant • Dinosaur
//
//	Elusive.
//	Play: You may exalt Umbra-Saurus. If you do, deal 3D to a creature.
var UmbraSaurus = set.New(
	"Umbra-Saurus",
	card.House.Saurian,
	card.Type.Creature,
	// TODO(variant): rarity relabelled from Variant to Special — handle manually
	card.Rarity.Special,
	card.Provenance(card.MM, "242"),
	card.WithPower(4),
	card.WithTraits(card.Traits.Mutant, card.Traits.Dinosaur),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
