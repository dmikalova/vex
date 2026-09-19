package engine

import "errors"

// ForEach resolves an effect once for each of a running count, so every repetition
// makes its own choices — Mothership Support deals 2 damage per friendly ready
// Mars creature and may pick a different creature each time. It is the
// choose-again counterpart to a Per clause, which multiplies one effect's amount
// against a single target.
type ForEach struct {
	// Times is how many repetitions to run.
	Times Count
	// Do is the effect resolved once per repetition.
	Do Effect
}

// Text renders the repetition as a leading "for each" clause over the effect's
// own phrase (rule 9).
func (e ForEach) Text() string { return forEach(e.Times, e.Do.Text()) }

// validate requires both halves of the repetition.
func (e ForEach) validate() error {
	if e.Times == nil {
		return errors.New("ForEach needs a Times count")
	}
	if e.Do == nil {
		return errors.New("ForEach needs an effect to Do")
	}
	return validateEffect(e.Do)
}

// Resolve runs the effect once per repetition, counting first so an effect that
// changes the board does not change how many times it runs.
func (e ForEach) Resolve(ctx *EffectContext) {
	for range e.Times.Value(ctx) {
		e.Do.Resolve(ctx)
	}
}

// RuleOfSix is the most times a card *name* can be used in one turn by one
// player. Every usage of a name shares one pool of six: a play, a discard, a use
// (reap, fight, or Action:), each Destroyed: resolution, each loop of a
// self-repeating ability past its first, and each chained Replicator-style trigger
// past the first. The use that fires an ability buys its first resolution free;
// what is left in the pool bounds the rest. Bait and Switch's "steal 1 Æmber ->
// repeat this effect" resolves the initial steal plus at most five repeats, so it
// steals six at most however far ahead the opponent is.
const RuleOfSix = 6

// Repeat resolves Do and repeats it as its Gate allows — automatically while a
// board condition holds (While), optionally at the controller's choice (MayWhile),
// or once if they pay by exalting a creature (ByExalting). The gate carries both
// the continue decision and the trailing "repeat" clause, so the axis a repeat
// varies along is one pluggable strategy rather than a family of near-duplicate
// nodes.
type Repeat struct {
	Do   Effect
	Gate RepeatGate
}

// validate checks the repeated effect and the gate.
func (e Repeat) validate() error {
	if e.Do == nil {
		return errors.New("Repeat needs an effect to Do")
	}
	if e.Gate == nil {
		return errors.New("Repeat needs a Gate")
	}
	if err := e.Gate.validate(); err != nil {
		return err
	}
	return validateEffect(e.Do)
}

// Text renders the effect, closing with the gate's own "repeat" clause.
func (e Repeat) Text() string { return e.Gate.text(e.Do) }

// Resolve runs Do and repeats it however the gate allows.
func (e Repeat) Resolve(ctx *EffectContext) { e.Gate.run(ctx, e.Do) }

// RepeatGate is the axis a Repeat varies along: it drives the loop that reruns the
// effect and renders the trailing clause that describes when it repeats. Three
// gates cover the shapes in play — While (automatic), MayWhile (optional), and
// ByExalting (paid, once).
type RepeatGate interface {
	// run resolves the effect once and then again for as long as the gate allows,
	// bounded by the Rule of Six.
	run(ctx *EffectContext, do Effect)
	// text renders the repeated effect's phrase plus the gate's "repeat" clause.
	text(do Effect) string
	// validate reports a misconfigured gate.
	validate() error
}

// While repeats the effect automatically while it keeps doing something and Cond
// holds — Numquid the Fair's "destroy an enemy creature -> if you are overwhelmed,
// repeat this effect." Do runs at least once; the loop stops as soon as Do does
// nothing (its gate is false) or Cond is not met. When Do cannot report progress
// the Rule of Six alone bounds the loop.
type While struct {
	Cond Condition
}

func (g While) validate() error { return nil }

func (g While) text(do Effect) string {
	return do.Text() + " -> " + g.Cond.CondText() + ", repeat this effect"
}

func (g While) run(ctx *EffectContext, do Effect) {
	// The first resolution rides free on the use that triggered it; each further
	// loop counts one usage against this card name's Rule-of-Six pool, so the loop
	// runs until the pool is spent.
	if !resolveGateOf(ctx, do) || !g.Cond.Met(ctx) {
		return
	}
	for !ctx.Resolver.AtRuleOfSix(ctx.Source) {
		ctx.Resolver.RecordUsage(ctx.Source)
		if !resolveGateOf(ctx, do) || !g.Cond.Met(ctx) {
			return
		}
	}
}

// MayWhile repeats the effect at the controller's choice for as long as Cond holds
// — Bouncing Deathquark's "-> if there is a friendly creature in play, you may
// repeat this effect." Do should make progress toward failing Cond so the loop can
// end. When Do leads with a single clickable choice, each repeat is driven by that
// choice — the controller keeps picking to repeat, or passes with Done — rather
// than a separate Yes/No question.
type MayWhile struct {
	Cond Condition
}

func (g MayWhile) validate() error { return nil }

func (g MayWhile) text(do Effect) string {
	return do.Text() + " -> " + g.Cond.CondText() + ", you may repeat this effect"
}

func (g MayWhile) run(ctx *EffectContext, do Effect) {
	// The first resolution rides free on the use that triggered it; each chosen
	// repeat counts one usage against this card name's Rule-of-Six pool.
	do.Resolve(ctx)
	d, byChoice := do.(declinableEffect)
	byChoice = byChoice && d.declinable()
	for !ctx.Resolver.AtRuleOfSix(ctx.Source) {
		if !g.Cond.Met(ctx) {
			return
		}
		if byChoice {
			if !d.resolveOptional(ctx) {
				return
			}
			ctx.Resolver.RecordUsage(ctx.Source)
			continue
		}
		if ctx.ChooseOption("Repeat this effect?", []string{"Yes", "No"}) != 0 {
			return
		}
		ctx.Resolver.RecordUsage(ctx.Source)
		do.Resolve(ctx)
	}
}

// ByExalting repeats the effect once if the controller exalts Creature to pay for
// it — Phalanx Strike's "You may exalt a friendly creature to repeat the preceding
// effect." "Repeat the preceding effect" repeats only the effect before the exalt
// clause, so the offer is made once and does not chain.
type ByExalting struct {
	// Creature names the creature exalted to pay for the single repeat.
	Creature Target
}

func (g ByExalting) validate() error {
	if !g.Creature.valid() {
		return errUnsetTarget("Repeat.ByExalting")
	}
	return nil
}

func (g ByExalting) text(do Effect) string {
	return punctuate(do.Text()) +
		" You may exalt " + g.Creature.Text() + " to repeat the preceding effect"
}

func (g ByExalting) run(ctx *EffectContext, do Effect) {
	do.Resolve(ctx)
	ids := g.exaltChoice(ctx)
	if len(ids) == 0 {
		return
	}
	for _, id := range ids {
		ctx.Resolver.AddAmberOn(id, 1)
		ctx.Resolver.Record(AemberExalted{
			Player:   ctx.Controller,
			Creature: id,
			Amount:   1,
		})
	}
	do.Resolve(ctx)
}

// exaltChoice offers the exalt that pays for the one repeat, or none to stop. A
// chosen target is its own declinable prompt. A back-reference like the chosen
// creature has no pool to pick from, but there is still one card to point at, so
// it is offered as a click on that creature — the same click-or-Done shape as
// every other optional single-card decision. Only a back-reference that resolves
// to several creatures has no single card to click, so it keeps the Yes/No.
func (g ByExalting) exaltChoice(ctx *EffectContext) []LocalID {
	if g.Creature.isChosen() {
		return g.Creature.SelectOptional(ctx)
	}
	ids := g.Creature.Select(ctx)
	if len(ids) == 0 {
		return nil
	}
	prompt := "Exalt " + g.Creature.Text() + " to repeat the preceding effect"
	if len(ids) == 1 {
		if _, ok := ctx.ChooseCardOptional(prompt, ids); !ok {
			return nil
		}
		return ids
	}
	if ctx.ChooseOption(prompt+"?", []string{"Yes", "No"}) != 0 {
		return nil
	}
	return ids
}
