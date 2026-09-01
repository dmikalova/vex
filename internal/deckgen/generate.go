package deckgen

import (
	"math/rand"

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
	deck := Deck{Set: set.Name, Seed: seed}
	for i := 0; i < PodCount && i < len(houses); i++ {
		deck.Pods[i] = g.expandConnections(g.fillPod(houses[i]))
	}
	return deck
}

// generator threads the single RNG and the deck-wide "already placed" set (for
// one-copy-per-deck) through the pipeline.
type generator struct {
	set    Set
	r      *rand.Rand
	placed map[string]bool
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

// expandConnections resolves a pod's connections: every non-maverick puller card
// ensures its connected partners are present, overwriting other (unprotected)
// slots with them. Puller and partner slots are protected, and the pass repeats
// to a fixpoint so a pulled partner's own connection also resolves. Cross-house
// (maverick) connections are off for v1, so a maverick puller is skipped. The
// loop always terminates: each productive pass consumes one of the finitely many
// unprotected slots, so placePartners eventually places nothing.
func (g *generator) expandConnections(pod HousePod) HousePod {
	protected := make([]bool, PodSize)
	for {
		required := map[string]bool{}
		var order []string
		for i := 0; i < PodSize; i++ {
			if pod.Slots[i].Maverick {
				continue
			}
			c, ok := g.set.byName[pod.Slots[i].Card.Name]
			if !ok || c.Profile.Connection.Empty() {
				continue
			}
			protected[i] = true
			for _, name := range c.Profile.Connection.Cards {
				if !required[name] {
					required[name] = true
					order = append(order, name)
				}
			}
		}
		if len(order) == 0 {
			return pod
		}
		present := map[string]bool{}
		for i := 0; i < PodSize; i++ {
			if required[pod.Slots[i].Card.Name] {
				present[pod.Slots[i].Card.Name] = true
				protected[i] = true
			}
		}
		if !g.placePartners(&pod, order, present, protected) {
			return pod
		}
	}
}

// placePartners overwrites unprotected slots with the missing connected partners,
// protecting each as it lands. It reports whether it placed anything (a fixpoint
// pass that places nothing ends the loop).
func (g *generator) placePartners(
	pod *HousePod,
	order []string,
	present map[string]bool,
	protected []bool,
) bool {
	placed := false
	for _, name := range order {
		if present[name] {
			continue
		}
		partner := g.set.byName[name]
		slot := freeSlot(protected)
		if slot < 0 {
			break
		}
		def := g.materialize(partner, SlotContext{House: pod.House, Rarity: partner.Def.Rarity})
		pod.Slots[slot] = Slot{Rarity: partner.Def.Rarity, Card: def}
		protected[slot] = true
		placed = true
	}
	return placed
}

// freeSlot returns the first unprotected slot index, or -1 when the pod is full.
func freeSlot(protected []bool) int {
	for i := range protected {
		if !protected[i] {
			return i
		}
	}
	return -1
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
// Special card adopts the pod's House.
func (g *generator) materialize(c Card, ctx SlotContext) engine.CardDefinition {
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
// placed. It reports false when nothing is eligible.
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
	return elig[g.r.Intn(len(elig))], true
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
