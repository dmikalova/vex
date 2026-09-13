package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Gleeful Mayhem
//
//	House:  Dis
//	Type:   Tactic
//	Rarity: Uncommon
//	Æmber:  1
//
//	Play: For each house, deal 5 damage to a Creature of that house.
var GleefulMayhem = card.New(
	"Gleeful Mayhem",
	card.House.Dis,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.WC, "090"),
	card.WithAemberBonus(1),
	card.WithAbility(card.Trigger.Play, card.DealDamagePerHouse{Amount: 5}),
)
