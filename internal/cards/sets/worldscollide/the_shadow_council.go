package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// The Shadow Council
//
//	House:  Shadows
//	Type:   Creature
//	Rarity: Rare
//	Power:  3
//	Traits: Elf • Leader • Thief
//
//	Elusive.
//	While The Shadow Council is in the center of your battleline, it gains, "Action: Steal 2 Æmber."
var TheShadowCouncil = card.New(
	"The Shadow Council",
	card.House.Shadows,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.WC, 283),
	card.WithPower(3),
	card.WithTraits(card.Traits.Elf, card.Traits.Leader, card.Traits.Thief),
	card.WithKeywords(card.Keyword.Elusive),
	card.WithConstant(card.ConstantAbility{
		Target:        card.Target.This,
		WhileInCenter: true,
		Granted: []card.Ability{{
			Trigger: card.Trigger.Action,
			Effect:  card.StealAember{Amount: 2},
		}},
	}),
)
