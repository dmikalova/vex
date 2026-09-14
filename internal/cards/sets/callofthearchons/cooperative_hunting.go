package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Cooperative Hunting
//
//	House:  Untamed
//	Type:   Tactic
//	Rarity: Common
//
//	Play: For each friendly Creature in play, deal 1 damage to a Creature.
var CooperativeHunting = set.New(
	"Cooperative Hunting",
	card.House.Untamed,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.CotA, "319"),
	card.WithAbility(
		card.Trigger.Play, card.DealDamage{
			Amount: 1,
			Per: card.InPlay{
				Player: card.Controller,
				Type:   card.Type.Creature,
			},
			Target: card.Target.Creature,
		}),
)
