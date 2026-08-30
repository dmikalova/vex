package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Punch
//
//	House:  Brobnar
//	Type:   Tactic
//	Rarity: Common
//	Æmber:  1
//
//	Play: Deal 3 damage to a creature.
var Punch = card.New(
	"Punch",
	card.House.Brobnar,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.CotA, 12),
	card.WithAemberBonus(1),
	card.WithAbility(
		card.Trigger.Play, card.DealDamage{
			Amount: 3,
			Target: card.Target.Creature,
		}),
)
