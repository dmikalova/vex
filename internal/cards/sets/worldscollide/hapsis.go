package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Hapsis
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Common
//	Power:  5
//	Traits: Mutant • Scientist
//
//	After a creature is destroyed in a fight with Hapsis, ward Hapsis. Draw a card.
var Hapsis = set.New(
	"Hapsis",
	card.House.Logos,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.WC, "134"),
	card.WithPower(5),
	card.WithTraits(card.Traits.Mutant, card.Traits.Scientist),
	card.WithAbility(
		card.Trigger.AfterDestroyedFighting, card.Sequence{
			Effects: []card.Effect{
				card.Ward{Target: card.Target.This},
				card.Draw{Amount: 1},
			},
		}),
)
