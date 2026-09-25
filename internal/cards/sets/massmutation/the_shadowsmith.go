package massmutation

import "github.com/dmikalova/vex/internal/card"

// The Shadowsmith
//
//	House:  Shadows
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  3
//	Traits: Mutant • Thief
//
//	Each Mutant creature gains elusive.
var TheShadowsmith = set.New(
	"The Shadowsmith",
	card.House.Shadows,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.MM, "276"),
	card.WithPower(3),
	card.WithTraits(card.Traits.Mutant, card.Traits.Thief),
	card.WithConstant(card.ConstantAbility{
		Keywords: card.Keywords(card.Keyword.Elusive),
		Target:   card.Target.EachCreature.WithTrait(card.Traits.Mutant),
	}),
)
