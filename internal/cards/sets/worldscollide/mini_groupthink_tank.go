package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Mini Groupthink Tank
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  3
//	Armor:  2
//	Traits: Robot • Experiment
//
//	Play/Fight/Reap: Deal 8 damage to a creature that shares a house with 2 of its neighbors.
var MiniGroupthinkTank = card.New(
	"Mini Groupthink Tank",
	card.House.Logos,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.WC, "155"),
	card.WithPower(3),
	card.WithArmor(2),
	card.WithTraits(card.Traits.Robot, card.Traits.Experiment),
	card.WithAbility(card.Trigger.PlayFightReap, card.DealDamage{
		Amount: 8,
		Target: card.Target.Creature.SharesHouseWithNeighbors(2),
	}),
)
