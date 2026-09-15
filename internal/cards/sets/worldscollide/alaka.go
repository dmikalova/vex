package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Alaka
//
//	House:  Brobnar
//	Type:   Creature
//	Rarity: Common
//	Power:  4
//	Traits: Giant
//
//	If you have used a creature to fight this turn, Alaka enters play ready.
var Alaka = set.New(
	"Alaka",
	card.House.Brobnar,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.WC, "1"),
	card.WithPower(4),
	card.WithTraits(card.Traits.Giant),
	card.WithEntersPlay(card.Conditional{
		Cond: card.UsedCreatureToFight{},
		Then: card.Ready{Target: card.Target.This},
	}),
)
