package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Subdue
//
//	House:  Star Alliance
//	Type:   Tactic
//	Rarity: Common
//	Bonus:  Æmber
//
//	Play: Deal 1 damage to a creature and stun it.
var Subdue = set.New(
	"Subdue",
	card.House.StarAlliance,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.MM, "314"),
	card.WithBonus(card.Bonus.Aember),
	card.WithAbility(
		card.Trigger.Play, card.DamageThen{
			Amount: 1,
			After:  card.Always,
			Target: card.Target.Creature,
			Then:   card.Stun{Target: card.Target.Triggering},
		}),
)
