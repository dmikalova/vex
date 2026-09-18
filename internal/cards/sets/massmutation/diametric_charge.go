package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Diametric Charge
//
//	House:  Logos
//	Type:   Tactic
//	Rarity: Common
//	Bonus:  Æmber
//
//	Play: Choose a creature. Deal 1 damage to the chosen creature and 2 damage to each of its neighbors.
var DiametricCharge = set.New(
	"Diametric Charge",
	card.House.Logos,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.MM, "070"),
	card.WithBonus(card.Bonus.Aember),
	card.WithAbility(
		card.Trigger.Play, card.DealDamage{Spread: card.CreatureAndNeighbors{
			Amount: 1,
			Splash: 2,
		}}),
)
