package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Explo-rover
//
//	House:  Star Alliance
//	Type:   Creature
//	Rarity: Common
//	Power:  3
//	Traits: Robot
//
//	Skirmish.
//	Explo-rover may be played as an upgrade instead of a creature, with the text: "This creature gains skirmish."
var ExploRover = set.New(
	"Explo-rover",
	card.House.StarAlliance,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.WC, "297"),
	card.WithPower(3),
	card.WithTraits(card.Traits.Robot),
	card.WithKeywords(card.Keyword.Skirmish),
	card.WithStatic(card.StaticModifier{
		Keywords: card.Keywords(card.Keyword.Skirmish),
	}),
	card.WithPlayableAsUpgrade(),
)
