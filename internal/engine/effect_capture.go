package engine

import (
	"fmt"
	"slices"
)

// Capturing Æmber moves it from a player's pool onto a capturing creature, where
// it counts for no player until that creature leaves play, at which point it goes
// to the pool of the capturing creature's controller's opponent. A creature can
// only capture what the Source pool holds. Target is the creature that captures
// (this creature by default); Source is the pool the Æmber comes from; Per repeats
// the capture, choosing a fresh Target each time (Hypnotic Command captures once
// for each friendly Mars creature).
type CaptureAember struct {
	// Amount is the fixed Æmber to capture.
	Amount int
	// All captures the whole Source pool instead of a fixed Amount.
	All bool
	// By captures a share of the Source pool instead of a fixed Amount
	// (By: AllBut(5) leaves the pool at exactly five).
	By Loss
	// Target is the creature that captures; the zero value is this creature.
	Target Target
	// Source is the pool the Æmber is taken from.
	Source Player
	// Per repeats the capture, choosing a fresh Target each time.
	Per Count
	// Distinct bars a creature an earlier repetition already picked from being
	// picked again, so a Per that repeats N times spreads the captures across N
	// different creatures (Unguarded Camp).
	Distinct bool
}

// validate requires an explicit Target and Source.
func (e CaptureAember) validate() error {
	if !e.Target.valid() {
		return errUnsetTarget("CaptureAember")
	}
	if e.Source == playerUnset {
		return errUnsetPlayer("CaptureAember")
	}
	if e.Amount != 0 && e.By != nil {
		return fmt.Errorf("CaptureAember: set Amount or By, not both (got Amount=%d)", e.Amount)
	}
	if e.Distinct && e.Per == nil {
		return fmt.Errorf("CaptureAember: Distinct is meaningless without a Per to repeat")
	}
	if e.Source == ItsOpponent && (e.All || e.By != nil) {
		return fmt.Errorf(
			"CaptureAember: ItsOpponent captures a fixed Amount, since a share of " +
				"\"its opponent's\" pool names a different pool for each capturer")
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
	switch {
	case e.All:
		body = fmt.Sprintf("%s captures all %s Æmber", capturer, e.poolPossessive())
	case e.By != nil:
		body = fmt.Sprintf(
			"%s captures %s from %s",
			capturer,
			e.By.object(e.poolPossessive()),
			e.fromText(),
		)
	default:
		body = fmt.Sprintf("%s captures %d Æmber from %s", capturer, e.Amount, e.fromText())
	}
	body = forEach(e.Per, body)
	if e.Distinct {
		body += fmt.Sprintf(
			". Each creature cannot capture more than %d Æmber this way",
			e.Amount,
		)
	}
	return body
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
	if e.Source == ItsOpponent {
		return "its opponent"
	}
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
	reps := 1
	if e.Per != nil {
		reps = e.Per.Value(ctx)
	}
	var captured []LocalID
	for i := 0; i < reps; i++ {
		ids := e.Target.selectWith(ctx, false, e.eligible(captured))
		if len(ids) == 0 {
			return
		}
		for _, id := range ids {
			// A creature the fight (or an earlier step) destroyed cannot capture: the
			// Æmber would sit on a card in a discard pile.
			if !resolverInPlay(ctx, id) {
				continue
			}
			captured = append(captured, id)
			pool := e.sourcePool(ctx, id)
			amt := e.Amount
			switch {
			case e.All:
				amt = ctx.Resolver.Aember(pool)
			case e.By != nil:
				amt = e.By.lose(ctx.Resolver.Aember(pool))
			}
			amt = min(amt, ctx.Resolver.Aember(pool))
			ctx.Resolver.SetAember(pool, ctx.Resolver.Aember(pool)-amt)
			ctx.Resolver.AddAmberOn(id, amt)
			ctx.Resolver.Record(AemberCaptured{Creature: id, Amount: amt})
		}
	}
}

// eligible narrows the candidates a repetition may choose from to the creatures
// no earlier repetition already had capture, but only when Distinct asks for it.
// A nil filter leaves the target's own candidate set untouched.
func (e CaptureAember) eligible(captured []LocalID) func(LocalID) bool {
	if !e.Distinct {
		return nil
	}
	return func(id LocalID) bool { return !slices.Contains(captured, id) }
}

// sourcePool names the pool one capture draws from. It is per capturer rather
// than per effect because ItsOpponent is relative to the capturing creature, so
// a single effect reaching every creature on the board takes from both pools.
func (e CaptureAember) sourcePool(ctx *EffectContext, capturer LocalID) int {
	if e.Source == ItsOpponent {
		return 1 - ctx.Resolver.Controller(capturer)
	}
	return ctx.PlayerFor(e.Source)
}
