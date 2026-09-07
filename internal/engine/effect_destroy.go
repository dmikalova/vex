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
	bonuses := make(map[LocalID]int, len(ids))
	for _, id := range ids {
		controllers[id] = ctx.Resolver.Controller(id)
		bonuses[id] = ctx.Resolver.AemberBonus(id)
	}
	ctx.Resolver.DestroyEachFrom(ctx.Controller, ctx.Source, ids)
	for _, id := range ids {
		if !resolverInPlay(ctx, id) {
			ctx.Produced.Destroyed[controllers[id]]++
			ctx.Produced.AemberBonusDestroyed += bonuses[id]
		}
	}
	return len(ids) > 0
}

// DestroyChosen destroys any number of creatures the controller picks from the
// Target pool, chosen one at a time and then destroyed together — Martyr's End
// destroys any number of friendly creatures. It tallies them into
// Produced.Destroyed so a following "gain 1 Æmber for each creature destroyed this
// way" can pay out.
type DestroyChosen struct {
	Target Target
}

// validate requires an explicit target pool to choose from.
func (e DestroyChosen) validate() error {
	if !e.Target.valid() {
		return errUnsetTarget("DestroyChosen")
	}
	return nil
}

// Text renders the effect, e.g. "destroy any number of friendly creatures".
func (e DestroyChosen) Text() string {
	return "destroy any number of " + singularNoun(e.Target.Text()) + "s"
}

// Resolve gathers the controller's picks one at a time, then destroys them all at
// once so their Destroyed abilities see each other still in play.
func (e DestroyChosen) Resolve(ctx *EffectContext) {
	picked := map[LocalID]bool{}
	var chosen []LocalID
	for {
		var cands []LocalID
		for _, id := range e.Target.Select(ctx) {
			if !picked[id] {
				cands = append(cands, id)
			}
		}
		if len(cands) == 0 {
			break
		}
		pick, ok := ctx.ChooseCardOptional("Choose a creature to destroy", cands)
		if !ok {
			break
		}
		picked[pick] = true
		chosen = append(chosen, pick)
	}
	Destroy{}.destroy(ctx, chosen)
}

// DestroyMostPowerfulUnlessReadyHouse destroys the most powerful creature
// controlled by each player who does not control a ready creature of House —
// Quicksand spares any player fielding a ready Untamed creature and destroys the
// most powerful creature of everyone else. When a player's largest creatures tie,
// the effect's controller chooses which one is destroyed.
type DestroyMostPowerfulUnlessReadyHouse struct {
	// House is the house whose ready creature spares its controller.
	House House
}

// validate requires an explicit house.
func (e DestroyMostPowerfulUnlessReadyHouse) validate() error {
	if e.House == HouseNone {
		return fmt.Errorf("DestroyMostPowerfulUnlessReadyHouse: House must be set")
	}
	return nil
}

// Text renders the effect, e.g. "destroy the most powerful creature controlled by
// each player who does not control a ready Untamed creature".
func (e DestroyMostPowerfulUnlessReadyHouse) Text() string {
	return "destroy the most powerful creature controlled by each player who " +
		"does not control a ready " + e.House.String() + " creature"
}

// Resolve destroys the most powerful creature of each player who lacks a ready
// creature of House, all at once so their Destroyed abilities see each other.
func (e DestroyMostPowerfulUnlessReadyHouse) Resolve(ctx *EffectContext) {
	var doomed []LocalID
	for p := 0; p < 2; p++ {
		if e.controlsReadyHouse(ctx, p) {
			continue
		}
		if ids := ctx.Resolver.Battleline(p); len(ids) > 0 {
			doomed = append(doomed, mostPowerfulN{n: 1}.refine(ctx, ids)...)
		}
	}
	Destroy{}.destroy(ctx, doomed)
}

// controlsReadyHouse reports whether player p controls a ready creature of House.
func (e DestroyMostPowerfulUnlessReadyHouse) controlsReadyHouse(ctx *EffectContext, p int) bool {
	for _, id := range ctx.Resolver.Battleline(p) {
		if ctx.Resolver.House(id) == e.House && !ctx.Resolver.Exhausted(id) {
			return true
		}
	}
	return false
}
