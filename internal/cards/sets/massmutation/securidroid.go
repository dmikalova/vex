package massmutation

import "github.com/dmikalova/vex/internal/card"

// Securi-Droid
//
//	House:  Star Alliance
//	Type:   Creature
//	Rarity: Common
//	Power:  4
//	Traits: Robot
//
//	Taunt.
//	Securi-Droid may be played as an upgrade instead of a creature, with the text: "This creature gains taunt."
var SecuriDroid = set.New(
	"Securi-Droid",
	card.House.StarAlliance,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.MM, "312"),
	card.WithPower(4),
	card.WithTraits(card.Traits.Robot),
	card.WithKeywords(card.Keyword.Taunt),
	card.WithStatic(card.StaticModifier{Keywords: card.Keywords(card.Keyword.Taunt)}),
	card.WithPlayableAsUpgrade(),
)
