package card

import "github.com/dmikalova/vactrol/internal/engine"

// Duration groups the spans a timed effect can last, e.g. card.Duration.NextTurn
// (see card.CannotFight). It mirrors the engine's duration.go.
var Duration = durations{
	EndOfTurn:           engine.EndOfTurn,
	NextTurn:            engine.NextTurn,
	UntilThisLeavesPlay: engine.UntilThisLeavesPlay,
	Forever:             engine.Forever,
}

type durations struct {
	// EndOfTurn lasts through the rest of the current turn, then lifts.
	EndOfTurn engine.Duration
	// NextTurn lasts through the affected player's next turn, then lifts.
	NextTurn engine.Duration
	// UntilThisLeavesPlay lasts until the card whose effect set it leaves play.
	UntilThisLeavesPlay engine.Duration
	// Forever never lifts; it lasts the rest of the game (latest ability wins).
	Forever engine.Duration
}
