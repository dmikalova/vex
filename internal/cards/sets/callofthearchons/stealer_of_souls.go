package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Stealer of Souls
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  6
//	Traits: Demon
//
//	After a creature is destroyed fighting Stealer of Souls, purge it, and gain 1 Æmber.
var StealerOfSouls = card.New(
	"Stealer of Souls",
	card.House.Dis,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, 98),
	card.WithPower(6),
	card.WithTraits(card.Traits.Demon),
	card.WithAbility(card.Trigger.AfterDestroyedFighting, card.Sequence{Effects: []card.Effect{
		card.PurgeCreature{Target: card.Target.Triggering},
		card.GainAember{
			Player: card.Controller,
			Amount: 1,
		},
	}}),
)
