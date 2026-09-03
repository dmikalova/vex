package engine

// ReadyIfFirstUse readies a creature only while the current use is its first use
// this turn. A creature is USED when it reaps, fights, or uses an Action: ability;
// the use count is advanced before Fight:/Reap: abilities resolve, so the first
// use is count 1 during this effect.
//
// Readying an exhausted creature stands it back up so it can be used again this
// turn.
type ReadyIfFirstUse struct {
	Target Target
}

// validate requires an explicit target.
func (e ReadyIfFirstUse) validate() error {
	if !e.Target.valid() {
		return errUnsetTarget("ReadyIfFirstUse")
	}
	return nil
}

// Text renders Rocket Boots' conditional readying clause.
func (e ReadyIfFirstUse) Text() string {
	return "if this is the first time " + e.Target.Text() + " was used this turn, ready it"
}

// Resolve readies each selected creature whose current use is its first this turn.
func (e ReadyIfFirstUse) Resolve(ctx *EffectContext) {
	for _, id := range e.Target.Select(ctx) {
		if ctx.Resolver.TimesUsedThisTurn(id) == 1 {
			ctx.Resolver.SetExhausted(id, false)
		}
	}
}

// ReadyCreatures readies creatures the controller chooses from its Target one at
// a time, up to a Max count — Commpod's "for each card revealed this way, ready a
// Mars creature." The controller picks which creature each time; it is mandatory
// while an eligible (exhausted) creature remains, and simply stops once none do.
//
// Readying an exhausted creature stands it back up so it can be used again this
// turn.
type ReadyCreatures struct {
	// Max caps how many creatures are readied; nil means one.
	Max Count
	// Target is the pool of creatures that may be readied. It should be an "each"
	// target (the candidates), not a chosen one — the controller picks within the
	// pool one at a time.
	Target Target
}

// validate requires an explicit target.
func (e ReadyCreatures) validate() error {
	if !e.Target.valid() {
		return errUnsetTarget("ReadyCreatures")
	}
	return nil
}

// Text renders the effect, e.g. "for each card revealed this way, ready a Mars
// creature".
func (e ReadyCreatures) Text() string {
	return forEach(e.Max, "ready a "+singularNoun(e.Target.Text()))
}

// Resolve readies up to Max creatures, one at a time, choosing among the exhausted
// members of the Target pool. It stops early only when no eligible creature is left.
func (e ReadyCreatures) Resolve(ctx *EffectContext) {
	n := 1
	if e.Max != nil {
		n = e.Max.Value(ctx)
	}
	for i := 0; i < n; i++ {
		var cands []LocalID
		for _, id := range e.Target.Select(ctx) {
			if ctx.Resolver.Exhausted(id) {
				cands = append(cands, id)
			}
		}
		if len(cands) == 0 {
			return
		}
		id, ok := ctx.ChooseCreature("Choose a creature to ready", cands)
		if !ok {
			return
		}
		ctx.Resolver.SetExhausted(id, false)
		ctx.Resolver.Record(CreatureReadied{Creature: id})
	}
}
