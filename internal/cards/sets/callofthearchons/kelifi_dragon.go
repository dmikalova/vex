package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Kelifi Dragon
//
//	House:  Brobnar
//	Type:   Creature
//	Rarity: Rare
//	Power:  12
//	Traits: Dragon
//
//	Kelifi Dragon cannot be played unless you have 7 Æmber or more.
//	Fight/Reap: Gain 1 Æmber, and deal 5 damage to a creature.
var KelifiDragon = card.New("Kelifi Dragon",
	card.House.Brobnar, card.Type.Creature, card.Rarity.Rare,
	card.Provenance(card.CotA, 37),
	card.WithPower(12),
	card.WithTraits("Dragon"),
	card.WithAemberThreshold(7),
	card.WithFightOrReap(card.Sequence{Effects: []card.Effect{
		card.GainAember{Player: card.Controller, Amount: 1},
		card.DealDamage{
			Target: card.Target.Creature,
			Amount: 5,
		},
	}}),
)
