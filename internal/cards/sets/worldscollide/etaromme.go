package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Etaromme
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Common
//	Power:  4
//	Traits: Demon
//
//	Reap: Destroy a Creature of the house with the most Creatures in play.
var Etaromme = card.New(
	"Etaromme",
	card.House.Dis,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.WC, "73"),
	card.WithPower(4),
	card.WithTraits(card.Traits.Demon),
	card.WithAbility(
		card.Trigger.Reap,
		card.Destroy{Target: card.Target.Creature.OfHouseWithMostCreatures()},
	),
)
