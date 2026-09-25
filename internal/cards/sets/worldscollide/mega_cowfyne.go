package worldscollide

import "github.com/dmikalova/vex/internal/card"

// Mega Cowfyne
//
//	House:  Brobnar
//	Type:   Creature
//	Rarity: Connected
//	Power:  7
//	Traits: Giant
//
//	Splash-attack 2.
var MegaCowfyne = set.New(
	"Mega Cowfyne",
	card.House.Brobnar,
	card.Type.Creature,
	card.Rarity.Connected,
	card.Provenance(card.WC, "55"),
	card.InCluster(card.Pulled(cowfynesBrewCluster, 1, 1.25)),
	card.WithPower(7),
	card.WithTraits(card.Traits.Giant),
	card.WithSplashAttack(2),
)
