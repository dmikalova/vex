//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// ComOfficerKirby
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Staralliance
//	Type:   Creature
//	Rarity: Common
//	Power:  3
//	Traits: Human
//
//	Play/Fight/Reap: You may play a non-Star Alliance artifact, upgrade, or action card this turn.
var ComOfficerKirby = card.New(
	"Com. Officer Kirby",
	card.House.Staralliance,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.WC, "295"),
	card.WithPower(3),
	card.WithTraits(card.Traits.Human),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
