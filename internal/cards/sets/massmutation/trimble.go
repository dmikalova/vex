package massmutation

import "github.com/dmikalova/vex/internal/card"

// Trimble
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  3
//	Traits: Mutant
//
//	Each Mutant creature gains skirmish.
var Trimble = set.New(
	"Trimble",
	card.House.Untamed,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.MM, "389"),
	card.WithPower(3),
	card.WithTraits(card.Traits.Mutant),
	card.WithConstant(card.ConstantAbility{
		Keywords: card.Keywords(card.Keyword.Skirmish),
		Target:   card.Target.EachCreature.WithTrait(card.Traits.Mutant),
	}),
)
