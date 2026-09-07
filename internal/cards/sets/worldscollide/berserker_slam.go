package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Berserker Slam
//
//	House:  Brobnar
//	Type:   Tactic
//	Rarity: Common
//	Æmber:  1
//
//	Play: Deal 4 damage to a flank creature. If this damage destroys that creature, its controller loses 1 Æmber.
var BerserkerSlam = card.New(
	"Berserker Slam",
	card.House.Brobnar,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.WC, 5),
	card.WithAemberBonus(1),
	card.WithAbility(
		card.Trigger.Play, card.DamageThenIfDestroyed{
			Amount: 4,
			Target: card.Target.Creature.OnFlank(),
			Then:   card.LoseAember{Player: card.ItsOwner, Amount: 1},
		}),
)
