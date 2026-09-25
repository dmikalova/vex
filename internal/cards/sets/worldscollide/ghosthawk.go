package worldscollide

import "github.com/dmikalova/vex/internal/card"

// Ghosthawk
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Common
//	Power:  2
//	Traits: Beast
//
//	Deploy.
//	Play: Reap with each of Ghosthawk's neighbors, one at a time.
var Ghosthawk = set.New(
	"Ghosthawk",
	card.House.Untamed,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.WC, "356"),
	card.WithPower(2),
	card.WithTraits(card.Traits.Beast),
	card.WithKeywords(card.Keyword.Deploy),
	card.WithAbility(
		card.Trigger.Play, card.OneAtATime{
			Target: card.Target.EachNeighbor,
			Verbs:  []card.CreatureVerb{card.ReapVerb{}},
		}),
)
