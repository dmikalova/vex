package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Mega Alaka
//
//	House:  Brobnar
//	Type:   Creature
//	Rarity: Connected
//	Power:  6
//	Traits: Giant
//
//	If you have used a Creature to fight this turn, Mega Alaka enters play ready.
var MegaAlaka = set.New(
	"Mega Alaka",
	card.House.Brobnar,
	card.Type.Creature,
	card.Rarity.Connected,
	card.Provenance(card.WC, "54"),
	card.InCluster(card.Pulled(alakasBrewCluster, 1, 1.25)),
	card.WithPower(6),
	card.WithTraits(card.Traits.Giant),
	card.WithEntersPlay(card.Conditional{
		Cond: card.UsedCreatureToFight{},
		Then: card.Ready{Target: card.Target.This},
	}),
)
