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
//	Explo-rover may be played as an Upgrade instead of a Creature, with the text: "This Creature gains skirmish."
var ExploRover = card.New(
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
