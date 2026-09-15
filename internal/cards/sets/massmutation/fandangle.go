//go:build todo

package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Fandangle
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Common
//	Power:  3
//	Traits: Mutant • Witch
//
//	While you have 4A or more, your non-Untamed creatures enter play ready.
var Fandangle = set.New(
	"Fandangle",
	card.House.Untamed,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.MM, "365"),
	card.WithPower(3),
	card.WithTraits(card.Traits.Mutant, card.Traits.Witch),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
