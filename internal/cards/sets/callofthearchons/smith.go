package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Smith
//
//	House:  Brobnar
//	Type:   Tactic
//	Rarity: Uncommon
//	Bonus:  Æmber
//
//	Play: If you control more creatures than your opponent, gain 2 Æmber.
var Smith = set.New(
	"Smith",
	card.House.Brobnar,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, "14"),
	card.WithBonus(card.Bonus.Aember),
	card.WithAbility(
		card.Trigger.Play, card.Conditional{
			Cond: card.ControlsMoreCreatures{},
			Then: card.GainAember{
				Player: card.Controller,
				Amount: 2,
			},
		}),
)
