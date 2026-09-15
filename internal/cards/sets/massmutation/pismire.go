//go:build todo

package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Pismire
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Common
//	Power:  4
//	Traits: Mutant
//
//	While there are more friendly Mutant creatures than enemy Mutant creatures, your opponent's keys cost +2A.
var Pismire = set.New(
	"Pismire",
	card.House.Untamed,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.MM, "372"),
	card.WithPower(4),
	card.WithTraits(card.Traits.Mutant),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
