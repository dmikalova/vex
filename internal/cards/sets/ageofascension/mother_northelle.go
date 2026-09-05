package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Mother Northelle
//
//	House:  Sanctum
//	Type:   Creature
//	Rarity: Rare
//	Power:  2
//	Traits: Human • Monk
//
//	Elusive.
//	Reap: Move 1 Æmber from a friendly creature to your pool.
var MotherNorthelle = card.New(
	"Mother Northelle",
	card.House.Sanctum,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.AoA, 257),
	card.WithPower(2),
	card.WithTraits(card.Traits.Human, card.Traits.Monk),
	card.WithKeywords(card.Keyword.Elusive),
	card.WithAbility(
		card.Trigger.Reap, card.MoveAember{
			Amount: 1,
			From:   card.Target.FriendlyCreature,
			To:     card.Controller,
		}),
)
