//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// SensorChiefGarcia
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
//	Play/Fight/Reap: Keys cost +2A during your opponent's next turn.
var SensorChiefGarcia = card.New(
	"Sensor Chief Garcia",
	card.House.Staralliance,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.WC, 305),
	card.WithPower(3),
	card.WithTraits(card.Traits.Human),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
