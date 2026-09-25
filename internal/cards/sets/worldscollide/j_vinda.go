package worldscollide

import "github.com/dmikalova/vex/internal/card"

// J. Vinda
//
//	House:  Shadows
//	Type:   Creature
//	Rarity: Common
//	Power:  2
//	Traits: Elf • Thief
//
//	Elusive.
//	Reap: Deal 1 damage to a creature. If this damage destroys that creature, steal 1 Æmber.
var JVinda = set.New(
	"J. Vinda",
	card.House.Shadows,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.WC, "242"),
	card.WithPower(2),
	card.WithTraits(card.Traits.Elf, card.Traits.Thief),
	card.WithKeywords(card.Keyword.Elusive),
	card.WithAbility(
		card.Trigger.Reap, card.DealDamage{
			Amount: 1,
			After:  card.IfDestroyed,
			Target: card.Target.Creature,
			Then:   card.StealAember{Amount: 1},
		}),
)
