package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Fertility Chant
//
//	House:  Untamed
//	Type:   Tactic
//	Rarity: Rare
//	Æmber:  4
//
//	Play: Your opponent gains 2 Æmber.
var FertilityChant = card.New(
	"Fertility Chant",
	card.House.Untamed,
	card.Type.Tactic,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 321),
	card.WithAemberBonus(4),
	card.WithAbility(
		card.Trigger.Play, card.GainAember{
			Player: card.Opponent,
			Amount: 2,
		}),
)
