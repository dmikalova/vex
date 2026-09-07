//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// FirstOfficerFrane
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Staralliance
//	Type:   Creature
//	Rarity: Common
//	Power:  4
//	Traits: Human
//
//	Play/Fight/Reap: A friendly creature captures 1A.
var FirstOfficerFrane = card.New(
	"First Officer Frane",
	card.House.Staralliance,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.WC, 298),
	card.WithPower(4),
	card.WithTraits(card.Traits.Human),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
