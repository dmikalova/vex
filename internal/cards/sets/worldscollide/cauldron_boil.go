package worldscollide

import "github.com/dmikalova/vex/internal/card"

// Cauldron Boil
//
//	House:  Untamed
//	Type:   Tactic
//	Rarity: Common
//	Bonus:  Æmber
//
//	Play: Deal 1 damage to each creature for each point of damage on it.
var CauldronBoil = set.New(
	"Cauldron Boil",
	card.House.Untamed,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.WC, "354"),
	card.WithBonus(card.Bonus.Aember),
	card.WithAbility(
		card.Trigger.Play, card.DealDamage{
			Amount:    1,
			Target:    card.Target.EachCreature,
			PerTarget: card.DamageOnIt,
		}),
)
