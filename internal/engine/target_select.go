package engine

import "slices"

// This file holds a Target's selection machinery: resolving a Target into
// concrete card ids (Select, SelectOptional, selectWith), narrowing them by the
// Target's per-card Filter (admitted), the base sets each Kind draws from
// (selectBase), and the battleline geometry the flank and neighbor axes read
// (onFlank, isNeighbor, neighbors, …). See target.go for the Target type,
// filter.go for the Filter it narrows by, and target_refinement.go for the
// Refinement strategies selectWith applies.

// Select resolves the target into concrete card ids, applying its filters. For a
// chosen kind it asks the controller to pick one of the filtered candidates
// (returning nil when there are none or the choice is declined).
func (t Target) Select(ctx *EffectContext) []LocalID {
	return t.selectWith(ctx, false, nil)
}

// wholeSide reports whether the target is a plain, unfiltered whole side — every
// friendly creature or every enemy creature with no narrowing — and which player
// that side belongs to. A side-wide, live-read status (Shield of Justice's
// damage immunity) can stand in for selecting those creatures one by one, which
// also covers creatures that arrive after the effect resolves. Any narrowing
// filter (a trait, a power bound, a flank) makes it not a whole side, so it falls
// back to selecting concrete creatures.
func (t Target) wholeSide(controller int) (player int, ok bool) {
	switch t.Kind {
	case TargetEachFriendlyCreature:
		player = controller
	case TargetEachEnemyCreature:
		player = 1 - controller
	default:
		return 0, false
	}
	if t != (Target{Kind: t.Kind}) {
		return 0, false
	}
	return player, true
}

// SelectOptional is Select inside a "you may": a chosen target is asked
// declinably, so the controller clicks the card they mean or passes, instead of
// answering a Yes/No and then being handed a pick they can no longer refuse. A
// target that chooses nothing has no decision to decline and behaves like Select.
func (t Target) SelectOptional(ctx *EffectContext) []LocalID {
	return t.selectWith(ctx, true, nil)
}

// candidates returns the filtered, refined creatures a chosen target would pick
// from, without making the choice — for an effect that repeats the choice itself,
// like a "for each" DealDamage that picks a creature per instance.
func (t Target) candidates(ctx *EffectContext) []LocalID {
	ids := t.admitted(ctx, t.selectBase(ctx))
	if t.Refinement != nil {
		ids = t.Refinement.refine(ctx, ids)
	}
	return ids
}

// empty reports that nothing matches this target, reading only the candidates a
// refinement would narrow: with nothing to narrow there is nothing to select, and
// unlike Select it asks the controller nothing.
func (t Target) empty(ctx *EffectContext) bool {
	return len(t.admitted(ctx, t.selectBase(ctx))) == 0
}

// couldSelect reports whether id is among the target's candidates without making
// the discretionary tie-break its own selection would — the tie-inclusive
// membership test conditions use (ItIsAmong). A refinement that ties at its
// cutoff counts every tied card as included.
func (t Target) couldSelect(ctx *EffectContext, id LocalID) bool {
	base := t.admitted(ctx, t.selectBase(ctx))
	if !slices.Contains(base, id) {
		return false
	}
	if t.Refinement == nil {
		return true
	}
	if mr, ok := t.Refinement.(membershipRefiner); ok {
		return mr.includes(ctx, base, id)
	}
	return slices.Contains(t.candidates(ctx), id)
}

// selectWith is the shared selection path; optional switches the chosen-kind
// prompt between a forced pick and a declinable one, and keep (when set) drops
// candidates the calling effect could not act on.
func (t Target) selectWith(
	ctx *EffectContext,
	optional bool,
	keep func(LocalID) bool,
) []LocalID {
	ids := t.admitted(ctx, t.selectBase(ctx))
	if keep != nil {
		kept := ids[:0:0]
		for _, id := range ids {
			if keep(id) {
				kept = append(kept, id)
			}
		}
		ids = kept
	}
	if t.Refinement != nil {
		ids = t.Refinement.refine(ctx, ids)
	}
	if !t.isChosen() {
		return t.expandNeighbors(ctx, ids)
	}
	if len(ids) == 0 {
		return nil
	}
	prompt := "Choose " + t.Text()
	var id LocalID
	var ok bool
	if optional {
		id, ok = ctx.ChooseCardOptional(prompt, ids)
	} else {
		id, ok = ctx.ChooseCreature(prompt, ids)
	}
	if !ok {
		return nil
	}
	return t.expandNeighbors(ctx, []LocalID{id})
}

// expandNeighbors applies the neighbour builders to an already-selected set:
// AndNeighbors keeps each selected creature and adds its battleline neighbors,
// NeighborsOf replaces the selection with them (Lord Golgotha hits the neighbors
// of the creature it fights, not that creature).
func (t Target) expandNeighbors(ctx *EffectContext, ids []LocalID) []LocalID {
	if !t.Neighbors.expands() {
		return ids
	}
	out := ids[:0:0]
	for _, id := range ids {
		if t.Neighbors.keepsSelected() {
			out = append(out, id)
		}
		ns := neighbors(ctx, id)
		// The fought creature may have left play in the fight that named it, so
		// "each neighbor of the fought creature" falls back to the neighbors the
		// fight snapshotted (Smite pops a warded neighbor even when the target dies).
		if t.Neighbors == NeighborsOnly && len(ns) == 0 && !resolverInPlay(ctx, id) {
			ns = ctx.Produced.Neighbors
		}
		out = append(out, ns...)
	}
	return out
}

// focus is the card an "another …" Target is other than: the card in context
// (ctx.It) when an effect has put one there, and otherwise the card doing the
// choosing. Falling back to the source is what keeps the exclusion honest —
// widening to every card instead would let "another creature" land on the very
// card it is defined against.
func focus(ctx *EffectContext) LocalID {
	if ctx.HasIt {
		return ctx.It
	}
	return ctx.Source
}

// excludesFocus reports whether the target is an "another …" one, i.e. defined
// by excluding the card in focus.
func (t Target) excludesFocus() bool {
	return t.Filter.Except == ExcludeFocus
}

// isChosen reports whether the Kind resolves to a single player-chosen creature.
func (t Target) isChosen() bool {
	return t.Kind == TargetChosenCreature || t.Kind == TargetChosenEnemyCreature ||
		t.Kind == TargetChosenFriendlyCreature ||
		t.Kind == TargetChosenArtifact || t.Kind == TargetChosenEnemyArtifact || t.Kind == TargetChosenFriendlyArtifact || t.Kind == TargetChosenCreatureOrArtifact ||
		t.Kind == TargetChosenUpgrade ||
		t.Kind == TargetChosenFriendlyCreatureOrArtifact ||
		t.Kind == TargetChosenEnemyCreatureOrArtifact
}

// isOfMostPopulousHouse reports whether id's house is (tied for) the house with
// the most creatures in play across both battlelines. Ties keep every tied house
// eligible so the chooser may pick a creature of any of them (Etaromme).
func isOfMostPopulousHouse(ctx *EffectContext, id LocalID) bool {
	counts := map[House]int{}
	for player := range 2 {
		for _, cid := range ctx.Resolver.Battleline(player) {
			counts[ctx.Resolver.House(cid)]++
		}
	}
	most := 0
	for _, n := range counts {
		if n > most {
			most = n
		}
	}
	return most > 0 && counts[ctx.Resolver.House(id)] == most
}

// admitted narrows ids to those the target's per-card Filter admits.
func (t Target) admitted(ctx *EffectContext, ids []LocalID) []LocalID {
	if !t.Filter.Narrows() {
		return ids
	}
	out := make([]LocalID, 0, len(ids))
	for _, id := range ids {
		if t.Filter.matches(ctx, id) {
			out = append(out, id)
		}
	}
	return out
}

// onFlank reports whether a creature is on a flank of its battleline (its
// leftmost or rightmost creature), or is considered one for the turn.
func onFlank(ctx *EffectContext, id LocalID) bool {
	if ctx.Resolver.ConsideredFlank(id) {
		return true
	}
	bl := battlelineContaining(ctx, id)
	return len(bl) > 0 &&
		(bl[0] == id || bl[len(bl)-1] == id)
}

// toSideOfSource reports whether id sits on the given side of src within src's
// battleline: dir +1 for the creatures to src's right, -1 for those to its left.
func toSideOfSource(ctx *EffectContext, src, id LocalID, dir int) bool {
	bl := battlelineContaining(ctx, src)
	si, ii := -1, -1
	for j, x := range bl {
		switch x {
		case src:
			si = j
		case id:
			ii = j
		}
	}
	if si < 0 || ii < 0 {
		return false
	}
	if dir > 0 {
		return ii > si
	}
	return ii < si
}

// isNeighbor reports whether id is one of src's battleline neighbors.
func isNeighbor(ctx *EffectContext, src, id LocalID) bool {
	return slices.Contains(neighbors(ctx, src), id)
}

// sharedHouseNeighbors counts how many of id's battleline neighbors share its
// house — the measure behind SharesHouseWithNeighbors (Groupthink Tank).
func sharedHouseNeighbors(ctx *EffectContext, id LocalID) int {
	house := ctx.Resolver.House(id)
	shared := 0
	for _, n := range neighbors(ctx, id) {
		if ctx.Resolver.House(n) == house {
			shared++
		}
	}
	return shared
}

// neighbors returns the creatures immediately adjacent to id in its controller's
// battleline — its left and right neighbors, when present. A card that is not in
// a battleline has no neighbors.
func neighbors(ctx *EffectContext, id LocalID) []LocalID {
	bl := battlelineContaining(ctx, id)
	i := -1
	for j, x := range bl {
		if x == id {
			i = j
			break
		}
	}
	if i < 0 {
		return nil
	}
	out := make([]LocalID, 0, 2)
	if i > 0 {
		out = append(out, bl[i-1])
	}
	if i < len(bl)-1 {
		out = append(out, bl[i+1])
	}
	return out
}

// pickCards has the controller choose cards one at a time from the pool avail
// returns, never repeating a pick, until they decline or the pool runs dry. limit
// caps how many are chosen; limit <= 0 is unbounded ("any number"). optional makes
// each prompt declinable, so the controller can stop early; a mandatory pick still
// stops when they pass or nothing matches. avail is re-read each round, so a pool
// that shifts as cards are picked stays current. It backs every "destroy/purge any
// number of ..." and "keep N ..." effect (Obsidian Forge, Destructive Analysis,
// Unnatural Selection, Tertiate).
func pickCards(
	ctx *EffectContext, prompt string, limit int, optional bool, avail func() []LocalID,
) []LocalID {
	picked := map[LocalID]bool{}
	var chosen []LocalID
	for limit <= 0 || len(chosen) < limit {
		var cands []LocalID
		for _, id := range avail() {
			if !picked[id] {
				cands = append(cands, id)
			}
		}
		if len(cands) == 0 {
			break
		}
		choose := ctx.ChooseCard
		if optional {
			choose = ctx.ChooseCardOptional
		}
		pick, ok := choose(prompt, cands)
		if !ok {
			break
		}
		picked[pick] = true
		chosen = append(chosen, pick)
	}
	return chosen
}

func battlelineContaining(ctx *EffectContext, id LocalID) []LocalID {
	for p := range 2 {
		bl := ctx.Resolver.Battleline(p)
		if slices.Contains(bl, id) {
			return bl
		}
	}
	return nil
}

// allCreatures is every creature in play, the controller's battleline first. The
// battlelines keep their own left-to-right order, which is the order a flank or
// neighbour target reads positions in, so a base set built here can still be
// narrowed by position.
func allCreatures(ctx *EffectContext) []LocalID {
	return append(
		ctx.Resolver.Battleline(ctx.Controller),
		ctx.Resolver.Battleline(ctx.Opponent())...)
}

// allArtifacts is every artifact in play, the controller's row first.
func allArtifacts(ctx *EffectContext) []LocalID {
	return append(
		ctx.Resolver.Artifacts(ctx.Controller),
		ctx.Resolver.Artifacts(ctx.Opponent())...)
}

// creaturesAndArtifacts is every creature in play followed by every artifact in
// play. Cards group by type before they group by player, so a creature of either
// side precedes every artifact. It is the pool for the target kinds that NAME the
// two types ("a creature or artifact"), which is why it stops at the two rows;
// a kind that says "card in play" reads resolverCardsInPlay and reaches upgrades.
func creaturesAndArtifacts(ctx *EffectContext) []LocalID {
	return append(allCreatures(ctx), allArtifacts(ctx)...)
}

// creaturesAndArtifactsOf is one player's creatures followed by their artifacts,
// the one-sided form of creaturesAndArtifacts.
func creaturesAndArtifactsOf(ctx *EffectContext, player int) []LocalID {
	return append(
		ctx.Resolver.Battleline(player),
		ctx.Resolver.Artifacts(player)...)
}

// allCardsInPlay is every card both players have in play — creatures, artifacts,
// and the upgrades on either — the controller's side first.
func allCardsInPlay(ctx *EffectContext) []LocalID {
	return append(
		resolverCardsInPlay(ctx, ctx.Controller),
		resolverCardsInPlay(ctx, ctx.Opponent())...)
}

// selectBase resolves the unfiltered base set chosen by Kind. Chosen kinds return
// the pool of candidates; Select applies filters and prompts for the choice.
func (t Target) selectBase(ctx *EffectContext) []LocalID {
	switch t.Kind {
	case TargetThisCreature:
		return []LocalID{ctx.Source}
	case TargetAttachedHost:
		if host, ok := ctx.Resolver.HostOf(ctx.Upgrade); ok {
			return []LocalID{host}
		}
		return nil
	case TargetGrantingCard:
		if ctx.HasGrantor {
			return []LocalID{ctx.Grantor}
		}
		return nil
	case TargetTriggeringCreature, TargetTheOtherCreature, TargetTheChosenCreature,
		TargetCreatureFought, TargetTheFoughtCreature, TargetTheSameCreature:
		if ctx.HasIt {
			return []LocalID{ctx.It}
		}
		return nil
	case TargetEachArtifact, TargetChosenArtifact:
		return allArtifacts(ctx)
	case TargetChosenEnemyArtifact:
		return ctx.Resolver.Artifacts(ctx.Opponent())
	case TargetChosenFriendlyArtifact:
		return ctx.Resolver.Artifacts(ctx.Controller)
	case TargetChosenUpgrade:
		var ups []LocalID
		for _, c := range allCreatures(ctx) {
			ups = append(ups, ctx.Resolver.Upgrades(c)...)
		}
		return ups
	case TargetEachEnemyArtifact:
		return ctx.Resolver.Artifacts(ctx.Opponent())
	case TargetEachFriendlyArtifact:
		return ctx.Resolver.Artifacts(ctx.Controller)
	case TargetEachCardInPlay:
		return allCardsInPlay(ctx)
	case TargetEachFriendlyCardInPlay:
		return resolverCardsInPlay(ctx, ctx.Controller)
	case TargetChosenCreatureOrArtifact:
		return creaturesAndArtifacts(ctx)
	case TargetChosenFriendlyCreatureOrArtifact:
		return creaturesAndArtifactsOf(ctx, ctx.Controller)
	case TargetChosenEnemyCreatureOrArtifact:
		return creaturesAndArtifactsOf(ctx, ctx.Opponent())
	case TargetEachCreature, TargetChosenCreature:
		return allCreatures(ctx)
	case TargetEachFriendlyCreature, TargetChosenFriendlyCreature:
		return ctx.Resolver.Battleline(ctx.Controller)
	case TargetEachEnemyCreature, TargetChosenEnemyCreature:
		return ctx.Resolver.Battleline(ctx.Opponent())
	case TargetFormerNeighbors:
		return ctx.Produced.Neighbors
	case TargetEachNeighbor:
		return neighbors(ctx, ctx.Source)
	case TargetEachUpgradeOnThis:
		return ctx.Resolver.Upgrades(ctx.Source)
	default:
		return nil
	}
}
