//go:build todo

package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Rockatiel
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Common
//	Power:  2
//	Traits: Beast • Mutant
//
//	Elusive. Hazardous 1.
//	Play/After Reap: Choose up to 2 creatures. Shuffle each chosen creature into its owner's deck.
var Rockatiel = set.New(
	"Rockatiel",
	card.House.Untamed,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.MoMu, "429"),
	card.WithPower(2),
	card.WithTraits(card.Traits.Beast, card.Traits.Mutant),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
