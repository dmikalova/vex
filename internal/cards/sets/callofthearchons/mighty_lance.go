package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Mighty Lance
//
//	House:  Sanctum
//	Type:   Tactic
//	Rarity: Rare
//
//	Play: Choose a creature. Deal 3 damage to the chosen creature and one of its neighbors.
var MightyLance = set.New(
	"Mighty Lance",
	card.House.Sanctum,
	card.Type.Tactic,
	card.Rarity.Rare,
	card.Provenance(card.CotA, "221"),
	card.WithAbility(card.Trigger.Play, card.DealDamage{Spread: card.CreatureAndNeighbors{
		Amount: 3,
		Splash: 3,
		Scope:  card.OneNeighbor,
	}}),
)
