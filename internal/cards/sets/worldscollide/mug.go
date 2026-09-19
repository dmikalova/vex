package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Mug
//
//	House:  Shadows
//	Type:   Tactic
//	Rarity: Common
//	Bonus:  Æmber
//
//	Play: Choose a creature. Move 1 Æmber from it to your pool. Deal 2 damage to it.
var Mug = set.New(
	"Mug",
	card.House.Shadows,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.WC, "244"),
	card.WithBonus(card.Bonus.Aember),
	card.WithAbility(
		card.Trigger.Play, card.ChooseCreatureThen{
			Target: card.Target.Creature,
			Then: card.Sequence{Effects: []card.Effect{
				card.MoveAember{
					Amount: 1,
					From:   card.Target.Triggering,
					To:     card.Controller,
				},
				card.DealDamage{
					Amount: 2,
					Target: card.Target.Triggering,
				},
			}},
		}),
)
