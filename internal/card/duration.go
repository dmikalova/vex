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
	EndOfTurn           engine.Duration
	NextTurn            engine.Duration
	UntilThisLeavesPlay engine.Duration
	Forever             engine.Duration
}
