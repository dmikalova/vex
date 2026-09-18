package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Booby Trap
//
//	House:  Shadows
//	Type:   Tactic
//	Rarity: Uncommon
//	Bonus:  Æmber
//
//	Play: Choose a creature that is not on a flank. Deal 4 damage to the chosen creature and 2 damage to each of its neighbors.
var BoobyTrap = set.New(
	"Booby Trap",
	card.House.Shadows,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, "268"),
	card.WithBonus(card.Bonus.Aember),
	card.WithAbility(
		card.Trigger.Play, card.DealDamage{Spread: card.CreatureAndNeighbors{
			Amount:     4,
			Splash:     2,
			NotOnFlank: true,
		}}),
)
