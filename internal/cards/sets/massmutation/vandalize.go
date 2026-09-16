package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Vandalize
//
//	House:  Shadows
//	Type:   Tactic
//	Rarity: Common
//	Bonus:  Æmber
//
//	Play: Reveal the top 3 cards of your opponent's deck. Discard 1. Put them back in any order.
var Vandalize = set.New(
	"Vandalize",
	card.House.Shadows,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.MM, "260"),
	card.WithBonus(card.Bonus.Aember),
	card.WithAbility(
		card.Trigger.Play, card.RevealTopOfDeck{
			Amount: 3,
			Player: card.Opponent,
			Then: []card.TopAct{
				card.ChooseAndMove{Count: 1, Dest: card.Into.Discard},
				card.ReorderRest{},
			},
		}),
)
