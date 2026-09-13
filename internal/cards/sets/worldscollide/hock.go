package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Hock
//
//	House:  Shadows
//	Type:   Tactic
//	Rarity: Common
//	Æmber:  1
//
//	Play: Destroy an Artifact -> gain 1 Æmber.
var Hock = card.New(
	"Hock",
	card.House.Shadows,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.WC, "239"),
	card.WithAemberBonus(1),
	card.WithAbility(
		card.Trigger.Play, card.Then{
			First:  card.Destroy{Target: card.Target.Artifact},
			Result: card.GainAember{Player: card.Controller, Amount: 1},
		}),
)
