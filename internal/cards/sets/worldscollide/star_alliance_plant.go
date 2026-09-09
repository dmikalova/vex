package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Star Alliance Plant
//
//	House:  Shadows
//	Type:   Creature
//	Rarity: Special
//	Power:  1
//	Traits: Elf • Thief
//
//	Elusive.
//	After a player chooses Star Alliance as their active house, gain 1 Æmber.
var StarAlliancePlant = card.New(
	"Star Alliance Plant",
	card.House.Shadows,
	card.Type.Creature,
	card.Rarity.Special,
	card.Provenance(card.WC, "290"),
	card.WithPower(1),
	card.WithTraits(card.Traits.Elf, card.Traits.Thief),
	card.WithKeywords(card.Keyword.Elusive),
	card.WithAbility(
		card.Trigger.AfterAnyPlayerChoosesHouse, card.Conditional{
			Cond: card.ChoseHouse{House: card.House.StarAlliance},
			Then: card.GainAember{Player: card.Controller, Amount: 1},
		}),
)
