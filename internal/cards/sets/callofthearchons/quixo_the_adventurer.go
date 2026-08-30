package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Quixo the "Adventurer"
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Common
//	Power:  3
//	Traits: Human • Scientist
//
//	Skirmish.
//	Fight: Draw a card.
var QuixoTheAdventurer = card.New(
	"Quixo the \"Adventurer\"",
	card.House.Logos,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.CotA, 144),
	card.WithPower(3),
	card.WithTraits("Human", "Scientist"),
	card.WithKeywords(card.Keyword.Skirmish),
	card.WithAbility(
		card.Trigger.Fight, card.Draw{Amount: 1}),
)
