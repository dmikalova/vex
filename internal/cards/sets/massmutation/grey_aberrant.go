package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Grey Aberrant
//
//	House:  Sanctum
//	Type:   Creature
//	Rarity: Common
//	Power:  3
//	Traits: Monk • Mutant
//
//	Each creature loses each of its traits.
var GreyAberrant = set.New(
	"Grey Aberrant",
	card.House.Sanctum,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.MoMu, "181"),
	card.WithPower(3),
	card.WithTraits(card.Traits.Monk, card.Traits.Mutant),
	card.WithConstant(card.ConstantAbility{
		Target:        card.Target.EachCreature,
		RemovesTraits: true,
	}),
)
