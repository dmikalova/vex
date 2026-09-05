package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Yxilo Bolter
//
//	House:  Mars
//	Type:   Creature
//	Rarity: Common
//	Power:  3
//	Traits: Martian • Soldier
//
//	Fight/Reap: Deal 2 damage to a creature. If this damage destroys that creature, purge it.
var YxiloBolter = card.New(
	"Yxilo Bolter",
	card.House.Mars,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.CotA, 204),
	card.WithPower(3),
	card.WithTraits(card.Traits.Martian, card.Traits.Soldier),
	card.WithFightOrReap(card.DamageThenIfDestroyed{
		Amount: 2,
		Target: card.Target.Creature,
		Then:   card.PurgeCreature{Target: card.Target.Triggering},
	}),
)
