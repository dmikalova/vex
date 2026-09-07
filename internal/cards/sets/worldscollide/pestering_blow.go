package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Pestering Blow
//
//	House:  Shadows
//	Type:   Tactic
//	Rarity: Common
//	Æmber:  1
//
//	Play: Deal 1 damage to a creature and enrage it.
var PesteringBlow = card.New(
	"Pestering Blow",
	card.House.Shadows,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.WC, 245),
	card.WithAemberBonus(1),
	card.WithAbility(
		card.Trigger.Play, card.DamageThen{
			Amount: 1,
			Target: card.Target.Creature,
			Then:   card.Enrage{Target: card.Target.Triggering},
		}),
)
