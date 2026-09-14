package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Twin Bolt Emission
//
//	House:  Logos
//	Type:   Tactic
//	Rarity: Common
//	Æmber:  1
//
//	Play: Deal 2 damage to a Creature and deal 2 damage to a different Creature.
var TwinBoltEmission = set.New(
	"Twin Bolt Emission",
	card.House.Logos,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.CotA, "124"),
	card.WithAemberBonus(1),
	card.WithAbility(
		card.Trigger.Play, card.DealDamage{Spread: card.DifferentCreatures{
			First:  2,
			Second: 2,
		}}),
)
