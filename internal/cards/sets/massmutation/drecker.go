//go:build todo

package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Drecker
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Common
//	Power:  4
//	Traits: Imp
//
//	Damage dealt to Drecker's neighbors during fights is also dealt to Drecker.
//	Reap: Steal 1A.
var Drecker = set.New(
	"Drecker",
	card.House.Dis,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.MM, "006"),
	card.WithPower(4),
	card.WithTraits(card.Traits.Imp),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
