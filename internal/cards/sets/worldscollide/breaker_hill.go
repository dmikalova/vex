package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Breaker Hill
//
//	House:  Shadows
//	Type:   Creature
//	Rarity: Common
//	Power:  1
//	Traits: Elf • Thief
//
//	Elusive.
//	Each neighboring creature gains, "Action: Steal 1 Æmber."
var BreakerHill = card.New(
	"Breaker Hill",
	card.House.Shadows,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.WC, "237"),
	card.WithPower(1),
	card.WithTraits(card.Traits.Elf, card.Traits.Thief),
	card.WithKeywords(card.Keyword.Elusive),
	card.WithConstant(card.ConstantAbility{
		Target: card.Target.EachCreature.Neighboring(),
		Granted: []card.Ability{{
			Trigger: card.Trigger.Action,
			Effect:  card.StealAember{Amount: 1},
		}},
	}),
)
