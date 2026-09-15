//go:build todo

package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// DinoKnight
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Sanctum
//	Type:   Creature
//	Rarity: Special
//	Power:  6
//	Armor:  2
//	Traits: Mutant • Knight
//
//	Play: You may exalt Dino-Knight. If you do, deal 3D to a creature.
var DinoKnight = set.New(
	"Dino-Knight",
	card.House.Sanctum,
	card.Type.Creature,
	// TODO(variant): rarity relabelled from Variant to Special — handle manually
	card.Rarity.Special,
	card.Provenance(card.MM, "178"),
	card.WithPower(6),
	card.WithArmor(2),
	card.WithTraits(card.Traits.Mutant, card.Traits.Knight),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
