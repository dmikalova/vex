package worldscollide

import "github.com/dmikalova/vex/internal/card"

// Gleeful Mayhem
//
//	House:  Dis
//	Type:   Tactic
//	Rarity: Uncommon
//	Bonus:  Æmber
//
//	Play: For each house, deal 5 damage to a creature of that house.
var GleefulMayhem = set.New(
	"Gleeful Mayhem",
	card.House.Dis,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.WC, "090"),
	card.WithBonus(card.Bonus.Aember),
	card.WithAbility(card.Trigger.Play, card.ForEachHouse{
		Do: card.DealDamage{
			Amount: 5,
			Target: card.Target.Creature.House(card.Houses.Each),
		},
	}),
)
