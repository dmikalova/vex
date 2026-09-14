package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Buzzle
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Common
//	Power:  3
//	Traits: Beast
//
//	Skirmish.
//	Play/Fight: You may purge a neighboring Creature -> ready Buzzle.
var Buzzle = set.New(
	"Buzzle",
	card.House.Dis,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.WC, "70"),
	card.WithPower(3),
	card.WithTraits(card.Traits.Beast),
	card.WithKeywords(card.Keyword.Skirmish),
	card.WithAbility(
		card.Trigger.Play, card.May{Do: card.Then{
			First:  card.PurgeCreature{Target: card.Target.Creature.Neighboring()},
			Result: card.Ready{Target: card.Target.This},
		}}),
	card.WithAbility(
		card.Trigger.Fight, card.May{Do: card.Then{
			First:  card.PurgeCreature{Target: card.Target.Creature.Neighboring()},
			Result: card.Ready{Target: card.Target.This},
		}}),
)
