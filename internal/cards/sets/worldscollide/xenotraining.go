package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Xenotraining
//
//	House:  Star Alliance
//	Type:   Tactic
//	Rarity: Uncommon
//	Æmber:  1
//
//	Play: For each house represented among friendly Creatures, a friendly Creature captures 1 Æmber from your opponent.
var Xenotraining = set.New(
	"Xenotraining",
	card.House.StarAlliance,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.WC, "323"),
	card.WithAemberBonus(1),
	card.WithAbility(
		card.Trigger.Play, card.CaptureAember{
			Amount: 1,
			Target: card.Target.FriendlyCreature,
			Source: card.Opponent,
			Times:  card.HousesAmong{Player: card.Controller, Type: card.Type.Creature},
		}),
)
