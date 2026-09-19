package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Unnatural Selection
//
//	House:  Untamed
//	Type:   Tactic
//	Rarity: Common
//	Bonus:  Æmber
//
//	Play: Choose 3 friendly creatures and 3 enemy creatures. Destroy each other creature.
var UnnaturalSelection = set.New(
	"Unnatural Selection",
	card.House.Untamed,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.WC, "367"),
	card.WithBonus(card.Bonus.Aember),
	card.WithAbility(
		card.Trigger.Play, card.Destroy{
			Target: card.Target.EachCreature.Refine(card.KeepPerSide(3)),
		}),
)
