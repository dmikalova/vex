package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Aember Imp
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Common
//	Power:  2
//	Traits: Imp
//
//	Elusive.
//	After a creature reaps, stun it.
var AemberImp = card.New(
	"Aember Imp",
	card.House.Dis,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.AoA, 53),
	card.WithPower(2),
	card.WithTraits(card.Traits.Imp),
	card.WithKeywords(card.Keyword.Elusive),
	card.WithAbility(
		card.Trigger.AfterCreatureReaps, card.Stun{Target: card.Target.Triggering}),
)
