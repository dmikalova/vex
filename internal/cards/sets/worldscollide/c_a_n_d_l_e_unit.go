package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// C.A.N.D.L.E. Unit
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  5
//	Armor:  1
//	Traits: Robot
//
//	After an enemy Creature reaps, draw a card.
//	Action: C.A.N.D.L.E. Unit captures 1 Æmber from your opponent.
var CANDLEUnit = set.New(
	"C.A.N.D.L.E. Unit",
	card.House.Logos,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.WC, "146"),
	card.WithPower(5),
	card.WithArmor(1),
	card.WithTraits(card.Traits.Robot),
	card.WithAbility(
		card.Trigger.AfterEnemyCreatureReaps, card.Draw{Amount: 1}),
	card.WithAbility(
		card.Trigger.Action, card.CaptureAember{
			Amount: 1,
			Target: card.Target.This,
			Source: card.Opponent,
		}),
)
