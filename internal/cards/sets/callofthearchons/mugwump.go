package callofthearchons

import "github.com/dmikalova/vex/internal/card"

// Mugwump
//
//	House:  Brobnar
//	Type:   Creature
//	Rarity: Rare
//	Power:  6
//	Traits: Giant
//
//	After a creature is destroyed in a fight with Mugwump, fully heal Mugwump. Give Mugwump a +1 power counter.
var Mugwump = set.New(
	"Mugwump",
	card.House.Brobnar,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.CotA, "42"),
	card.WithPower(6),
	card.WithTraits(card.Traits.Giant),
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
