package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Marmo Swarm
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Common
//	Power:  2
//	Traits: Beast
//
//	Marmo Swarm gains +1 power for each Æmber in your pool.
var MarmoSwarm = card.New(
	"Marmo Swarm",
	card.House.Untamed,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.AoA, 327),
	card.WithPower(2),
	card.WithTraits(card.Traits.Beast),
	card.WithConstant(card.ConstantAbility{
		Target:     card.Target.This,
		PowerBonus: 1,
		Per:        card.AemberInPool{Player: card.Controller},
	}),
)
