//go:build todo

package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// SubjectKirby
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Staralliance
//	Type:   Creature
//	Rarity: Common
//	Power:  2
//	Traits: Mutant
//
//	Play/Fight/Reap: You may play a non-Star Alliance creature this turn.
var SubjectKirby = set.New(
	"Subject Kirby",
	card.House.Staralliance,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.MM, "315"),
	card.WithPower(2),
	card.WithTraits(card.Traits.Mutant),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
