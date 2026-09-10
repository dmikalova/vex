package engine

// This file holds a Target's selection machinery: resolving a Target into
// concrete card ids (Select, SelectOptional, selectWith), narrowing them by the
// Target's filters (filter), the base sets each Kind draws from (selectBase), and
// the battleline geometry the flank and neighbor filters read (onFlank,
// isNeighbor, neighbors, …). See target.go for the Target type and its filter
// builders, and target_refinement.go for the Refinement strategies selectWith applies.

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

// empty reports that nothing matches this target, reading only the candidates a
// refinement would narrow: with nothing to narrow there is nothing to select, and
// unlike Select it asks the controller nothing.
func (t Target) empty(ctx *EffectContext) bool {
	return len(t.filter(ctx, t.selectBase(ctx))) == 0
}

// selectWith is the shared selection path; optional switches the chosen-kind
// prompt between a forced pick and a declinable one, and keep (when set) drops
// candidates the calling effect could not act on.
func (t Target) selectWith(
	ctx *EffectContext,
	optional bool,
	keep func(LocalID) bool,
) []LocalID {
	ids := t.filter(ctx, t.selectBase(ctx))
	if keep != nil {
		kept := ids[:0:0]
		for _, id := range ids {
			if keep(id) {
				kept = append(kept, id)
			}
		}
		ids = kept
	}
	if t.refinement != nil {
		ids = t.refinement.refine(ctx, ids)
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
	if !t.withNeighbors && !t.neighborsOf {
		return ids
	}
	out := ids[:0:0]
	for _, id := range ids {
		if t.withNeighbors {
			out = append(out, id)
		}
		out = append(out, neighbors(ctx, id)...)
	}
	return out
}

// isChosen reports whether the Kind resolves to a single player-chosen creature.
func (t Target) isChosen() bool {
	return t.Kind == TargetChosenCreature || t.Kind == TargetChosenEnemyCreature ||
		t.Kind == TargetChosenFriendlyCreature || t.Kind == TargetChosenOtherFriendlyCreature ||
		t.Kind == TargetChosenOtherCreature ||
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
	for player := 0; player < 2; player++ {
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

// filter narrows ids to those matching the target's trait, power, damaged, and
// flank filters.
func (t Target) filter(ctx *EffectContext, ids []LocalID) []LocalID {
	if t.trait == traitUnset &&
		t.exceptTrait == traitUnset &&
		t.house == HouseNone &&
		t.exceptHouse == HouseNone &&
		!t.chosenHouse &&
		!t.activeHouse &&
		!t.contextualHouse &&
		!t.houseWithMostCreatures &&
		!t.sharesTrait &&
		!t.hasMaxPower &&
		!t.hasMinPower &&
		!t.hasExactPower &&
		!t.oddPower &&
		!t.evenPower &&
		!t.damaged &&
		!t.undamaged &&
		!t.stunned &&
		!t.ready &&
		!t.withAember &&
		!t.withoutAember &&
		!t.withCounter.valid() &&
		!t.withArmor &&
		!t.withUpgrade &&
		t.sharesHouseNeighbors == 0 &&
		t.keyword == keywordUnset &&
		!t.onFlank &&
		!t.notOnFlank &&
		!t.neighboring &&
		!t.toRightOfSource &&
		!t.toLeftOfSource &&
		!t.other &&
		t.named == "" {
		return ids
	}
	out := make([]LocalID, 0, len(ids))
	for _, id := range ids {
		if t.trait != traitUnset && !ctx.Resolver.HasTrait(id, t.trait) {
			continue
		}
		if t.exceptTrait != traitUnset && ctx.Resolver.HasTrait(id, t.exceptTrait) {
			continue
		}
		if t.house != HouseNone && ctx.Resolver.House(id) != t.house {
			continue
		}
		if t.exceptHouse != HouseNone && ctx.Resolver.House(id) == t.exceptHouse {
			continue
		}
		if t.chosenHouse && ctx.Resolver.House(id) != ctx.ChosenHouse {
			continue
		}
		if t.activeHouse && ctx.Resolver.House(id) != ctx.Resolver.ActiveHouse() {
			continue
		}
		if t.contextualHouse &&
			(!ctx.HasIt || ctx.Resolver.House(id) != ctx.Resolver.House(ctx.It)) {
			continue
		}
		if t.houseWithMostCreatures && !isOfMostPopulousHouse(ctx, id) {
			continue
		}
		if t.sharesTrait && (!ctx.HasIt || !ctx.Resolver.SharesTrait(ctx.It, id)) {
			continue
		}
		if t.hasMaxPower && ctx.Resolver.Power(id) > t.maxPower {
			continue
		}
		if t.hasMinPower && ctx.Resolver.Power(id) < t.minPower {
			continue
		}
		if t.hasExactPower && ctx.Resolver.Power(id) != t.exactPower {
			continue
		}
		if t.oddPower && ctx.Resolver.Power(id)%2 == 0 {
			continue
		}
		if t.evenPower && ctx.Resolver.Power(id)%2 != 0 {
			continue
		}
		if t.damaged && ctx.Resolver.Damage(id) == 0 {
			continue
		}
		if t.undamaged && ctx.Resolver.Damage(id) != 0 {
			continue
		}
		if t.withAember && ctx.Resolver.AmberOn(id) == 0 {
			continue
		}
		if t.withoutAember && ctx.Resolver.AmberOn(id) != 0 {
			continue
		}
		if t.withCounter.valid() && ctx.Resolver.CountersOn(id, t.withCounter) == 0 {
			continue
		}
		if t.withArmor && ctx.Resolver.Armor(id) == 0 {
			continue
		}
		if t.withUpgrade && len(ctx.Resolver.Upgrades(id)) == 0 {
			continue
		}
		if t.sharesHouseNeighbors > 0 &&
			sharedHouseNeighbors(ctx, id) < t.sharesHouseNeighbors {
			continue
		}
		if t.keyword.valid() && !ctx.Resolver.HasKeyword(id, t.keyword) {
			continue
		}
		if t.stunned && !ctx.Resolver.Stunned(id) {
			continue
		}
		if t.ready && ctx.Resolver.Exhausted(id) {
			continue
		}
		if t.onFlank && ctx.Resolver.IsCreature(id) && !onFlank(ctx, id) {
			continue
		}
		if t.notOnFlank && onFlank(ctx, id) {
			continue
		}
		if t.neighboring && !isNeighbor(ctx, ctx.Source, id) {
			continue
		}
		if t.toRightOfSource && !toSideOfSource(ctx, ctx.Source, id, +1) {
			continue
		}
		if t.toLeftOfSource && !toSideOfSource(ctx, ctx.Source, id, -1) {
			continue
		}
		if t.other && id == ctx.Source {
			continue
		}
		if t.named != "" && ctx.Resolver.Name(id) != t.named {
			continue
		}
		out = append(out, id)
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
	for _, n := range neighbors(ctx, src) {
		if n == id {
			return true
		}
	}
	return false
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

// creaturesExcept returns every creature in play except one, walking both
// players' battlelines in order. It backs effects that target "a different
// creature" or "another creature" than one already chosen.
func creaturesExcept(ctx *EffectContext, exclude LocalID) []LocalID {
	var out []LocalID
	for p := 0; p < 2; p++ {
		for _, id := range ctx.Resolver.Battleline(p) {
			if id != exclude {
				out = append(out, id)
			}
		}
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
	for p := 0; p < 2; p++ {
		bl := ctx.Resolver.Battleline(p)
		for _, x := range bl {
			if x == id {
				return bl
			}
		}
	}
	return nil
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
	case TargetGrantingArtifact:
		if ctx.HasGrantor {
			return []LocalID{ctx.Grantor}
		}
		return nil
	case TargetTriggeringCreature, TargetTheOtherCreature, TargetTheChosenCreature,
		TargetCreatureFought, TargetTheFoughtCreature:
		if ctx.HasIt {
			return []LocalID{ctx.It}
		}
		return nil
	case TargetEachArtifact, TargetChosenArtifact:
		return append(
			ctx.Resolver.Artifacts(ctx.Controller),
			ctx.Resolver.Artifacts(ctx.Opponent())...)
	case TargetChosenEnemyArtifact:
		return ctx.Resolver.Artifacts(ctx.Opponent())
	case TargetChosenFriendlyArtifact:
		return ctx.Resolver.Artifacts(ctx.Controller)
	case TargetChosenUpgrade:
		var ups []LocalID
		for _, c := range append(
			ctx.Resolver.Battleline(ctx.Controller),
			ctx.Resolver.Battleline(ctx.Opponent())...) {
			ups = append(ups, ctx.Resolver.Upgrades(c)...)
		}
		return ups
	case TargetEachEnemyArtifact:
		return ctx.Resolver.Artifacts(ctx.Opponent())
	case TargetEachFriendlyArtifact:
		return ctx.Resolver.Artifacts(ctx.Controller)
	case TargetEachCardInPlay, TargetChosenCreatureOrArtifact:
		ids := ctx.Resolver.Battleline(ctx.Controller)
		ids = append(ids, ctx.Resolver.Battleline(ctx.Opponent())...)
		ids = append(ids, ctx.Resolver.Artifacts(ctx.Controller)...)
		ids = append(ids, ctx.Resolver.Artifacts(ctx.Opponent())...)
		return ids
	case TargetEachFriendlyCardInPlay, TargetChosenFriendlyCreatureOrArtifact:
		return append(
			ctx.Resolver.Battleline(ctx.Controller),
			ctx.Resolver.Artifacts(ctx.Controller)...)
	case TargetChosenEnemyCreatureOrArtifact:
		return append(
			ctx.Resolver.Battleline(ctx.Opponent()),
			ctx.Resolver.Artifacts(ctx.Opponent())...)
	case TargetEachCreature, TargetChosenCreature:
		return append(
			ctx.Resolver.Battleline(ctx.Controller),
			ctx.Resolver.Battleline(ctx.Opponent())...)
	case TargetEachFriendlyCreature, TargetChosenFriendlyCreature:
		return ctx.Resolver.Battleline(ctx.Controller)
	case TargetEachEnemyCreature, TargetChosenEnemyCreature:
		return ctx.Resolver.Battleline(ctx.Opponent())
	case TargetEachOtherFriendlyCreature, TargetChosenOtherFriendlyCreature:
		out := make([]LocalID, 0)
		for _, id := range ctx.Resolver.Battleline(ctx.Controller) {
			if id != ctx.Source {
				out = append(out, id)
			}
		}
		return out
	case TargetChosenOtherCreature:
		if !ctx.HasIt {
			return append(
				ctx.Resolver.Battleline(ctx.Controller),
				ctx.Resolver.Battleline(ctx.Opponent())...)
		}
		return creaturesExcept(ctx, ctx.It)
	case TargetFormerNeighbors:
		return ctx.Produced.Neighbors
	default:
		return nil
	}
}
