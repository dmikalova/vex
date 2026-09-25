package worldscollide

import "github.com/dmikalova/vex/internal/card"

// Barn Razing
//
//	House:  Brobnar
//	Type:   Tactic
//	Rarity: Common
//
//	Play: For the remainder of the turn, each time a friendly creature fights, your opponent loses 1 Æmber.
var BarnRazing = set.New(
	"Barn Razing",
	card.House.Brobnar,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.WC, "4"),
	card.WithAbility(
		card.Trigger.Play, card.ForRemainderOfTurn{
			On: card.Event.Fight,
			Do: card.LoseAember{
				Player: card.Opponent,
				Amount: 1,
			},
		}),
)
