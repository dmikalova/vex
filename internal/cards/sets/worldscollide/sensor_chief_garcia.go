package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Sensor Chief Garcia
//
//	House:  Star Alliance
//	Type:   Creature
//	Rarity: Common
//	Power:  3
//	Traits: Human
//
//	Play/Fight/Reap: Keys cost +2 Æmber during your opponent's next turn.
var SensorChiefGarcia = card.New(
	"Sensor Chief Garcia",
	card.House.StarAlliance,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.WC, "305"),
	card.WithPower(3),
	card.WithTraits(card.Traits.Human),
	card.WithAbility(card.Trigger.PlayFightReap, card.RaiseKeyCost{
		Player:   card.Opponent,
		Amount:   2,
		Duration: card.Duration.OpponentNextTurn,
	}),
)
