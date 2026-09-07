package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Panpaca, Jaga
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  3
//	Traits: Beast
//
//	Skirmish.
//	Each creature to the left of Panpaca, Jaga gains skirmish.
var PanpacaJaga = card.New(
	"Panpaca, Jaga",
	card.House.Untamed,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.AoA, "348"),
	card.WithPower(3),
	card.WithTraits(card.Traits.Beast),
	card.WithKeywords(card.Keyword.Skirmish),
	card.WithConstant(card.ConstantAbility{
		Target:   card.Target.EachCreature.ToLeftOfSource(),
		Keywords: card.Keywords(card.Keyword.Skirmish),
	}),
)
