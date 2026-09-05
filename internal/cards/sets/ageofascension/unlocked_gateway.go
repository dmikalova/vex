package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Unlocked Gateway
//
//	House:  Dis
//	Type:   Tactic
//	Rarity: Common
//
//	Omega.
//	Play: Destroy each creature.
var UnlockedGateway = card.New(
	"Unlocked Gateway",
	card.House.Dis,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.AoA, 67),
	card.WithKeywords(card.Keyword.Omega),
	card.WithAbility(
		card.Trigger.Play, card.Destroy{Target: card.Target.EachCreature}),
)
