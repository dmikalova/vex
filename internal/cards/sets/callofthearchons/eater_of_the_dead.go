package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Eater of the Dead
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Rare
//	Power:  4
//	Traits: Demon
//
//	Fight/Reap: Purge a creature from a discard pile -> give Eater of the Dead a +1 power counter.
var EaterOfTheDead = card.New(
	"Eater of the Dead",
	card.House.Dis,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 84),
	card.WithPower(4),
	card.WithTraits(card.Traits.Demon),
	card.WithFightOrReap(card.Then{
		First: card.PurgeCard{
			Zone: card.Discard,
			Type: card.Type.Creature,
		},
		Result: card.AddPowerCounter{
			Target: card.Target.This,
			Amount: 1,
		},
	}),
)
