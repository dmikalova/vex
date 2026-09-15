package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Look Over There!
//
//	House:  Shadows
//	Type:   Tactic
//	Rarity: Common
//
//	Play: Deal 2 damage to a creature. If it is not destroyed, steal 1 Æmber.
var LookOverThere = set.New(
	"Look Over There!",
	card.House.Shadows,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.MM, "250"),
	card.WithAbility(
		card.Trigger.Play, card.DamageThen{
			Amount: 2,
			After:  card.IfSurvives,
			Target: card.Target.Creature,
			Then:   card.StealAember{Amount: 1},
		}),
)
