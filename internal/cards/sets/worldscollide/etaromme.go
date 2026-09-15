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
//	Reap: Destroy a creature of the house with the most creatures in play.
var Etaromme = set.New(
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
