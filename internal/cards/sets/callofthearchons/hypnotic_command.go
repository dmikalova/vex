package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Hypnotic Command
//
//	House:  Mars
//	Type:   Tactic
//	Rarity: Rare
//
//	Play: For each friendly Mars creature, an enemy creature captures 1 Æmber from their own side.
var HypnoticCommand = card.New(
	"Hypnotic Command",
	card.House.Mars,
	card.Type.Tactic,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 164),
	card.WithAbility(
		card.Trigger.Play, card.CaptureAember{
			Amount: 1,
			Target: card.Target.EnemyCreature,
			Source: card.Opponent,
			Per: card.InPlay{
				Player: card.Controller,
				Type:   card.Type.Creature,
				House:  card.House.Self,
			},
		}),
)
