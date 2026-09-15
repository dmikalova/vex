//go:build todo

package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// CeraSevera
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Saurian
//	Type:   Creature
//	Rarity: Common
//	Power:  4
//	Traits: Mutant
//
//	After Fight/After Reap: Capture 1 Aember.
//	Destroyed: Choose an enemy creature. Deal 1 Damage to that creature for each Aember on Cera Severa
var CeraSevera = set.New(
	"Cera Severa",
	card.House.Saurian,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.MoMu, "232"),
	card.WithPower(4),
	card.WithTraits(card.Traits.Mutant),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
