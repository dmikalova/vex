package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// The Grim Reaper
//
//	House:  Brobnar
//	Type:   Creature
//	Rarity: Special
//	Power:  4
//	Traits: Robot • Specter
//
//	If you are haunted, The Grim Reaper enters play ready.
//	Reap: Purge an enemy Creature, and purge a friendly Creature.
var TheGrimReaper = card.New(
	"The Grim Reaper",
	card.House.Brobnar,
	card.Type.Creature,
	card.Rarity.Special,
	card.Provenance(card.WC, "A07"),
	card.WithPower(4),
	card.WithTraits(card.Traits.Robot, card.Traits.Specter),
	card.WithEntersPlay(card.Conditional{
		Cond: card.Haunted{},
		Then: card.Ready{Target: card.Target.This},
	}),
	card.WithAbility(
		card.Trigger.Reap, card.Sequence{
			Effects: []card.Effect{
				card.PurgeCreature{Target: card.Target.EnemyCreature},
				card.PurgeCreature{Target: card.Target.FriendlyCreature},
			},
		}),
)
