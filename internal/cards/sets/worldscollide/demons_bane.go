package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Demons' Bane
//
//	House:  Dis
//	Type:   Tactic
//	Rarity: Rare
//	Æmber:  1
//
//	Play: Destroy a Demon creature.
var DemonsBane = card.New(
	"Demons' Bane",
	card.House.Dis,
	card.Type.Tactic,
	card.Rarity.Rare,
	card.Provenance(card.WC, 123),
	card.WithAemberBonus(1),
	card.WithAbility(
		card.Trigger.Play, card.Destroy{
			Target: card.Target.Creature.WithTrait(card.Traits.Demon),
		}),
)
