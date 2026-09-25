package ageofascension

import "github.com/dmikalova/vex/internal/card"

// Mars Needs Aember
//
//	House:  Mars
//	Type:   Tactic
//	Rarity: Common
//	Bonus:  Æmber
//
//	Play: Each enemy damaged non-Mars creature captures 1 Æmber from your opponent.
var MarsNeedsAember = set.New(
	"Mars Needs Aember",
	card.House.Mars,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.AoA, "166"),
	card.WithBonus(card.Bonus.Aember),
	card.WithAbility(
		card.Trigger.Play, card.CaptureAember{
			Amount: 1,
			Target: card.Target.EachEnemyCreature.Damaged().
				House(card.Houses.Except(card.House.Self)),
			Source: card.Opponent,
		}),
)
