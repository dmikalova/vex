package engine

import "fmt"

// A chain is a penalty a card can inflict on its controller: while a player holds
// chains they draw fewer cards each turn — one fewer for every 6 chains — until the
// chains are shed, one on each turn the reduction blocks a draw. Gaining a chain is
// the cost some strong effects charge, so a card's power is paid for by a slower
// hand refill (see Game.drawStep).
type GainChains struct {
	// Player takes the chains; an unset zero value means the controller, so the
	// common "gain N chains" cards need not name it.
	Player Player
	// Amount is how many chains to gain.
	Amount int
}

// Text renders the effect, e.g. "gain 1 chain" or "your opponent gains 2 chains".
func (e GainChains) Text() string {
	if e.Player == Opponent {
		return fmt.Sprintf("your opponent gains %d %s", e.Amount, chainNoun(e.Amount))
	}
	return fmt.Sprintf("gain %d %s", e.Amount, chainNoun(e.Amount))
}

// Resolve adds the chains to the selected player, defaulting to the controller.
func (e GainChains) Resolve(ctx *EffectContext) {
	ctx.Resolver.GainChains(ctx.PlayerFor(e.chained()), e.Amount)
}

// chained resolves the target player, treating the unset zero value as the
// controller.
func (e GainChains) chained() Player {
	if e.Player == playerUnset {
		return Controller
	}
	return e.Player
}

// chainNoun renders "chain" or "chains" for a count.
func chainNoun(n int) string {
	if n == 1 {
		return "chain"
	}
	return "chains"
}
