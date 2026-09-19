package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Skixuno
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Rare
//	Power:  1
//	Traits: Demon
//
//	Omega.
//	Play: Destroy each other creature. For each creature destroyed this way, give Skixuno a +1 power counter.
var Skixuno = set.New(
	"Skixuno",
	card.House.Dis,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.MM, "048"),
	card.WithPower(1),
	card.WithTraits(card.Traits.Demon),
	card.WithKeywords(card.Keyword.Omega),
	card.WithAbility(
		card.Trigger.Play, card.Sequence{Effects: []card.Effect{
			card.Destroy{Target: card.Target.EachCreature.Other()},
			card.AddPowerCounter{
				Target: card.Target.This,
				Amount: 1,
				Per:    card.CreaturesDestroyed{},
			},
		}},
	),
)
