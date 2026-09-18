package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Festering Touch
//
//	House:  Dis
//	Type:   Tactic
//	Rarity: Common
//	Bonus:  Æmber
//
//	Play: Choose up to 2 creatures. Deal 1 damage to each chosen creature. Deal 3 damage instead to each chosen creature that was already damaged.
var FesteringTouch = set.New(
	"Festering Touch",
	card.House.Dis,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.WC, "75"),
	card.WithBonus(card.Bonus.Aember),
	card.WithAbility(
		card.Trigger.Play,
		card.DealDamage{
			Spread: card.UpToCreatures{Count: 2, Amount: 1, WhenDamaged: 3},
		},
	),
)
