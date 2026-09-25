package massmutation

import "github.com/dmikalova/vex/internal/card"

// Brabble
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Common
//	Power:  2
//	Traits: Imp
//
//	Destroyed: If it is your turn, your opponent loses 1 Æmber. Otherwise, your opponent loses 3 Æmber.
var Brabble = set.New(
	"Brabble",
	card.House.Dis,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.MM, "003"),
	card.WithPower(2),
	card.WithTraits(card.Traits.Imp),
	card.WithAbility(
		card.Trigger.Destroyed, card.Conditional{
			Cond: card.ItIsYourTurn{},
			Then: card.LoseAember{
				Player: card.Opponent,
				Amount: 1,
			},
			Else: card.LoseAember{
				Player: card.Opponent,
				Amount: 3,
			},
		}),
)
