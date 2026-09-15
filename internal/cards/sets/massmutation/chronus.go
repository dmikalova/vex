//go:build todo

package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Chronus
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  3
//	Traits: Mutant
//
//	Enhance Draw Draw.
//	After you resolve a R bonus icon, you may archive a card.
var Chronus = set.New(
	"Chronus",
	card.House.Logos,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.MM, "084"),
	card.WithPower(3),
	card.WithTraits(card.Traits.Mutant),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
