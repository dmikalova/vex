package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Fangtooth Cavern
//
//	House:  Untamed
//	Type:   Artifact
//	Rarity: Uncommon
//	Æmber:  1
//	Traits: Location
//
//	At the end of your turn, destroy the least powerful Creature.
var FangtoothCavern = set.New(
	"Fangtooth Cavern",
	card.House.Untamed,
	card.Type.Artifact,
	card.Rarity.Uncommon,
	card.Provenance(card.WC, "370"),
	card.WithAemberBonus(1),
	card.WithTraits(card.Traits.Location),
	card.WithAbility(
		card.Trigger.EndOfTurn, card.Destroy{
			Target: card.Target.EachCreature.Refine(card.LeastPowerful),
		}),
)
