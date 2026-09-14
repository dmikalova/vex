package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Grimlocus Dux
//
//	House:  Saurian
//	Type:   Creature
//	Rarity: Rare
//	Power:  11
//	Armor:  2
//	Traits: Dinosaur • Soldier
//
//	Taunt.
//	Play: Exalt Grimlocus Dux 2 times.
var GrimlocusDux = set.New(
	"Grimlocus Dux",
	card.House.Saurian,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.WC, "221"),
	card.WithPower(11),
	card.WithArmor(2),
	card.WithTraits(card.Traits.Dinosaur, card.Traits.Soldier),
	card.WithKeywords(card.Keyword.Taunt),
	card.WithAbility(
		card.Trigger.Play, card.Exalt{
			Target: card.Target.This,
			Amount: 2,
		}),
)
