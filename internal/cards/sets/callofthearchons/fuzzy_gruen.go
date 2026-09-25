package callofthearchons

import "github.com/dmikalova/vex/internal/card"

// Fuzzy Gruen
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Rare
//	Power:  5
//	Bonus:  Æmber Æmber
//	Traits: Beast
//
//	Play: Your opponent gains 1 Æmber.
var FuzzyGruen = set.New(
	"Fuzzy Gruen",
	card.House.Untamed,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.CotA, "353"),
	card.WithBonus(card.Bonus.Aember, card.Bonus.Aember),
	card.WithPower(5),
	card.WithTraits(card.Traits.Beast),
	card.WithAbility(
		card.Trigger.Play, card.GainAember{
			Player: card.Opponent,
			Amount: 1,
		}),
)
