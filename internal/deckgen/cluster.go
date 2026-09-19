package deckgen

import (
	"fmt"
	"sort"

	"github.com/dmikalova/vactrol/internal/engine"
)

// ClusterStrategy is how a cluster fills out once a member triggers it
// (ADR 0036). The zero value is invalid, caught at NewSet.
type ClusterStrategy uint8

const (
	clusterStrategyInvalid ClusterStrategy = iota
	// WholePool places every member of the cluster (the four Horsemen).
	WholePool
	// RandomCount places a random number of distinct members in [Min, Max] (the
	// seven sins, 3–7).
	RandomCount
	// SelfPull places a random number of copies of the single triggering member
	// itself (Plague Rat pulls more Plague Rats). The count is Min + Poisson(Mean −
	// Min), capped at PodSize: at least Min, averaging about Mean, with a thin tail
	// that reaches a full pod only very rarely. Its cluster has exactly one member.
	SelfPull
	// PullExact places one copy of each non-lead member per lead instance in the
	// pod — two Timetravellers pull two Help from Future Self, never one. It is
	// always ByLead: the lead is the puller, the other members its exact partners.
	PullExact
	// Pull places a per-partner random count of each non-lead member when the lead
	// rolls in — Troop Call pulls a couple of Niffle Apes and, less often, a Niffle
	// Queen. It is always ByLead. Unlike the other strategies each partner carries
	// its own rate (Min copies at least, averaging about Mean, on a Poisson tail),
	// set with card.Pulled, so one lead can pull two partners at different rates.
	Pull
	// OnePerHouse places one member in each of the deck's Houses. It is the only
	// deck-wide strategy and the only one gated complete-by-construction — every
	// House the set can deck must have a member, or NewSet panics. Shards use it.
	OnePerHouse
	// PerGigantic places one member, chosen at random, into each gigantic base's
	// pod, stamped to that pod's House. It fires off deck structure (a gigantic
	// present), not a member being drawn, so its members are Houseless reservoir
	// cards drawn only by the pull. The tutors use it (ADR 0044).
	PerGigantic
)

// ClusterTrigger is what fires a cluster. The zero value is invalid, caught at
// NewSet.
type ClusterTrigger uint8

const (
	clusterTriggerInvalid ClusterTrigger = iota
	// ByLead fires the cluster only when its designated lead member is placed
	// (Horseman of Pestilence pulls the other Horsemen).
	ByLead
	// ByAnyMember fires the cluster when any member is placed (any sin, any
	// Shard).
	ByAnyMember
)

// ClusterMembership marks a card a member of a named cluster and carries the
// cluster's strategy and trigger, so the generator can resolve the whole family
// from any one member (ADR 0036). Every member of a cluster carries identical
// Strategy/Trigger — card.InCluster copies them from a single shared card.Cluster
// declaration, so they cannot drift. Min/Max/Mean are the fill's rate: shared by
// the family for RandomCount and SelfPull, but per-partner for Pull (each pulled
// member sets its own with card.Pulled). Lead marks the one member that fires a
// ByLead cluster.
type ClusterMembership struct {
	Name     string
	Strategy ClusterStrategy
	Trigger  ClusterTrigger
	Min, Max int
	// Mean is the average count a Poisson pull aims for (SelfPull copies of self,
	// Pull copies of a partner); unused by the other strategies.
	Mean float64
	Lead bool
}

// Empty reports whether the card belongs to no cluster.
func (m ClusterMembership) Empty() bool { return m.Name == "" }

// clusterIndex is a cluster resolved to its members: the strategy that fills it,
// the trigger that fires it, its members (by pool entry and by House), and the
// name of its lead member for a ByLead cluster.
type clusterIndex struct {
	strategy ClusterStrategy
	trigger  ClusterTrigger
	min, max int
	mean     float64
	members  []Card
	byHouse  map[engine.House]Card
	lead     string
}

// ClusterPool is the catalog-wide cluster index: the clusters of the whole card
// catalog, built once from every registered card and attached to each base Set
// with WithClusters, so a deck-wide OnePerHouse cluster (the Shards) resolves
// across every House a deck can reach — including a House an errant pod brings in
// from another set — not only the Houses of the drawing set (ADR 0036).
type ClusterPool struct {
	clusters map[string]clusterIndex
}

// NewClusterPool builds a ClusterPool from a flat list of every catalog card. It
// indexes the same way NewSet does, but over the whole catalog rather than one
// set and with no draw-pool filtering: a Connected member (a reservoir Shard) is a
// cluster member here just as it is in a set's own clusters.
func NewClusterPool(cards []Card) *ClusterPool {
	return &ClusterPool{clusters: buildClusters(cards)}
}

// Draftable reports whether a pool Card would enter a Set's draw pool: it is
// housed, neither Houseless nor Connected, and not a reservoir card. A set whose
// whole pool is undraftable is a reservoir set — its cards enter decks only
// through a cross-set cluster — and gets no Set of its own (ADR 0036).
func Draftable(c Card) bool {
	if c.Profile.Houseless || c.Profile.Reservoir {
		return false
	}
	return c.Def.House != engine.HouseNone && c.Def.Rarity != engine.Connected
}

// buildClusters groups a set's cards into clusters by membership name. Each
// member contributes its identical strategy/trigger; a member's own House keys it
// in byHouse (for OnePerHouse placement); the Lead flag records the ByLead lead.
// Members are sorted by name so the index is deterministic regardless of package
// init order.
func buildClusters(cards []Card) map[string]clusterIndex {
	idx := map[string]clusterIndex{}
	for i := range cards {
		c := cards[i]
		m := c.Profile.Cluster
		if m.Empty() {
			continue
		}
		ci := idx[m.Name]
		ci.strategy = m.Strategy
		ci.trigger = m.Trigger
		ci.min, ci.max = m.Min, m.Max
		ci.mean = m.Mean
		if ci.byHouse == nil {
			ci.byHouse = map[engine.House]Card{}
		}
		ci.members = append(ci.members, c)
		ci.byHouse[c.Def.House] = c
		if m.Lead {
			ci.lead = c.Def.Name
		}
		idx[m.Name] = ci
	}
	for name, ci := range idx {
		sort.Slice(
			ci.members,
			func(i, j int) bool { return ci.members[i].Def.Name < ci.members[j].Def.Name },
		)
		idx[name] = ci
	}
	return idx
}

// validateClusters fails loudly on an authoring-time contradiction: an invalid
// strategy or trigger, a ByLead cluster without exactly one lead, a RandomCount
// range that cannot be satisfied, or — for OnePerHouse — a House the set can deck
// that has no member (the complete-by-construction gate, ADR 0036 / ADR 0018). A
// missing House is satisfied explicitly by a //go:build todo stub card, so the
// gap is visible rather than a silent short deck.
func (s Set) validateClusters() {
	for name, ci := range s.clusters {
		s.validateCluster(name, ci)
	}
}

// validateCluster panics if a single named cluster is malformed: it has a strategy
// and trigger, a ByLead cluster names a lead, each strategy's own shape holds, the
// cluster can actually be rolled to fire, and a OnePerHouse cluster covers every
// house. PerGigantic fires off deck structure (a gigantic present), not a member
// being drawn, so it needs no trigger, lead, or rollable member.
func (s Set) validateCluster(name string, ci clusterIndex) {
	if ci.strategy == clusterStrategyInvalid {
		panic(fmt.Sprintf("deckgen: cluster %q in set %q has no strategy", name, s.Name))
	}
	if ci.strategy == PerGigantic {
		return
	}
	if ci.trigger == clusterTriggerInvalid {
		panic(fmt.Sprintf("deckgen: cluster %q in set %q has no trigger", name, s.Name))
	}
	if ci.trigger == ByLead && ci.lead == "" {
		panic(
			fmt.Sprintf(
				"deckgen: ByLead cluster %q in set %q has no lead member",
				name,
				s.Name,
			),
		)
	}
	switch ci.strategy {
	case RandomCount:
		validateRandomCountCluster(name, s.Name, ci)
	case SelfPull:
		validateSelfPullCluster(name, s.Name, ci)
	case PullExact:
		validatePullExactCluster(name, s.Name, ci)
	case Pull:
		validatePullCluster(name, s.Name, ci)
	}
	if !clusterCanFire(ci) {
		panic(fmt.Sprintf(
			"deckgen: cluster %q in set %q can never fire — %s, so nothing rolls to "+
				"trigger it (give a triggering member a rollable rarity)",
			name, s.Name, clusterTriggerDesc(ci.trigger),
		))
	}
	if ci.strategy == OnePerHouse {
		s.validateOnePerHouseCluster(name, ci)
	}
}

// validateRandomCountCluster panics if a RandomCount cluster's [min,max] does not
// fit its placeable pool. A ByLead RandomCount cluster plants its lead and places a
// count of the other members, so the placeable pool is the non-lead members (Dark
// Harbinger pulls its Mutations). A ByAnyMember one has no lead, so every member is
// placeable (the sins).
func validateRandomCountCluster(name, setName string, ci clusterIndex) {
	placeable := len(ci.members)
	if ci.lead != "" {
		placeable--
	}
	if ci.min < 1 || ci.max < ci.min || ci.max > placeable {
		panic(fmt.Sprintf(
			"deckgen: RandomCount cluster %q in set %q wants [%d,%d] of %d placeable members",
			name,
			setName,
			ci.min,
			ci.max,
			placeable,
		))
	}
}

// validateSelfPullCluster panics if a SelfPull cluster is not a single member with
// a well-formed 1 ≤ min ≤ mean ≤ PodSize count.
func validateSelfPullCluster(name, setName string, ci clusterIndex) {
	if len(ci.members) != 1 || ci.min < 1 || ci.mean < float64(ci.min) ||
		ci.mean > PodSize {
		panic(fmt.Sprintf(
			"deckgen: SelfPull cluster %q in set %q wants min %d mean %g of %d members "+
				"(needs exactly one member, 1 ≤ min ≤ mean ≤ %d)",
			name, setName, ci.min, ci.mean, len(ci.members), PodSize,
		))
	}
}

// validatePullExactCluster panics if a PullExact cluster is not ByLead with at
// least a lead and one other member.
func validatePullExactCluster(name, setName string, ci clusterIndex) {
	if ci.trigger != ByLead || len(ci.members) < 2 {
		panic(fmt.Sprintf(
			"deckgen: PullExact cluster %q in set %q needs ByLead and a partner "+
				"(a lead plus at least one other member)",
			name, setName,
		))
	}
}

// validatePullCluster panics if a Pull cluster is not ByLead with a partner, or any
// non-lead member's pulled count is malformed.
func validatePullCluster(name, setName string, ci clusterIndex) {
	if ci.trigger != ByLead || len(ci.members) < 2 {
		panic(fmt.Sprintf(
			"deckgen: Pull cluster %q in set %q needs ByLead and a partner "+
				"(a lead plus at least one other member)",
			name, setName,
		))
	}
	for i := range ci.members {
		m := ci.members[i]
		r := m.Profile.Cluster
		if m.Def.Name == ci.lead {
			continue
		}
		if r.Min < 0 || r.Mean < float64(r.Min) {
			panic(fmt.Sprintf(
				"deckgen: Pull cluster %q in set %q pulls %q at min %d mean %g "+
					"(needs 0 ≤ min ≤ mean; set it with card.Pulled)",
				name, setName, m.Def.Name, r.Min, r.Mean,
			))
		}
	}
}

// validateOnePerHouseCluster panics if a OnePerHouse cluster is missing a member
// for any house the set uses.
func (s Set) validateOnePerHouseCluster(name string, ci clusterIndex) {
	for _, h := range s.houses {
		if _, ok := ci.byHouse[h]; !ok {
			panic(fmt.Sprintf(
				"deckgen: OnePerHouse cluster %q in set %q has no member for House %s "+
					"(stub one with //go:build todo)",
				name, s.Name, h,
			))
		}
	}
}

// clusterCanFire reports whether a cluster has a triggering member that rolls in
// the pool. A cluster whose every triggering member is Rarity.Connected can never
// be drawn to fire itself. For ByLead only the lead fires it, so the lead must
// roll; for ByAnyMember any non-Connected member suffices.
func clusterCanFire(ci clusterIndex) bool {
	if ci.trigger == ByLead {
		for i := range ci.members {
			m := ci.members[i]
			if m.Def.Name == ci.lead {
				return m.Def.Rarity != engine.Connected
			}
		}
	}
	for i := range ci.members {
		m := ci.members[i]
		if m.Def.Rarity != engine.Connected {
			return true
		}
	}
	return false
}

// clusterTriggerDesc names the triggering member for a validation message.
func clusterTriggerDesc(t ClusterTrigger) string {
	if t == ByLead {
		return "its lead member is Rarity.Connected"
	}
	return "every member is Rarity.Connected"
}

// validateCrossClusters gates the attached catalog ClusterPool against this set:
// every OnePerHouse cluster must have a member for each House the set can deck —
// its native Houses and the foreign Houses an errant pod can bring in — or it
// panics, the same complete-by-construction gate validateClusters applies to a
// set's own OnePerHouse clusters (ADR 0036 / ADR 0018), widened to the Houses an
// errant pod reaches. It is a no-op when no ClusterPool is attached.
func (s Set) validateCrossClusters() {
	if s.crossClusters == nil {
		return
	}
	houses := append(append([]engine.House(nil), s.houses...), s.errantHouses...)
	for name, ci := range s.crossClusters.clusters {
		if ci.strategy != OnePerHouse {
			continue
		}
		for _, h := range houses {
			if _, ok := ci.byHouse[h]; !ok {
				panic(fmt.Sprintf(
					"deckgen: cross-set OnePerHouse cluster %q has no member for House %s "+
						"deckable by set %q (native or errant)",
					name, h, s.Name,
				))
			}
		}
	}
}

// onePerHouseClusters returns the OnePerHouse clusters to resolve deck-wide, in
// name order: from the attached catalog ClusterPool when one is set, else the
// set's own clusters. The catalog pool lets a deck-wide cluster complete across
// every House a deck can reach (including an errant House), where the set's own
// clusters cover only its native Houses.
func (s Set) onePerHouseClusters() []clusterIndex {
	src := s.clusters
	if s.crossClusters != nil {
		src = s.crossClusters.clusters
	}
	names := make([]string, 0, len(src))
	for name, ci := range src {
		if ci.strategy == OnePerHouse {
			names = append(names, name)
		}
	}
	sort.Strings(names)
	out := make([]clusterIndex, 0, len(names))
	for _, name := range names {
		out = append(out, src[name])
	}
	return out
}

// perGiganticClusters returns the PerGigantic clusters to resolve when a gigantic
// is placed, in name order, drawn from the attached catalog ClusterPool when one
// is set so a tutor family shared across sets (the tutors) is reachable by every
// gigantic-bearing set, else from the set's own clusters.
func (s Set) perGiganticClusters() []clusterIndex {
	src := s.clusters
	if s.crossClusters != nil {
		src = s.crossClusters.clusters
	}
	names := make([]string, 0, len(src))
	for name, ci := range src {
		if ci.strategy == PerGigantic {
			names = append(names, name)
		}
	}
	sort.Strings(names)
	out := make([]clusterIndex, 0, len(names))
	for _, name := range names {
		out = append(out, src[name])
	}
	return out
}
