//go:build todo

package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// AgentHooMan
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Mars
//	Type:   Creature
//	Rarity: Common
//	Power:  2
//	Traits: Martian • Agent
//
//	Elusive. (The first time this creature is attacked each turn, no damage is dealt.)
//	Reap: Choose a friendly non-Mars creature and an enemy non-Mars creature. Stun the chosen creatures.
var AgentHooMan = card.New(
	"Agent Hoo-man",
	card.House.Mars,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.AoA, 160),
	card.WithPower(2),
	card.WithTraits(card.Traits.Martian, card.Traits.Agent),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
