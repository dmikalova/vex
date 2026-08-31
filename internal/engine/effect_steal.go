package engine

import "fmt"

// Stealing Aember moves it from the opponent's pool into your own. You can only
// steal as much Aember as the opponent actually has.
//
//rulebook:effect Steal Aember
type StealAember struct {
	Amount int
}

// Text renders the effect, e.g. "steal 1 Æmber".
func (e StealAember) Text() string { return fmt.Sprintf("steal %d Æmber", e.Amount) }

// Resolve moves the Æmber from the opponent's pool to the controller's.
func (e StealAember) Resolve(ctx *EffectContext) { e.resolveGate(ctx) }

// resolveGate performs the steal and reports whether any Æmber actually moved, so
// it can gate a following effect (Nerve Blast's "steal 1 -> deal 2 damage").
func (e StealAember) resolveGate(ctx *EffectContext) bool {
	opp := ctx.Opponent()
	if ctx.Resolver.AemberProtected(opp) {
		return false
	}
	amt := min(e.Amount, ctx.Resolver.Aember(opp))
	ctx.Resolver.SetAember(opp, ctx.Resolver.Aember(opp)-amt)
	ctx.Resolver.SetAember(ctx.Controller, ctx.Resolver.Aember(ctx.Controller)+amt)
	ctx.Resolver.Logf("%s steals %d Æmber from %s", ctx.Resolver.PlayerName(ctx.Controller), amt, ctx.Resolver.PlayerName(opp))
	return amt > 0
}
