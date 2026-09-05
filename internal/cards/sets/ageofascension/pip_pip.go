package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Pip Pip
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Common
//	Power:  2
//	Traits: Human • Scientist
//
//	After an enemy creature reaps, stun it.
var PipPip = card.New(
	"Pip Pip",
	card.House.Logos,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.AoA, 116),
	card.WithPower(2),
	card.WithTraits(card.Traits.Human, card.Traits.Scientist),
	card.WithAbility(
		card.Trigger.AfterEnemyCreatureReaps, card.Stun{Target: card.Target.Triggering}),
)
