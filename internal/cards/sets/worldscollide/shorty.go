package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Shorty
//
//	House:  Brobnar
//	Type:   Creature
//	Rarity: Common
//	Power:  4
//	Traits: Giant
//
//	Assault 4.
//	Reap: Enrage Shorty.
var Shorty = card.New(
	"Shorty",
	card.House.Brobnar,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.WC, "13"),
	card.WithPower(4),
	card.WithTraits(card.Traits.Giant),
	card.WithAssault(4),
	card.WithAbility(
		card.Trigger.Reap, card.Enrage{Target: card.Target.This}),
)
