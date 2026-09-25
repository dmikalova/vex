package worldscollide

import "github.com/dmikalova/vex/internal/card"

// Imperium
//
//	House:  Saurian
//	Type:   Tactic
//	Rarity: Common
//	Bonus:  Æmber
//
//	Play: Ward 2 friendly creatures.
var Imperium = set.New(
	"Imperium",
	card.House.Saurian,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.WC, "186"),
	card.WithBonus(card.Bonus.Aember),
	card.WithAbility(
		card.Trigger.Play, card.Ward{
			Target: card.Target.EachFriendlyCreature,
			Amount: 2,
		}),
)
