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
func (e Stun) Resolve(ctx *EffectContext) { e.stun(ctx, e.Target.Select(ctx)) }

// declinable reports that the stun is a single clickable creature.
func (e Stun) declinable() bool { return e.Target.isChosen() }

// vacuous reports that there is no creature to stun, so a "you may" wrapping it
// need not ask.
func (e Stun) vacuous(ctx *EffectContext) bool { return e.Target.empty(ctx) }

// resolveOptional asks for the creature declinably, so "you may stun a creature"
// is answered by clicking that creature rather than by a separate Yes/No.
func (e Stun) resolveOptional(ctx *EffectContext) bool {
	return e.stun(ctx, e.Target.SelectOptional(ctx))
}

// stun applies the status to an already-selected set and reports whether any
// creature was chosen.
func (e Stun) stun(ctx *EffectContext, ids []LocalID) bool {
	for _, id := range ids {
		if ctx.Resolver.Stunned(id) {
			ctx.Resolver.Record(CreatureStunned{Creature: id, By: ctx.Source, AlreadyStunned: true})
			continue
		}
		ctx.Resolver.SetStunned(id, true)
		ctx.Resolver.Record(CreatureStunned{Creature: id, By: ctx.Source})
	}
	return len(ids) > 0
}

// Unstunning a creature removes the stun status from each creature the effect
// targets, freeing it to act normally instead of having to shake the stun off.
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

// Resolve clears the stun on each selected creature and logs the ones it freed as
// one act, so a card that unstuns several creatures reads as one line.
func (e Unstun) Resolve(ctx *EffectContext) {
	var freed []LocalID
	for _, id := range e.Target.Select(ctx) {
		if !ctx.Resolver.Stunned(id) {
			continue
		}
		ctx.Resolver.SetStunned(id, false)
		freed = append(freed, id)
	}
	if len(freed) > 0 {
		ctx.Resolver.Record(CreaturesUnstunned{Player: ctx.Controller, Creatures: freed})
	}
}
