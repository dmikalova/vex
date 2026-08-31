package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Mighty Lance
//
//	House:  Sanctum
//	Type:   Tactic
//	Rarity: Rare
//
//	Play: Deal 3 damage to a creature and 3 damage to a neighbor of that creature.
var MightyLance = card.New(
	"Mighty Lance",
	card.House.Sanctum,
	card.Type.Tactic,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 221),
	card.WithAbility(card.Trigger.Play, card.DamageCreatureAndNeighbor{
		Amount:         3,
		NeighborAmount: 3,
	}),
)
