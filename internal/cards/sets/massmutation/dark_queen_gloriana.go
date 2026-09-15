//go:build todo

package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// DarkQueenGloriana
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Rare
//	Power:  3
//	Traits: Mutant
//
//	Enhance Aember Aember.
//	Play: Return a friendly non-Untamed creature to your hand.
var DarkQueenGloriana = set.New(
	"Dark Queen Gloriana",
	card.House.Untamed,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.MM, "397"),
	card.WithPower(3),
	card.WithTraits(card.Traits.Mutant),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
