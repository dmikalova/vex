package worldscollide

import "github.com/dmikalova/vex/internal/card"

// Fangtooth Cavern
//
//	House:  Untamed
//	Type:   Artifact
//	Rarity: Uncommon
//	Bonus:  Æmber
//	Traits: Location
//
//	At the end of your turn, destroy the least powerful creature.
var FangtoothCavern = set.New(
	"Fangtooth Cavern",
	card.House.Untamed,
	card.Type.Artifact,
	card.Rarity.Uncommon,
	card.Provenance(card.WC, "370"),
	card.WithBonus(card.Bonus.Aember),
	card.WithTraits(card.Traits.Location),
	card.WithAbility(
		card.Trigger.EndOfTurn, card.Destroy{
			Target: card.Target.EachCreature.Refine(card.LeastPowerful),
		}),
)
