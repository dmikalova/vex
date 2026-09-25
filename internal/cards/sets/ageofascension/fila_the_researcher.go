package ageofascension

import "github.com/dmikalova/vex/internal/card"

// Fila the Researcher
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  1
//	Traits: Human • Scientist
//
//	Elusive.
//	After a creature is played adjacent to Fila the Researcher, draw a card.
var FilaTheResearcher = set.New(
	"Fila the Researcher",
	card.House.Logos,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.AoA, "129"),
	card.WithPower(1),
	card.WithTraits(card.Traits.Human, card.Traits.Scientist),
	card.WithKeywords(card.Keyword.Elusive),
	card.WithAbility(
		card.Trigger.AfterCreaturePlayedAdjacent, card.Draw{Amount: 1}),
)
