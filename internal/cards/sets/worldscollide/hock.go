package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Hock
//
//	House:  Shadows
//	Type:   Tactic
//	Rarity: Common
//	Bonus:  Æmber
//
//	Play: Destroy an artifact -> gain 1 Æmber.
var Hock = set.New(
	"Hock",
	card.House.Shadows,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.WC, "239"),
	card.WithBonus(card.Bonus.Aember),
	card.WithAbility(
		card.Trigger.Play, card.Then{
			First: card.Destroy{Target: card.Target.Artifact},
			Result: card.GainAember{
				Player: card.Controller,
				Amount: 1,
			},
		}),
)
