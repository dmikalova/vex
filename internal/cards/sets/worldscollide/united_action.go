package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// United Action
//
//	House:  Star Alliance
//	Type:   Tactic
//	Rarity: Rare
//
//	Alpha.
//	Play: For the remainder of the turn, you may play cards from any house for which you have a card in play. You cannot use cards this turn.
var UnitedAction = card.New(
	"United Action",
	card.House.StarAlliance,
	card.Type.Tactic,
	card.Rarity.Rare,
	card.Provenance(card.WC, "343"),
	card.WithKeywords(card.Keyword.Alpha),
	card.WithAbility(
		card.Trigger.Play, card.Sentences{Effects: []card.Effect{
			card.MayPlayOffHouse{
				Controlled: true,
				Grant:      card.GrantPlay,
			},
			card.CannotUse{
				Player:   card.Controller,
				Duration: card.Duration.RemainderOfPlayerTurn,
			},
		}}),
)
