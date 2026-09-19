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
	g := &generator{
		set:    set,
		r:      rand.New(rand.NewSource(seed)),
		placed: map[string]bool{},
	}
	houses := set.pickHouses(g.r)
	plans := g.planPods(houses)
	for i := 0; i < PodCount && i < len(plans); i++ {
		g.deckHouses[i] = plans[i].house
	}
	deck := Deck{
		Set:  set.Name,
		Seed: seed,
	}
	for i := 0; i < PodCount && i < len(plans); i++ {
		deck.Pods[i] = g.placeGiganticTutor(
			g.placeGiganticArt(g.expandPodClusters(g.fillPodPlan(plans[i]))))
	}
	g.expandClusters(&deck)
	g.expandFilteredClusters(&deck)
	deck.validate()
	g.applyEnhancements(&deck)
	return deck
}

// podPlan is a pod's fill recipe: its House and how the whole pod is drawn. An
// interloper pod draws every slot as a same-House card from another Set; an errant
// pod does the same but for a foreign House (one not native to this Set), which
// has replaced the pod's picked native House. A plain pod is neither and draws
// from this Set's own pool.
type podPlan struct {
	house      engine.House
	interloper bool
	errant     bool
}

// planPods turns the picked Houses into per-pod recipes, rolling the House-level
// overlays. Each pod rolls errant first, then interloper. The errant roll only
// fires when the Set has a legacy pool carrying a foreign House and a non-zero
// rate; the interloper roll only when the Set has a legacy pool and a non-zero
// rate — so a single-set build and every Tuning that leaves both rates at zero
// draw exactly the RNG stream they did before the overlays existed.
func (g *generator) planPods(houses []engine.House) []podPlan {
	plans := make([]podPlan, len(houses))
	t := g.set.Tuning
	errantOK := t.ErrantRate > 0 && g.set.legacy != nil && len(g.set.errantHouses) > 0
	interloperOK := t.InterloperRate > 0 && g.set.legacy != nil
	used := map[engine.House]bool{}
	for _, h := range houses {
		used[h] = true
	}
	for i, h := range houses {
		plans[i] = podPlan{house: h}
		if errantOK && g.chance(t.ErrantRate) {
			if fh, ok := g.pickErrantHouse(used); ok {
				used[fh] = true
				plans[i].house = fh
				plans[i].errant = true
				continue
			}
		}
		if interloperOK && g.chance(t.InterloperRate) {
			plans[i].interloper = true
		}
	}
	return plans
}

// pickErrantHouse picks a foreign House for an errant pod at random from the Set's
// errant Houses, excluding those already used by another pod so a deck's Houses
// stay distinct. It reports false when every errant House is taken.
func (g *generator) pickErrantHouse(used map[engine.House]bool) (engine.House, bool) {
	avail := make([]engine.House, 0, len(g.set.errantHouses))
	for _, h := range g.set.errantHouses {
		if !used[h] {
			avail = append(avail, h)
		}
	}
	if len(avail) == 0 {
		return engine.HouseNone, false
	}
	return avail[g.r.Intn(len(avail))], true
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
	return g.fillPodPlan(podPlan{house: house})
}

func (g *generator) fillPodPlan(plan podPlan) HousePod {
	pod := HousePod{House: plan.house}
	placed := make([]placedCard, 0, PodSize)
	for i := range PodSize {
		slot, pc := g.fillSlot(plan, placed)
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
		if ci.strategy == OnePerHouse || ci.strategy == PerGigantic ||
			!podClusterFires(pod, ci) {
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

// placeGiganticArt places the art half of every gigantic base half in a pod into
// a free slot of the same pod, so a pod that drew a gigantic's base half holds
// both halves (ADR 0042). The base carries its synthetic art half on its
// generation profile; deck generation is where the second physical card enters
// the deck. Each base spends one ordinary (non-half) slot on its art, so a pod
// with a gigantic holds one fewer other card.
func (g *generator) placeGiganticArt(pod HousePod) HousePod {
	for i := range pod.Slots {
		if pod.Slots[i].Card.GiganticRole != engine.GiganticBase {
			continue
		}
		art := g.set.byName[pod.Slots[i].Card.Name].Profile.GiganticArt
		if art == nil {
			continue
		}
		slot := freeGiganticSlot(&pod)
		if slot < 0 {
			continue
		}
		pod.Slots[slot] = Slot{
			Rarity: art.Rarity,
			Card:   *art,
		}
	}
	return pod
}

// freeGiganticSlot returns the first slot holding an ordinary card — one that is
// not a half of a gigantic — so placing an art half never clobbers a base half or
// an art half already placed. It returns -1 when every slot holds a gigantic half.
func freeGiganticSlot(pod *HousePod) int {
	for i := range pod.Slots {
		if pod.Slots[i].Card.GiganticRole == engine.GiganticNone {
			return i
		}
	}
	return -1
}

// placeGiganticTutor pulls one tutor into the pod of every gigantic base, chosen
// at random and stamped to the pod's House, so a deck that runs a gigantic also
// runs a tutor that fetches its halves (ADR 0044). It runs after placeGiganticArt,
// so each base has already spent one ordinary slot on its art; the tutor spends a
// second, and a pod with a gigantic holds base, art, and tutor. It gathers the
// free ordinary slots once, up front, so two gigantic bases in one pod each get a
// distinct slot rather than the second tutor clobbering the first. The tutors are
// a PerGigantic cluster read from the catalog pool, so any gigantic-bearing set
// reaches them and a new tutor joins every such set's pull with no per-set change.
func (g *generator) placeGiganticTutor(pod HousePod) HousePod {
	tutors := g.giganticTutors()
	if len(tutors) == 0 {
		return pod
	}
	bases := 0
	free := make([]int, 0, PodSize)
	for i := range pod.Slots {
		switch pod.Slots[i].Card.GiganticRole {
		case engine.GiganticBase:
			bases++
		case engine.GiganticNone:
			free = append(free, i)
		}
	}
	for b := 0; b < bases && b < len(free); b++ {
		tutor := tutors[g.r.Intn(len(tutors))]
		ctx := SlotContext{
			House:   pod.House,
			Rarity:  tutor.Def.Rarity,
			Special: true,
		}
		pod.Slots[free[b]] = Slot{
			Rarity: tutor.Def.Rarity,
			Card:   g.materialize(tutor, ctx),
		}
	}
	return pod
}

// giganticTutors gathers every member of the set's PerGigantic clusters into the
// one pool each gigantic's tutor pull chooses from.
func (g *generator) giganticTutors() []Card {
	var out []Card
	for _, ci := range g.set.perGiganticClusters() {
		out = append(out, ci.members...)
	}
	return out
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
	for _, s := range &pod.Slots {
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
// RandomCount cluster, by shuffling the members and taking that many. A ByLead
// cluster excludes its lead from the pick: the lead is already planted in the pod
// (its roll fired the cluster), so it is never placed again as one of the count —
// only its Connected partners are (Dark Harbinger pulls its Mutations, not itself).
func (g *generator) randomMembers(ci clusterIndex) []Card {
	pool := nonLeadMembers(ci)
	k := ci.min
	if ci.max > ci.min {
		k += g.r.Intn(ci.max - ci.min + 1)
	}
	shuffled := append([]Card(nil), pool...)
	g.r.Shuffle(len(shuffled), func(i, j int) {
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	})
	return shuffled[:k]
}

// nonLeadMembers returns a cluster's members with its lead removed, or all members
// when the cluster has no lead (a ByAnyMember cluster).
func nonLeadMembers(ci clusterIndex) []Card {
	if ci.lead == "" {
		return ci.members
	}
	out := make([]Card, 0, len(ci.members))
	for _, m := range ci.members {
		if m.Def.Name != ci.lead {
			out = append(out, m)
		}
	}
	return out
}

// selfPullCount rolls a SelfPull copy count: Min + Poisson(Mean − Min), capped at
// PodSize — at least Min, averaging about Mean, reaching a whole pod only on the
// thin tail.
func (g *generator) selfPullCount(ci clusterIndex) int {
	n := min(ci.min+poisson(g.r, ci.mean-float64(ci.min)), PodSize)
	return n
}

// pullCount rolls a Pull partner count from that partner's own rate: Min +
// Poisson(Mean − Min), capped at PodSize. A partner with Min 0 (Niffle Queen)
// often rolls none.
func (g *generator) pullCount(m ClusterMembership) int {
	n := min(m.Min+poisson(g.r, m.Mean-float64(m.Min)), PodSize)
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
			SlotContext{
				House:    pod.House,
				Rarity:   m.Def.Rarity,
				Maverick: maverick,
			},
		)
		if m.Profile.OneCopyPerDeck {
			g.placed[m.Def.Name] = true
		}
		pod.Slots[slot] = Slot{
			Rarity:   m.Def.Rarity,
			Maverick: maverick,
			Card:     def,
		}
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
	for _, s := range &pod.Slots {
		if s.Card.Name == name {
			n++
		}
	}
	return n
}

// expandClusters resolves the deck-wide cluster strategies once every pod is
// filled. Today that is OnePerHouse (the Shards): when any Shard has been drawn
// into the deck, each pod is given its House's Shard, so a single Shard pulls the
// whole cycle in. The clusters come from the attached catalog ClusterPool when one
// is set (so an errant pod's foreign House gets its Shard, drawn from another set),
// else from the set's own clusters. The pod-local strategies (WholePool,
// RandomCount) resolve in the pod pass, not here.
func (g *generator) expandClusters(deck *Deck) {
	for _, ci := range g.set.onePerHouseClusters() {
		if !g.clusterTriggered(deck, ci) {
			continue
		}
		for i := range PodCount {
			pod := &deck.Pods[i]
			if pod.House == engine.HouseNone {
				continue
			}
			// validateClusters (own) or validateCrossClusters (catalog) guarantees
			// every deckable House has a member, so the pod's House indexes one.
			g.ensureClusterMember(pod, ci, ci.byHouse[pod.House])
		}
	}
}

// clusterTriggered reports whether a cluster has fired: for ByAnyMember, that any
// member has been drawn into the deck; for ByLead, that the lead member has.
func (g *generator) clusterTriggered(deck *Deck, ci clusterIndex) bool {
	for i := range PodCount {
		for _, s := range &deck.Pods[i].Slots {
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
	for _, s := range &pod.Slots {
		if inCluster(ci, s.Card.Name) {
			return
		}
	}
	place, maverick := member, false
	if g.chance(g.set.Tuning.MaverickRate) {
		place, maverick = g.otherClusterMember(ci, member), true
	}
	ctx := SlotContext{
		House:    pod.House,
		Rarity:   place.Def.Rarity,
		Maverick: maverick,
	}
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
// by an optional duplicate-pull or a fresh (possibly maverick) draw. In an
// interloper or errant pod every slot draws from the legacy pool, so the legacy
// branch is forced rather than rolled.
func (g *generator) fillSlot(plan podPlan, placed []placedCard) (Slot, placedCard) {
	house := plan.house
	t := g.set.Tuning
	if len(g.set.special) > 0 && g.chance(t.SpecialRate) {
		if c, ok := g.pick(g.set.special); ok {
			return g.commit(c, SlotContext{
				House:   house,
				Rarity:  c.Def.Rarity,
				Special: true,
			})
		}
	}

	rarity := g.rollRarity()
	if plan.interloper || plan.errant || g.chance(t.LegacyRate) {
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
		return g.commit(c, SlotContext{
			House:  house,
			Rarity: rarity,
		})
	}

	maverick := g.chance(t.MaverickRate)
	c, ok := g.draw(house, rarity, maverick)
	if !ok {
		return Slot{Rarity: rarity}, placedCard{rarity: rarity}
	}
	ctx := SlotContext{
		House:    house,
		Rarity:   rarity,
		Maverick: maverick && c.Def.House != house,
	}
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
	return slot, placedCard{
		card:   c,
		rarity: ctx.Rarity,
	}
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
		def = engine.Rehouse(def, ctx.House)
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
