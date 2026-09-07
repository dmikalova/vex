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
//	Fight: Discard a random card from your opponent's archives.
var Tantadlin = card.New(
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
		card.Trigger.Fight, card.DiscardRandomFromArchives{Player: card.Opponent}),
)
