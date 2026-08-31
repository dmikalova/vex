package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Loot the Bodies
//
//	House:  Brobnar
//	Type:   Tactic
//	Rarity: Common
//
//	Play: For the remainder of the turn, each time an enemy creature is destroyed, gain 1 Æmber.
var LootTheBodies = card.New(
	"Loot the Bodies",
	card.House.Brobnar,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.CotA, 10),
	card.WithAbility(
		card.Trigger.Play, card.ForRemainderOfTurn{
			On: card.Event.EnemyCreatureDestroyed,
			Do: card.GainAember{
				Player: card.Controller,
				Amount: 1,
			},
		}),
)
