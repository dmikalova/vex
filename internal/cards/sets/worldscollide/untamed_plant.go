package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Untamed Plant
//
//	House:  Shadows
//	Type:   Creature
//	Rarity: Special
//	Power:  1
//	Traits: Elf • Thief
//
//	Elusive.
//	After a player chooses Untamed as their active house, gain 1 Æmber.
var UntamedPlant = card.New(
	"Untamed Plant",
	card.House.Shadows,
	card.Type.Creature,
	card.Rarity.Special,
	card.Provenance(card.WC, "291"),
	card.WithPower(1),
	card.WithTraits(card.Traits.Elf, card.Traits.Thief),
	card.WithKeywords(card.Keyword.Elusive),
	card.WithAbility(
		card.Trigger.AfterAnyPlayerChoosesHouse, card.Conditional{
			Cond: card.ChoseHouse{House: card.House.Untamed},
			Then: card.GainAember{Player: card.Controller, Amount: 1},
		}),
)
