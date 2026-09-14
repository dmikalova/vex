package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Ulyq Megamouth
//
//	House:  Mars
//	Type:   Creature
//	Rarity: Common
//	Power:  3
//	Traits: Martian • Scientist
//
//	Fight/Reap: Use a friendly non-Mars Creature.
var UlyqMegamouth = set.New(
	"Ulyq Megamouth",
	card.House.Mars,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.CotA, "200"),
	card.WithPower(3),
	card.WithTraits(card.Traits.Martian, card.Traits.Scientist),
	card.WithAbility(card.Trigger.FightReap, card.OnChooseCreature{
		Target: card.Target.FriendlyCreature.House(card.Houses.Except(card.House.Self)),
		Verbs:  []card.CreatureVerb{card.UseVerb{}},
	}),
)
