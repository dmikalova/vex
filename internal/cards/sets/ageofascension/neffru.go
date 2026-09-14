package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Neffru
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Rare
//	Power:  4
//	Traits: Demon
//
//	After a Creature is destroyed, its owner gains 1 Æmber.
var Neffru = set.New(
	"Neffru",
	card.House.Dis,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.AoA, "94"),
	card.WithPower(4),
	card.WithTraits(card.Traits.Demon),
	card.WithAbility(
		card.Trigger.AfterCreatureDestroyed, card.GainAember{
			Player: card.ItsOwner,
			Amount: 1,
		}),
)
