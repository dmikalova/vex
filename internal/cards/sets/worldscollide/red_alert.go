package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Red Alert
//
//	House:  Star Alliance
//	Type:   Tactic
//	Rarity: Common
//
//	Play: For each creature your opponent controls in excess of you, deal 1 damage to each enemy creature.
var RedAlert = card.New(
	"Red Alert",
	card.House.StarAlliance,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.WC, 303),
	card.WithAbility(
		card.Trigger.Play, card.DealDamage{
			Amount: 1,
			Per:    card.ExcessCreatures{Player: card.Opponent},
			Target: card.Target.EachEnemyCreature,
		}),
)
