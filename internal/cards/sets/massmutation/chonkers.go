//go:build todo

package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Chonkers
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Rare
//	Power:  1
//	Traits: Mutant
//
//	Skirmish.
//	After an enemy creature is destroyed fighting Chonkers, double the number of +1 power counters on Chonkers.
//	Play: Give Chonkers a +1 power counter.
var Chonkers = set.New(
	"Chonkers",
	card.House.Untamed,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.MM, "396"),
	card.WithPower(1),
	card.WithTraits(card.Traits.Mutant),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
