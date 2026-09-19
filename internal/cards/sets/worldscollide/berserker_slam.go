package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Berserker Slam
//
//	House:  Brobnar
//	Type:   Tactic
//	Rarity: Common
//	Bonus:  Æmber
//
//	Play: Deal 4 damage to a flank creature. If this damage destroys that creature, its controller loses 1 Æmber.
var BerserkerSlam = set.New(
	"Berserker Slam",
	card.House.Brobnar,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.WC, "5"),
	card.WithBonus(card.Bonus.Aember),
	card.WithAbility(
		card.Trigger.Play, card.DealDamage{
			Amount: 4,
			After:  card.IfDestroyed,
			Target: card.Target.Creature.OnFlank(),
			Then: card.LoseAember{
				Player: card.ItsOwner,
				Amount: 1,
			},
		}),
)
