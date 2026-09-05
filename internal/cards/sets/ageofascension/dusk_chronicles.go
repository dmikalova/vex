package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Dusk Chronicles
//
//	House:  Shadows
//	Type:   Tactic
//	Rarity: Common
//	Æmber:  1
//
//	Play: If your opponent has more Æmber than you, draw a card. If you have more Æmber than your opponent, archive a card from your hand.
var DuskChronicles = card.New(
	"Dusk Chronicles",
	card.House.Shadows,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.AoA, 268),
	card.WithAemberBonus(1),
	card.WithAbility(
		card.Trigger.Play, card.Sentences{
			Effects: []card.Effect{
				card.Conditional{
					Cond: card.OpponentAember{Is: card.MoreThanYou},
					Then: card.Draw{Amount: 1},
				},
				card.Conditional{
					Cond: card.YourAember{Is: card.MoreThanOpponent},
					Then: card.ArchiveFromHand{Amount: 1},
				},
			},
		}),
)
