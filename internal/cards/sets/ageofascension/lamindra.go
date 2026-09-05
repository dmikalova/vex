package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Lamindra
//
//	House:  Shadows
//	Type:   Creature
//	Rarity: Common
//	Power:  1
//	Traits: Elf • Thief
//
//	Deploy, Elusive.
//	Each neighboring creature gains elusive.
var Lamindra = card.New(
	"Lamindra",
	card.House.Shadows,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.AoA, 272),
	card.WithPower(1),
	card.WithTraits(card.Traits.Elf, card.Traits.Thief),
	card.WithKeywords(card.Keyword.Deploy, card.Keyword.Elusive),
	card.WithConstant(card.ConstantAbility{
		Keywords: card.Keywords(card.Keyword.Elusive),
		Target:   card.Target.EachCreature.Neighboring(),
	}),
)
