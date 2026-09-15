//go:build todo

package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// TitanEngineer
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Common
//	Power:  6
//	Traits: Cyborg • Scientist
//
//	While Titan Engineer is not on a flank, keys cost +1 A.
var TitanEngineer = set.New(
	"Titan Engineer",
	card.House.Logos,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.MM, "081"),
	card.WithPower(6),
	card.WithTraits(card.Traits.Cyborg, card.Traits.Scientist),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
