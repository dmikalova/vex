package ageofascension

import "github.com/dmikalova/vex/internal/card"

// Po's Pixies
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Rare
//	Power:  1
//	Traits: Faerie
//
//	Elusive.
//	Æmber stolen or captured from your pool is taken from the common supply instead.
var PosPixies = set.New(
	"Po's Pixies",
	card.House.Untamed,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.AoA, "362"),
	card.WithPower(1),
	card.WithTraits(card.Traits.Faerie),
	card.WithKeywords(card.Keyword.Elusive),
	card.WithReplaces(card.Instead{
		Of:     card.Event.AemberTakenFromPool,
		Player: card.Controller,
		With:   card.FromCommonSupply,
	}),
)
