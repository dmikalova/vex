package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Hugger-Mugger
//
//	House:  Shadows
//	Type:   Creature
//	Rarity: Common
//	Power:  3
//	Traits: Elf • Thief
//
//	Elusive.
//	Play: Hugger-Mugger captures 1 Æmber from your opponent. If your opponent has more forged keys than you, steal 1 Æmber.
var HuggerMugger = set.New(
	"Hugger-Mugger",
	card.House.Shadows,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.WC, "240"),
	card.WithPower(3),
	card.WithTraits(card.Traits.Elf, card.Traits.Thief),
	card.WithKeywords(card.Keyword.Elusive),
	card.WithAbility(
		card.Trigger.Play, card.Sequence{Effects: []card.Effect{
			card.CaptureAember{
				Amount: 1,
				Target: card.Target.This,
				Source: card.Opponent,
			},
			card.Conditional{
				Cond: card.HasMoreForgedKeys{Player: card.Opponent},
				Then: card.StealAember{Amount: 1},
			},
		}}),
)
