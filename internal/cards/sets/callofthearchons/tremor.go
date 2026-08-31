package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Tremor
//
//	House:  Brobnar
//	Type:   Tactic
//	Rarity: Common
//
//	Play: Stun a creature and each of its neighbors.
var Tremor = card.New(
	"Tremor",
	card.House.Brobnar,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.CotA, 16),
	card.WithAbility(
		card.Trigger.Play, card.Stun{Target: card.Target.Creature.AndNeighbors()}),
)
