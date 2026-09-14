package deckgen

import (
	"math/rand"
	"sort"

	"github.com/dmikalova/vactrol/internal/engine"
)

// Card is one pool entry: a card definition, its deck-building profile, and an
// optional Materializer. A nil Materializer means a concrete card, materialized
// by identity (plus any maverick rehousing).
type Card struct {
	Def          engine.CardDefinition
	Profile      GenerationProfile
	Materializer Materializer
}

// Set is the pool a Deck is generated from: cards bucketed by House and rarity,
// the houseless Special pool, the selectable Houses, and the Tuning. Build one
// with NewSet.
type Set struct {
	Name   string
	Tuning Tuning

	houses  []engine.House
	pool    map[engine.House]map[engine.Rarity][]Card
	byHouse map[engine.House][]Card
	byName  map[string]Card
	special []Card

	// clusters groups the set's cards into card families placed by strategy
	// (ADR 0036): the Horsemen (WholePool), the sins (RandomCount), the Shards
	// (OnePerHouse). It is keyed by cluster name and built from members' profiles.
	clusters map[string]clusterIndex

	// filtered holds the deck-wide filtered pulls the set's cards lead (ADR 0036):
	// a pull tops a deck up to a floor of cards matching a predicate (Chief
	// Engineer Walls guaranteeing Upgrades and Robots), keyed by cluster name.
	filtered map[string]FilteredCluster

	// legacy is the shared cross-set pool a slot may draw from at
	// Tuning.LegacyRate, keeping the pod's House. It is nil for a single-set build
	// and attached by WithLegacy; the same *Legacy is shared by every Set, which
	// draws only the entries whose set differs from its own Name.
	legacy *Legacy

	// errantHouses are the Houses present in the legacy pool but not native to
	// this Set — the foreign Houses an errant pod can bring into a deck (ADR 0036).
	// Computed by WithLegacy; empty without a legacy pool.
	errantHouses []engine.House

	// crossClusters is the catalog-wide cluster index used to resolve deck-wide
	// (OnePerHouse) clusters across every House a deck can reach, including an
	// errant House (ADR 0036). Nil for a single-set build, which resolves deck-wide
	// clusters from its own clusters; attached by WithClusters.
	crossClusters *ClusterPool
}

// NewSet builds a Set from a flat list of pool entries, bucketing them by House
// and rarity. Houseless cards go to the Special pool; cards with no House are
// dropped; Connected cards are indexed by name but never pooled (they enter a
// deck only when a cluster places them). The selectable Houses are those
// present in the pool, sorted by name.
func NewSet(name string, cards []Card, tuning Tuning) Set {
	s := Set{
		Name:    name,
		Tuning:  tuning,
		pool:    map[engine.House]map[engine.Rarity][]Card{},
		byHouse: map[engine.House][]Card{},
		byName:  map[string]Card{},
	}
	seen := map[engine.House]bool{}
	for _, c := range cards {
		if c.Def.Name != "" {
			s.byName[c.Def.Name] = c
		}
		if c.Profile.Houseless {
			s.special = append(s.special, c)
			continue
		}
		h := c.Def.House
		if h == engine.HouseNone {
			continue
		}
		if c.Def.Rarity == engine.Connected {
			continue
		}
		if s.pool[h] == nil {
			s.pool[h] = map[engine.Rarity][]Card{}
		}
		s.pool[h][c.Def.Rarity] = append(s.pool[h][c.Def.Rarity], c)
		s.byHouse[h] = append(s.byHouse[h], c)
		seen[h] = true
	}
	for h := range seen {
		s.houses = append(s.houses, h)
	}
	sort.Slice(s.houses, func(i, j int) bool { return s.houses[i].String() < s.houses[j].String() })
	s.clusters = buildClusters(cards)
	s.filtered = buildFilteredClusters(cards)
	s.validateClusters()
	s.validateFilteredClusters()
	return s
}

// LegacyEntry pairs a pool Card with the name of the source set it belongs to.
// The shared Legacy is built from these so that, when a Set draws a legacy card,
// it can exclude the entries that belong to the Set itself.
type LegacyEntry struct {
	Card Card
	Set  string
}

// Legacy is the cross-set pool shared by every Set: every set's cards bucketed by
// House and rarity, each tagged with the set it came from. One Legacy is built for
// the whole catalog with NewLegacy and referenced by every Set, rather than each
// Set holding its own copy of every other set's cards. A Set draws only the
// entries whose set differs from its own, so a legacy slot always pulls a card
// printed in another set.
type Legacy struct {
	// byHouseRarity buckets entries by House then rarity so a slot draws a legacy
	// card of its own rolled rarity; byHouse is the flat per-House fallback for a
	// rolled rarity that House has no legacy card of.
	byHouseRarity map[engine.House]map[engine.Rarity][]LegacyEntry
	byHouse       map[engine.House][]LegacyEntry
}

// NewLegacy buckets legacy entries by House and rarity. Only housed, non-Connected
// cards are pooled — a legacy card keeps its own House and rarity, since it is not
// rehoused the way a maverick is; Houseless Specials, Connected cards, and cards
// with no House are skipped for drawing.
func NewLegacy(entries []LegacyEntry) *Legacy {
	l := &Legacy{
		byHouseRarity: map[engine.House]map[engine.Rarity][]LegacyEntry{},
		byHouse:       map[engine.House][]LegacyEntry{},
	}
	for _, e := range entries {
		c := e.Card
		if c.Profile.Houseless || c.Def.Rarity == engine.Connected {
			continue
		}
		h := c.Def.House
		if h == engine.HouseNone {
			continue
		}
		l.byHouse[h] = append(l.byHouse[h], e)
		if l.byHouseRarity[h] == nil {
			l.byHouseRarity[h] = map[engine.Rarity][]LegacyEntry{}
		}
		l.byHouseRarity[h][c.Def.Rarity] = append(l.byHouseRarity[h][c.Def.Rarity], e)
	}
	return l
}

// candidates returns the entries' cards excluding those from the drawing set, so a
// Set never draws one of its own cards as a legacy card.
func (l *Legacy) candidates(entries []LegacyEntry, exclude string) []Card {
	out := make([]Card, 0, len(entries))
	for _, e := range entries {
		if e.Set == exclude {
			continue
		}
		out = append(out, e.Card)
	}
	return out
}

// houses returns the Houses the legacy pool can draw, sorted by name.
func (l *Legacy) houses() []engine.House {
	hs := make([]engine.House, 0, len(l.byHouse))
	for h := range l.byHouse {
		hs = append(hs, h)
	}
	sort.Slice(hs, func(i, j int) bool { return hs[i].String() < hs[j].String() })
	return hs
}

// WithLegacy attaches the shared Legacy pool to the set: cards printed in other
// sets that a slot may draw instead of one of this set's own, at the set's
// Tuning.LegacyRate. The same *Legacy is shared by every Set; this Set draws only
// the entries whose set differs from its Name. It also computes the set's errant
// Houses (the pool Houses not native to the set, ADR 0036). It returns the set
// with the pool attached so it reads as a builder step.
func (s Set) WithLegacy(l *Legacy) Set {
	s.legacy = l
	s.errantHouses = foreignHouses(s.houses, l.houses())
	return s
}

// foreignHouses returns the pool Houses not among the native ones, preserving the
// pool's order — the Houses an errant pod can bring into a deck from the cross-set
// legacy pool (ADR 0036).
func foreignHouses(native, pool []engine.House) []engine.House {
	isNative := make(map[engine.House]bool, len(native))
	for _, h := range native {
		isNative[h] = true
	}
	var out []engine.House
	for _, h := range pool {
		if !isNative[h] {
			out = append(out, h)
		}
	}
	return out
}

// WithClusters attaches the catalog-wide ClusterPool used to resolve deck-wide
// (OnePerHouse) clusters (ADR 0036). It gates the pool against this set: every
// OnePerHouse cluster must have a member for each House the set can deck — its
// native Houses and the foreign Houses an errant pod can bring in — or it panics,
// so a set that can deck a House with no member is caught at assembly, not with a
// short deck. Run it after WithLegacy, which computes the errant Houses.
func (s Set) WithClusters(cp *ClusterPool) Set {
	s.crossClusters = cp
	s.validateCrossClusters()
	return s
}

// Houses returns the Set's selectable Houses, sorted by name.
func (s Set) Houses() []engine.House { return append([]engine.House(nil), s.houses...) }

// member reports whether the set prints a card of this name — either natively or
// as a reprint it folds into its own pool (ADR 0021). A legacy draw uses it to
// avoid tagging a card the set already prints as legacy: the legacy pool still
// holds such cards, but one drawn from a legacy slot is one of the set's own
// cards, not a guest from another set.
func (s Set) member(name string) bool {
	_, ok := s.byName[name]
	return ok
}

// pickHouses selects PodCount distinct Houses, weighted and honoring exclusions,
// and returns them sorted by name.
func (s Set) pickHouses(r *rand.Rand) []engine.House {
	remaining := append([]engine.House(nil), s.houses...)
	excluded := map[engine.House]bool{}
	picked := make([]engine.House, 0, PodCount)
	for len(picked) < PodCount {
		cands := make([]engine.House, 0, len(remaining))
		for _, h := range remaining {
			if !excluded[h] {
				cands = append(cands, h)
			}
		}
		if len(cands) == 0 {
			break
		}
		h := s.weightedHouse(cands, r)
		picked = append(picked, h)
		remaining = removeHouse(remaining, h)
		for _, pair := range s.Tuning.HouseExclusions {
			switch h {
			case pair[0]:
				excluded[pair[1]] = true
			case pair[1]:
				excluded[pair[0]] = true
			}
		}
	}
	return picked
}

// weightedHouse draws one House from cands by weight, deterministically from r.
func (s Set) weightedHouse(cands []engine.House, r *rand.Rand) engine.House {
	total := 0.0
	for _, h := range cands {
		total += s.houseWeight(h)
	}
	x := r.Float64() * total
	for i := 0; i < len(cands)-1; i++ {
		if x -= s.houseWeight(cands[i]); x < 0 {
			return cands[i]
		}
	}
	return cands[len(cands)-1]
}

func (s Set) houseWeight(h engine.House) float64 {
	if w, ok := s.Tuning.HouseWeights[h]; ok && w > 0 {
		return w
	}
	return 1
}

func removeHouse(hs []engine.House, drop engine.House) []engine.House {
	out := hs[:0]
	for _, h := range hs {
		if h != drop {
			out = append(out, h)
		}
	}
	return out
}
