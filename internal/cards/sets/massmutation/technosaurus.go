//go:build todo

package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// TechnoSaurus
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Saurian
//	Type:   Creature
//	Rarity: Special
//	Power:  5
//	Traits: Mutant • Dinosaur
//
//	Play: You may exalt Techno-Saurus. If you do, deal 3D to a creature.
//	Reap: Discard a card from your hand. If you do, draw a card.
var TechnoSaurus = set.New(
	"Techno-Saurus",
	card.House.Saurian,
	card.Type.Creature,
	// TODO(variant): rarity relabelled from Variant to Special — handle manually
	card.Rarity.Special,
	card.Provenance(card.MM, "241"),
	card.WithPower(5),
	card.WithTraits(card.Traits.Mutant, card.Traits.Dinosaur),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
