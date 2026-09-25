package callofthearchons

import "github.com/dmikalova/vex/internal/card"

// Scrambler Storm
//
//	House:  Logos
//	Type:   Tactic
//	Rarity: Uncommon
//	Bonus:  Æmber
//
//	Play: Your opponent cannot play tactics during their next turn.
var ScramblerStorm = set.New(
	"Scrambler Storm",
	card.House.Logos,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, "122"),
	card.WithBonus(card.Bonus.Aember),
	card.WithAbility(card.Trigger.Play, card.CannotPlay{
		Player:   card.Opponent,
		Type:     card.Type.Tactic,
		Duration: card.Duration.OpponentNextTurn,
	}),
)
