package engine

// This file holds the tide: a game-wide state, neutral at the start of the game,
// that a card raises to become high for the raiser and low for their opponent. No
// implemented card raises it yet — the raising cards arrive in a later set — so
// this is a stub: the state and the "is it high/low for a player" reads exist so a
// card that only checks the tide (Valoocanth is barred while the tide is low) can
// ship now, but with nothing to raise it the tide stays neutral and every check
// reads false.

// Tide is the game-wide tide state from a single point of view: neutral for
// everyone, or high for one player and low for the other. It is a flat, comparable
// value so it fits the pointerless GameState (ADR 0005).
type Tide uint8

const (
	// TideNeutral is the starting state: the tide is neither high nor low for
	// either player.
	TideNeutral Tide = iota
	// TideHighForP0 is the tide high for player 0 and low for player 1.
	TideHighForP0
	// TideHighForP1 is the tide high for player 1 and low for player 0.
	TideHighForP1
)

// TideIsHigh reports whether the tide is high for the given player.
func (g *Game) TideIsHigh(player int) bool {
	return (g.State.Tide == TideHighForP0 && player == 0) ||
		(g.State.Tide == TideHighForP1 && player == 1)
}

// TideIsLow reports whether the tide is low for the given player. The tide is low
// for a player exactly when it is high for their opponent; it is neither while the
// tide is neutral.
func (g *Game) TideIsLow(player int) bool {
	return (g.State.Tide == TideHighForP0 && player == 1) ||
		(g.State.Tide == TideHighForP1 && player == 0)
}
