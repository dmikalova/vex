package engine

import "fmt"

// A chain is a penalty a card can inflict on its controller: while a player holds
// chains they draw fewer cards each turn — one fewer for every 6 chains — until the
// chains are shed, one on each turn the reduction blocks a draw. Gaining a chain is
// the cost some strong effects charge, so a card's power is paid for by a slower
// hand refill (see Game.drawStep).
type GainChains struct {
	Amount int
}

// Text renders the effect, e.g. "gain 1 chain" or "gain 2 chains".
func (e GainChains) Text() string {
	return fmt.Sprintf("gain %d %s", e.Amount, chainNoun(e.Amount))
}

// Resolve adds the chains to the controller.
func (e GainChains) Resolve(ctx *EffectContext) {
	ctx.Resolver.GainChains(ctx.Controller, e.Amount)
}

// chainNoun renders "chain" or "chains" for a count.
func chainNoun(n int) string {
	if n == 1 {
		return "chain"
	}
	return "chains"
}
