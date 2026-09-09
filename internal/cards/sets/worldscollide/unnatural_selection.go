package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Unnatural Selection
//
//	House:  Untamed
//	Type:   Tactic
//	Rarity: Common
//	Æmber:  1
//
//	Play: Choose 3 friendly creatures and 3 enemy creatures. Destroy each other creature.
var UnnaturalSelection = card.New(
	"Unnatural Selection",
	card.House.Untamed,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.WC, "367"),
	card.WithAemberBonus(1),
	card.WithAbility(
		card.Trigger.Play, card.DestroyAllExceptChosen{
			FriendlyKept: 3,
			EnemyKept:    3,
		}),
)
