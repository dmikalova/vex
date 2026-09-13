package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Mars Needs Aember
//
//	House:  Mars
//	Type:   Tactic
//	Rarity: Common
//	Æmber:  1
//
//	Play: Each enemy damaged non-Mars Creature captures 1 Æmber from your opponent.
var MarsNeedsAember = card.New(
	"Mars Needs Aember",
	card.House.Mars,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.AoA, "166"),
	card.WithAemberBonus(1),
	card.WithAbility(
		card.Trigger.Play, card.CaptureAember{
			Amount: 1,
			Target: card.Target.EachEnemyCreature.Damaged().ExceptHouse(card.House.Self),
			Source: card.Opponent,
		}),
)
