package worldscollide

import "github.com/dmikalova/vex/internal/card"

// Mega Shorty
//
//	House:  Brobnar
//	Type:   Creature
//	Rarity: Connected
//	Power:  6
//	Traits: Giant
//
//	Assault 4.
//	Reap: Enrage Mega Shorty.
var MegaShorty = set.New(
	"Mega Shorty",
	card.House.Brobnar,
	card.Type.Creature,
	card.Rarity.Connected,
	card.Provenance(card.WC, "61"),
	card.InCluster(card.Pulled(shortysBrewCluster, 1, 1.25)),
	card.WithPower(6),
	card.WithTraits(card.Traits.Giant),
	card.WithAssault(4),
	card.WithAbility(
		card.Trigger.Reap, card.Enrage{Target: card.Target.This}),
)
