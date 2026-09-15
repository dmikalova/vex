package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Poison Wave
//
//	House:  Shadows
//	Type:   Tactic
//	Rarity: Common
//	Bonus:  Æmber
//
//	Play: Deal 2 damage to each creature.
var PoisonWave = set.New(
	"Poison Wave",
	card.House.Shadows,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.CotA, "280"),
	card.WithBonus(card.Bonus.Aember),
	card.WithAbility(
		card.Trigger.Play, card.DealDamage{
			Amount: 2,
			Target: card.Target.EachCreature,
		}),
)
