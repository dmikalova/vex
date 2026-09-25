package worldscollide

import "github.com/dmikalova/vex/internal/card"

// Groupthink Tank
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  4
//	Armor:  3
//	Traits: Robot • Experiment
//
//	Action: Deal 4 damage to each creature that shares a house with at least 1 of its neighbors.
var GroupthinkTank = set.New(
	"Groupthink Tank",
	card.House.Logos,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.WC, "151"),
	card.WithPower(4),
	card.WithArmor(3),
	card.WithTraits(card.Traits.Robot, card.Traits.Experiment),
	card.WithAbility(card.Trigger.Action, card.DealDamage{
		Amount: 4,
		Target: card.Target.EachCreature.SharesHouseWithNeighbors(1),
	}),
)
