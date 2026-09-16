package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Reckless Rizzo
//
//	House:  Shadows
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  1
//	Traits: Elf • Thief
//
//	Elusive.
//	Action: Steal 2 Æmber. Until the start of your next turn, Reckless Rizzo loses elusive.
var RecklessRizzo = set.New(
	"Reckless Rizzo",
	card.House.Shadows,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.MM, "273"),
	card.WithPower(1),
	card.WithTraits(card.Traits.Elf, card.Traits.Thief),
	card.WithKeywords(card.Keyword.Elusive),
	card.WithAbility(
		card.Trigger.Action, card.Sentences{Effects: []card.Effect{
			card.StealAember{Amount: 2},
			card.LoseKeywordsUntilNextTurn{
				Target:   card.Target.This,
				Keywords: []card.KeywordValue{card.Keyword.Elusive},
			},
		}}),
)
