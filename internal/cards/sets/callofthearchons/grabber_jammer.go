package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Grabber Jammer
//
//	House:  Mars
//	Type:   Creature
//	Rarity: Common
//	Power:  4
//	Armor:  1
//	Traits: Robot
//
//	Your opponent's keys cost +1 Æmber.
//	Fight/Reap: Grabber Jammer captures 1 Æmber from your opponent.
var GrabberJammer = card.New(
	"Grabber Jammer",
	card.House.Mars,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.CotA, 193),
	card.WithPower(4),
	card.WithArmor(1),
	card.WithTraits(card.Traits.Robot),
	card.WithKeyCost(card.KeyCostChange(card.Opponent, 1)),
	card.WithFightOrReap(card.CaptureAember{
		Amount: 1,
		Target: card.Target.This,
		Source: card.Opponent,
	}),
)
