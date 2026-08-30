package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Ammonia Clouds
//
//	House:  Mars
//	Type:   Tactic
//	Rarity: Common
//
//	Play: Deal 3 damage to each creature.
var AmmoniaClouds = card.New(
	"Ammonia Clouds",
	card.House.Mars,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.CotA, 160),
	card.WithAbility(
		card.Trigger.Play, card.DealDamage{
			Amount: 3,
			Target: card.Target.EachCreature,
		}),
)
