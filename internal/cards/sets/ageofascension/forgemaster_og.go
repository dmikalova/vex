package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Forgemaster Og
//
//	House:  Brobnar
//	Type:   Creature
//	Rarity: Rare
//	Power:  4
//	Traits: Giant
//
//	After a player forges a key, that player loses all their Æmber.
var ForgemasterOg = set.New(
	"Forgemaster Og",
	card.House.Brobnar,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.AoA, "38"),
	card.WithPower(4),
	card.WithTraits(card.Traits.Giant),
	card.WithAbility(
		card.Trigger.AfterPlayerForgesKey,
		card.LoseAember{Player: card.ThatPlayer, By: card.AllAember},
	),
)
