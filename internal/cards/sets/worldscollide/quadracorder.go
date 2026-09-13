package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Quadracorder
//
//	House:  Star Alliance
//	Type:   Upgrade
//	Rarity: Uncommon
//
//	This Creature gains, "Your opponent's keys cost +1 Æmber for each house represented among friendly Creatures."
var Quadracorder = card.New(
	"Quadracorder",
	card.House.StarAlliance,
	card.Type.Upgrade,
	card.Rarity.Uncommon,
	card.Provenance(card.WC, "316"),
	card.WithStatic(card.StaticModifier{
		KeyCostChange: card.KeyCostChange(card.Opponent, 1).Per(card.HousesAmong{
			Player: card.Controller,
			Type:   card.Type.Creature,
		}),
	}),
)
