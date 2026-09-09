package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Festering Touch
//
//	House:  Dis
//	Type:   Tactic
//	Rarity: Common
//	Æmber:  1
//
//	Play: Deal 1 damage to up to 2 creatures, dealing 3 damage instead to a chosen creature that was already damaged.
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
