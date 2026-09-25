package massmutation

import "github.com/dmikalova/vex/internal/card"

// Boss Zarek
//
//	House:  Shadows
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  3
//	Traits: Mutant • Thief
//
//	Each friendly creature with Æmber on it gains elusive.
//	Enhance Capture Capture Capture.
var BossZarek = set.New(
	"Boss Zarek",
	card.House.Shadows,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.MM, "264"),
	card.WithEnhance(card.Bonus.Capture, card.Bonus.Capture, card.Bonus.Capture),
	card.WithPower(3),
	card.WithTraits(card.Traits.Mutant, card.Traits.Thief),
	card.WithConstant(card.ConstantAbility{
		Target:   card.Target.EachFriendlyCreature.WithAember(),
		Keywords: card.Keywords(card.Keyword.Elusive),
	}),
)
