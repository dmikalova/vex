package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Foggify
//
//	House:  Logos
//	Type:   Tactic
//	Rarity: Common
//	Æmber:  1
//
//	Play: During your opponent's next turn, after an enemy Creature is used to fight, stun it.
var Foggify = set.New(
	"Foggify",
	card.House.Logos,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.CotA, "110"),
	card.WithAemberBonus(1),
	card.WithAbility(
		card.Trigger.Play, card.StunEnemyFighters{
			Duration: card.Duration.OpponentNextTurn,
		}),
)
