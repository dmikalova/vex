//go:build todo

package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// DarkCenturion
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Saurian
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  5
//	Traits: Mutant • Soldier
//
//	Enhance Capture Capture.
//	Action: Move 1A from a creature to the common supply. If you do, ward that creature.
var DarkCenturion = set.New(
	"Dark Centurion",
	card.House.Saurian,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.MM, "203"),
	card.WithPower(5),
	card.WithTraits(card.Traits.Mutant, card.Traits.Soldier),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
