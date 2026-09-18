package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Relentless Whispers
//
//	House:  Shadows
//	Type:   Tactic
//	Rarity: Common
//	Bonus:  Æmber
//
//	Play: Deal 2 damage to a creature. If this damage destroys that creature, steal 1 Æmber.
var RelentlessWhispers = set.New(
	"Relentless Whispers",
	card.House.Shadows,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.CotA, "281"),
	card.WithBonus(card.Bonus.Aember),
	card.WithAbility(
		card.Trigger.Play, card.DealDamage{
			Amount: 2,
			After:  card.IfDestroyed,
			Target: card.Target.Creature,
			Then:   card.StealAember{Amount: 1},
		}),
)
