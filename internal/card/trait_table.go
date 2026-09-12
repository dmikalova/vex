package card

import (
	"hash/fnv"
	"math/rand"
	"sort"
	"strings"
	"sync"

	"github.com/dmikalova/vactrol/internal/cards/provenance"
	"github.com/dmikalova/vactrol/internal/engine"
)

var (
	baneTraitOnce sync.Once
	baneTrait     map[engine.House]engine.Trait
)

// MostCommonCreatureTrait returns the trait carried by the most creatures of the
// given House across the whole KeyForge card pool (the provenance catalogs, not
// just the cards vactrol has implemented), computed once and cached. Sourcing the
// full pool means every House has a well-defined answer even before its creatures
// are implemented — Mars is Martian and Sanctum is Knight from the first card.
// When several traits tie for the most creatures, the winner is chosen by a
// shuffle seeded from the names of the creatures carrying a tied trait — so the
// choice is deterministic, independent of catalog order, and not biased toward
// alphabetically early traits. It returns the zero Trait for a House whose
// creatures carry no enumerated trait.
func MostCommonCreatureTrait(h engine.House) engine.Trait {
	baneTraitOnce.Do(buildBaneTraits)
	return baneTrait[h]
}

func buildBaneTraits() {
	counts := map[engine.House]map[engine.Trait]int{}
	names := map[engine.House]map[engine.Trait][]string{}
	for _, set := range provenance.Sets() {
		for _, c := range set.Cards {
			if c.Type != "creature" {
				continue
			}
			house, ok := parseHouseSlug(c.House)
			if !ok {
				continue
			}
			if counts[house] == nil {
				counts[house] = map[engine.Trait]int{}
				names[house] = map[engine.Trait][]string{}
			}
			for _, slug := range c.Traits {
				tr, ok := engine.ParseTrait(slug)
				if !ok {
					continue
				}
				counts[house][tr]++
				names[house][tr] = append(names[house][tr], c.Name)
			}
		}
	}
	baneTrait = make(map[engine.House]engine.Trait, len(counts))
	for h, c := range counts {
		baneTrait[h] = topTrait(c, names[h])
	}
}

// parseHouseSlug resolves a provenance house slug (lowercase, space-stripped, e.g.
// "staralliance") to its House. Later-set houses vactrol does not model return
// false and are skipped.
func parseHouseSlug(slug string) (engine.House, bool) {
	for h := engine.HouseNone + 1; int(h) < engine.NumHouses; h++ {
		if strings.EqualFold(strings.ReplaceAll(h.String(), " ", ""), slug) {
			return h, true
		}
	}
	return engine.HouseNone, false
}

// topTrait returns the trait with the highest count, breaking a tie with a
// shuffle seeded from the sorted names of the creatures carrying a tied trait.
func topTrait(counts map[engine.Trait]int, names map[engine.Trait][]string) engine.Trait {
	most := 0
	for _, n := range counts {
		if n > most {
			most = n
		}
	}
	var tied []engine.Trait
	for tr, n := range counts {
		if n == most {
			tied = append(tied, tr)
		}
	}
	sort.Slice(tied, func(i, j int) bool { return tied[i].String() < tied[j].String() })
	if len(tied) == 1 {
		return tied[0]
	}
	var seed []string
	for _, tr := range tied {
		seed = append(seed, names[tr]...)
	}
	sort.Strings(seed)
	h := fnv.New64a()
	for _, n := range seed {
		h.Write([]byte(n))
		h.Write([]byte{0})
	}
	r := rand.New(rand.NewSource(int64(h.Sum64())))
	r.Shuffle(len(tied), func(i, j int) { tied[i], tied[j] = tied[j], tied[i] })
	return tied[0]
}
