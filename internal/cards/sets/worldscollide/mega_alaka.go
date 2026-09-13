package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Mega Alaka
//
//	House:  Brobnar
//	Type:   Creature
//	Rarity: Special
//	Power:  6
//	Traits: Giant
//
//	If you have used a Creature to fight this turn, Mega Alaka enters play ready.
var MegaAlaka = card.New(
	"Mega Alaka",
	card.House.Brobnar,
	card.Type.Creature,
	card.Rarity.Special,
	card.Provenance(card.WC, "54"),
	card.WithPower(6),
	card.WithTraits(card.Traits.Giant),
	card.WithEntersPlay(card.Conditional{
		Cond: card.UsedCreatureToFight{},
		Then: card.Ready{Target: card.Target.This},
	}),
)
