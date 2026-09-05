package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Roxador
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Rare
//	Power:  4
//	Traits: Beast
//
//	Skirmish.
//	Roxador deals 2 Damage when fighting.
//	Fight: Stun the creature Roxador fights.
var Roxador = card.New(
	"Roxador",
	card.House.Untamed,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.AoA, 365),
	card.WithPower(4),
	card.WithTraits(card.Traits.Beast),
	card.WithKeywords(card.Keyword.Skirmish),
	card.WithAttackDamage(card.AttackDamage{
		Amount: 2,
		Fixed:  true,
	}),
	card.WithAbility(
		card.Trigger.Fight, card.Stun{Target: card.Target.CreatureFought}),
)
