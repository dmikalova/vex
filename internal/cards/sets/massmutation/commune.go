package massmutation

import "github.com/dmikalova/vex/internal/card"

// Commune
//
//	House:  Untamed
//	Type:   Tactic
//	Rarity: Uncommon
//
//	Omega.
//	Play: Lose all your Æmber. Gain 4 Æmber.
var Commune = set.New(
	"Commune",
	card.House.Untamed,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.MM, "380"),
	card.WithKeywords(card.Keyword.Omega),
	card.WithAbility(
		card.Trigger.Play, card.Sequence{
			Effects: []card.Effect{
				card.LoseAember{
					Player: card.Controller,
					By:     card.AllAember,
				},
				card.GainAember{
					Player: card.Controller,
					Amount: 4,
				},
			},
		}),
)
