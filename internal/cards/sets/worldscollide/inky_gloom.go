package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Inky Gloom
//
//	House:  Shadows
//	Type:   Tactic
//	Rarity: Common
//
//	Play: Your opponent cannot use Creatures to reap during their next turn.
var InkyGloom = set.New(
	"Inky Gloom",
	card.House.Shadows,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.WC, "241"),
	card.WithAbility(
		card.Trigger.Play, card.CannotReap{
			Player:   card.Opponent,
			Duration: card.Duration.OpponentNextTurn,
		}),
)
