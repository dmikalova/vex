package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Optio Gorkus
//
//	House:  Saurian
//	Type:   Creature
//	Rarity: Rare
//	Power:  3
//	Armor:  3
//	Traits: Dinosaur • Soldier
//
//	Elusive.
//	Each of {self}'s neighbors gains, "Destroyed: Move all Æmber from this creature to Optio Gorkus."
var OptioGorkus = set.New(
	"Optio Gorkus",
	card.House.Saurian,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.MM, "226"),
	card.WithPower(3),
	card.WithArmor(3),
	card.WithTraits(card.Traits.Dinosaur, card.Traits.Soldier),
	card.WithKeywords(card.Keyword.Elusive),
	card.WithConstant(card.ConstantAbility{
		Target: card.Target.EachNeighbor,
		Granted: []card.Ability{{
			Trigger: card.Trigger.Destroyed,
			Effect: card.MoveAember{
				All:  true,
				From: card.Target.This,
				Onto: card.Target.GrantingCard,
			},
		}},
	}),
)
