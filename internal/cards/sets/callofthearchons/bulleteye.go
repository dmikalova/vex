package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Bulleteye
//
//	House:  Shadows
//	Type:   Creature
//	Rarity: Rare
//	Power:  2
//	Traits: Elf • Thief
//
//	Elusive.
//	Reap: Destroy a flank creature.
var Bulleteye = card.New(
	"Bulleteye",
	card.House.Shadows,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 297),
	card.WithPower(2),
	card.WithTraits(card.Traits.Elf, card.Traits.Thief),
	card.WithKeywords(card.Keyword.Elusive),
	card.WithAbility(
		card.Trigger.Reap, card.Destroy{Target: card.Target.Creature.OnFlank()}),
)
