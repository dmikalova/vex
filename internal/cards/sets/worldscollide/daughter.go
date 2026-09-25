package worldscollide

import "github.com/dmikalova/vex/internal/card"

// Daughter
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Common
//	Power:  2
//	Traits: Cyborg • Scientist
//
//	Elusive.
//	Your hand size is 1 more.
var Daughter = set.New(
	"Daughter",
	card.House.Logos,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.WC, "131"),
	card.WithPower(2),
	card.WithTraits(card.Traits.Cyborg, card.Traits.Scientist),
	card.WithKeywords(card.Keyword.Elusive),
	card.WithDrawModifier(card.Controller, 1),
)
