package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Scientists' Bane
//
//	House:  Dis
//	Type:   Tactic
//	Rarity: Special
//	Æmber:  1
//
//	Play: Destroy a Scientist creature.
var ScientistsBane = card.New(
	"Scientists' Bane",
	card.House.Dis,
	card.Type.Tactic,
	card.Rarity.Special,
	card.Provenance(card.WC, "127"),
	card.WithAemberBonus(1),
	card.WithAbility(
		card.Trigger.Play, card.Destroy{
			Target: card.Target.Creature.WithTrait(card.Traits.Scientist),
		}),
)
