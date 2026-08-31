package engine

import "fmt"

// Capturing Aember moves it from a player's pool onto a capturing creature, where
// it counts for no player until that creature leaves play, at which point it goes
// to the pool of the capturing creature's controller's opponent. A creature can
// only capture what the Source pool holds. Target is the creature that captures
// (this creature by default); Source is the pool the Aember comes from; Per repeats
// the capture, choosing a fresh Target each time (Hypnotic Command captures once
// for each friendly Mars creature).
//
//rulebook:effect Capture Aember
type CaptureAember struct {
	Amount int
	// All captures the whole Source pool instead of a fixed Amount.
	All bool
	// Target is the creature that captures; the zero value is this creature.
	Target Target
	// Source is the pool the Æmber is taken from.
	Source Player
	// Per repeats the capture, choosing a fresh Target each time.
	Per Count
}

// validate requires an explicit Target and Source.
func (e CaptureAember) validate() error {
	if !e.Target.valid() {
		return errUnsetTarget("CaptureAember")
	}
	if e.Source == playerUnset {
		return errUnsetPlayer("CaptureAember")
	}
	return nil
}

// Text renders the effect, e.g. "{self} captures 3 Æmber from your opponent",
// "{self} captures all your opponent's Æmber", or (for a chosen enemy capturer)
// "an enemy creature captures 1 Æmber from their own side".
func (e CaptureAember) Text() string {
	capturer := SelfName
	if e.Target.Kind != TargetThisCreature {
		capturer = e.Target.Text()
	}
	var body string
	if e.All {
		body = fmt.Sprintf("%s captures all %s Æmber", capturer, e.poolPossessive())
	} else {
		body = fmt.Sprintf("%s captures %d Æmber from %s", capturer, e.Amount, e.fromText())
	}
	return forEach(e.Per, body)
}

// poolPossessive names the Source pool as a possessive for the "all" wording,
// e.g. "your opponent's".
func (e CaptureAember) poolPossessive() string {
	if e.Source == Controller {
		return "your"
	}
	return "your opponent's"
}

// fromText names the Source pool relative to the capturer: "your opponent" when a
// friendly creature captures the opponent's pool, or "their own side" when an
// enemy creature captures the opponent's (its own) pool.
func (e CaptureAember) fromText() string {
	enemyCapturer := e.Target.Kind == TargetChosenEnemyCreature
	poolIsControllers := e.Source == Controller
	if enemyCapturer != poolIsControllers { // the capturer captures its own pool
		if enemyCapturer {
			return "their own side"
		}
		return "your own side"
	}
	if poolIsControllers {
		return "you"
	}
	return "your opponent"
}

// Resolve moves Æmber from the Source pool onto each capturing creature, repeating
// Per times and choosing a fresh Target each time. Capturing stops early if the
// Target selects nothing (no eligible creature, or the choice is declined).
func (e CaptureAember) Resolve(ctx *EffectContext) {
	pool := ctx.PlayerFor(e.Source)
	reps := 1
	if e.Per != nil {
		reps = e.Per.Value(ctx)
	}
	for i := 0; i < reps; i++ {
		ids := e.Target.Select(ctx)
		if len(ids) == 0 {
			return
		}
		for _, id := range ids {
			amt := e.Amount
			if e.All {
				amt = ctx.Resolver.Aember(pool)
			}
			amt = min(amt, ctx.Resolver.Aember(pool))
			ctx.Resolver.SetAember(pool, ctx.Resolver.Aember(pool)-amt)
			ctx.Resolver.AddAmberOn(id, amt)
			ctx.Resolver.Logf("%s captures %d Æmber", ctx.Resolver.Name(id), amt)
		}
	}
}
