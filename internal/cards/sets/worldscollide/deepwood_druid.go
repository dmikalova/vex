package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Deepwood Druid
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Common
//	Power:  3
//	Traits: Elf • Witch
//
//	Deploy.
//	Play/Reap: Fully heal a neighboring creature.
var DeepwoodDruid = set.New(
	"Deepwood Druid",
	card.House.Untamed,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.WC, "355"),
	card.WithPower(3),
	card.WithTraits(card.Traits.Elf, card.Traits.Witch),
	card.WithKeywords(card.Keyword.Deploy),
	card.WithAbility(card.Trigger.PlayReap, card.Heal{
		Fully:  true,
		Target: card.Target.Creature.Neighboring(),
	}),
)
