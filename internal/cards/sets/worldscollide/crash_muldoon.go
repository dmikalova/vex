package worldscollide

import "github.com/dmikalova/vex/internal/card"

// Crash Muldoon
//
//	House:  Star Alliance
//	Type:   Creature
//	Rarity: Rare
//	Power:  3
//	Traits: Human • Pilot
//
//	Deploy.
//	Crash Muldoon enters play ready.
//	Action: Use a neighboring non-Star Alliance creature.
var CrashMuldoon = set.New(
	"Crash Muldoon",
	card.House.StarAlliance,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.WC, "327"),
	card.WithPower(3),
	card.WithTraits(card.Traits.Human, card.Traits.Pilot),
	card.WithKeywords(card.Keyword.Deploy),
	card.WithEntersPlay(card.Ready{Target: card.Target.This}),
	card.WithAbility(
		card.Trigger.Action, card.Use{
			Max: 1,
			Target: card.Target.EachFriendlyCreature.Neighboring().
				House(card.Houses.Except(card.House.Self)),
		}),
)
