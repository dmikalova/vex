package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Take that, Smartypants
//
//	House:  Brobnar
//	Type:   Tactic
//	Rarity: Rare
//	Bonus:  Æmber
//
//	Play: If there are 3 or more enemy Logos cards in play, steal 2 Æmber.
var TakeThatSmartypants = set.New(
	"Take that, Smartypants",
	card.House.Brobnar,
	card.Type.Tactic,
	card.Rarity.Rare,
	card.Provenance(card.CotA, "11"),
	card.WithBonus(card.Bonus.Aember),
	card.WithAbility(card.Trigger.Play, card.Conditional{
		Cond: card.InPlay{
			Player: card.Opponent,
			House:  card.Houses.Named(card.House.Logos),
			Amount: 3,
		},
		Then: card.StealAember{Amount: 2},
	}),
)
