package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Niffle Queen
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  6
//	Traits: Beast • Niffle
//
//	Each other friendly Beast trait creature gains +1 power.
//	Each other friendly Niffle trait creature gains +1 power.
var NiffleQueen = card.New(
	"Niffle Queen",
	card.House.Untamed,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, 364),
	card.WithPower(6),
	card.WithTraits(card.Traits.Beast, card.Traits.Niffle),
	card.WithConstant(card.ConstantAbility{
		Target:     card.Target.EachOtherFriendlyCreature.WithTrait(card.Traits.Beast),
		PowerBonus: 1,
	}),
	card.WithConstant(card.ConstantAbility{
		Target:     card.Target.EachOtherFriendlyCreature.WithTrait(card.Traits.Niffle),
		PowerBonus: 1,
	}),
)
