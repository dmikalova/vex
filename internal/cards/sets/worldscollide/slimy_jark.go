package worldscollide

import "github.com/dmikalova/vex/internal/card"

// Slimy Jark
//
//	House:  Brobnar
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  1
//	Traits: Goblin
//
//	Skirmish, Elusive.
//	Fight: Enrage the creature Slimy Jark fought.
var SlimyJark = set.New(
	"Slimy Jark",
	card.House.Brobnar,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.WC, "30"),
	card.WithPower(1),
	card.WithTraits(card.Traits.Goblin),
	card.WithKeywords(card.Keyword.Skirmish, card.Keyword.Elusive),
	card.WithAbility(
		card.Trigger.Fight, card.Enrage{Target: card.Target.CreatureFought}),
)
