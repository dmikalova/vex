package engine

import "fmt"

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
// "choose a creature - destroy …" when the target's refinement leads with a choice.
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
	powers := make(map[LocalID]int, len(ids))
	for _, id := range ids {
		controllers[id] = ctx.Resolver.Controller(id)
		powers[id] = ctx.Resolver.Power(id)
	}
	ctx.Resolver.DestroyEachFrom(ctx.Controller, ctx.Source, ids)
	if len(ids) == 1 {
		// Bind the destroyed card so a following effect can read it even after it
		// leaves play — its controller ("its owner discards") or its printed Æmber
		// bonus (Rustgnawer). Those are printed properties, so they survive.
		ctx.It, ctx.HasIt = ids[0], true
		ctx.ItController = controllers[ids[0]]
	}
	for _, id := range ids {
		if !resolverInPlay(ctx, id) {
			ctx.Produced.Destroyed[controllers[id]]++
			ctx.Produced.DestroyedPower += powers[id]
		}
	}
	return len(ids) > 0
}

// DestroyChosen destroys creatures the controller picks from the Target pool,
// chosen one at a time and then destroyed together — Martyr's End destroys any
// number of friendly creatures. Amount fixes the count when set: Ritual of Tognath
// destroys exactly 2 friendly creatures; the zero value destroys any number. It
// tallies them into Produced.Destroyed so a following "gain 1 Æmber for each
// creature destroyed this way" can pay out.
type DestroyChosen struct {
	Target Target
	Amount int
}

// validate requires an explicit target pool to choose from.
func (e DestroyChosen) validate() error {
	if !e.Target.valid() {
		return errUnsetTarget("DestroyChosen")
	}
	if e.Amount < 0 {
		return fmt.Errorf("DestroyChosen: negative Amount %d", e.Amount)
	}
	return nil
}

// Text renders the effect, e.g. "destroy any number of friendly creatures" or,
// with a fixed Amount, "destroy 2 friendly creatures".
func (e DestroyChosen) Text() string {
	if e.Amount > 0 {
		return fmt.Sprintf("destroy %d %ss", e.Amount, singularNoun(e.Target.Text()))
	}
	return "destroy any number of " + singularNoun(e.Target.Text()) + "s"
}

// Resolve gathers the controller's picks one at a time, then destroys them all at
// once so their Destroyed abilities see each other still in play. With a fixed
// Amount the picks are mandatory up to that many; otherwise the controller takes
// any number.
func (e DestroyChosen) Resolve(ctx *EffectContext) {
	chosen := pickCards(
		ctx,
		"Choose a creature to destroy",
		e.Amount,
		e.Amount == 0,
		func() []LocalID {
			return e.Target.Select(ctx)
		},
	)
	Destroy{}.destroy(ctx, chosen)
}

// Gather selects the set of creatures a BatchDestroy destroys and renders the noun
// phrase that names them. Modeling "which creatures" as a strategy lets
// BatchDestroy keep the one simultaneous-destruction batch while the selection axis
// varies per card.
type Gather interface {
	gather(ctx *EffectContext) []LocalID
	gatherText() string
	validate() error
}

// BatchDestroy destroys a computed set of creatures in one simultaneous batch, so
// their Destroyed abilities see each other still in play (the KeyForge
// simultaneous-destruction rule). Its Gather strategy picks the set and names it;
// BatchDestroy carries out the single destruction. Reach for it whenever an effect
// must destroy a set it has to compute — one per player, a power ranking — rather
// than a plain Target the runtime already batches: two sequential Destroy ops would
// break the shared batch.
type BatchDestroy struct {
	Gather Gather
}

// validate requires a gather strategy and defers to its configuration.
func (e BatchDestroy) validate() error {
	if e.Gather == nil {
		return fmt.Errorf("BatchDestroy: Gather must be set")
	}
	return e.Gather.validate()
}

// Text renders "destroy <gathered set>".
func (e BatchDestroy) Text() string { return "destroy " + e.Gather.gatherText() }

// Resolve destroys the gathered set all at once so their Destroyed abilities see
// each other still in play.
func (e BatchDestroy) Resolve(ctx *EffectContext) {
	Destroy{}.destroy(ctx, e.Gather.gather(ctx))
}

// EachPlayerUnless gathers, from each player not spared by a per-player board
// condition, the creatures a refinement keeps of that player's battleline —
// Quicksand takes each unspared player's most powerful creature. Spare is read
// from each player's own perspective, so it names its board as the controller's
// ("friendly"): a player fielding a match is skipped. Take then keeps the doomed
// creatures of everyone else, with ties broken by the effect's controller.
type EachPlayerUnless struct {
	// Spare is the board condition that, read from a player's own perspective,
	// exempts that player. It must be phrased as the controller's board (Player:
	// Controller) because it is re-based onto each player in turn.
	Spare InPlay
	// Take keeps the doomed creatures of an unspared player's battleline.
	Take Refinement
}

// validate requires a Take refinement and a Spare phrased from the controller's
// perspective, since Spare is re-based onto each player in turn.
func (g EachPlayerUnless) validate() error {
	if g.Take == nil {
		return fmt.Errorf("EachPlayerUnless: Take must be set")
	}
	if g.Spare.Player != Controller {
		return fmt.Errorf("EachPlayerUnless: Spare must be phrased from the " +
			"controller's perspective (Player: Controller)")
	}
	return nil
}

// gatherText renders the noun phrase, e.g. "the most powerful creature controlled
// by each player who does not have a friendly ready Untamed creature in play".
func (g EachPlayerUnless) gatherText() string {
	return g.Take.clause("each creature") +
		" controlled by each player who does not have " +
		indefinite(g.Spare.noun()) + " in play"
}

// gather takes the doomed creatures of each player the Spare condition does not
// exempt. Spare is evaluated from each player's own perspective by re-basing the
// context's controller onto that player; ties are broken by the effect's actual
// controller, so refine reads the original context.
func (g EachPlayerUnless) gather(ctx *EffectContext) []LocalID {
	var doomed []LocalID
	for p := 0; p < 2; p++ {
		pctx := *ctx
		pctx.Controller = p
		if g.Spare.Met(&pctx) {
			continue
		}
		doomed = append(doomed, g.Take.refine(ctx, ctx.Resolver.Battleline(p))...)
	}
	return doomed
}

// DestroyEachCreatureAtEndOfTurn schedules "destroy each creature" to resolve in
// the active player's end-of-turn window rather than now — Ragnarok wipes the board
// only once the turn it is played is ending, after its owner has spent the turn
// fighting for Æmber. The wipe orders alongside the in-play "at the end of your
// turn" abilities (ADR 0013); the schedule survives the ready phase, which runs
// earlier.
type DestroyEachCreatureAtEndOfTurn struct{}

// Text renders the effect, e.g. "at the end of the turn, destroy each creature".
func (e DestroyEachCreatureAtEndOfTurn) Text() string {
	return "at the end of the turn, destroy each creature"
}

// Resolve schedules the board wipe into the end-of-turn window; that window carries
// it out.
func (e DestroyEachCreatureAtEndOfTurn) Resolve(ctx *EffectContext) {
	ctx.Resolver.ScheduleAtEndOfTurn(ctx.Source, schedDestroyEachCreature)
}
