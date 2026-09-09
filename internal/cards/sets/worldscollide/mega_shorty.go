package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Mega Shorty
//
//	House:  Brobnar
//	Type:   Creature
//	Rarity: Special
//	Power:  6
//	Traits: Giant
//
//	Assault 4.
//	Reap: Enrage Mega Shorty.
var MegaShorty = card.New(
	"Mega Shorty",
	card.House.Brobnar,
	card.Type.Creature,
	card.Rarity.Special,
	card.Provenance(card.WC, "61"),
	card.WithPower(6),
	card.WithTraits(card.Traits.Giant),
	card.WithAssault(4),
	card.WithAbility(
		card.Trigger.Reap, card.Enrage{Target: card.Target.This}),
)
