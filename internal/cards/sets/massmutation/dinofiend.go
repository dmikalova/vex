//go:build todo

package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// DinoFiend
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Special
//	Power:  4
//	Traits: Mutant • Demon
//
//	Play: You may exalt Dino-Fiend. If you do, deal 3D to a creature.
//	Destroyed: Steal 1A.
var DinoFiend = set.New(
	"Dino-Fiend",
	card.House.Dis,
	card.Type.Creature,
	// TODO(variant): rarity relabelled from Variant to Special — handle manually
	card.Rarity.Special,
	card.Provenance(card.MM, "055"),
	card.WithPower(4),
	card.WithTraits(card.Traits.Mutant, card.Traits.Demon),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
