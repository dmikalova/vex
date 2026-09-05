package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Ixxyxli Fixfinger
//
//	House:  Mars
//	Type:   Creature
//	Rarity: Common
//	Power:  2
//	Armor:  2
//	Traits: Martian • Scientist
//
//	Elusive.
//	Each other friendly Mars creature gains +1 armor.
var IxxyxliFixfinger = card.New(
	"Ixxyxli Fixfinger",
	card.House.Mars,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.AoA, 164),
	card.WithPower(2),
	card.WithArmor(2),
	card.WithTraits(card.Traits.Martian, card.Traits.Scientist),
	card.WithKeywords(card.Keyword.Elusive),
	card.WithConstant(card.ConstantAbility{
		ArmorBonus: 1,
		Target:     card.Target.EachOtherFriendlyCreature.OfHouse(card.House.Self),
	}),
)
