package engine

// Destroying a creature removes it from play. When an effect destroys several
// creatures they are destroyed simultaneously: every one is tagged for
// destruction and stays in play while their "Destroyed:" abilities resolve, in an
// order the controller chooses, so each ability sees the others still present;
// only then does each creature still in play move to the discard pile, along with
// its upgrades. A destroy effect can target every creature or only those matching
// a filter, such as "each creature with power 3 or lower".
//
//rulebook:effect Destroy
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
	ids := e.Target.Select(ctx)
	ctx.Resolver.DestroyEach(ctx.Controller, ids)
	for _, id := range ids {
		if !resolverInPlay(ctx, id) {
			ctx.Produced.Destroyed++
		}
	}
	return len(ids) > 0
}
