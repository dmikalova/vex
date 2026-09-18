package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// A. Vinda
//
//	House:  Shadows
//	Type:   Creature
//	Rarity: Common
//	Power:  4
//	Traits: Elf • Thief
//
//	Reap: Deal 1 damage to a creature. If this damage destroys that creature, your opponent discards a random card from their hand.
var AVinda = set.New(
	"A. Vinda",
	card.House.Shadows,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.WC, "235"),
	card.WithPower(4),
	card.WithTraits(card.Traits.Elf, card.Traits.Thief),
	card.WithAbility(
		card.Trigger.Reap, card.DealDamage{
			Amount: 1,
			After:  card.IfDestroyed,
			Target: card.Target.Creature,
			Then: card.DiscardCard{
				Player:    card.Opponent,
				Zones:     []card.Zone{card.Hand},
				Selection: card.Random{Count: 1},
			},
		}),
)
