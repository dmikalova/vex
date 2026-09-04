package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Piranha Monkeys
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Rare
//	Power:  2
//	Traits: Beast
//
//	Play/Reap: Deal 2 damage to each other creature.
var PiranhaMonkeys = card.New(
	"Piranha Monkeys",
	card.House.Untamed,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 365),
	card.WithPower(2),
	card.WithTraits(card.Traits.Beast),
	card.WithPlayReap(card.DealDamage{
		Amount: 2,
		Target: card.Target.EachCreature.Other(),
	}),
)
