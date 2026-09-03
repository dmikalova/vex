package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Shatter Storm
//
//	House:  Mars
//	Type:   Tactic
//	Rarity: Rare
//
//	Play: Lose all your Æmber, and for each Æmber you lost this way, your opponent loses 3 Æmber.
var ShatterStorm = card.New(
	"Shatter Storm",
	card.House.Mars,
	card.Type.Tactic,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 176),
	card.WithAbility(
		card.Trigger.Play, card.Sequence{Effects: []card.Effect{
			card.LoseAember{
				Player: card.Controller,
				By:     card.AllAember,
			},
			card.LoseAember{
				Player: card.Opponent,
				Amount: 3,
				Per:    card.AemberLostThisWay{Player: card.Controller},
			},
		}}),
)
