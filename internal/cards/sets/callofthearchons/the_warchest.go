package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// The Warchest
//
//	House:  Brobnar
//	Type:   Artifact
//	Rarity: Uncommon
//	Traits: Item
//
//	Action: For each enemy Creature that was destroyed in a fight this turn, gain 1 Æmber.
var TheWarchest = set.New(
	"The Warchest",
	card.House.Brobnar,
	card.Type.Artifact,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, "27"),
	card.WithTraits(card.Traits.Item),
	card.WithAbility(
		card.Trigger.Action, card.GainAember{
			Player: card.Controller,
			Amount: 1,
			Per: card.TurnCount{
				Player: card.Controller,
				Of:     card.TurnStat.EnemyCreaturesFightKilled,
			},
		}),
)
