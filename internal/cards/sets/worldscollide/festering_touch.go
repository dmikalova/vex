package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Festering Touch
//
//	House:  Dis
//	Type:   Tactic
//	Rarity: Common
//	Æmber:  1
//
//	Play: Choose up to 2 Creatures. Deal 1 damage to each chosen Creature. If that Creature was already damaged, deal 3 damage instead.
var FesteringTouch = card.New(
	"Festering Touch",
	card.House.Dis,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.WC, "75"),
	card.WithAemberBonus(1),
	card.WithAbility(
		card.Trigger.Play,
		card.DealDamage{
			Spread: card.UpToCreatures{Count: 2, Amount: 1, WhenDamaged: 3},
		},
	),
)
