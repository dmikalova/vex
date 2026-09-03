package deckgen

import (
	"fmt"
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
}

// NewSet builds a Set from a flat list of pool entries, bucketing them by House
// and rarity. Houseless cards go to the Special pool; cards with no House are
// dropped; Connected cards are indexed by name but never pooled (they enter a
// deck only when a connection pulls them). The selectable Houses are those
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
	s.validateConnections()
	return s
}

// validateConnections fails loudly if any card connects to a card absent from the
// set, or pulls it an impossible number of times or at an impossible rate: a
// connection to a card that does not exist is an authoring error, never a link to
// silently drop at generation time.
func (s Set) validateConnections() {
	for _, cards := range s.byName {
		for _, cc := range cards.Profile.Connection.Cards {
			if _, ok := s.byName[cc.Name]; !ok {
				panic(fmt.Sprintf(
					"deckgen: card %q connects to %q, which is not in set %q",
					cards.Def.Name, cc.Name, s.Name,
				))
			}
			if cc.Copies < 1 {
				panic(fmt.Sprintf(
					"deckgen: card %q pulls %d copies of %q; a connection pulls at least one",
					cards.Def.Name, cc.Copies, cc.Name,
				))
			}
			if cc.Chance <= 0 || cc.Chance > 1 {
				panic(fmt.Sprintf(
					"deckgen: card %q pulls %q at chance %v; a connection fires in (0, 1]",
					cards.Def.Name, cc.Name, cc.Chance,
				))
			}
		}
	}
}

// Houses returns the Set's selectable Houses, sorted by name.
func (s Set) Houses() []engine.House { return append([]engine.House(nil), s.houses...) }

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
