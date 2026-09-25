package deckgen

import (
	"fmt"
	"sort"

	"github.com/dmikalova/vex/internal/engine"
)

// FilteredCluster is a deck-wide pull keyed to a lead card that fills by a
// predicate instead of by named members (ADR 0036). When the lead card is in the
// deck, generation guarantees at least Floor cards matching Match somewhere in the
// deck, topping up from the set pool in any House (rehoused as mavericks
// off-House). Cards already in the deck that Match count toward the Floor, so a
// deck that rolls enough on its own pulls nothing. Unlike a ClusterMembership it
// names no members — the pulled cards are whatever the pool offers that Match — so
// a card can lead a filtered pull in addition to belonging to a named cluster.
type FilteredCluster struct {
	// Name identifies the cluster for diagnostics.
	Name string
	// Lead is the card whose presence in a deck fires the pull, set from the
	// declaring card at NewSet.
	Lead string
	// Floor is the minimum number of matching cards guaranteed in the deck.
	Floor int
	// Mean is the average number of matching cards the pull aims for. Each deck
	// rolls a target of Floor + Poisson(Mean − Floor), so the count is at least
	// Floor and averages about Mean on a Poisson tail. Mean equal to Floor (the
	// zero-Mean default too) pulls exactly to the Floor, matching a flat guarantee.
	Mean float64
	// Match reports whether a pool card satisfies the pull. It is a static
	// predicate over a definition, so it lives in deckgen data with no resolver.
	Match func(engine.CardDefinition) bool
}

// buildFilteredClusters indexes the filtered pulls a set's cards lead, keyed by
// cluster name and stamped with the declaring card as the Lead. A card marks
// itself a lead with card.PullsMatching.
func buildFilteredClusters(cards []Card) map[string]FilteredCluster {
	idx := map[string]FilteredCluster{}
	for i := range cards {
		c := cards[i]
		if c.Profile.Leads == nil {
			continue
		}
		fc := *c.Profile.Leads
		fc.Lead = c.Def.Name
		idx[fc.Name] = fc
	}
	return idx
}

// matchingPool returns the set's draftable pool cards whose definition satisfies
// the filtered cluster's predicate, sorted by name so the fill is deterministic.
func (s Set) matchingPool(fc FilteredCluster) []Card {
	var out []Card
	for _, h := range s.houses {
		cardsByHouse := s.byHouse[h]
		for i := range cardsByHouse {
			c := cardsByHouse[i]
			if fc.Match(c.Def) {
				out = append(out, c)
			}
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Def.Name < out[j].Def.Name })
	return out
}

// validateFilteredClusters fails loudly on an authoring contradiction: a nil
// predicate, a non-positive Floor, or a pool that cannot satisfy the Floor — the
// complete-by-construction gate (ADR 0036 / ADR 0018) applied to a filtered pull,
// so a deck can never silently come up short.
func (s Set) validateFilteredClusters() {
	for name, fc := range s.filtered {
		if fc.Match == nil {
			panic(fmt.Sprintf(
				"deckgen: filtered cluster %q in set %q has no predicate", name, s.Name,
			))
		}
		if fc.Floor < 1 {
			panic(fmt.Sprintf(
				"deckgen: filtered cluster %q in set %q wants floor %d (needs ≥ 1)",
				name, s.Name, fc.Floor,
			))
		}
		if n := len(s.matchingPool(fc)); n < fc.Floor {
			panic(fmt.Sprintf(
				"deckgen: filtered cluster %q in set %q has %d pool cards matching its "+
					"predicate, needs %d",
				name, s.Name, n, fc.Floor,
			))
		}
		if fc.Mean != 0 && fc.Mean < float64(fc.Floor) {
			panic(fmt.Sprintf(
				"deckgen: filtered cluster %q in set %q wants mean %g below floor %d",
				name, s.Name, fc.Mean, fc.Floor,
			))
		}
	}
}

// expandFilteredClusters resolves the deck-wide filtered pulls once every pod is
// filled and the OnePerHouse clusters are placed. For each pull whose lead is in
// the deck it counts the cards already matching the predicate and, if short of the
// Floor, tops up from the pool with distinct matching cards not already in the
// deck. Pulls are visited in name order so the fill is deterministic.
func (g *generator) expandFilteredClusters(deck *Deck) {
	for _, name := range g.filteredNames() {
		fc := g.set.filtered[name]
		if !deckHasCard(deck, fc.Lead) {
			continue
		}
		need := g.filteredTarget(fc) - deckCountMatching(deck, fc.Match)
		cands := g.filteredCandidates(deck, fc)
		for i := range cands {
			if need <= 0 {
				break
			}
			if g.placeFilteredMatch(deck, fc, cands[i]) {
				need--
			}
		}
	}
}

// filteredTarget rolls a pull's per-deck target count: Floor + Poisson(Mean −
// Floor), so it is at least the Floor and averages about the Mean. Existing
// matching cards count against it in expandFilteredClusters, so the roll is the
// total the deck aims for, not the number pulled.
func (g *generator) filteredTarget(fc FilteredCluster) int {
	return fc.Floor + poisson(g.r, fc.Mean-float64(fc.Floor))
}

// filteredNames returns the set's filtered-cluster names in sorted order.
func (g *generator) filteredNames() []string {
	names := make([]string, 0, len(g.set.filtered))
	for name := range g.set.filtered {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// deckHasCard reports whether any pod holds a card of the given name.
func deckHasCard(deck *Deck, name string) bool {
	for i := range PodCount {
		for j := range deck.Pods[i].Slots {
			if deck.Pods[i].Slots[j].Card.Name == name {
				return true
			}
		}
	}
	return false
}

// deckCountMatching counts the deck's slots whose card satisfies the predicate.
func deckCountMatching(deck *Deck, match func(engine.CardDefinition) bool) int {
	n := 0
	for i := range PodCount {
		for j := range deck.Pods[i].Slots {
			if match(deck.Pods[i].Slots[j].Card) {
				n++
			}
		}
	}
	return n
}

// filteredCandidates returns the pool cards that could be pulled to satisfy a
// filtered cluster: those matching its predicate, not already in the deck, and not
// a one-copy-per-deck card already placed. They are shuffled so the top-up varies
// by seed.
func (g *generator) filteredCandidates(deck *Deck, fc FilteredCluster) []Card {
	inDeck := map[string]bool{}
	for i := range PodCount {
		for j := range deck.Pods[i].Slots {
			inDeck[deck.Pods[i].Slots[j].Card.Name] = true
		}
	}
	var cands []Card
	for i := range g.set.matchingPool(fc) {
		c := g.set.matchingPool(fc)[i]
		if inDeck[c.Def.Name] || (c.Profile.OneCopyPerDeck && g.placed[c.Def.Name]) {
			continue
		}
		cands = append(cands, c)
	}
	g.r.Shuffle(len(cands), func(i, j int) { cands[i], cands[j] = cands[j], cands[i] })
	return cands
}

// protectedNames is the set of card names a filtered pull must not overwrite: the
// members of every named cluster and every deck-wide OnePerHouse cluster, so
// topping up an upgrade/robot floor never displaces a Shard, a Horseman, or a
// blaster's signature creature.
func (g *generator) protectedNames() map[string]bool {
	names := map[string]bool{}
	for _, ci := range g.set.clusters {
		for i := range ci.members {
			m := ci.members[i]
			names[m.Def.Name] = true
		}
	}
	for _, ci := range g.set.onePerHouseClusters() {
		for i := range ci.members {
			m := ci.members[i]
			names[m.Def.Name] = true
		}
	}
	return names
}

// placeFilteredMatch places one copy of cand by overwriting a slot that holds no
// protected card and does not already match the pull, so the fill never displaces
// the lead, a cluster member, or a card already counting toward the Floor. It
// prefers a pod of the card's own House (a natural, non-maverick fit) and falls
// back to any pod, rehousing the card as a maverick. It reports false when no such
// slot exists anywhere in the deck.
func (g *generator) placeFilteredMatch(deck *Deck, fc FilteredCluster, cand Card) bool {
	protected := g.protectedNames()
	for _, pi := range podOrder(deck, cand.Def.House) {
		pod := &deck.Pods[pi]
		if pod.House == engine.HouseNone {
			continue
		}
		for si := range pod.Slots {
			name := pod.Slots[si].Card.Name
			if name == fc.Lead || protected[name] || fc.Match(pod.Slots[si].Card) {
				continue
			}
			maverick := cand.Def.House != pod.House
			def := g.materialize(
				cand,
				SlotContext{
					House:    pod.House,
					Rarity:   cand.Def.Rarity,
					Maverick: maverick,
				},
			)
			if cand.Profile.OneCopyPerDeck {
				g.placed[cand.Def.Name] = true
			}
			pod.Slots[si] = Slot{
				Rarity:   cand.Def.Rarity,
				Maverick: maverick,
				Card:     def,
			}
			return true
		}
	}
	return false
}

// podOrder returns the deck's pod indices with the pods native to house first, so
// a pulled card lands in its own House as a normal card before falling back to a
// maverick slot elsewhere.
func podOrder(deck *Deck, house engine.House) []int {
	order := make([]int, 0, PodCount)
	for i := range PodCount {
		if deck.Pods[i].House == house {
			order = append(order, i)
		}
	}
	for i := range PodCount {
		if deck.Pods[i].House != house {
			order = append(order, i)
		}
	}
	return order
}
