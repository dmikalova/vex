package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Scowly Caper
//
//	House:  Shadows
//	Type:   Creature
//	Rarity: Rare
//	Power:  2
//	Traits: Elf • Thief
//
//	Skirmish, Treachery, Versatile.
//	At the end of your turn, destroy a neighboring Creature.
var ScowlyCaper = set.New(
	"Scowly Caper",
	card.House.Shadows,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.AoA, "313"),
	card.WithPower(2),
	card.WithTraits(card.Traits.Elf, card.Traits.Thief),
	card.WithKeywords(card.Keyword.Skirmish, card.Keyword.Treachery, card.Keyword.Versatile),
	card.WithAbility(
		card.Trigger.EndOfTurn, card.Destroy{
			Target: card.Target.Creature.Neighboring(),
		}),
)
