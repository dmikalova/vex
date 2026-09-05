package engine

import "fmt"

// Exhaust turns a creature sideways so it cannot be used again until it readies.
// Like Stun and Unstun it is a simple "verb the target" effect, so it folds
// together with a neighbouring status change on the same target when they share a
// Sequence (see combinable) — e.g. "stun and exhaust this creature".

// Exhausting a creature turns it sideways so it cannot be used again until it
// readies at the end of its controller's turn. It exhausts each creature the
// effect targets.
type Exhaust struct {
	Target Target
}

// validate requires an explicit target.
func (e Exhaust) validate() error {
	if !e.Target.valid() {
		return errUnsetTarget("Exhaust")
	}
	return nil
}

func (e Exhaust) verb() string       { return "exhaust" }
func (e Exhaust) targetText() string { return e.Target.Text() }

// Text renders the effect, e.g. "exhaust this creature".
func (e Exhaust) Text() string { return e.verb() + " " + e.targetText() }

// Resolve exhausts each selected creature.
func (e Exhaust) Resolve(ctx *EffectContext) {
	for _, id := range e.Target.Select(ctx) {
		ctx.Resolver.SetExhausted(id, true)
	}
}

// ExhaustCreatures exhausts up to Max creatures the controller chooses one at a
// time from the Target pool — Nocturnal Maneuver's "exhaust up to 3 creatures".
// Already-exhausted creatures are not offered.
type ExhaustCreatures struct {
	Max    int
	Target Target
}

// validate requires a target and a positive maximum.
func (e ExhaustCreatures) validate() error {
	if !e.Target.valid() {
		return errUnsetTarget("ExhaustCreatures")
	}
	if e.Max < 1 {
		return fmt.Errorf("ExhaustCreatures: Max must be positive")
	}
	return nil
}

// Text renders the effect, e.g. "exhaust up to 3 creatures".
func (e ExhaustCreatures) Text() string {
	return fmt.Sprintf("exhaust up to %d %s", e.Max, singularNoun(e.Target.Text())+"s")
}

// Resolve exhausts up to Max ready creatures from the pool, one at a time, letting
// the controller stop early.
func (e ExhaustCreatures) Resolve(ctx *EffectContext) {
	for i := 0; i < e.Max; i++ {
		var cands []LocalID
		for _, id := range e.Target.Select(ctx) {
			if !ctx.Resolver.Exhausted(id) {
				cands = append(cands, id)
			}
		}
		chosen, ok := ctx.ChooseCardOptional("Choose a creature to exhaust", cands)
		if !ok {
			return
		}
		ctx.Resolver.SetExhausted(chosen, true)
		ctx.Resolver.Record(CreatureExhausted{Creature: chosen})
	}
}

// Readying a creature turns it upright again so it can be used, the opposite of
// exhausting. It readies each creature the effect targets.
type Ready struct {
	Target Target
}

// validate requires an explicit target.
func (e Ready) validate() error {
	if !e.Target.valid() {
		return errUnsetTarget("Ready")
	}
	return nil
}

func (e Ready) verb() string       { return "ready" }
func (e Ready) targetText() string { return e.Target.Text() }

// Text renders the effect, e.g. "ready this creature".
func (e Ready) Text() string { return e.verb() + " " + e.targetText() }

// Resolve readies each selected creature.
func (e Ready) Resolve(ctx *EffectContext) {
	for _, id := range e.Target.Select(ctx) {
		ctx.Resolver.SetExhausted(id, false)
	}
}
