package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Grovekeeper
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Common
//	Power:  3
//	Traits: Human • Witch
//
//	At the end of your turn, give each neighboring creature a +1 power counter.
var Grovekeeper = card.New(
	"Grovekeeper",
	card.House.Untamed,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.AoA, 324),
	card.WithPower(3),
	card.WithTraits(card.Traits.Human, card.Traits.Witch),
	card.WithAbility(
		card.Trigger.EndOfTurn, card.AddPowerCounter{
			Target: card.Target.EachCreature.Neighboring(),
			Amount: 1,
		}),
)
