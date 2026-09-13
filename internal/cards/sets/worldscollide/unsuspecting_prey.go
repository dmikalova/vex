package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Unsuspecting Prey
//
//	House:  Untamed
//	Type:   Tactic
//	Rarity: Common
//
//	Play: Deal 2 damage to up to 3 undamaged Creatures.
var UnsuspectingPrey = card.New(
	"Unsuspecting Prey",
	card.House.Untamed,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.WC, "368"),
	card.WithAbility(
		card.Trigger.Play, card.DealDamage{
			Spread: card.UpToCreatures{Count: 3, Amount: 2, Undamaged: true},
		}),
)
