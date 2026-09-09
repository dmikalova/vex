package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Imperium
//
//	House:  Saurian
//	Type:   Tactic
//	Rarity: Common
//	Æmber:  1
//
//	Play: Ward 2 friendly creatures.
var Imperium = card.New(
	"Imperium",
	card.House.Saurian,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.WC, "186"),
	card.WithAemberBonus(1),
	card.WithAbility(
		card.Trigger.Play, card.Ward{
			Target: card.Target.EachFriendlyCreature,
			Amount: 2,
		}),
)
