package deckgen

import (
	"math"
	"math/rand"
	"sort"

	"github.com/dmikalova/vactrol/internal/engine"
)

// Generate builds a Deck from a Set and a seed, deterministically. The same
// (Set, seed) always yields the same Deck within one version of the Set's pool.
func Generate(set Set, seed int64) Deck {
	if set.Tuning.RarityWeights == nil {
		set.Tuning = DefaultTuning()
	}
	g := &generator{set: set, r: rand.New(rand.NewSource(seed)), placed: map[string]bool{}}
	houses := set.pickHouses(g.r)
	for i := 0; i < PodCount && i < len(houses); i++ {
		g.deckHouses[i] = houses[i]
	}
	deck := Deck{Set: set.Name, Seed: seed}
	for i := 0; i < PodCount && i < len(houses); i++ {
		deck.Pods[i] = g.expandPodClusters(g.fillPod(houses[i]))
	}
	g.expandClusters(&deck)
	return deck
}

// generator threads the single RNG and the deck-wide "already placed" set (for
// one-copy-per-deck) through the pipeline.
type generator struct {
	set        Set
	r          *rand.Rand
	placed     map[string]bool
	deckHouses [PodCount]engine.House
}

// placedCard remembers a filled slot's pool entry and rolled rarity so the
// duplicate-pull can copy an earlier same-rarity card in the same pod.
type placedCard struct {
	card   Card
	rarity engine.Rarity
}

func (g *generator) fillPod(house engine.House) HousePod {
	pod := HousePod{House: house}
	placed := make([]placedCard, 0, PodSize)
	for i := 0; i < PodSize; i++ {
		slot, pc := g.fillSlot(house, placed)
		pod.Slots[i] = slot
		placed = append(placed, pc)
	}
	return pod
}

// expandPodClusters resolves the pod-local cluster strategies for one filled pod:
// WholePool (place every member), RandomCount (a random count of distinct
// members), SelfPull (a random count of copies of the single triggering member),
// PullExact (one of each partner per lead instance), and Pull (a per-partner
// random count of each partner when the lead rolls in). A cluster fires when its
// trigger is met in the pod — ByLead when its lead is present, ByAnyMember when any
// member is. Each strategy tops the pod up by overwriting slots that hold no member
// of the cluster, so the fired member and any siblings already placed stay put.
// OnePerHouse is deck-wide and resolved in expandClusters, not here. Clusters are
// visited in name order so the fill is deterministic when a set has several.
func (g *generator) expandPodClusters(pod HousePod) HousePod {
	for _, name := range g.clusterNames() {
		ci := g.set.clusters[name]
		if ci.strategy == OnePerHouse || !podClusterFires(pod, ci) {
			continue
		}
		switch ci.strategy {
		case WholePool:
			for _, m := range ci.members {
				g.placeMemberCopies(&pod, ci, m, 1)
			}
		case RandomCount:
			for _, m := range g.randomMembers(ci) {
				g.placeMemberCopies(&pod, ci, m, 1)
			}
		case SelfPull:
			g.placeMemberCopies(&pod, ci, ci.members[0], g.selfPullCount(ci))
		case PullExact:
			leads := countMember(pod, ci.lead)
			for _, m := range ci.members {
				if m.Def.Name != ci.lead {
					g.placeMemberCopies(&pod, ci, m, leads)
				}
			}
		case Pull:
			for _, m := range ci.members {
				if m.Def.Name != ci.lead {
					g.placeMemberCopies(&pod, ci, m, g.pullCount(m.Profile.Cluster))
				}
			}
		}
	}
	return pod
}

// clusterNames returns the set's cluster names in sorted order, so pod-local
// cluster expansion is deterministic when a set has several.
func (g *generator) clusterNames() []string {
	names := make([]string, 0, len(g.set.clusters))
	for name := range g.set.clusters {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// podClusterFires reports whether a cluster's trigger is met in the pod: for
// ByLead, that its lead member is present; for ByAnyMember, that any member is.
func podClusterFires(pod HousePod, ci clusterIndex) bool {
	for _, s := range pod.Slots {
		name := s.Card.Name
		if ci.trigger == ByLead {
			if name == ci.lead {
				return true
			}
			continue
		}
		if inCluster(ci, name) {
			return true
		}
	}
	return false
}

// randomMembers picks a random count of distinct members in [min, max] for a
// RandomCount cluster, by shuffling the members and taking that many.
func (g *generator) randomMembers(ci clusterIndex) []Card {
	k := ci.min
	if ci.max > ci.min {
		k += g.r.Intn(ci.max - ci.min + 1)
	}
	shuffled := append([]Card(nil), ci.members...)
	g.r.Shuffle(len(shuffled), func(i, j int) {
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	})
	return shuffled[:k]
}

// selfPullCount rolls a SelfPull copy count: Min + Poisson(Mean − Min), capped at
// PodSize — at least Min, averaging about Mean, reaching a whole pod only on the
// thin tail.
func (g *generator) selfPullCount(ci clusterIndex) int {
	n := ci.min + poisson(g.r, ci.mean-float64(ci.min))
	if n > PodSize {
		n = PodSize
	}
	return n
}

// pullCount rolls a Pull partner count from that partner's own rate: Min +
// Poisson(Mean − Min), capped at PodSize. A partner with Min 0 (Niffle Queen)
// often rolls none.
func (g *generator) pullCount(m ClusterMembership) int {
	n := m.Min + poisson(g.r, m.Mean-float64(m.Min))
	if n > PodSize {
		n = PodSize
	}
	return n
}

// poisson draws from a Poisson distribution with the given mean using Knuth's
// algorithm; a mean of 0 always returns 0.
func poisson(r *rand.Rand, lambda float64) int {
	l := math.Exp(-lambda)
	k, p := 0, 1.0
	for {
		k++
		p *= r.Float64()
		if p <= l {
			return k - 1
		}
	}
}

// placeMemberCopies tops the pod up to want copies of member m, overwriting slots
// that hold no member of the cluster so the triggering member and its placed
// siblings are never displaced. A member whose native House differs from the pod's
// is rehoused as a maverick; a OneCopyPerDeck member is recorded as placed.
func (g *generator) placeMemberCopies(pod *HousePod, ci clusterIndex, m Card, want int) {
	for countMember(*pod, m.Def.Name) < want {
		slot := freeClusterSlot(pod, ci)
		if slot < 0 {
			return
		}
		maverick := m.Def.House != pod.House
		def := g.materialize(
			m,
			SlotContext{House: pod.House, Rarity: m.Def.Rarity, Maverick: maverick},
		)
		if m.Profile.OneCopyPerDeck {
			g.placed[m.Def.Name] = true
		}
		pod.Slots[slot] = Slot{Rarity: m.Def.Rarity, Maverick: maverick, Card: def}
	}
}

// freeClusterSlot returns the first slot holding no member of the cluster, or -1
// when every slot already holds one.
func freeClusterSlot(pod *HousePod, ci clusterIndex) int {
	for i := range pod.Slots {
		if !inCluster(ci, pod.Slots[i].Card.Name) {
			return i
		}
	}
	return -1
}

// countMember counts how many of the pod's slots hold the named card.
func countMember(pod HousePod, name string) int {
	n := 0
	for _, s := range pod.Slots {
		if s.Card.Name == name {
			n++
		}
	}
	return n
}

// expandClusters resolves the deck-wide cluster strategies once every pod is
// filled. Today that is OnePerHouse (the Shards): when any Shard has been drawn
// into the deck, each pod is given its House's Shard, so a single Shard pulls the
// whole cycle in. The pod-local strategies (WholePool, RandomCount) resolve in
// the pod pass, not here.
func (g *generator) expandClusters(deck *Deck) {
	for _, ci := range g.set.clusters {
		if ci.strategy != OnePerHouse {
			continue
		}
		if !g.clusterTriggered(deck, ci) {
			continue
		}
		for i := 0; i < PodCount; i++ {
			pod := &deck.Pods[i]
			if pod.House == engine.HouseNone {
				continue
			}
			// validateClusters guarantees every deckable House has a member, so
			// the pod's House indexes one.
			g.ensureClusterMember(pod, ci, ci.byHouse[pod.House])
		}
	}
}

// clusterTriggered reports whether a cluster has fired: for ByAnyMember, that any
// member has been drawn into the deck; for ByLead, that the lead member has.
func (g *generator) clusterTriggered(deck *Deck, ci clusterIndex) bool {
	for i := 0; i < PodCount; i++ {
		for _, s := range deck.Pods[i].Slots {
			name := s.Card.Name
			if ci.trigger == ByLead {
				if name == ci.lead {
					return true
				}
				continue
			}
			if inCluster(ci, name) {
				return true
			}
		}
	}
	return false
}

// ensureClusterMember gives a pod its House's cluster member, unless the pod
// already holds a member of the cluster (the Shard that fired the cycle stays put
// in its own pod). It overwrites a random slot. At the Set's MaverickRate the pod
// instead receives another House's member, rehoused as a maverick — the printed-
// house rule KeyForge uses for a maverick card.
func (g *generator) ensureClusterMember(pod *HousePod, ci clusterIndex, member Card) {
	for _, s := range pod.Slots {
		if inCluster(ci, s.Card.Name) {
			return
		}
	}
	place, maverick := member, false
	if g.chance(g.set.Tuning.MaverickRate) {
		place, maverick = g.otherClusterMember(ci, member), true
	}
	ctx := SlotContext{House: pod.House, Rarity: place.Def.Rarity, Maverick: maverick}
	def := g.materialize(place, ctx)
	if place.Profile.OneCopyPerDeck {
		g.placed[place.Def.Name] = true
	}
	pod.Slots[g.r.Intn(PodSize)] = Slot{
		Rarity:   place.Def.Rarity,
		Maverick: maverick,
		Card:     def,
	}
}

// otherClusterMember picks a member of the cluster whose name differs from the
// excluded one, for a maverick substitution. It is only reached for a cluster
// with more than one member (a triggered OnePerHouse cluster deckable in a
// multi-House deck always has at least two), so the pick is never empty.
func (g *generator) otherClusterMember(ci clusterIndex, exclude Card) Card {
	alts := make([]Card, 0, len(ci.members))
	for _, m := range ci.members {
		if m.Def.Name != exclude.Def.Name {
			alts = append(alts, m)
		}
	}
	return alts[g.r.Intn(len(alts))]
}

// inCluster reports whether a card name is a member of the cluster.
func inCluster(ci clusterIndex, name string) bool {
	for _, m := range ci.members {
		if m.Def.Name == name {
			return true
		}
	}
	return false
}

// fillSlot resolves one slot: a rare Special overlay, else a rarity roll followed
// by an optional duplicate-pull or a fresh (possibly maverick) draw.
func (g *generator) fillSlot(house engine.House, placed []placedCard) (Slot, placedCard) {
	t := g.set.Tuning
	if len(g.set.special) > 0 && g.chance(t.SpecialRate) {
		if c, ok := g.pick(g.set.special); ok {
			return g.commit(c, SlotContext{House: house, Rarity: c.Def.Rarity, Special: true})
		}
	}

	rarity := g.rollRarity()
	if g.chance(t.LegacyRate) {
		if c, ok := g.drawLegacy(house, rarity); ok {
			// A legacy slot may land on a card this set also prints as a reprint
			// (ADR 0021 keeps such cards in the legacy pool). A card the set prints
			// is one of its own, so it is not tagged legacy even when a legacy slot
			// drew it.
			return g.commit(c, SlotContext{
				House:  house,
				Rarity: c.Def.Rarity,
				Legacy: !g.set.member(c.Def.Name),
			})
		}
	}
	if c, ok := g.tryDuplicate(rarity, placed); ok {
		return g.commit(c, SlotContext{House: house, Rarity: rarity})
	}

	maverick := g.chance(t.MaverickRate)
	c, ok := g.draw(house, rarity, maverick)
	if !ok {
		return Slot{Rarity: rarity}, placedCard{rarity: rarity}
	}
	ctx := SlotContext{House: house, Rarity: rarity, Maverick: maverick && c.Def.House != house}
	return g.commit(c, ctx)
}

// commit materializes the card, records a one-copy-per-deck placement, and builds
// the slot and its placed record.
func (g *generator) commit(c Card, ctx SlotContext) (Slot, placedCard) {
	def := g.materialize(c, ctx)
	if c.Profile.OneCopyPerDeck {
		g.placed[c.Def.Name] = true
	}
	slot := Slot{
		Rarity:   ctx.Rarity,
		Maverick: ctx.Maverick,
		Legacy:   ctx.Legacy,
		Special:  ctx.Special,
		Card:     def,
	}
	return slot, placedCard{card: c, rarity: ctx.Rarity}
}

// materialize produces the final playable definition. A template binds itself via
// its Materializer; a concrete card is used as-is. Either way a Maverick or
// Special card adopts the pod's House. The deck's Houses are threaded in so a
// template can bind a partner house (ADR 0036).
func (g *generator) materialize(c Card, ctx SlotContext) engine.CardDefinition {
	ctx.DeckHouses = g.deckHouses
	def := c.Def
	if c.Materializer != nil {
		def = c.Materializer.Materialize(ctx, g.r)
	}
	if ctx.Maverick || ctx.Special {
		def.House = ctx.House
	}
	return def
}

// draw picks a card for the pod House at the rolled rarity. A maverick draw comes
// from a different House of the same Set; it falls back to any rarity in the
// source House, then to any card of the pod House, so a slot always fills.
func (g *generator) draw(house engine.House, rarity engine.Rarity, maverick bool) (Card, bool) {
	source := house
	if maverick {
		if alt, ok := g.otherHouse(house); ok {
			source = alt
		}
	}
	if c, ok := g.pick(g.set.pool[source][rarity]); ok {
		return c, true
	}
	if c, ok := g.pick(g.set.byHouse[source]); ok {
		return c, true
	}
	return g.pick(g.set.byHouse[house])
}

// drawLegacy picks a legacy card for the pod House at the rolled rarity from the
// shared cross-set pool, excluding this set's own cards, and falling back to any
// rarity in that House so a legacy slot still fills when the House has no legacy
// card of that rarity. It reports false when no legacy pool is attached, or the
// House has no legacy card from another set.
func (g *generator) drawLegacy(house engine.House, rarity engine.Rarity) (Card, bool) {
	l := g.set.legacy
	if l == nil {
		return Card{}, false
	}
	if c, ok := g.pick(l.candidates(l.byHouseRarity[house][rarity], g.set.Name)); ok {
		return c, true
	}
	return g.pick(l.candidates(l.byHouse[house], g.set.Name))
}

// tryDuplicate copies an already-placed same-pod, same-rarity card with the
// per-rarity duplicate probability. One-copy-per-deck cards are never copied.
func (g *generator) tryDuplicate(rarity engine.Rarity, placed []placedCard) (Card, bool) {
	rate := g.set.Tuning.DuplicateRate[rarity]
	if rate <= 0 || !g.chance(rate) {
		return Card{}, false
	}
	elig := make([]Card, 0, len(placed))
	for _, pc := range placed {
		if pc.rarity != rarity || pc.card.Def.Name == "" || pc.card.Profile.OneCopyPerDeck {
			continue
		}
		elig = append(elig, pc.card)
	}
	if len(elig) == 0 {
		return Card{}, false
	}
	return elig[g.r.Intn(len(elig))], true
}

// pick returns a random eligible card, skipping one-copy-per-deck cards already
// placed. Eligible cards are drawn weighted by their RarityWeight (default 1),
// so a card can be made rarer within its rarity than its peers. It reports false
// when nothing is eligible.
func (g *generator) pick(cards []Card) (Card, bool) {
	elig := make([]Card, 0, len(cards))
	for _, c := range cards {
		if c.Profile.OneCopyPerDeck && g.placed[c.Def.Name] {
			continue
		}
		elig = append(elig, c)
	}
	if len(elig) == 0 {
		return Card{}, false
	}
	total := 0.0
	for _, c := range elig {
		total += drawWeight(c)
	}
	x := g.r.Float64() * total
	for i := 0; i < len(elig)-1; i++ {
		if x -= drawWeight(elig[i]); x < 0 {
			return elig[i], true
		}
	}
	return elig[len(elig)-1], true
}

// drawWeight is a card's relative draw weight within its house+rarity bucket,
// defaulting to 1 when its RarityWeight is unset (0 or less).
func drawWeight(c Card) float64 {
	if c.Profile.RarityWeight > 0 {
		return c.Profile.RarityWeight
	}
	return 1
}

// otherHouse picks a Set House other than house.
func (g *generator) otherHouse(house engine.House) (engine.House, bool) {
	others := make([]engine.House, 0, len(g.set.houses))
	for _, h := range g.set.houses {
		if h != house {
			others = append(others, h)
		}
	}
	if len(others) == 0 {
		return engine.HouseNone, false
	}
	return others[g.r.Intn(len(others))], true
}

// rarityOrder fixes the iteration order of rarities so a weighted roll is
// deterministic regardless of map ordering.
var rarityOrder = []engine.Rarity{
	engine.Common, engine.Uncommon, engine.Rare, engine.Special,
}

func (g *generator) rollRarity() engine.Rarity {
	w := g.set.Tuning.RarityWeights
	var order []engine.Rarity
	total := 0.0
	for _, rr := range rarityOrder {
		if w[rr] > 0 {
			order = append(order, rr)
			total += w[rr]
		}
	}
	if len(order) == 0 {
		return engine.Common
	}
	x := g.r.Float64() * total
	for i := 0; i < len(order)-1; i++ {
		if x -= w[order[i]]; x < 0 {
			return order[i]
		}
	}
	return order[len(order)-1]
}

func (g *generator) chance(p float64) bool { return g.r.Float64() < p }
