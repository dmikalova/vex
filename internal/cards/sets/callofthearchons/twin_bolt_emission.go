package callofthearchons

import "github.com/dmikalova/vex/internal/card"

// Twin Bolt Emission
//
//	House:  Logos
//	Type:   Tactic
//	Rarity: Common
//	Bonus:  Æmber
//
//	Play: Deal 2 damage to a creature and 2 damage to a different creature.
var TwinBoltEmission = set.New(
	"Twin Bolt Emission",
	card.House.Logos,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.CotA, "124"),
	card.WithBonus(card.Bonus.Aember),
	card.WithAbility(
		card.Trigger.Play, card.DealDamage{Spread: card.DifferentCreatures{
			First:  2,
			Second: 2,
		}}),
)
