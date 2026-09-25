package worldscollide

import "github.com/dmikalova/vex/internal/card"

// Prescriptive Grammarbot
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Rare
//	Power:  3
//	Traits: Robot
//
//	Taunt, Hazardous 3.
//	Reap: Enrage a creature.
var PrescriptiveGrammarbot = set.New(
	"Prescriptive Grammarbot",
	card.House.Logos,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.WC, "173"),
	card.WithPower(3),
	card.WithTraits(card.Traits.Robot),
	card.WithKeywords(card.Keyword.Taunt),
	card.WithHazardous(3),
	card.WithAbility(
		card.Trigger.Reap, card.Enrage{Target: card.Target.Creature}),
)
