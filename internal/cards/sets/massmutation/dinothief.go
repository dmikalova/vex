//go:build todo

package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// DinoThief
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Shadows
//	Type:   Creature
//	Rarity: Special
//	Power:  4
//	Traits: Mutant • Thief
//
//	Elusive.
//	Play: You may exalt Dino-Thief. If you do, deal 3D to a creature.
var DinoThief = set.New(
	"Dino-Thief",
	card.House.Shadows,
	card.Type.Creature,
	// TODO(variant): rarity relabelled from Variant to Special — handle manually
	card.Rarity.Special,
	card.Provenance(card.MM, "297"),
	card.WithPower(4),
	card.WithTraits(card.Traits.Mutant, card.Traits.Thief),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
