package card

import "github.com/dmikalova/vactrol/internal/engine"

// Duration groups the spans a timed effect can last, e.g.
// card.Duration.OpponentNextTurn (see card.Restrict). It mirrors the
// engine's duration.go.
var Duration = durations{
	RemainderOfPlayerTurn: engine.RemainderOfPlayerTurn,
	OpponentNextTurn:      engine.OpponentNextTurn,
	StartOfPlayerNextTurn: engine.StartOfPlayerNextTurn,
	EndOfPlayerNextTurn:   engine.EndOfPlayerNextTurn,
	UntilThisLeavesPlay:   engine.UntilThisLeavesPlay,
	Forever:               engine.Forever,
}

type durations struct {
	// RemainderOfPlayerTurn lasts through the rest of the current turn, then lifts.
	RemainderOfPlayerTurn engine.Duration
	// OpponentNextTurn is dormant this turn and bites only during the affected
	// player's next turn, lifting when that turn ends.
	OpponentNextTurn engine.Duration
	// StartOfPlayerNextTurn is live now, survives the opponent's turn, and lifts at
	// the start of the caster's next turn.
	StartOfPlayerNextTurn engine.Duration
	// EndOfPlayerNextTurn lasts from now through the end of the affected player's next turn.
	EndOfPlayerNextTurn engine.Duration
	// UntilThisLeavesPlay lasts until the card whose effect set it leaves play.
	UntilThisLeavesPlay engine.Duration
	// Forever never lifts; it lasts the rest of the game (latest ability wins).
	Forever engine.Duration
}
