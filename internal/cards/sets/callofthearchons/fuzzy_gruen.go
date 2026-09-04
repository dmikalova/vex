package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Fuzzy Gruen
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Rare
//	Power:  5
//	Æmber:  2
//	Traits: Beast
//
//	Play: Your opponent gains 1 Æmber.
var FuzzyGruen = card.New(
	"Fuzzy Gruen",
	card.House.Untamed,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 353),
	card.WithPower(5),
	card.WithAemberBonus(2),
	card.WithTraits(card.Traits.Beast),
	card.WithAbility(
		card.Trigger.Play, card.GainAember{
			Player: card.Opponent,
			Amount: 1,
		}),
)
