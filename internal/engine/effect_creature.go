package engine

import (
	"fmt"
	"slices"
	"strings"
)

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

// A narrowingVerb can only act on some of the creatures its target names, and
// so narrows the choice offered without changing the printed text: "use a
// friendly creature" never meant an exhausted one, since using it would do
// nothing at all.
type narrowingVerb interface {
	CreatureVerb
	canApplyTo(ctx *EffectContext, target LocalID) bool
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
	ctx.Resolver.Record(CreatureReadied{Creature: target})
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
	fightAmong(ctx, target, ctx.Resolver.Battleline(1-owner))
}

// fightAmong has attacker fight one of enemies, chosen by attacker's controller,
// and reports which enemy it fought. Taking the candidates as an argument is what
// lets RepeatedFight narrow them to the enemies no earlier fight has used.
func fightAmong(ctx *EffectContext, attacker LocalID, enemies []LocalID) (LocalID, bool) {
	owner := ctx.Resolver.Owner(attacker)
	enemy, ok := ctx.Resolver.ChooseCreature(
		owner,
		ctx.Source,
		"Choose a creature to fight",
		enemies,
	)
	if !ok {
		ctx.Resolver.Record(NoCreatureToFight{Creature: attacker})
		return 0, false
	}
	ctx.Resolver.FightWith(attacker, enemy)
	return enemy, true
}

// ChooseCreatureThen asks the controller to choose a creature from Target,
// records it on the effect context as "it" (read by a Target.Triggering inside
// Then), and resolves Then unconditionally — unlike a Then result gate, Then
// always runs once a creature is chosen, whether or not it does anything
// (Protectrix protects a creature it heals even when there was no damage to
// heal). It models "Choose a creature. <do something to that creature>."
type ChooseCreatureThen struct {
	Target Target
	Then   Effect
}

// validate requires an explicit target and a valid Then.
func (e ChooseCreatureThen) validate() error {
	if !e.Target.valid() {
		return errUnsetTarget("ChooseCreatureThen")
	}
	return validateEffect(e.Then)
}

// Text renders the effect, e.g. "choose a creature - fully heal it".
func (e ChooseCreatureThen) Text() string {
	return "choose " + e.Target.Text() + " - " + e.Then.Text()
}

// Resolve asks for a creature, records it as "it", then resolves Then. A Target
// that chooses nothing (no candidate) leaves Then unresolved.
func (e ChooseCreatureThen) Resolve(ctx *EffectContext) {
	ids := e.Target.Select(ctx)
	if len(ids) == 0 {
		return
	}
	for _, id := range ids {
		ctx.It, ctx.HasIt = id, true
	}
	e.Then.Resolve(ctx)
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
	e.applyTo(ctx, e.Target.selectWith(ctx, false, e.actionable(ctx)))
}

// actionable is the predicate that drops the creatures no verb could act on, or
// nil when every verb takes any creature. Verbs apply in order, so a readying
// verb answers the readiness a later verb wants: "ready and use a friendly
// creature" may perfectly well pick an exhausted one.
func (e OnChooseCreature) actionable(ctx *EffectContext) func(LocalID) bool {
	var narrowing []narrowingVerb
	readied := false
	for _, v := range e.Verbs {
		if _, ok := v.(ReadyVerb); ok {
			readied = true
			continue
		}
		if n, ok := v.(narrowingVerb); ok && !readied {
			narrowing = append(narrowing, n)
		}
	}
	if len(narrowing) == 0 {
		return nil
	}
	return func(id LocalID) bool {
		for _, n := range narrowing {
			if !n.canApplyTo(ctx, id) {
				return false
			}
		}
		return true
	}
}

// declinable reports that the verbs hang off a single clickable creature.
func (e OnChooseCreature) declinable() bool { return e.Target.isChosen() }

// resolveOptional is Resolve under a May: the creature is asked declinably, so
// "you may ready and fight with a neighboring creature" is answered by clicking
// that neighbor rather than by a separate Yes/No.
func (e OnChooseCreature) resolveOptional(ctx *EffectContext) bool {
	return e.applyTo(ctx, e.Target.selectWith(ctx, true, e.actionable(ctx)))
}

// applyTo runs every verb over each selected creature and reports whether any
// creature was acted on. A creature that has left play takes no more verbs: an
// earlier sentence can destroy the very creature a later one names (Transposition
// Sandals swaps a creature off the flank that was keeping it alive, then says to
// use it), and a verb can destroy the creature the next verb would act on.
func (e OnChooseCreature) applyTo(ctx *EffectContext, ids []LocalID) bool {
	acted := false
	for _, id := range ids {
		for _, v := range e.Verbs {
			if !ctx.Resolver.InPlay(id) {
				break
			}
			v.Apply(ctx, id)
			acted = true
		}
	}
	return acted
}

// applyToChosen is one declinable pass that skips the creatures already spent,
// returning the ones it acted on. It is how OneAtATime keeps each pass on a
// different creature without the Target itself having to know about the others.
func (e OnChooseCreature) applyToChosen(ctx *EffectContext, spent []LocalID) []LocalID {
	actionable := e.actionable(ctx)
	ids := e.Target.selectWith(ctx, true, func(id LocalID) bool {
		if slices.Contains(spent, id) {
			return false
		}
		return actionable == nil || actionable(id)
	})
	e.applyTo(ctx, ids)
	return ids
}

// UseVerb uses the chosen creature. The controller picks how to use it — reap,
// fight, or its "Action:" ability — and that use resolves completely, nesting any
// further uses it triggers, before control returns. This is "Use a friendly
// creature": a creature can only be used while ready, so only ready creatures are
// offered and an ability with none to offer asks nothing.
type UseVerb struct{}

// VerbText returns the verb phrase.
func (UseVerb) VerbText() string { return "use" }

// canApplyTo offers only ready creatures: a creature can only be used while
// ready, so choosing an exhausted one would spend the ability on nothing.
func (UseVerb) canApplyTo(ctx *EffectContext, target LocalID) bool {
	return !ctx.Resolver.Exhausted(target)
}

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
	if ctx.Resolver.HasTrigger(target, TriggerAction) {
		labels = append(labels, "use its action")
		uses = append(uses, func() { ctx.Resolver.UseActionOf(ctx.Controller, target) })
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

// OneAtATime repeats a chosen-creature action several passes over, each pass on
// a creature no earlier pass took — Relentless Assault readies and fights with
// up to 3 different friendly creatures, one at a time. Each pass resolves fully
// (including any fight it triggers) before the next choice is offered, which is
// what "one at a time" means: the controller sees the board each pass and may
// stop early by declining.
type OneAtATime struct {
	// Times is how many passes the controller may take at most.
	Times  int
	Target Target
	Verbs  []CreatureVerb
}

// validate requires an explicit target and a positive number of passes.
func (e OneAtATime) validate() error {
	if !e.Target.valid() {
		return errUnsetTarget("OneAtATime")
	}
	if e.Times <= 0 {
		return fmt.Errorf("OneAtATime: Times must be positive")
	}
	return nil
}

// each is the single pass this effect repeats.
func (e OneAtATime) each() OnChooseCreature {
	return OnChooseCreature{Target: e.Target, Verbs: e.Verbs}
}

// Text renders the effect, e.g.
// "ready and fight with up to 3 different friendly creatures, one at a time".
func (e OneAtATime) Text() string {
	verbs := make([]string, 0, len(e.Verbs))
	for _, v := range e.Verbs {
		verbs = append(verbs, v.VerbText())
	}
	return fmt.Sprintf(
		"%s up to %d different %ss, one at a time",
		strings.Join(verbs, " and "),
		e.Times,
		singularNoun(e.Target.Text()),
	)
}

// Resolve takes up to Times passes, stopping as soon as a pass acts on nobody —
// the pool ran dry or the controller declined.
func (e OneAtATime) Resolve(ctx *EffectContext) {
	pass := e.each()
	var used []LocalID
	for range e.Times {
		chosen := pass.applyToChosen(ctx, used)
		if len(chosen) == 0 {
			return
		}
		used = append(used, chosen...)
	}
}

// RepeatedFight readies a creature and fights with it several times over, each
// fight against an enemy no earlier fight of this effect used — One Stood Against
// Many fights 3 times, each time against a different enemy creature. Where
// OneAtATime makes the *actor* differ each pass, RepeatedFight makes the *enemy*
// differ: the same creature may be chosen to fight every time. Each fight resolves
// fully before the next choice is offered.
type RepeatedFight struct {
	// Times is how many fights the effect takes at most.
	Times int
	// Target chooses the creature that readies and fights, once per fight.
	Target Target
}

// validate requires an explicit target and a positive number of fights.
func (e RepeatedFight) validate() error {
	if !e.Target.valid() {
		return errUnsetTarget("RepeatedFight")
	}
	if e.Times <= 0 {
		return fmt.Errorf("RepeatedFight: Times must be positive")
	}
	return nil
}

// Text renders the effect, e.g. "ready and fight with a friendly creature 3 times,
// each time against a different enemy creature. Resolve these fights one at a
// time".
func (e RepeatedFight) Text() string {
	return fmt.Sprintf(
		"ready and fight with %s %d times, each time against a different enemy "+
			"creature. Resolve these fights one at a time",
		e.Target.Text(),
		e.Times,
	)
}

// Resolve takes up to Times fights, stopping as soon as there is no untouched
// enemy left, no creature to fight with, or the controller declines a choice.
func (e RepeatedFight) Resolve(ctx *EffectContext) {
	var fought []LocalID
	for range e.Times {
		enemies := e.remaining(ctx, fought)
		if len(enemies) == 0 {
			return
		}
		ids := e.Target.Select(ctx)
		if len(ids) == 0 {
			return
		}
		attacker := ids[0]
		ReadyVerb{}.Apply(ctx, attacker)
		enemy, ok := fightAmong(ctx, attacker, enemies)
		if !ok {
			return
		}
		fought = append(fought, enemy)
	}
}

// remaining lists the enemy creatures no earlier fight of this effect used.
func (e RepeatedFight) remaining(ctx *EffectContext, fought []LocalID) []LocalID {
	var enemies []LocalID
	for _, id := range ctx.Resolver.Battleline(ctx.Opponent()) {
		if !slices.Contains(fought, id) {
			enemies = append(enemies, id)
		}
	}
	return enemies
}
