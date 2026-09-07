package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Beasts' Bane
//
//	House:  Dis
//	Type:   Tactic
//	Rarity: Special
//	Æmber:  1
//
//	Play: Destroy a Beast creature.
var BeastsBane = card.New(
	"Beasts' Bane",
	card.House.Dis,
	card.Type.Tactic,
	card.Rarity.Special,
	card.Provenance(card.WC, "122"),
	card.WithAemberBonus(1),
	card.WithAbility(
		card.Trigger.Play, card.Destroy{
			Target: card.Target.Creature.WithTrait(card.Traits.Beast),
		}),
)
