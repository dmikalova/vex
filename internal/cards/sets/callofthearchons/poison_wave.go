package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Poison Wave
//
//	House:  Shadows
//	Type:   Tactic
//	Rarity: Common
//	Æmber:  1
//
//	Play: Deal 2 damage to each creature.
var PoisonWave = card.New(
	"Poison Wave",
	card.House.Shadows,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.CotA, 280),
	card.WithAemberBonus(1),
	card.WithAbility(
		card.Trigger.Play, card.DealDamage{
			Amount: 2,
			Target: card.Target.EachCreature,
		}),
)
