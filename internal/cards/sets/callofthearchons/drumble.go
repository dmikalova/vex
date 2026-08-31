package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Drumble
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  2
//	Traits: Imp
//
//	Elusive.
//	Play: If your opponent has 7 Æmber or more, Drumble captures all your opponent's Æmber.
var Drumble = card.New(
	"Drumble",
	card.House.Dis,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, 82),
	card.WithPower(2),
	card.WithTraits("Imp"),
	card.WithKeywords(card.Keyword.Elusive),
	card.WithAbility(
		card.Trigger.Play, card.Conditional{
			Cond: card.OpponentAember{Is: card.AtLeast, Amount: 7},
			Then: card.CaptureAember{
				All:    true,
				Target: card.Target.This,
				Source: card.Opponent,
			},
		}),
)
