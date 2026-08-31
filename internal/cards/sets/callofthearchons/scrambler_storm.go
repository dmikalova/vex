package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Scrambler Storm
//
//	House:  Logos
//	Type:   Tactic
//	Rarity: Uncommon
//	Æmber:  1
//
//	Play: Your opponent cannot play action cards during their next turn.
var ScramblerStorm = card.New(
	"Scrambler Storm",
	card.House.Logos,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, 122),
	card.WithAemberBonus(1),
	card.WithAbility(card.Trigger.Play, card.CannotPlayNextTurn{
		Player: card.Opponent,
		Type:   card.Type.Tactic,
	}),
)
