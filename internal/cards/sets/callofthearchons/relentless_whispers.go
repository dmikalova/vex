package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Relentless Whispers
//
//	House:  Shadows
//	Type:   Tactic
//	Rarity: Common
//	Æmber:  1
//
//	Play: Deal 2 damage to a Creature. If this damage destroys that Creature, steal 1 Æmber.
var RelentlessWhispers = card.New(
	"Relentless Whispers",
	card.House.Shadows,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.CotA, "281"),
	card.WithAemberBonus(1),
	card.WithAbility(
		card.Trigger.Play, card.DamageThen{
			Amount: 2,
			After:  card.IfDestroyed,
			Target: card.Target.Creature,
			Then:   card.StealAember{Amount: 1},
		}),
)
