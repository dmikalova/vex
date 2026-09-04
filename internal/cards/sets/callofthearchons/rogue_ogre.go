package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Rogue Ogre
//
//	House:  Brobnar
//	Type:   Creature
//	Rarity: Rare
//	Power:  6
//	Traits: Giant • Mutant
//
//	At the end of your turn, if you played exactly 1 card this turn, heal 2 damage from Rogue Ogre, and Rogue Ogre captures 1 Æmber from your opponent.
var RogueOgre = card.New(
	"Rogue Ogre",
	card.House.Brobnar,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 45),
	card.WithPower(6),
	card.WithTraits(card.Traits.Giant, card.Traits.Mutant),
	card.WithAbility(
		card.Trigger.EndOfTurn, card.Conditional{
			Cond: card.CountIs{
				Count:  card.CardsPlayed{Player: card.Controller},
				Is:     card.Exactly,
				Amount: 1,
			},
			Then: card.Sequence{Effects: []card.Effect{
				card.Heal{
					Amount: 2,
					Target: card.Target.This,
				},
				card.CaptureAember{
					Amount: 1,
					Target: card.Target.This,
					Source: card.Opponent,
				},
			}},
		}),
)
