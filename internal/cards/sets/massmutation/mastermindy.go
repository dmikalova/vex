package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Mastermindy
//
//	House:  Shadows
//	Type:   Creature
//	Rarity: Rare
//	Power:  2
//	Traits: Elf • Thief
//
//	Elusive.
//	At the end of your turn, put a scheme counter on Mastermindy.
//	Action: For each scheme counter on Mastermindy, steal 1 Æmber. Remove each scheme counter from Mastermindy.
var Mastermindy = set.New(
	"Mastermindy",
	card.House.Shadows,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.MM, "285"),
	card.WithPower(2),
	card.WithTraits(card.Traits.Elf, card.Traits.Thief),
	card.WithKeywords(card.Keyword.Elusive),
	card.WithAbility(
		card.Trigger.EndOfTurn, card.PlaceCounter{
			Kind:   card.Counter.Scheme,
			Target: card.Target.This,
		}),
	card.WithAbility(
		card.Trigger.Action, card.Sentences{Effects: []card.Effect{
			card.StealAember{
				Amount: 1,
				Per:    card.CountersOnThis{Kind: card.Counter.Scheme},
			},
			card.RemoveCounters{Kind: card.Counter.Scheme, Target: card.Target.This},
		}}),
)
