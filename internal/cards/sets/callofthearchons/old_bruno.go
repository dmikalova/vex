package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Old Bruno
//
//	House:  Shadows
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  3
//	Traits: Elf • Thief
//
//	Elusive.
//	Play: Old Bruno captures 3 Æmber from your opponent.
var OldBruno = card.New(
	"Old Bruno",
	card.House.Shadows,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, 307),
	card.WithPower(3),
	card.WithTraits(card.Traits.Elf, card.Traits.Thief),
	card.WithKeywords(card.Keyword.Elusive),
	card.WithAbility(
		card.Trigger.Play, card.CaptureAember{
			Amount: 3,
			Target: card.Target.This,
			Source: card.Opponent,
		}),
)
