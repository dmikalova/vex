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
			// Po's Pixies: the source pool keeps its Æmber and the capture is drawn
			// from the common supply instead, so only the creature's Æmber grows.
			fromSupply := ctx.Resolver.AemberTakenFromSupply(pool)
			if !fromSupply {
				ctx.Resolver.SetAember(pool, held-amt)
			}
			ctx.Resolver.AddAmberOn(id, amt)
			ctx.Resolver.Record(AemberCaptured{
				Creature:   id,
				Amount:     amt,
				Source:     ctx.Source,
				FromSupply: fromSupply,
			})
			ctx.It, ctx.HasIt = id, true
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

// CaptureFromAnyPlayer captures up to Amount Æmber onto this creature, drawn from
// both players' pools in whatever split the controller chooses ("capture N Æmber
// from any combination of players", Crassosaurus). Each unit is taken from a pool
// the controller names when both still hold Æmber, so a rational controller strips
// the opponent first but may spend their own pool too. Capturing stops once both
// pools are empty, so a capture is capped by what the pools actually hold.
//
// The split is decided here; the moves themselves delegate to CaptureAember, one
// capture per contributing pool, so the capture, supply-redirect, and log behavior
// stay identical to every other capture.
type CaptureFromAnyPlayer struct {
	// Amount is the total Æmber to capture across both pools.
	Amount int
}

// validate requires a positive Amount.
func (e CaptureFromAnyPlayer) validate() error {
	if e.Amount <= 0 {
		return fmt.Errorf("CaptureFromAnyPlayer needs a positive Amount")
	}
	return nil
}

// Text renders the effect, e.g. "{self} captures 10 Æmber from any combination of
// players".
func (e CaptureFromAnyPlayer) Text() string {
	return fmt.Sprintf(
		"%s captures %d Æmber from any combination of players", SelfName, e.Amount,
	)
}

// Resolve asks the controller, one Æmber at a time, which pool to draw from while
// both hold Æmber, then captures each pool's share onto this creature. Reading the
// pool totals up front and decrementing local budgets keeps a pool from being
// picked past what it holds.
func (e CaptureFromAnyPlayer) Resolve(ctx *EffectContext) {
	ownLeft := ctx.Resolver.Aember(ctx.Controller)
	oppLeft := ctx.Resolver.Aember(ctx.Opponent())
	fromOwn, fromOpp := 0, 0
	for n := 0; n < e.Amount && (ownLeft > 0 || oppLeft > 0); n++ {
		takeOwn := ownLeft > 0
		if ownLeft > 0 && oppLeft > 0 {
			takeOwn = ctx.ChooseOption(
				"Capture 1 Æmber from which pool?",
				[]string{"your pool", "your opponent's pool"},
			) == 0
		}
		if takeOwn {
			fromOwn, ownLeft = fromOwn+1, ownLeft-1
		} else {
			fromOpp, oppLeft = fromOpp+1, oppLeft-1
		}
	}
	e.captureFrom(ctx, Controller, fromOwn)
	e.captureFrom(ctx, Opponent, fromOpp)
}

// captureFrom captures amt Æmber onto this creature from one pool, delegating to
// CaptureAember so the move behaves exactly like any other capture. A zero share
// is skipped so an untouched pool never logs a capture.
func (e CaptureFromAnyPlayer) captureFrom(ctx *EffectContext, source Player, amt int) {
	if amt <= 0 {
		return
	}
	CaptureAember{
		Amount: amt,
		Target: Target{Kind: TargetThisCreature},
		Source: source,
	}.Resolve(ctx)
}

// MoveAemberToSupply removes Æmber sitting on a creature and returns it to the
// common supply, the reverse of a capture — Aubade the Grim discards one of its
// own captured Æmber each time it reaps. A creature holding fewer than Amount is
// simply emptied rather than driven negative.
type MoveAemberToSupply struct {
	// Amount is the Æmber to remove from each target. Ignored when All is set.
	Amount int
	// All removes every Æmber on each target rather than a fixed Amount, rendering
	// "move each Æmber on <target> to the common supply" (Imperial Scutum,
	// Praefectus Ludo).
	All bool
	// Target is the creature the Æmber is removed from; the zero value is this
	// creature.
	Target Target
}

// validate requires an explicit Target and either a positive Amount or the All
// flag, but not both.
func (e MoveAemberToSupply) validate() error {
	if !e.Target.valid() {
		return errUnsetTarget("MoveAemberToSupply")
	}
	if err := errAmountOr("MoveAemberToSupply", "All", e.Amount, e.All); err != nil {
		return err
	}
	if !e.All && e.Amount <= 0 {
		return fmt.Errorf("MoveAemberToSupply: Amount must be positive")
	}
	return nil
}

// Text renders the effect, e.g. "move 1 Æmber from {self} to the common supply",
// or "move each Æmber on {self} to the common supply" in All mode.
func (e MoveAemberToSupply) Text() string {
	target := SelfName
	if e.Target.Kind != TargetThisCreature {
		target = e.Target.Text()
	}
	if e.All {
		return fmt.Sprintf("move each \u00c6mber on %s to the common supply", target)
	}
	return fmt.Sprintf("move %d Æmber from %s to the common supply", e.Amount, target)
}

// Resolve removes the Æmber from each target — Amount, or all of it in All mode —
// returning it to the common supply. A target holding none is skipped.
func (e MoveAemberToSupply) Resolve(ctx *EffectContext) {
	for _, id := range e.Target.Select(ctx) {
		have := ctx.Resolver.AmberOn(id)
		remove := have
		if !e.All {
			remove = min(e.Amount, have)
		}
		if remove <= 0 {
			continue
		}
		ctx.Resolver.AddAmberOn(id, -remove)
		ctx.Resolver.Record(AemberMovedToCommonSupply{Creature: id, Amount: remove})
	}
}
