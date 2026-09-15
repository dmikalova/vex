//go:build todo

package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// XenoSaurus
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
//	Play: You may exalt Xeno-Saurus. If you do, deal 3D to a creature.
//	Fight: Look at the top 3 cards of your deck. Put 1 into your hand and 1 on the bottom of your deck.
var XenoSaurus = set.New(
	"Xeno-Saurus",
	card.House.Saurian,
	card.Type.Creature,
	// TODO(variant): rarity relabelled from Variant to Special — handle manually
	card.Rarity.Special,
	card.Provenance(card.MM, "243"),
	card.WithPower(5),
	card.WithTraits(card.Traits.Mutant, card.Traits.Dinosaur),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
