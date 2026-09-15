//go:build todo

package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// DinoBeast
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Special
//	Power:  5
//	Traits: Mutant • Beast
//
//	Skirmish.
//	Play: You may exalt Dino-Beast. If you do, deal 3D to a creature.
var DinoBeast = set.New(
	"Dino-Beast",
	card.House.Untamed,
	card.Type.Creature,
	// TODO(variant): rarity relabelled from Variant to Special — handle manually
	card.Rarity.Special,
	card.Provenance(card.MM, "417"),
	card.WithPower(5),
	card.WithTraits(card.Traits.Mutant, card.Traits.Beast),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
