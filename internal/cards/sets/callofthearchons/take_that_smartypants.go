package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Take that, Smartypants
//
//	House:  Brobnar
//	Type:   Tactic
//	Rarity: Rare
//	Æmber:  1
//
//	Play: If there are 3 enemy Logos cards in play, steal 2 Æmber.
var TakeThatSmartypants = card.New(
	"Take that, Smartypants",
	card.House.Brobnar,
	card.Type.Tactic,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 11),
	card.WithAemberBonus(1),
	card.WithAbility(card.Trigger.Play, card.Conditional{
		Cond: card.InPlay{
			Player: card.Opponent,
			House:  card.House.Logos,
			Amount: 3,
		},
		Then: card.StealAember{Amount: 2},
	}),
)
