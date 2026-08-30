package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Mugwump
//
//	House:  Brobnar
//	Type:   Creature
//	Rarity: Rare
//	Power:  6
//	Traits: Giant
//
//	After a creature is destroyed fighting Mugwump, fully heal Mugwump, and give Mugwump a +1 power counter.
var Mugwump = card.New(
	"Mugwump",
	card.House.Brobnar,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 42),
	card.WithPower(6),
	card.WithTraits("Giant"),
	card.WithAbility(
		card.Trigger.AfterDestroyedFighting, card.Sequence{
			Effects: []card.Effect{
				card.Heal{
					Fully:  true,
					Target: card.Target.This,
				},
				card.AddPowerCounter{
					Target: card.Target.This,
					Amount: 1,
				},
			},
		}),
)
