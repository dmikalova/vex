//go:build todo

package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// BossZarek
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Shadows
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  3
//	Traits: Mutant • Thief
//
//	Enhance Capture Capture Capture.
//	Each friendly creature with A on it gains elusive.
var BossZarek = set.New(
	"Boss Zarek",
	card.House.Shadows,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.MM, "264"),
	card.WithPower(3),
	card.WithTraits(card.Traits.Mutant, card.Traits.Thief),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
