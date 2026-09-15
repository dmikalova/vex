//go:build todo

package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Gladiodontus
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Saurian
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  15
//	Traits: Mutant
//
//	Gladiodontus enters play stunned.
//	Gladiodontus only deals 5D when fighting.
//	Fight/Reap: If this is the first time Gladiodontus has been used this turn, ready and enrage it.
var Gladiodontus = set.New(
	"Gladiodontus",
	card.House.Saurian,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.MM, "206"),
	card.WithPower(15),
	card.WithTraits(card.Traits.Mutant),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
