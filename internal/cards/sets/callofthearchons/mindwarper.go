package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Mindwarper
//
//	House:  Mars
//	Type:   Creature
//	Rarity: Common
//	Power:  2
//	Traits: Martian • Scientist
//
//	Elusive.
//	Action: An enemy creature captures 1 Æmber from their own side.
var Mindwarper = card.New(
	"Mindwarper",
	card.House.Mars,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.CotA, 196),
	card.WithPower(2),
	card.WithTraits("Martian", "Scientist"),
	card.WithKeywords(card.Keyword.Elusive),
	card.WithAbility(
		card.Trigger.Action, card.CaptureAember{
			Amount: 1,
			Target: card.Target.EnemyCreature,
			Source: card.Opponent,
		}),
)
