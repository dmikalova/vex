package worldscollide

import "github.com/dmikalova/vex/internal/card"

// Phalanx Strike
//
//	House:  Saurian
//	Type:   Tactic
//	Rarity: Common
//	Bonus:  Æmber
//
//	Play: For each friendly creature in play, deal 1 damage to a creature. You may exalt a friendly creature to repeat the preceding effect.
var PhalanxStrike = set.New(
	"Phalanx Strike",
	card.House.Saurian,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.WC, "189"),
	card.WithBonus(card.Bonus.Aember),
	card.WithAbility(
		card.Trigger.Play, card.Repeat{
			Do: card.DealDamage{
				Amount: 1,
				Per: card.CardsInPlay{
					Player: card.Controller,
					Type:   card.Type.Creature,
				},
				Target: card.Target.Creature,
			},
			Gate: card.ByExalting{Creature: card.Target.FriendlyCreature},
		}),
)
