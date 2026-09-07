package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Giants' Bane
//
//	House:  Dis
//	Type:   Tactic
//	Rarity: Uncommon
//	Æmber:  1
//
//	Play: Destroy a Giant creature.
var GiantsBane = card.New(
	"Giants' Bane",
	card.House.Dis,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.WC, "89"),
	card.WithAemberBonus(1),
	card.WithAbility(
		card.Trigger.Play, card.Destroy{
			Target: card.Target.Creature.WithTrait(card.Traits.Giant),
		}),
)
