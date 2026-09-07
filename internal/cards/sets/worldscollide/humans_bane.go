package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Humans' Bane
//
//	House:  Dis
//	Type:   Tactic
//	Rarity: Rare
//	Æmber:  1
//
//	Play: Destroy a Human creature.
var HumansBane = card.New(
	"Humans' Bane",
	card.House.Dis,
	card.Type.Tactic,
	card.Rarity.Rare,
	card.Provenance(card.WC, 126),
	card.WithAemberBonus(1),
	card.WithAbility(
		card.Trigger.Play, card.Destroy{
			Target: card.Target.Creature.WithTrait(card.Traits.Human),
		}),
)
