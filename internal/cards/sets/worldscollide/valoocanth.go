package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Valoocanth
//
//	House:  Brobnar
//	Type:   Creature
//	Rarity: Special
//	Power:  6
//	Traits: Aquan
//
//	While the tide is low, Valoocanth cannot be used.
//	Fight/Reap: Exhaust an enemy creature and each of its neighbors.
var Valoocanth = card.New(
	"Valoocanth",
	card.House.Brobnar,
	card.Type.Creature,
	card.Rarity.Special,
	card.Provenance(card.WC, "A10"),
	card.WithPower(6),
	card.WithTraits(card.Traits.Aquan),
	card.WithCannotBeUsedWhile(card.TideIsLow{}),
	card.WithAbility(card.Trigger.FightReap,
		card.Exhaust{Target: card.Target.EnemyCreature.AndNeighbors()}),
)
