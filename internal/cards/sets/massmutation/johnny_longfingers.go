package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Johnny Longfingers
//
//	House:  Shadows
//	Type:   Creature
//	Rarity: Rare
//	Power:  4
//	Traits: Mutant • Thief
//
//	Each friendly Mutant creature gains, "Destroyed: Steal 1 Æmber."
var JohnnyLongfingers = set.New(
	"Johnny Longfingers",
	card.House.Shadows,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.MM, "283"),
	card.WithPower(4),
	card.WithTraits(card.Traits.Mutant, card.Traits.Thief),
	card.WithConstant(card.ConstantAbility{
		Target: card.Target.EachFriendlyCreature.WithTrait(card.Traits.Mutant),
		Granted: []card.Ability{{
			Trigger: card.Trigger.Destroyed,
			Effect:  card.StealAember{Amount: 1},
		}},
	}),
)
