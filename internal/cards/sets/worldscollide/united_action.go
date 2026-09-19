package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// United Action
//
//	House:  Star Alliance
//	Type:   Tactic
//	Rarity: Rare
//
//	Alpha.
//	Play: For the remainder of the turn, you may play cards from any house for which you have a card in play. You cannot use any cards for the remainder of the turn.
var UnitedAction = set.New(
	"United Action",
	card.House.StarAlliance,
	card.Type.Tactic,
	card.Rarity.Rare,
	card.Provenance(card.WC, "343"),
	card.WithKeywords(card.Keyword.Alpha),
	card.WithAbility(
		card.Trigger.Play, card.Sequence{Effects: []card.Effect{
			card.MayPlayOrUse{
				Houses: card.GrantHouses.Controlled,
				Grant:  card.GrantPlay,
			},
			card.Restrict{
				Player:   card.Controller,
				Action:   card.Restricted.Use,
				Duration: card.Duration.RemainderOfPlayerTurn,
			},
		}}),
)
