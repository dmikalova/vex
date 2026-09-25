package massmutation

import "github.com/dmikalova/vex/internal/card"

// Stealthster
//
//	House:  Star Alliance
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  3
//	Traits: Robot
//
//	Elusive.
//	Stealthster may be played as an upgrade instead of a creature, with the text: "This creature gains elusive."
var Stealthster = set.New(
	"Stealthster",
	card.House.StarAlliance,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.MM, "329"),
	card.WithPower(3),
	card.WithTraits(card.Traits.Robot),
	card.WithKeywords(card.Keyword.Elusive),
	card.WithStatic(card.StaticModifier{Keywords: card.Keywords(card.Keyword.Elusive)}),
	card.WithPlayableAsUpgrade(),
)
