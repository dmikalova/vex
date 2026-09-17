package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Tantadlin
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Common
//	Power:  9
//	Traits: Tree
//
//	Tantadlin deals 2 Damage when fighting.
//	Fight: Your opponent discards a random card from their archives.
var Tantadlin = set.New(
	"Tantadlin",
	card.House.Untamed,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.AoA, "333"),
	card.WithPower(9),
	card.WithTraits(card.Traits.Tree),
	card.WithAttackDamage(card.AttackDamage{
		Amount: 2,
		Fixed:  true,
	}),
	card.WithAbility(
		card.Trigger.Fight, card.DiscardCard{
			Player:    card.Opponent,
			Zones:     []card.Zone{card.Archives},
			Selection: card.Random{Count: 1},
		}),
)
