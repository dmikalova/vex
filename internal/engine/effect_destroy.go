package engine

// Destroying a creature removes it from play. When an effect destroys several
// creatures they are destroyed simultaneously: every one is tagged for
// destruction and stays in play while their "Destroyed:" abilities resolve, in an
// order the controller chooses, so each ability sees the others still present;
// only then does each creature still in play move to the discard pile, along with
// its upgrades. A destroy effect can target every creature or only those matching
// a filter, such as "each creature with power 3 or lower".
type Destroy struct {
	Target Target
}

// validate requires an explicit target.
func (e Destroy) validate() error {
	if !e.Target.valid() {
		return errUnsetTarget("Destroy")
	}
	return nil
}

func (e Destroy) verb() string       { return "destroy" }
func (e Destroy) targetText() string { return e.Target.Text() }

// Text renders the effect, e.g. "destroy each creature with power 3 or lower", or
// "choose a creature - destroy …" when the target's selector leads with a choice.
func (e Destroy) Text() string {
	body := e.verb() + " " + e.targetText()
	if lead, ok := e.Target.leadIn(); ok {
		return lead + " - " + body
	}
	return body
}

// Resolve destroys each selected creature, letting the controller order them.
func (e Destroy) Resolve(ctx *EffectContext) { e.resolveGate(ctx) }

// resolveGate destroys the selected creatures simultaneously and reports whether
// any were, so Destroy can be the first half of a Then ("destroy a creature ->
// ..."). It tallies how many actually left play on the context (read by
// CardsDestroyedFewerThan), counting after the batch so a save (Armageddon Cloak)
// is not counted.
func (e Destroy) resolveGate(ctx *EffectContext) bool {
	return e.destroy(ctx, e.Target.Select(ctx))
}

// declinable reports that the destruction is a single clickable creature.
func (e Destroy) declinable() bool { return e.Target.isChosen() }

// vacuous reports that there is nothing here to destroy, so a "you may" wrapping
// it need not ask.
func (e Destroy) vacuous(ctx *EffectContext) bool { return e.Target.empty(ctx) }

// resolveOptional is resolveGate under a May: the creature is asked declinably, so
// "you may destroy another friendly creature" is answered by clicking that
// creature rather than by a separate Yes/No.
func (e Destroy) resolveOptional(ctx *EffectContext) bool {
	return e.destroy(ctx, e.Target.SelectOptional(ctx))
}

// destroy carries out the destruction of an already-selected set.
func (e Destroy) destroy(ctx *EffectContext, ids []LocalID) bool {
	controllers := make(map[LocalID]int, len(ids))
	for _, id := range ids {
		controllers[id] = ctx.Resolver.Controller(id)
	}
	ctx.Resolver.DestroyEachFrom(ctx.Controller, ctx.Source, ids)
	for _, id := range ids {
		if !resolverInPlay(ctx, id) {
			ctx.Produced.Destroyed[controllers[id]]++
		}
	}
	return len(ids) > 0
}
