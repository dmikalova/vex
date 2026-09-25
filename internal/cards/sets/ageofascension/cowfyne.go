package ageofascension

import "github.com/dmikalova/vex/internal/card"

// Cowfyne
//
//	House:  Brobnar
//	Type:   Creature
//	Rarity: Common
//	Power:  5
//	Traits: Giant
//
//	Splash-attack 2.
var Cowfyne = set.New(
	"Cowfyne",
	card.House.Brobnar,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.AoA, "5"),
	card.WithPower(5),
	card.WithTraits(card.Traits.Giant),
	card.WithSplashAttack(2),
)
