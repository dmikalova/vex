package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Carlo Phantom
//
//	House:  Shadows
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  1
//	Traits: Elf • Thief
//
//	Elusive, Skirmish.
//	After you play an artifact, steal 1 Æmber.
var CarloPhantom = card.New(
	"Carlo Phantom",
	card.House.Shadows,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, 298),
	card.WithPower(1),
	card.WithTraits(card.Traits.Elf, card.Traits.Thief),
	card.WithKeywords(card.Keyword.Elusive, card.Keyword.Skirmish),
	card.WithAbility(card.Trigger.AfterCardPlayed, card.Conditional{
		Cond: card.ItIs{Type: card.Type.Artifact},
		Then: card.StealAember{Amount: 1},
	}),
)
