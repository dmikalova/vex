package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Duskwitch
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Common
//	Power:  1
//	Traits: Human • Witch
//
//	Omega, Elusive.
//	Your creatures enter play ready.
var Duskwitch = card.New(
	"Duskwitch",
	card.House.Untamed,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.AoA, "320"),
	card.WithPower(1),
	card.WithTraits(card.Traits.Human, card.Traits.Witch),
	card.WithKeywords(card.Keyword.Omega, card.Keyword.Elusive),
	card.WithFriendlyEntersPlayReady(card.Type.Creature),
)
