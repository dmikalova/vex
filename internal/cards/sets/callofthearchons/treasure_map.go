package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Treasure Map
//
//	House:  Shadows
//	Type:   Tactic
//	Rarity: Rare
//	Æmber:  1
//
//	Play: If you played exactly 1 card this turn, gain 3 Æmber, and you cannot play cards for the remainder of the turn.
var TreasureMap = card.New(
	"Treasure Map",
	card.House.Shadows,
	card.Type.Tactic,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 284),
	card.WithAemberBonus(1),
	card.WithAbility(
		card.Trigger.Play, card.Sequence{
			Effects: []card.Effect{
				card.Conditional{
					Cond: card.CountIs{
						Count:  card.CardsPlayed{Player: card.Controller},
						Is:     card.Exactly,
						Amount: 1,
					},
					Then: card.GainAember{
						Player: card.Controller,
						Amount: 3,
					},
				},
				card.CannotPlay{
					Player:   card.Controller,
					Duration: card.Duration.EndOfTurn,
				},
			},
		}),
)
