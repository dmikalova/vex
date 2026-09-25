package callofthearchons

import "github.com/dmikalova/vex/internal/card"

// Fertility Chant
//
//	House:  Untamed
//	Type:   Tactic
//	Rarity: Rare
//	Bonus:  Æmber Æmber Æmber Æmber
//
//	Play: Your opponent gains 2 Æmber.
var FertilityChant = set.New(
	"Fertility Chant",
	card.House.Untamed,
	card.Type.Tactic,
	card.Rarity.Rare,
	card.Provenance(card.CotA, "321"),
	card.WithBonus(card.Bonus.Aember, card.Bonus.Aember, card.Bonus.Aember, card.Bonus.Aember),
	card.WithAbility(
		card.Trigger.Play, card.GainAember{
			Player: card.Opponent,
			Amount: 2,
		}),
)
