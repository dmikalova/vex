package worldscollide

import "github.com/dmikalova/vex/internal/card"

// Pestering Blow
//
//	House:  Shadows
//	Type:   Tactic
//	Rarity: Common
//	Bonus:  Æmber
//
//	Play: Deal 1 damage to a creature and enrage it.
var PesteringBlow = set.New(
	"Pestering Blow",
	card.House.Shadows,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.WC, "245"),
	card.WithBonus(card.Bonus.Aember),
	card.WithAbility(
		card.Trigger.Play, card.DealDamage{
			Amount: 1,
			After:  card.Always,
			Target: card.Target.Creature,
			Then:   card.Enrage{Target: card.Target.Triggering},
		}),
)
