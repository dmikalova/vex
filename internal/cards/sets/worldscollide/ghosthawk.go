package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Ghosthawk
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Common
//	Power:  2
//	Traits: Beast
//
//	Deploy.
//	Play: You may reap with up to 2 different neighboring Creatures, one at a time.
var Ghosthawk = card.New(
	"Ghosthawk",
	card.House.Untamed,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.WC, "356"),
	card.WithPower(2),
	card.WithTraits(card.Traits.Beast),
	card.WithKeywords(card.Keyword.Deploy),
	card.WithAbility(
		card.Trigger.Play, card.May{Do: card.OneAtATime{
			Times:  card.Fixed(2),
			Target: card.Target.Creature.Neighboring(),
			Verbs:  []card.CreatureVerb{card.ReapVerb{}},
		}}),
)
