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

// ExceptMostPowerful is a Refinement that drops the single most powerful creature
// from the set, letting the controller choose which one to keep when several tie
// for most powerful. A set of one or none keeps its (only, most powerful) member,
// so nothing is selected.
var ExceptMostPowerful Refinement = exceptMostPowerful{}

// exceptMostPowerful implements the ExceptMostPowerful refinement.
type exceptMostPowerful struct{}

// clause renders "<phrase> except the most powerful <noun>", e.g. "each enemy
// creature except the most powerful enemy creature".
func (exceptMostPowerful) clause(phrase string) string {
	return phrase + " except the most powerful " + strings.TrimPrefix(phrase, "each ")
}

// refine returns ids without the single most powerful creature, letting the
// controller choose which to keep when several tie. A set of one or none keeps
// its (most powerful) member, so nothing is selected.
func (exceptMostPowerful) refine(ctx *EffectContext, ids []LocalID) []LocalID {
	if len(ids) <= 1 {
		return nil
	}
	highest := ctx.Resolver.Power(ids[0])
	for _, id := range ids[1:] {
		if p := ctx.Resolver.Power(id); p > highest {
			highest = p
		}
	}
	mostPowerful := make([]LocalID, 0, len(ids))
	for _, id := range ids {
		if ctx.Resolver.Power(id) == highest {
			mostPowerful = append(mostPowerful, id)
		}
	}
	spared := mostPowerful[0]
	if len(mostPowerful) > 1 {
		if chosen, ok := ctx.ChooseCreature(
			"Choose the most powerful creature to keep",
			mostPowerful,
		); ok {
			spared = chosen
		}
	}
	out := make([]LocalID, 0, len(ids)-1)
	for _, id := range ids {
		if id != spared {
			out = append(out, id)
		}
	}
	return out
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
	var friendly, enemy []LocalID
	for _, id := range ids {
		if ctx.Resolver.Controller(id) == ctx.Controller {
			friendly = append(friendly, id)
		} else {
			enemy = append(enemy, id)
		}
	}
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

// MostPowerful returns a Refinement that keeps the n most powerful creatures of a
// set, e.g. card.Target.EachCreature.Refine(card.MostPowerful(3)) (Three Fates).
// When more creatures tie at the cutoff than there are remaining slots, the
// controller chooses which of the tied creatures to include.
func MostPowerful(n int) Refinement { return mostPowerfulN{n: n} }

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

// LowestAndHighestPower is a Refinement that keeps every creature tied for the
// lowest power together with every creature tied for the highest power, selecting
// both extremes at once so they resolve simultaneously (Standardized Testing
// destroys each creature with the lowest power and each with the highest).
var LowestAndHighestPower Refinement = lowestAndHighestPower{}

// lowestAndHighestPower implements the LowestAndHighestPower refinement.
type lowestAndHighestPower struct{}

// clause renders "<phrase> with the lowest power and <phrase> with the highest
// power", e.g. "each creature" -> "each creature with the lowest power and each
// creature with the highest power".
func (lowestAndHighestPower) clause(phrase string) string {
	return phrase + " with the lowest power and " + phrase + " with the highest power"
}

// refine keeps every creature whose power equals the set's minimum or maximum, so
// a set with a single power keeps all of it. An empty set selects nothing.
func (lowestAndHighestPower) refine(ctx *EffectContext, ids []LocalID) []LocalID {
	if len(ids) == 0 {
		return nil
	}
	low, high := ctx.Resolver.Power(ids[0]), ctx.Resolver.Power(ids[0])
	for _, id := range ids[1:] {
		switch p := ctx.Resolver.Power(id); {
		case p < low:
			low = p
		case p > high:
			high = p
		}
	}
	kept := make([]LocalID, 0, len(ids))
	for _, id := range ids {
		if p := ctx.Resolver.Power(id); p == low || p == high {
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

// selfHouseResolved resolves the SelfHouse sentinel inside the count, which lives
// in an unexported field reflection cannot reach on its own.
func (p powerLessThan) selfHouseResolved(house House) any {
	p.limit = resolvedIn(p.limit, house)
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
