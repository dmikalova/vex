package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// GleefulMayhem
//
//	House:  Dis
//	Type:   Action
//	Rarity: Uncommon
//	Æmber:  1
//
//	Play: For each house, deal 5D to a creature of that house.
var GleefulMayhem = card.New(
	"Gleeful Mayhem",
	card.House.Dis,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.WC, "090"),
	card.WithAemberBonus(1),
	card.WithAbility(card.Trigger.Play, card.DealDamagePerHouse{Amount: 5}),
)
