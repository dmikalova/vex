package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Zap
//
//	House:  Star Alliance
//	Type:   Tactic
//	Rarity: Common
//	Æmber:  1
//
//	Play: For each house represented among Creatures in play, deal 1 damage to a Creature.
var Zap = card.New(
	"Zap",
	card.House.StarAlliance,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.WC, "307"),
	card.WithAemberBonus(1),
	card.WithAbility(
		card.Trigger.Play, card.DealDamage{
			Amount: 1,
			Per: card.HousesAmong{
				Player: card.EachPlayer,
				Type:   card.Type.Creature,
			},
			Target: card.Target.Creature,
		}),
)
