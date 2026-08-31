package engine

import "strings"

// Some abilities have you choose a single creature and then do one or more things
// to it — "Ready and fight with a friendly creature." OnChooseCreature models
// that: it picks one creature (via its Target) and applies an ordered list of
// CreatureVerbs to it, so the verbs read as one sentence sharing a target.

// CreatureVerb is a single action applied to a chosen creature. Verbs are the
// building blocks of OnChooseCreature and compose into natural card text.
type CreatureVerb interface {
	VerbText() string
	Apply(ctx *EffectContext, target LocalID)
}

// ReadyVerb readies the chosen creature.
//
// Readying an exhausted creature stands it back up so it can be used again this
// turn (to reap, fight, or use an action).
type ReadyVerb struct{}

// VerbText returns the verb phrase.
func (ReadyVerb) VerbText() string { return "ready" }

// Apply readies the creature.
func (ReadyVerb) Apply(ctx *EffectContext, target LocalID) {
	ctx.Resolver.SetExhausted(target, false)
	ctx.Resolver.Logf("%s is readied", ctx.Resolver.Name(target))
}

// FightVerb makes the chosen creature fight an enemy creature.
//
// Using a creature to fight has it attack an enemy creature: both deal their
// power as damage to each other simultaneously (see Game.fight). The creature's
// controller chooses which enemy it fights.
type FightVerb struct{}

// VerbText returns the verb phrase.
func (FightVerb) VerbText() string { return "fight with" }

// Apply has the creature fight a chosen enemy creature.
func (FightVerb) Apply(ctx *EffectContext, target LocalID) {
	owner := ctx.Resolver.Owner(target)
	enemies := ctx.Resolver.Battleline(1 - owner)
	victim, ok := ctx.Resolver.ChooseCreature(
		owner,
		ctx.Source,
		"Choose a creature to fight",
		enemies,
	)
	if !ok {
		ctx.Resolver.Logf("%s has no creature to fight", ctx.Resolver.Name(target))
		return
	}
	ctx.Resolver.FightWith(target, victim)
}

// OnChooseCreature picks a single creature named by its Target and applies one or
// more verbs to it. The verbs share the single chosen target, which lets card
// text read naturally, e.g. "Ready and fight with a friendly creature."
type OnChooseCreature struct {
	Target Target
	Verbs  []CreatureVerb
}

// validate requires an explicit target.
func (e OnChooseCreature) validate() error {
	if !e.Target.valid() {
		return errUnsetTarget("OnChooseCreature")
	}
	return nil
}

// Text joins the verbs and the shared target, e.g.
// "ready and fight with a friendly creature".
func (e OnChooseCreature) Text() string {
	parts := make([]string, 0, len(e.Verbs))
	for _, v := range e.Verbs {
		parts = append(parts, v.VerbText())
	}
	return strings.Join(parts, " and ") + " " + e.Target.Text()
}

// Resolve chooses the creature (through its Target) once, then applies each verb
// to it. A Target that selects nothing (no candidate or a declined choice) simply
// applies no verbs.
func (e OnChooseCreature) Resolve(ctx *EffectContext) {
	for _, id := range e.Target.Select(ctx) {
		for _, v := range e.Verbs {
			v.Apply(ctx, id)
		}
	}
}

// UseVerb uses the chosen creature. The controller picks how to use it — reap,
// fight, or its "Action:" ability — and that use resolves completely, nesting any
// further uses it triggers, before control returns. This is "Use a friendly
// creature": a creature can only be used while ready, so an already-exhausted
// creature may be chosen but nothing happens when it is used.
type UseVerb struct{}

// VerbText returns the verb phrase.
func (UseVerb) VerbText() string { return "use" }

// Apply offers the controller the target's available uses and resolves the chosen
// one.
func (UseVerb) Apply(ctx *EffectContext, target LocalID) {
	owner := ctx.Resolver.Owner(target)
	labels := []string{"reap"}
	uses := []func(){func() { ctx.Resolver.ReapWith(target) }}
	if len(ctx.Resolver.Battleline(1-owner)) > 0 {
		labels = append(labels, "fight")
		uses = append(uses, func() { FightVerb{}.Apply(ctx, target) })
	}
	if ctx.Resolver.HasAction(target) {
		labels = append(labels, "use its action")
		uses = append(uses, func() { ctx.Resolver.UseActionOf(target) })
	}
	idx := ctx.Resolver.ChooseOption(
		owner,
		ctx.Source,
		"Choose how to use "+ctx.Resolver.Name(target),
		labels,
	)
	if idx < 0 || idx >= len(uses) {
		return
	}
	uses[idx]()
}

// StunVerb stuns the chosen creature. As a verb it shares a single chosen target
// with the other verbs of an OnChooseCreature, so "stun and exhaust a creature"
// reads and resolves as one target rather than two separate choices.
type StunVerb struct{}

// VerbText returns the verb phrase.
func (StunVerb) VerbText() string { return "stun" }

// Apply stuns the creature.
func (StunVerb) Apply(ctx *EffectContext, target LocalID) {
	ctx.Resolver.SetStunned(target, true)
}

// ExhaustVerb exhausts the chosen creature. Like StunVerb it is a verb so it can
// share a single chosen target with other verbs of an OnChooseCreature.
type ExhaustVerb struct{}

// VerbText returns the verb phrase.
func (ExhaustVerb) VerbText() string { return "exhaust" }

// Apply exhausts the creature.
func (ExhaustVerb) Apply(ctx *EffectContext, target LocalID) {
	ctx.Resolver.SetExhausted(target, true)
}
