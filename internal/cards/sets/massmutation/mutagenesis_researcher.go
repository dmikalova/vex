//go:build todo

package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// MutagenesisResearcher
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Staralliance
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  3
//	Traits: Mutant • Scientist
//
//	Enhance Aember Capture Damage Draw.
var MutagenesisResearcher = set.New(
	"Mutagenesis Researcher",
	card.House.Staralliance,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.MM, "325"),
	card.WithPower(3),
	card.WithTraits(card.Traits.Mutant, card.Traits.Scientist),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
