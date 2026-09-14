package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Phalanx Strike
//
//	House:  Saurian
//	Type:   Tactic
//	Rarity: Common
//	Æmber:  1
//
//	Play: For each friendly Creature in play, deal 1 damage to a Creature. You may exalt a friendly Creature to repeat the preceding effect.
var PhalanxStrike = set.New(
	"Phalanx Strike",
	card.House.Saurian,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.WC, "189"),
	card.WithAemberBonus(1),
	card.WithAbility(
		card.Trigger.Play, card.Repeat{
			Do: card.DealDamage{
				Amount: 1,
				Per: card.InPlay{
					Player: card.Controller,
					Type:   card.Type.Creature,
				},
				Target: card.Target.Creature,
			},
			Gate: card.ByExalting{Creature: card.Target.FriendlyCreature},
		}),
)
