package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Daughter
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Common
//	Power:  2
//	Traits: Cyborg • Scientist
//
//	Elusive.
//	During your "draw cards" phase, refill your hand to 1 additional card.
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
