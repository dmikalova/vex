package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Anahita the Trader
//
//	House:  Sanctum
//	Type:   Creature
//	Rarity: Rare
//	Power:  2
//	Traits: Human • Merchant
//
//	Reap: Your opponent gains control of a friendly artifact -> steal 2 Æmber.
var AnahitaTheTrader = set.New(
	"Anahita the Trader",
	card.House.Sanctum,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.AoA, "248"),
	card.WithPower(2),
	card.WithTraits(card.Traits.Human, card.Traits.Merchant),
	card.WithAbility(
		card.Trigger.Reap, card.Then{
			First: card.TakeControl{
				Target:     card.Target.FriendlyArtifact,
				Duration:   card.Duration.UntilCardLeavesPlay,
				ToOpponent: true,
			},
			Result: card.StealAember{Amount: 2},
		}),
)
