package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Mogghunter
//
//	House:  Brobnar
//	Type:   Creature
//	Rarity: Common
//	Power:  6
//	Traits: Giant
//
//	Fight: Deal 2 damage to a flank creature.
var Mogghunter = card.New(
	"Mogghunter",
	card.House.Brobnar,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.WC, 11),
	card.WithPower(6),
	card.WithTraits(card.Traits.Giant),
	card.WithAbility(
		card.Trigger.Fight, card.DealDamage{
			Target: card.Target.Creature.OnFlank(),
			Amount: 2,
		}),
)
