package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Looter Goblin
//
//	House:  Brobnar
//	Type:   Creature
//	Rarity: Rare
//	Power:  2
//	Traits: Goblin
//
//	Elusive.
//	Reap: For the remainder of the turn, each time an enemy creature is destroyed, gain 1 Æmber.
var LooterGoblin = card.New(
	"Looter Goblin",
	card.House.Brobnar,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 41),
	card.WithPower(2),
	card.WithTraits(card.Traits.Goblin),
	card.WithKeywords(card.Keyword.Elusive),
	card.WithAbility(
		card.Trigger.Reap, card.ForRemainderOfTurn{
			On: card.Event.EnemyCreatureDestroyed,
			Do: card.GainAember{
				Player: card.Controller,
				Amount: 1,
			},
		}),
)
