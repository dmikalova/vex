package worldscollide

import "github.com/dmikalova/vex/internal/card"

// Imprinted Murmook
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Common
//	Power:  3
//	Traits: Beast
//
//	Elusive.
//	Your keys cost -1 Æmber.
var ImprintedMurmook = set.New(
	"Imprinted Murmook",
	card.House.Untamed,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.WC, "358"),
	card.WithPower(3),
	card.WithTraits(card.Traits.Beast),
	card.WithKeywords(card.Keyword.Elusive),
	card.WithKeyCost(card.KeyCostChange(card.Controller, -1)),
)
