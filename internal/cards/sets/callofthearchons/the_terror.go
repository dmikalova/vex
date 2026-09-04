package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// The Terror
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Common
//	Power:  5
//	Traits: Demon • Knight
//
//	Play: If your opponent has no Æmber, gain 2 Æmber.
var TheTerror = card.New(
	"The Terror",
	card.House.Dis,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.CotA, 101),
	card.WithPower(5),
	card.WithTraits(card.Traits.Demon, card.Traits.Knight),
	card.WithAbility(
		card.Trigger.Play, card.Conditional{
			Cond: card.OpponentAember{Is: card.Exactly, Amount: 0},
			Then: card.GainAember{
				Player: card.Controller,
				Amount: 2,
			},
		}),
)
