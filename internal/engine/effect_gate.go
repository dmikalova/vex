package engine

// A result gate resolves one action and then a follow-up, but only when the first
// action actually happened — written A -> B (destroy a creature -> steal 1 Æmber;
// purge a creature -> give a +1 power counter). The follow-up never runs when the
// gate does nothing: no valid target, an empty zone, or a declined choice. It is
// distinct from a conditional, which turns on a fact about the board rather than an
// action succeeding.
// GatingEffect is an effect that can be the first half of a Then: besides
// resolving, it reports whether it did anything. The report is an unexported
// method, so only engine effects (Destroy, Purge) can be a gate.
type GatingEffect interface {
	Effect
	// resolveGate resolves the effect and reports whether it did anything.
	resolveGate(ctx *EffectContext) bool
}

// resolveGateOf resolves any effect and reports whether it did something. A
// GatingEffect answers for itself; an effect that cannot report progress always
// counts as progress, leaving the caller's own bound (the Rule of Six) to end the
// loop.
func resolveGateOf(ctx *EffectContext, e Effect) bool {
	if g, ok := e.(GatingEffect); ok {
		return g.resolveGate(ctx)
	}
	e.Resolve(ctx)
	return true
}

// Then is the "A -> B" result gate: it resolves First and, only when First did
// something, resolves Result. When First does nothing and an Else is set, the gate
// takes the Else arm instead — the two-verb "you may A. If you do, B. Otherwise, C."
// branch (Novu Dynamo gains Æmber or destroys itself; Auto-Vac 5150 taxes keys or
// archives a card). This is the one case a gate carries an otherwise (rule 5).
type Then struct {
	First  GatingEffect
	Result Effect
	// Else, when set, resolves in place of Result whenever First did nothing —
	// because the controller declined the optional first half or it had no legal
	// target. Leaving it nil keeps the plain "A -> B" gate with no otherwise.
	Else Effect
}

// Text renders the gate, e.g. "destroy a damaged creature -> steal 1 Æmber", or
// with an Else arm "discard a card -> gain 1 Æmber. Otherwise, destroy {self}".
func (e Then) Text() string {
	body := e.First.Text() + " -> " + e.Result.Text()
	if e.Else == nil {
		return body
	}
	return body + ". Otherwise, " + e.Else.Text()
}

// Resolve runs First and then Result if First did something, or the Else arm (when
// set) if it did not.
func (e Then) Resolve(ctx *EffectContext) {
	if e.First.resolveGate(ctx) {
		e.Result.Resolve(ctx)
		return
	}
	if e.Else != nil {
		e.Else.Resolve(ctx)
	}
}

// declinable reports whether the gate's first half is itself one card choice; only
// then can the whole gate be answered by clicking a card.
func (e Then) declinable() bool {
	d, ok := e.First.(declinableEffect)
	return ok && d.declinable()
}

// resolveOptional runs the gate under a May: when First can be declined by card
// click it is asked that way, so "you may destroy another friendly creature ->
// fully heal Chuff Ape" is one click, and Result still hangs off First happening.
func (e Then) resolveOptional(ctx *EffectContext) bool {
	if !e.First.(declinableEffect).resolveOptional(ctx) {
		return false
	}
	e.Result.Resolve(ctx)
	return true
}

// validate surfaces a configuration error in any arm of the gate.
func (e Then) validate() error {
	if err := validateEffect(e.First); err != nil {
		return err
	}
	if err := validateEffect(e.Result); err != nil {
		return err
	}
	if e.Else != nil {
		return validateEffect(e.Else)
	}
	return nil
}
