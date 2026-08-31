package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Booby Trap
//
//	House:  Shadows
//	Type:   Tactic
//	Rarity: Uncommon
//	Æmber:  1
//
//	Play: Deal 4 damage to a creature that is not on a flank and 2 damage to each of its neighbors.
var BoobyTrap = card.New(
	"Booby Trap",
	card.House.Shadows,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, 268),
	card.WithAemberBonus(1),
	card.WithAbility(
		card.Trigger.Play, card.DealDamage{Spread: card.CreatureAndNeighbors{
			Amount: 4,
			Splash: 2,
		}}),
)
