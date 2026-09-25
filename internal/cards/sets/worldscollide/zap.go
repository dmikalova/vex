package worldscollide

import "github.com/dmikalova/vex/internal/card"

// Zap
//
//	House:  Star Alliance
//	Type:   Tactic
//	Rarity: Common
//	Bonus:  Æmber
//
//	Play: For each house represented among creatures in play, deal 1 damage to a creature.
var Zap = set.New(
	"Zap",
	card.House.StarAlliance,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.WC, "307"),
	card.WithBonus(card.Bonus.Aember),
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
