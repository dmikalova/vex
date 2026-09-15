//go:build todo

package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// SpareArmCarmine
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Shadows
//	Type:   Creature
//	Rarity: Common
//	Power:  2
//	Traits: Mutant
//
//	After Reap: If there are more friendly Mutant creatures than enemy Mutant creatures, steal 2 Aember. Otherwise, steal 1 Aember.
var SpareArmCarmine = set.New(
	"Spare Arm Carmine",
	card.House.Shadows,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.MoMu, "307"),
	card.WithPower(2),
	card.WithTraits(card.Traits.Mutant),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
