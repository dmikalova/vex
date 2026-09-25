package massmutation

import "github.com/dmikalova/vex/internal/card"

// Drecker
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Common
//	Power:  4
//	Traits: Imp
//
//	Damage dealt to Drecker's neighbors during fights is also dealt to Drecker.
//	Reap: Steal 1 Æmber.
var Drecker = set.New(
	"Drecker",
	card.House.Dis,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.MM, "006"),
	card.WithPower(4),
	card.WithTraits(card.Traits.Imp),
	card.WithAlsoTakesNeighborFightDamage(),
	card.WithAbility(
		card.Trigger.Reap, card.StealAember{Amount: 1}),
)
