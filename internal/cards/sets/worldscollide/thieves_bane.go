package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Thieves' Bane
//
//	House:  Dis
//	Type:   Tactic
//	Rarity: Rare
//	Æmber:  1
//
//	Play: Destroy a Thief creature.
var ThievesBane = card.New(
	"Thieves' Bane",
	card.House.Dis,
	card.Type.Tactic,
	card.Rarity.Rare,
	card.Provenance(card.WC, 128),
	card.WithAemberBonus(1),
	card.WithAbility(
		card.Trigger.Play, card.Destroy{
			Target: card.Target.Creature.WithTrait(card.Traits.Thief),
		}),
)
