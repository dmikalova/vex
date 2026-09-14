package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Cauldron Boil
//
//	House:  Untamed
//	Type:   Tactic
//	Rarity: Common
//	Æmber:  1
//
//	Play: Deal 1 damage to each Creature for each point of damage on it.
var CauldronBoil = set.New(
	"Cauldron Boil",
	card.House.Untamed,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.WC, "354"),
	card.WithAemberBonus(1),
	card.WithAbility(
		card.Trigger.Play, card.DealDamage{
			Amount:    1,
			Target:    card.Target.EachCreature,
			PerTarget: card.DamageOnIt,
		}),
)
