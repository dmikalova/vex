package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Positron Bolt
//
//	House:  Logos
//	Type:   Tactic
//	Rarity: Uncommon
//	Æmber:  1
//
//	Play: Choose a flank creature. Deal 3 damage to it, 2 damage to its neighbor, and 1 damage to the neighbor's other neighbor.
var PositronBolt = card.New(
	"Positron Bolt",
	card.House.Logos,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, 118),
	card.WithAemberBonus(1),
	card.WithAbility(
		card.Trigger.Play, card.DealDamage{Spread: card.FlankWalk{
			Amounts: []int{3, 2, 1},
		}}),
)
