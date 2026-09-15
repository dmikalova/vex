//go:build todo

package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// SacroSaurus
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Saurian
//	Type:   Creature
//	Rarity: Special
//	Power:  6
//	Armor:  2
//	Traits: Mutant • Dinosaur
//
//	Play: You may exalt Sacro-Saurus. If you do, deal 3D to a creature.
var SacroSaurus = set.New(
	"Sacro-Saurus",
	card.House.Saurian,
	card.Type.Creature,
	// TODO(variant): rarity relabelled from Variant to Special — handle manually
	card.Rarity.Special,
	card.Provenance(card.MM, "240"),
	card.WithPower(6),
	card.WithArmor(2),
	card.WithTraits(card.Traits.Mutant, card.Traits.Dinosaur),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
