//go:build todo

package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// CrewmanJorg
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Staralliance
//	Type:   Creature
//	Rarity: Common
//	Power:  3
//	Traits: Human • Thief
//
//	Enhance Capture.
//	Action: If Crewman Jorg has no Star Alliance neighbors, steal 1A.
var CrewmanJorg = set.New(
	"Crewman Jorg",
	card.House.Staralliance,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.MM, "305"),
	card.WithPower(3),
	card.WithTraits(card.Traits.Human, card.Traits.Thief),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
