package engine

import (
	"fmt"
	"slices"
	"sort"
	"strings"
)

// This file holds the Refinement strategies a Target can carry: the Refinement
// interface, its optional leadingRefinement capability, and every concrete
// implementation. Each renders its own clause (and optional lead) and refines the
// selected set. See target.go for the Target type and its filters, and
// target_select.go for the selection machinery that applies a refinement.

// A Refinement refines a Target's selected set relative to the whole set — a rule
// that compares the candidates to each other rather than testing each on its own,
// such as "except the most powerful creature". It both narrows the ids (refine)
// and contributes a clause to the target's printed phrase (clause), so these
// niche rules compose onto any Target without a field per rule.
type Refinement interface {
	refine(ctx *EffectContext, ids []LocalID) []LocalID
	clause(phrase string) string
}

// membershipRefiner is the optional capability of a Refinement that can report
// whether a card is among its result without making the discretionary tie-break
// its refine makes — so a condition (ItIsAmong) can ask "is it one of these"
// without prompting. A tie at the cutoff counts every tied card as included.
type membershipRefiner interface {
	Refinement
	includes(ctx *EffectContext, ids []LocalID, id LocalID) bool
}

// leadingRefinement is the optional capability of a Refinement whose choice reads
// before the effect's verb, so the phrase runs left to right — the choice stated
// first, then its consequence (SamePowerAsChosen's "choose a creature - destroy
// …"). A Refinement without it contributes only a trailing clause.
type leadingRefinement interface {
	lead() string
}

// leadIn returns the leading clause a Refinement contributes ahead of the effect's
// verb, and whether the target has one — so an effect can render "choose a
// creature - destroy …" left to right instead of burying the choice.
func (t Target) leadIn() (string, bool) {
	if l, ok := t.Refinement.(leadingRefinement); ok {
		return l.lead(), true
	}
	return "", false
}

// AllExcept is a Refinement that keeps every creature its inner Refinement drops —
// the complement over the selected set. It composes with the singular selectors to
// spare one creature and take the rest: AllExcept{Inner: MostPowerful} is "each
// creature except the most powerful creature" (Champion's Challenge). The inner
// refinement runs first, so any choice it makes (which tied creature is the most
// powerful) decides which creature it spares.
//
// It is named for what it keeps rather than "Except", which the filter vocabulary
// already spends on the one card a target leaves out (Filter.Except) and on the
// house a matcher bars (Houses.Except).
type AllExcept struct {
	// Inner is the refinement whose result is complemented.
	Inner Refinement
}

// clause renders "<phrase> except <inner clause>", e.g. AllExcept{Inner:
// MostPowerful} over "each enemy creature" -> "each enemy creature except the most
// powerful enemy creature".
func (n AllExcept) clause(phrase string) string {
	return phrase + " except " + n.Inner.clause(phrase)
}

// refine keeps the creatures the inner refinement does not, preserving the
// original order. An empty result selects nothing.
func (n AllExcept) refine(ctx *EffectContext, ids []LocalID) []LocalID {
	dropped := n.Inner.refine(ctx, ids)
	drop := make(map[LocalID]bool, len(dropped))
	for _, id := range dropped {
		drop[id] = true
	}
	kept := make([]LocalID, 0, len(ids))
	for _, id := range ids {
		if !drop[id] {
			kept = append(kept, id)
		}
	}
	if len(kept) == 0 {
		return nil
	}
	return kept
}

// AnyOf is a Refinement that keeps every creature any one of its Refinements
// keeps — the union over the selected set. It composes the tier selectors to take
// several extremes at once: AnyOf{Refinements: []Refinement{LowestPower{},
// HighestPower{}}} is "each creature with the lowest or highest power"
// (Standardized Testing). The union preserves the original order and keeps each
// creature once.
type AnyOf struct {
	// Refinements are the members whose results are unioned.
	Refinements []Refinement
}

// clause states a shared frame once when every member is framed alike, joining
// only the varying words with "or" — inside one noun phrase "and" would read as a
// creature meeting both rules. Members that cannot share a frame render in full
// and join with "and", where each is already a complete noun phrase and "or"
// would read as a choice between the two groups.
func (a AnyOf) clause(phrase string) string {
	if folded, ok := foldFramedClauses(a.Refinements, phrase); ok {
		return folded
	}
	clauses := make([]string, len(a.Refinements))
	for i, r := range a.Refinements {
		clauses[i] = r.clause(phrase)
	}
	return strings.Join(clauses, " and ")
}

// framedClause is the optional capability of a Refinement whose clause is a fixed
// frame around one varying word — "each creature with the" + "lowest" + "power" —
// so a combinator can state the frame once and list only the words that differ.
type framedClause interface {
	clauseFrame(phrase string) (head, tail string)
	clauseWord() string
}

// framedClauseText assembles a framed Refinement's own clause, so its frame is
// written once and cannot drift from what a combinator folds it into.
func framedClauseText(f framedClause, phrase string) string {
	head, tail := f.clauseFrame(phrase)
	return head + " " + f.clauseWord() + " " + tail
}

// foldFramedClauses states a shared frame once and lists only the varying words,
// so AnyOf over LowestPower and HighestPower reads "each creature with the lowest
// or highest power" rather than repeating the subject per member. It declines unless
// there are at least two members and every one is framed alike; members framed
// differently have nothing to share, so they render in full.
// Pinned by TestAnyOfFoldsSharedClauseFrame.
func foldFramedClauses(refinements []Refinement, phrase string) (string, bool) {
	if len(refinements) < 2 {
		return "", false
	}
	var head, tail string
	words := make([]string, 0, len(refinements))
	for i, r := range refinements {
		f, ok := r.(framedClause)
		if !ok {
			return "", false
		}
		h, t := f.clauseFrame(phrase)
		if i == 0 {
			head, tail = h, t
		} else if h != head || t != tail {
			return "", false
		}
		words = append(words, f.clauseWord())
	}
	return head + " " + joinOr(words) + " " + tail, true
}

// refine keeps every creature any member keeps, preserving the original order and
// keeping each creature once. An empty result selects nothing.
func (a AnyOf) refine(ctx *EffectContext, ids []LocalID) []LocalID {
	keep := make(map[LocalID]bool)
	for _, r := range a.Refinements {
		for _, id := range r.refine(ctx, ids) {
			keep[id] = true
		}
	}
	kept := make([]LocalID, 0, len(ids))
	for _, id := range ids {
		if keep[id] {
			kept = append(kept, id)
		}
	}
	if len(kept) == 0 {
		return nil
	}
	return kept
}

// SamePowerAsChosen keeps every creature sharing the power of one the controller
// chooses from the set — the chosen creature included — so a Destroy paired with
// it wipes out a whole power bracket (Dance of Doom). A declined choice selects
// nothing. It leads with "choose a creature" so the printed phrase reads left to
// right, the choice before its consequence.
var SamePowerAsChosen Refinement = SamePowerAs{Sides: []CreatureSide{EitherSide}}

// SamePowerAsEitherChosen is the two-pick form of SamePowerAsChosen: the controller
// chooses one friendly and one enemy creature, and it keeps every creature sharing
// the power of either — both chosen creatures included — so a Destroy paired with
// it wipes out both power brackets at once (Quintrino Flux). Both picks are made
// before the set narrows, so the kept set is computed from the board as it stands
// before any destruction. A declined pick contributes no power. It leads with the
// pair of choices so the printed phrase reads left to right.
var SamePowerAsEitherChosen Refinement = SamePowerAs{
	Sides: []CreatureSide{FriendlySide, EnemySide},
}

// SamePowerAs keeps every creature whose power matches one of a set the controller
// chooses. Each pick names the battleline its candidates come from, so one pick
// over both sides is Dance of Doom (SamePowerAsChosen) and a friendly-plus-enemy
// pair is Quintrino Flux (SamePowerAsEitherChosen). Every pick is made before the
// set narrows, so the kept set unions the chosen powers over the board as it
// stands, and a declined pick contributes no power.
type SamePowerAs struct {
	// Sides is the battleline each pick draws its candidates from, one entry per
	// pick, in the order the picks are made.
	Sides []CreatureSide
}

// lead renders the picks ahead of the effect's verb, e.g. "choose a creature" or
// "choose a friendly creature and an enemy creature".
func (s SamePowerAs) lead() string {
	nouns := make([]string, len(s.Sides))
	for i, side := range s.Sides {
		nouns[i] = side.noun()
	}
	return "choose " + strings.Join(nouns, " and ")
}

// clause names the creatures the kept power is measured against — singular for one
// pick, "either of the chosen creatures" for more.
func (s SamePowerAs) clause(phrase string) string {
	if len(s.Sides) == 1 {
		return phrase + " with the same power as the chosen creature"
	}
	return phrase + " with the same power as either of the chosen creatures"
}

// refine makes every pick, then keeps each creature whose power matches any chosen
// creature's. The picks are made first, so the kept set unions their power brackets
// from the pre-narrowing board.
func (s SamePowerAs) refine(ctx *EffectContext, ids []LocalID) []LocalID {
	powers := map[int]bool{}
	for _, side := range s.Sides {
		if chosen, ok := ctx.ChooseCreature(side.prompt(), side.creatures(ctx, ids)); ok {
			powers[ctx.Resolver.Power(chosen)] = true
		}
	}
	if len(powers) == 0 {
		return nil
	}
	out := make([]LocalID, 0, len(ids))
	for _, id := range ids {
		if powers[ctx.Resolver.Power(id)] {
			out = append(out, id)
		}
	}
	return out
}

// CreatureSide names the battleline a SamePowerAs pick draws its candidates from.
type CreatureSide uint8

const (
	// EitherSide draws from both battlelines (Dance of Doom's single choice). It is
	// the zero value: a pick that names no side draws from the whole set.
	EitherSide CreatureSide = iota
	// FriendlySide draws from the controller's battleline (Quintrino Flux).
	FriendlySide
	// EnemySide draws from the opponent's battleline (Quintrino Flux).
	EnemySide
)

// CreatureSides lists every side in declaration order, so the census enumerates
// them rather than keeping a list a new side could fall out of.
func CreatureSides() []CreatureSide {
	return []CreatureSide{EitherSide, FriendlySide, EnemySide}
}

// noun renders the side in a "choose ..." lead, e.g. "a creature", "a friendly
// creature", or "an enemy creature".
func (s CreatureSide) noun() string {
	switch s {
	case FriendlySide:
		return "a friendly creature"
	case EnemySide:
		return "an enemy creature"
	default:
		return "a creature"
	}
}

// prompt is what the chooser sees when picking from this side. It is the lead's
// own noun, so the prompt and the printed phrase cannot drift apart.
func (s CreatureSide) prompt() string { return "Choose " + s.noun() }

// creatures narrows the candidate ids to the side this pick draws from.
func (s CreatureSide) creatures(ctx *EffectContext, ids []LocalID) []LocalID {
	friendly, enemy := creaturesBySide(ctx, ids)
	switch s {
	case FriendlySide:
		return friendly
	case EnemySide:
		return enemy
	default:
		return ids
	}
}

// creaturesBySide splits ids into the controller's creatures and the opponent's,
// preserving each side's original order — the both-battlelines refinements
// (SamePowerAsEitherChosen, KeepPerSide, PortionPerSide) act on each side in turn.
func creaturesBySide(ctx *EffectContext, ids []LocalID) (friendly, enemy []LocalID) {
	for _, id := range ids {
		if ctx.Resolver.Controller(id) == ctx.Controller {
			friendly = append(friendly, id)
		} else {
			enemy = append(enemy, id)
		}
	}
	return friendly, enemy
}

// KeepPerSide is a Refinement that spares a chosen number of creatures on each
// battleline and selects every other creature — Unnatural Selection keeps 3 per
// side and destroys the rest, as Refine on card.Target.EachCreature. It splits the
// set by side, has the controller choose up to N keepers on each, and returns the
// creatures left over. A side no larger than the keep count is spared whole with
// no prompt, since keeping cannot change which creatures survive. It leads with
// the pair of choices so the printed phrase reads left to right.
type KeepPerSide struct {
	// N is how many creatures the controller keeps on each battleline.
	N int
}

// lead renders the pair of keep choices ahead of the effect's verb, e.g. "choose 3
// friendly creatures and 3 enemy creatures".
func (k KeepPerSide) lead() string {
	return fmt.Sprintf("choose %d friendly creatures and %d enemy creatures", k.N, k.N)
}

// clause renders the leftover set, e.g. "each creature" -> "each other creature".
func (k KeepPerSide) clause(phrase string) string {
	return strings.Replace(phrase, "each ", "each other ", 1)
}

// refine keeps up to n creatures on each side, chosen by the controller, and
// returns every other creature. The sides are read from the pre-narrowing board so
// the whole leftover set selects at once — the controller keeps the friendly side
// first, then the enemy side.
func (k KeepPerSide) refine(ctx *EffectContext, ids []LocalID) []LocalID {
	friendly, enemy := creaturesBySide(ctx, ids)
	return append(k.leftover(ctx, friendly), k.leftover(ctx, enemy)...)
}

// leftover has the controller keep up to n of side, returning the rest. A side no
// larger than the keep count is spared whole without a prompt: every creature is
// kept, so the choice would be vacuous.
func (k KeepPerSide) leftover(ctx *EffectContext, side []LocalID) []LocalID {
	if len(side) <= k.N {
		return nil
	}
	kept := map[LocalID]bool{}
	for _, id := range pickCards(ctx, "Choose a creature to keep", k.N, false, func() []LocalID {
		return side
	}) {
		kept[id] = true
	}
	var rest []LocalID
	for _, id := range side {
		if !kept[id] {
			rest = append(rest, id)
		}
	}
	return rest
}

// PortionPerSide is a Refinement that selects a fraction of the creatures on each
// battleline, chosen by the controller — Tertiate takes one third of each side,
// rounding up, as Refine on card.Target.EachCreature. It reads both sides from the
// pre-narrowing board and selects the whole chosen set at once.
type PortionPerSide struct {
	// Share is the fraction of each battleline the controller chooses.
	Share Fraction
}

// clause renders the printed phrase from the fraction, ignoring the base noun,
// e.g. "one third of all enemy creatures and one third of all friendly creatures,
// rounding up each time". The rounding trails both sides as one clause, so the
// rule is stated once rather than repeated per side.
func (p PortionPerSide) clause(string) string {
	return fmt.Sprintf(
		"one %s of all enemy creatures and one %s of all friendly creatures, %s each time",
		p.Share.word(), p.Share.word(), p.Share.roundingPhrase(),
	)
}

// refine selects the fraction of each side, the enemy side first, chosen by the
// controller. Both sides are read from the pre-narrowing board so the whole chosen
// set selects at once.
func (p PortionPerSide) refine(ctx *EffectContext, ids []LocalID) []LocalID {
	friendly, enemy := creaturesBySide(ctx, ids)
	return append(p.share(ctx, enemy), p.share(ctx, friendly)...)
}

// share has the controller choose the fraction of side, returning the chosen ids.
func (p PortionPerSide) share(ctx *EffectContext, side []LocalID) []LocalID {
	return pickCards(
		ctx,
		"Choose a creature to destroy",
		p.Share.of(len(side)),
		false,
		func() []LocalID {
			return side
		},
	)
}

// LeastPowerful is a Refinement that keeps only the single least powerful creature
// of a set, e.g. card.Target.EachCreature.Refine(card.Refine.LeastPowerful)
// (Horseman of Famine). When several tie for least powerful the controller
// chooses which.
type LeastPowerful struct{}

// clause renders "the least powerful <noun>", e.g. "each creature" -> "the least
// powerful creature".
func (LeastPowerful) clause(phrase string) string {
	return "the least powerful " + strings.TrimPrefix(phrase, "each ")
}

// refine returns the single least powerful creature, letting the controller
// choose which to keep when several tie. An empty set selects nothing.
func (LeastPowerful) refine(ctx *EffectContext, ids []LocalID) []LocalID {
	if len(ids) == 0 {
		return nil
	}
	lowest := ctx.Resolver.Power(ids[0])
	for _, id := range ids[1:] {
		if p := ctx.Resolver.Power(id); p < lowest {
			lowest = p
		}
	}
	tied := make([]LocalID, 0, len(ids))
	for _, id := range ids {
		if ctx.Resolver.Power(id) == lowest {
			tied = append(tied, id)
		}
	}
	pick := tied[0]
	if len(tied) > 1 {
		if chosen, ok := ctx.ChooseCreature("Choose the least powerful creature", tied); ok {
			pick = chosen
		}
	}
	return []LocalID{pick}
}

// MostPowerful is a Refinement that keeps the single most powerful creature of a
// set, letting the controller choose which to keep when several tie — the
// highest-power mirror of LeastPowerful, e.g.
// card.Target.EachCreature.Refine(card.Refine.MostPowerful). It is the one-creature
// form of MostPowerfulN, named because most cards want exactly it.
var MostPowerful Refinement = MostPowerfulN{N: 1}

// MostPowerfulN is a Refinement that keeps the N most powerful creatures of a set,
// e.g. card.Target.EachCreature.Refine(card.MostPowerfulN{N: 3}) (Three Fates).
// When more creatures tie at the cutoff than there are remaining slots, the
// controller chooses which of the tied creatures to include.
type MostPowerfulN struct {
	// N is how many of the most powerful creatures the refinement keeps.
	N int
}

// clause renders "the N most powerful <noun>s", e.g. "the 3 most powerful
// creatures", and the singular "the most powerful <noun>" when N is 1.
func (m MostPowerfulN) clause(phrase string) string {
	noun := strings.TrimPrefix(phrase, "each ")
	if m.N == 1 {
		return "the most powerful " + noun
	}
	return fmt.Sprintf("the %d most powerful %ss", m.N, noun)
}

// refine keeps the N highest-power creatures, letting the controller break ties
// at the cutoff. A set no larger than N keeps all of it.
func (m MostPowerfulN) refine(ctx *EffectContext, ids []LocalID) []LocalID {
	if len(ids) <= m.N {
		return ids
	}
	sorted := make([]LocalID, len(ids))
	copy(sorted, ids)
	sort.SliceStable(sorted, func(i, j int) bool {
		return ctx.Resolver.Power(sorted[i]) > ctx.Resolver.Power(sorted[j])
	})
	threshold := ctx.Resolver.Power(sorted[m.N-1])
	chosen := make([]LocalID, 0, m.N)
	tied := make([]LocalID, 0)
	for _, id := range sorted {
		switch p := ctx.Resolver.Power(id); {
		case p > threshold:
			chosen = append(chosen, id)
		case p == threshold:
			tied = append(tied, id)
		}
	}
	for slots := m.N - len(chosen); slots > 0 && len(tied) > 0; slots-- {
		if len(tied) == slots {
			chosen = append(chosen, tied...)
			break
		}
		pick := tied[0]
		if c, ok := ctx.ChooseCreature("Choose one of the most powerful creatures", tied); ok {
			pick = c
		}
		chosen = append(chosen, pick)
		for i, id := range tied {
			if id == pick {
				tied = append(tied[:i], tied[i+1:]...)
				break
			}
		}
	}
	return chosen
}

// includes reports whether id ties for or exceeds the N-th highest power in the
// set — the tie-inclusive membership a condition wants, without the tie-break
// prompt refine makes. A set no larger than N includes every member.
func (m MostPowerfulN) includes(ctx *EffectContext, ids []LocalID, id LocalID) bool {
	if !slices.Contains(ids, id) {
		return false
	}
	if len(ids) <= m.N {
		return true
	}
	sorted := make([]LocalID, len(ids))
	copy(sorted, ids)
	sort.SliceStable(sorted, func(i, j int) bool {
		return ctx.Resolver.Power(sorted[i]) > ctx.Resolver.Power(sorted[j])
	})
	threshold := ctx.Resolver.Power(sorted[m.N-1])
	return ctx.Resolver.Power(id) >= threshold
}

// HighestPower is a Refinement that keeps every creature tied for the highest
// power of a set — a tier, not a single creature, so it makes no choice (contrast
// MostPowerful, which keeps exactly one). An empty set selects nothing.
type HighestPower struct{}

// clause renders "<phrase> with the highest power".
func (r HighestPower) clause(phrase string) string { return framedClauseText(r, phrase) }

// clauseFrame frames the power tier so AnyOf can fold it with another tier.
func (HighestPower) clauseFrame(phrase string) (head, tail string) {
	return phrase + " with the", "power"
}

// clauseWord names the tier this refinement keeps.
func (HighestPower) clauseWord() string { return "highest" }

// refine keeps every creature whose power equals the set's maximum, so a set with
// a single power keeps all of it. An empty set selects nothing.
func (HighestPower) refine(ctx *EffectContext, ids []LocalID) []LocalID {
	if len(ids) == 0 {
		return nil
	}
	high := ctx.Resolver.Power(ids[0])
	for _, id := range ids[1:] {
		if p := ctx.Resolver.Power(id); p > high {
			high = p
		}
	}
	kept := make([]LocalID, 0, len(ids))
	for _, id := range ids {
		if ctx.Resolver.Power(id) == high {
			kept = append(kept, id)
		}
	}
	return kept
}

// LowestPower is a Refinement that keeps every creature tied for the lowest power
// of a set — a tier, not a single creature, so it makes no choice (contrast
// LeastPowerful, which keeps exactly one). An empty set selects nothing.
type LowestPower struct{}

// clause renders "<phrase> with the lowest power".
func (r LowestPower) clause(phrase string) string { return framedClauseText(r, phrase) }

// clauseFrame frames the power tier so AnyOf can fold it with another tier.
func (LowestPower) clauseFrame(phrase string) (head, tail string) {
	return phrase + " with the", "power"
}

// clauseWord names the tier this refinement keeps.
func (LowestPower) clauseWord() string { return "lowest" }

// refine keeps every creature whose power equals the set's minimum, so a set with
// a single power keeps all of it. An empty set selects nothing.
func (LowestPower) refine(ctx *EffectContext, ids []LocalID) []LocalID {
	if len(ids) == 0 {
		return nil
	}
	low := ctx.Resolver.Power(ids[0])
	for _, id := range ids[1:] {
		if p := ctx.Resolver.Power(id); p < low {
			low = p
		}
	}
	kept := make([]LocalID, 0, len(ids))
	for _, id := range ids {
		if ctx.Resolver.Power(id) == low {
			kept = append(kept, id)
		}
	}
	return kept
}

// PowerLessThan is a Refinement that keeps every creature of a set whose power is
// below a running count — Exterminate! Exterminate! destroys each non-Mars
// creature with power less than the number of friendly Mars creatures you
// control. The threshold is read when the effect resolves, so it reflects the
// board at that moment.
//
// A threshold that is a printed number, or none at all, is a comparable Filter
// axis instead (PowerBound); this refinement is for the live Count.
type PowerLessThan struct {
	// Limit is the running count a creature's power must stay below.
	Limit Count
}

// clause renders "<phrase> with power less than the number of …".
func (p PowerLessThan) clause(phrase string) string {
	return phrase + " with power less than " + cardinalCountText(p.Limit)
}

// refine keeps the creatures whose power is below the count's value.
func (p PowerLessThan) refine(ctx *EffectContext, ids []LocalID) []LocalID {
	limit := p.Limit.Value(ctx)
	kept := make([]LocalID, 0, len(ids))
	for _, id := range ids {
		if ctx.Resolver.Power(id) < limit {
			kept = append(kept, id)
		}
	}
	return kept
}
