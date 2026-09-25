package ageofascension

import "github.com/dmikalova/vex/internal/card"

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
var Duskwitch = set.New(
	"Duskwitch",
	card.House.Untamed,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.AoA, "320"),
	card.WithPower(1),
	card.WithTraits(card.Traits.Human, card.Traits.Witch),
	card.WithKeywords(card.Keyword.Omega, card.Keyword.Elusive),
	card.WithFriendlyEntersPlayReady(card.EntersReadyGrant{Type: card.Type.Creature}),
)
