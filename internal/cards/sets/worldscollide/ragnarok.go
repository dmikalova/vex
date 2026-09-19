package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Ragnarok
//
//	House:  Brobnar
//	Type:   Tactic
//	Rarity: Rare
//
//	Alpha.
//	Play: You cannot use creatures to reap for the remainder of the turn. For the remainder of the turn, each time a friendly creature fights, gain 1 Æmber. At the end of the turn, destroy each creature.
var Ragnarok = set.New(
	"Ragnarok",
	card.House.Brobnar,
	card.Type.Tactic,
	card.Rarity.Rare,
	card.Provenance(card.WC, "47"),
	card.WithKeywords(card.Keyword.Alpha),
	card.WithAbility(
		card.Trigger.Play, card.Sequence{
			Effects: []card.Effect{
				card.Restrict{
					Player:   card.Controller,
					Action:   card.Restricted.Reaping,
					Duration: card.Duration.RemainderOfPlayerTurn,
				},
				card.ForRemainderOfTurn{
					On: card.Event.Fight,
					Do: card.GainAember{
						Player: card.Controller,
						Amount: 1,
					},
				},
				card.DestroyEachCreatureAtEndOfTurn{},
			},
		}),
)
