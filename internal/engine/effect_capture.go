package engine

import (
	"fmt"
	"slices"
)

// Capturing Æmber moves it from a player's pool onto a capturing creature, where
// it counts for no player until that creature leaves play, at which point it goes
// to the pool of the capturing creature's controller's opponent. A creature can
// only capture what the Source pool holds. Target is the creature that captures
// (this creature by default); Source is the pool the Æmber comes from. The two
// count axes are distinct: Per scales how much one capturer takes (Yxili Marauder
// captures 1 per friendly ready Mars creature onto itself), while Times repeats
// the capture, choosing a fresh Target each time (Hypnotic Command captures once
// for each friendly Mars creature).
type CaptureAember struct {
	// Amount is the fixed Æmber to capture; Per scales it "for each ...".
	Amount int
	// Per multiplies the Amount one capturer takes by a running count.
	Per Count
	// All captures the whole Source pool instead of a fixed Amount.
	All bool
	// By captures a share of the Source pool instead of a fixed Amount
	// (By: AllBut(5) leaves the pool at exactly five).
	By Loss
	// Target is the creature that captures; the zero value is this creature.
	Target Target
	// Source is the pool the Æmber is taken from.
	Source Player
	// Times repeats the capture, choosing a fresh Target each time.
	Times Count
	// Distinct bars a creature an earlier repetition already picked from being
	// picked again, so a Times that repeats N times spreads the captures across N
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
	if err := errAmountOr("CaptureAember", "By", e.Amount, e.By != nil); err != nil {
		return err
	}
	if e.Distinct && e.Times == nil {
		return fmt.Errorf("CaptureAember: Distinct is meaningless without a Times to repeat")
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
	default:
		body = fmt.Sprintf(
			"%s captures %s from %s",
			capturer,
			aemberObject(e.Amount, e.By, e.poolPossessive()),
			e.fromText(),
		)
	}
	// A card sets at most one count axis; whichever is present leads the "for
	// each" clause. Per and Times render the same phrasing but differ in resolution.
	count := e.Times
	if count == nil {
		count = e.Per
	}
	body = forEach(count, body)
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
// Times times and choosing a fresh Target each time, and scaling one capturer's
// take by Per. Capturing stops early if the Target selects nothing (no eligible
// creature, or the choice is declined).
func (e CaptureAember) Resolve(ctx *EffectContext) {
	reps := 1
	if e.Times != nil {
		reps = e.Times.Value(ctx)
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
			// All captures the whole pool, which the AllAember share already expresses.
			by := e.By
			if e.All {
				by = AllAember
			}
			held := ctx.Resolver.Aember(pool)
			amt := min(poolAmount(scaled(e.Amount, e.Per, ctx), by, nil, ctx, held), held)
			ctx.Resolver.SetAember(pool, held-amt)
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

// MoveAemberToCommonSupply removes Æmber sitting on a creature and returns it to the
// common supply, the reverse of a capture — Aubade the Grim discards one of its
// own captured Æmber each time it reaps. A creature holding fewer than Amount is
// simply emptied rather than driven negative.
type MoveAemberToCommonSupply struct {
	// Amount is the Æmber to remove from each target.
	Amount int
	// Target is the creature the Æmber is removed from; the zero value is this
	// creature.
	Target Target
}

// validate requires an explicit Target and a positive Amount.
func (e MoveAemberToCommonSupply) validate() error {
	if !e.Target.valid() {
		return errUnsetTarget("MoveAemberToCommonSupply")
	}
	if e.Amount <= 0 {
		return fmt.Errorf("MoveAemberToCommonSupply: Amount must be positive")
	}
	return nil
}

// Text renders the effect, e.g. "move 1 Æmber from {self} to the common supply".
func (e MoveAemberToCommonSupply) Text() string {
	target := SelfName
	if e.Target.Kind != TargetThisCreature {
		target = e.Target.Text()
	}
	return fmt.Sprintf("move %d Æmber from %s to the common supply", e.Amount, target)
}

// Resolve removes up to Amount Æmber from each target, returning it to the common
// supply. A target holding none is skipped.
func (e MoveAemberToCommonSupply) Resolve(ctx *EffectContext) {
	for _, id := range e.Target.Select(ctx) {
		remove := min(e.Amount, ctx.Resolver.AmberOn(id))
		if remove <= 0 {
			continue
		}
		ctx.Resolver.AddAmberOn(id, -remove)
		ctx.Resolver.Record(AemberMovedToCommonSupply{Creature: id, Amount: remove})
	}
}
