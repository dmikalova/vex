package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Project Z.Y.X.
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Rare
//	Power:  5
//	Armor:  1
//	Traits: Cyborg • Mutant
//
//	Fight/Reap: You may play a card from your archives.
var ProjectZYX = card.New(
	"Project Z.Y.X.",
	card.House.Logos,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.AoA, "152"),
	card.WithPower(5),
	card.WithArmor(1),
	card.WithTraits(card.Traits.Cyborg, card.Traits.Mutant),
	card.WithAbility(card.Trigger.FightReap, card.May{Do: card.PlayFrom{From: card.Archives}}),
)
