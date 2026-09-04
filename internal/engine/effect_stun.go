package engine

// Stun and Unstun apply and remove the stun status. Each is a simple "verb the
// target" effect, so a stun that runs beside another status change on the same
// target (such as an Exhaust) folds into one phrase in a Sequence (see
// combinable).

// A stun is a status placed on a creature. A stunned creature must shake off the
// stun before it can do anything else: the next time it is used to reap, fight,
// or use an "Action:" ability, it is exhausted and the stun is removed instead of
// that action happening. Its constant abilities and any effect that does not
// require using it keep working while it is stunned. Stunning applies this status
// to each creature the effect targets.
//
//rulebook:effect Stun
type Stun struct {
	Target Target
}

// validate requires an explicit target.
func (e Stun) validate() error {
	if !e.Target.valid() {
		return errUnsetTarget("Stun")
	}
	return nil
}

func (e Stun) verb() string       { return "stun" }
func (e Stun) targetText() string { return e.Target.Text() }

// Text renders the effect, e.g. "stun each friendly creature".
func (e Stun) Text() string { return e.verb() + " " + e.targetText() }

// Resolve stuns each selected creature. A creature already stunned still gets a
// log line — the source still had to choose it — just without a state change.
func (e Stun) Resolve(ctx *EffectContext) {
	for _, id := range e.Target.Select(ctx) {
		if ctx.Resolver.Stunned(id) {
			ctx.Resolver.Record(CreatureStunned{Creature: id, By: ctx.Source, AlreadyStunned: true})
			continue
		}
		ctx.Resolver.SetStunned(id, true)
		ctx.Resolver.Record(CreatureStunned{Creature: id, By: ctx.Source})
	}
}

// Unstunning a creature removes the stun status from each creature the effect
// targets, freeing it to act normally instead of having to shake the stun off.
//
//rulebook:effect Unstun
type Unstun struct {
	Target Target
}

// validate requires an explicit target.
func (e Unstun) validate() error {
	if !e.Target.valid() {
		return errUnsetTarget("Unstun")
	}
	return nil
}

func (e Unstun) verb() string       { return "unstun" }
func (e Unstun) targetText() string { return e.Target.Text() }

// Text renders the effect, e.g. "unstun each friendly creature".
func (e Unstun) Text() string { return e.verb() + " " + e.targetText() }

// Resolve clears the stun on each selected creature.
func (e Unstun) Resolve(ctx *EffectContext) {
	for _, id := range e.Target.Select(ctx) {
		ctx.Resolver.SetStunned(id, false)
	}
}
