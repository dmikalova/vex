package worldscollide

import "github.com/dmikalova/vex/internal/card"

// Inky Gloom
//
//	House:  Shadows
//	Type:   Tactic
//	Rarity: Common
//
//	Play: Your opponent cannot use creatures to reap during their next turn.
var InkyGloom = set.New(
	"Inky Gloom",
	card.House.Shadows,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.WC, "241"),
	card.WithAbility(
		card.Trigger.Play, card.Restrict{
			Player:   card.Opponent,
			Action:   card.Restricted.Reaping,
			Duration: card.Duration.OpponentNextTurn,
		}),
)
