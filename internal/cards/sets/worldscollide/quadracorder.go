package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Quadracorder
//
//	House:  Star Alliance
//	Type:   Upgrade
//	Rarity: Uncommon
//
//	This creature gains, "Your opponent's keys cost +1 Æmber for each house represented among friendly creatures (to a maximum of 3)."
var Quadracorder = card.New(
	"Quadracorder",
	card.House.StarAlliance,
	card.Type.Upgrade,
	card.Rarity.Uncommon,
	card.Provenance(card.WC, 316),
	card.WithStatic(card.StaticModifier{
		KeyCostChange: card.KeyCostChange(card.Opponent, 1).Per(card.HousesAmong{
			Player: card.Controller,
			Type:   card.Type.Creature,
			Max:    3,
		}),
	}),
)
