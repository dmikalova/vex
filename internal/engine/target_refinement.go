package engine

import (
	"fmt"
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
	if l, ok := t.refinement.(leadingRefinement); ok {
		return l.lead(), true
	}
	return "", false
}

// Except is a Refinement that keeps every creature its inner Refinement drops —
// the complement over the selected set. It composes with the singular selectors to
// spare one creature and take the rest: Except(MostPowerful) is "each creature
// except the most powerful creature" (Champion's Challenge). The inner refinement
// runs first, so any choice it makes (which tied creature is the most powerful)
// decides which creature Except spares.
func Except(inner Refinement) Refinement { return except{inner: inner} }

// except implements the Except combinator.
type except struct{ inner Refinement }

// clause renders "<phrase> except <inner clause>", e.g. Except(MostPowerful) over
// "each enemy creature" -> "each enemy creature except the most powerful enemy
// creature".
func (n except) clause(phrase string) string {
	return phrase + " except " + n.inner.clause(phrase)
}

// refine keeps the creatures the inner refinement does not, preserving the
// original order. An empty result selects nothing.
func (n except) refine(ctx *EffectContext, ids []LocalID) []LocalID {
	dropped := n.inner.refine(ctx, ids)
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
// several extremes at once: AnyOf(LowestPower, HighestPower) is "each creature
// with the lowest power and each creature with the highest power" (Standardized
// Testing). The union preserves the original order and keeps each creature once.
func AnyOf(refinements ...Refinement) Refinement {
	return anyOf{refinements: refinements}
}

// anyOf implements the AnyOf combinator.
type anyOf struct{ refinements []Refinement }

// clause joins each member's clause with " and ", e.g. "each creature with the
// lowest power and each creature with the highest power".
func (a anyOf) clause(phrase string) string {
	clauses := make([]string, len(a.refinements))
	for i, r := range a.refinements {
		clauses[i] = r.clause(phrase)
	}
	return strings.Join(clauses, " and ")
}

// refine keeps every creature any member keeps, preserving the original order and
// keeping each creature once. An empty result selects nothing.
func (a anyOf) refine(ctx *EffectContext, ids []LocalID) []LocalID {
	keep := make(map[LocalID]bool)
	for _, r := range a.refinements {
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

// SamePowerAsChosen is a Refinement that keeps every creature sharing the power of
// one the controller chooses from the set — the chosen creature included — so a
// Destroy paired with it wipes out a whole power bracket (Dance of Doom). A
// declined choice selects nothing. It leads with "choose a creature" so the
// printed phrase reads left to right, the choice before its consequence.
var SamePowerAsChosen Refinement = samePowerAsChosen{}

// samePowerAsChosen implements the SamePowerAsChosen refinement.
type samePowerAsChosen struct{}

// lead renders the choice ahead of the effect's verb.
func (samePowerAsChosen) lead() string { return "choose a creature" }

// clause renders "<phrase> with the same power as the chosen creature".
func (samePowerAsChosen) clause(phrase string) string {
	return phrase + " with the same power as the chosen creature"
}

// refine picks a creature and keeps every one in the set with matching power.
func (samePowerAsChosen) refine(ctx *EffectContext, ids []LocalID) []LocalID {
	chosen, ok := ctx.ChooseCreature("Choose a creature", ids)
	if !ok {
		return nil
	}
	power := ctx.Resolver.Power(chosen)
	out := make([]LocalID, 0, len(ids))
	for _, id := range ids {
		if ctx.Resolver.Power(id) == power {
			out = append(out, id)
		}
	}
	return out
}

// SamePowerAsEitherChosen is the two-choice sibling of SamePowerAsChosen: the
// controller chooses one friendly and one enemy creature from the set, and it
// keeps every creature sharing the power of either — both chosen creatures
// included — so a Destroy paired with it wipes out both power brackets at once
// (Quintrino Flux). Both creatures are chosen before the set is narrowed, so the
// kept set is computed from the board as it stands before any destruction. A
// declined choice contributes no power. It leads with the pair of choices so the
// printed phrase reads left to right, the choices before their consequence.
var SamePowerAsEitherChosen Refinement = samePowerAsEitherChosen{}

// samePowerAsEitherChosen implements the SamePowerAsEitherChosen refinement.
type samePowerAsEitherChosen struct{}

// lead renders both choices ahead of the effect's verb.
func (samePowerAsEitherChosen) lead() string {
	return "choose a friendly creature and an enemy creature"
}

// clause renders "<phrase> with the same power as either of the chosen creatures".
func (samePowerAsEitherChosen) clause(phrase string) string {
	return phrase + " with the same power as either of the chosen creatures"
}

// refine chooses one friendly and one enemy creature, then keeps every creature
// in the set matching either chosen power. Both choices are made first, so the
// kept set unions the two power brackets from the pre-narrowing board.
func (samePowerAsEitherChosen) refine(ctx *EffectContext, ids []LocalID) []LocalID {
	friendly, enemy := creaturesBySide(ctx, ids)
	powers := map[int]bool{}
	if fc, ok := ctx.ChooseCreature("Choose a friendly creature", friendly); ok {
		powers[ctx.Resolver.Power(fc)] = true
	}
	if ec, ok := ctx.ChooseCreature("Choose an enemy creature", enemy); ok {
		powers[ctx.Resolver.Power(ec)] = true
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

// KeepPerSide returns a Refinement that spares a chosen number of creatures on
// each battleline and selects every other creature — Unnatural Selection keeps 3
// per side and destroys the rest, as Refine on card.Target.EachCreature. It splits
// the set by side, has the controller choose up to n keepers on each, and returns
// the creatures left over. A side no larger than the keep count is spared whole
// with no prompt, since keeping cannot change which creatures survive. It leads
// with the pair of choices so the printed phrase reads left to right.
func KeepPerSide(n int) Refinement { return keepPerSide{n: n} }

// keepPerSide implements the KeepPerSide refinement.
type keepPerSide struct{ n int }

// lead renders the pair of keep choices ahead of the effect's verb, e.g. "choose 3
// friendly creatures and 3 enemy creatures".
func (k keepPerSide) lead() string {
	return fmt.Sprintf("choose %d friendly creatures and %d enemy creatures", k.n, k.n)
}

// clause renders the leftover set, e.g. "each creature" -> "each other creature".
func (k keepPerSide) clause(phrase string) string {
	return strings.Replace(phrase, "each ", "each other ", 1)
}

// refine keeps up to n creatures on each side, chosen by the controller, and
// returns every other creature. The sides are read from the pre-narrowing board so
// the whole leftover set selects at once — the controller keeps the friendly side
// first, then the enemy side.
func (k keepPerSide) refine(ctx *EffectContext, ids []LocalID) []LocalID {
	friendly, enemy := creaturesBySide(ctx, ids)
	return append(k.leftover(ctx, friendly), k.leftover(ctx, enemy)...)
}

// leftover has the controller keep up to n of side, returning the rest. A side no
// larger than the keep count is spared whole without a prompt: every creature is
// kept, so the choice would be vacuous.
func (k keepPerSide) leftover(ctx *EffectContext, side []LocalID) []LocalID {
	if len(side) <= k.n {
		return nil
	}
	kept := map[LocalID]bool{}
	for _, id := range pickCards(ctx, "Choose a creature to keep", k.n, false, func() []LocalID {
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

// PortionPerSide returns a Refinement that selects a fraction of the creatures on
// each battleline, chosen by the controller — Tertiate takes one third of each
// side, rounding up, as Refine on card.Target.EachCreature. It reads both sides
// from the pre-narrowing board and selects the whole chosen set at once.
func PortionPerSide(f Fraction) Refinement { return portionPerSide{f: f} }

// portionPerSide implements the PortionPerSide refinement.
type portionPerSide struct{ f Fraction }

// clause renders the printed phrase from the fraction, ignoring the base noun,
// e.g. "one third of all enemy creatures and one third of all friendly creatures
// (rounding up each time)".
func (p portionPerSide) clause(string) string {
	return fmt.Sprintf(
		"one %s of all enemy creatures and one %s of all friendly creatures (%s each time)",
		p.f.word(), p.f.word(), p.f.roundingPhrase(),
	)
}

// refine selects the fraction of each side, the enemy side first, chosen by the
// controller. Both sides are read from the pre-narrowing board so the whole chosen
// set selects at once.
func (p portionPerSide) refine(ctx *EffectContext, ids []LocalID) []LocalID {
	friendly, enemy := creaturesBySide(ctx, ids)
	return append(p.share(ctx, enemy), p.share(ctx, friendly)...)
}

// share has the controller choose the fraction of side, returning the chosen ids.
func (p portionPerSide) share(ctx *EffectContext, side []LocalID) []LocalID {
	return pickCards(
		ctx,
		"Choose a creature to destroy",
		p.f.of(len(side)),
		false,
		func() []LocalID {
			return side
		},
	)
}

// LeastPowerful is a Refinement that keeps only the single least powerful creature
// of a set, e.g. card.Target.EachCreature.Refine(card.LeastPowerful) (Horseman
// of Famine). When several tie for least powerful the controller chooses which.
var LeastPowerful Refinement = leastPowerful{}

// leastPowerful implements the LeastPowerful refinement.
type leastPowerful struct{}

// clause renders "the least powerful <noun>", e.g. "each creature" -> "the least
// powerful creature".
func (leastPowerful) clause(phrase string) string {
	return "the least powerful " + strings.TrimPrefix(phrase, "each ")
}

// refine returns the single least powerful creature, letting the controller
// choose which to keep when several tie. An empty set selects nothing.
func (leastPowerful) refine(ctx *EffectContext, ids []LocalID) []LocalID {
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
// card.Target.EachCreature.Refine(card.MostPowerful). For the "keep the top N"
// form use MostPowerfulN.
var MostPowerful Refinement = mostPowerfulN{n: 1}

// MostPowerfulN returns a Refinement that keeps the n most powerful creatures of a
// set, e.g. card.Target.EachCreature.Refine(card.MostPowerfulN(3)) (Three Fates).
// When more creatures tie at the cutoff than there are remaining slots, the
// controller chooses which of the tied creatures to include.
func MostPowerfulN(n int) Refinement { return mostPowerfulN{n: n} }

// mostPowerfulN implements the MostPowerful refinement.
type mostPowerfulN struct{ n int }

// clause renders "the N most powerful <noun>s", e.g. "the 3 most powerful
// creatures", and the singular "the most powerful <noun>" when n is 1.
func (m mostPowerfulN) clause(phrase string) string {
	noun := strings.TrimPrefix(phrase, "each ")
	if m.n == 1 {
		return "the most powerful " + noun
	}
	return fmt.Sprintf("the %d most powerful %ss", m.n, noun)
}

// refine keeps the n highest-power creatures, letting the controller break ties
// at the cutoff. A set no larger than n keeps all of it.
func (m mostPowerfulN) refine(ctx *EffectContext, ids []LocalID) []LocalID {
	if len(ids) <= m.n {
		return ids
	}
	sorted := append([]LocalID(nil), ids...)
	sort.SliceStable(sorted, func(i, j int) bool {
		return ctx.Resolver.Power(sorted[i]) > ctx.Resolver.Power(sorted[j])
	})
	threshold := ctx.Resolver.Power(sorted[m.n-1])
	chosen := make([]LocalID, 0, m.n)
	tied := make([]LocalID, 0)
	for _, id := range sorted {
		switch p := ctx.Resolver.Power(id); {
		case p > threshold:
			chosen = append(chosen, id)
		case p == threshold:
			tied = append(tied, id)
		}
	}
	for slots := m.n - len(chosen); slots > 0 && len(tied) > 0; slots-- {
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

// HouseWithAtLeast returns a Refinement that keeps only creatures whose house has
// at least n creatures in play, counting that house across both players'
// battlelines — a house is a house regardless of who controls its creatures. No
// Safety in Numbers deals its damage to each creature that belongs to a house
// with 3 or more creatures in play.
func HouseWithAtLeast(n int) Refinement { return houseWithAtLeast{n: n} }

// houseWithAtLeast implements the HouseWithAtLeast refinement.
type houseWithAtLeast struct{ n int }

// clause renders "<phrase> that belongs to a house that has N or more creatures
// in play", e.g. "each creature that belongs to a house that has 3 or more
// creatures in play".
func (h houseWithAtLeast) clause(phrase string) string {
	return fmt.Sprintf(
		"%s that belongs to a house that has %d or more creatures in play",
		phrase, h.n,
	)
}

// refine keeps only creatures whose house has at least n creatures in play,
// counting each house across both battlelines so both players' creatures of a
// house count toward its total.
func (h houseWithAtLeast) refine(ctx *EffectContext, ids []LocalID) []LocalID {
	counts := map[House]int{}
	for player := 0; player < 2; player++ {
		for _, cid := range ctx.Resolver.Battleline(player) {
			counts[ctx.Resolver.House(cid)]++
		}
	}
	kept := make([]LocalID, 0, len(ids))
	for _, id := range ids {
		if counts[ctx.Resolver.House(id)] >= h.n {
			kept = append(kept, id)
		}
	}
	return kept
}

// WithoutSharedTrait returns a Refinement that keeps only creatures that share no
// trait with any other creature in the same controller's battleline — a "loner".
// A creature with no traits shares no trait, so it is kept. Good of the Many
// destroys each creature that does not share a trait with another creature in its
// controller's battleline.
func WithoutSharedTrait() Refinement { return withoutSharedTrait{} }

// withoutSharedTrait implements the WithoutSharedTrait refinement.
type withoutSharedTrait struct{}

// clause renders "<phrase> that does not share a trait with another creature in
// its controller's battleline".
func (withoutSharedTrait) clause(phrase string) string {
	return phrase + " that does not share a trait with another creature in its controller's battleline"
}

// refine keeps a creature only when no other creature its controller controls
// shares a trait with it, comparing each candidate against its own battleline
// mates alone (an enemy sharing a trait does not save it).
func (withoutSharedTrait) refine(ctx *EffectContext, ids []LocalID) []LocalID {
	kept := make([]LocalID, 0, len(ids))
	for _, id := range ids {
		shared := false
		for _, mate := range ctx.Resolver.Battleline(ctx.Resolver.Controller(id)) {
			if mate != id && ctx.Resolver.SharesTrait(id, mate) {
				shared = true
				break
			}
		}
		if !shared {
			kept = append(kept, id)
		}
	}
	return kept
}

// HighestPower is a Refinement that keeps every creature tied for the highest
// power of a set — a tier, not a single creature, so it makes no choice (contrast
// MostPowerful, which keeps exactly one). An empty set selects nothing.
var HighestPower Refinement = highestPower{}

// highestPower implements the HighestPower refinement.
type highestPower struct{}

// clause renders "<phrase> with the highest power".
func (highestPower) clause(phrase string) string {
	return phrase + " with the highest power"
}

// refine keeps every creature whose power equals the set's maximum, so a set with
// a single power keeps all of it. An empty set selects nothing.
func (highestPower) refine(ctx *EffectContext, ids []LocalID) []LocalID {
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
var LowestPower Refinement = lowestPower{}

// lowestPower implements the LowestPower refinement.
type lowestPower struct{}

// clause renders "<phrase> with the lowest power".
func (lowestPower) clause(phrase string) string {
	return phrase + " with the lowest power"
}

// refine keeps every creature whose power equals the set's minimum, so a set with
// a single power keeps all of it. An empty set selects nothing.
func (lowestPower) refine(ctx *EffectContext, ids []LocalID) []LocalID {
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
func PowerLessThan(limit Count) Refinement { return powerLessThan{limit: limit} }

// powerLessThan implements the PowerLessThan refinement.
type powerLessThan struct{ limit Count }

// houseReplaced resolves the SelfHouse sentinel inside the count (or rehouses it
// for a Maverick), which lives in an unexported field reflection cannot reach on
// its own.
func (p powerLessThan) houseReplaced(from, to House) any {
	p.limit = replacedIn(p.limit, from, to)
	return p
}

// clause renders "<phrase> with power less than the number of …".
func (p powerLessThan) clause(phrase string) string {
	return phrase + " with power less than " + cardinalCountText(p.limit)
}

// refine keeps the creatures whose power is below the count's value.
func (p powerLessThan) refine(ctx *EffectContext, ids []LocalID) []LocalID {
	limit := p.limit.Value(ctx)
	kept := make([]LocalID, 0, len(ids))
	for _, id := range ids {
		if ctx.Resolver.Power(id) < limit {
			kept = append(kept, id)
		}
	}
	return kept
}
