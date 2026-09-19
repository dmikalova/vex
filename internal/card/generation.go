package card

import (
	"math/rand"

	"github.com/dmikalova/vactrol/internal/deckgen"
)

// Deck-generation metadata re-exported for authoring. These attach to a card the
// same way Provenance does — recorded in the registry, not on the engine
// definition — and are read by the deck generator (see internal/deckgen and
// ADR 0004).
type (
	// Materializer turns a card template into a concrete card at generation time.
	Materializer = deckgen.Materializer
	// SlotContext is what a Materializer is given: the pod's house and the slot's
	// rolled rarity and provenance flags.
	SlotContext = deckgen.SlotContext
	// Cluster declares a card family deck generation places together (ADR 0036):
	// its name, the strategy that fills it, and the trigger that fires it. Declare
	// one shared value and hand it to every member with card.InCluster, so all
	// members carry identical placement rules that cannot drift.
	//
	//	var shardCluster = card.Cluster{
	//	  Name:     "Shard",
	//	  Strategy: card.ClusterStrategy.OnePerHouse,
	//	  Trigger:  card.ClusterTrigger.ByAnyMember,
	//	}
	Cluster = deckgen.ClusterMembership
	// FilteredCluster is a deck-wide pull keyed to a lead card that fills by a
	// predicate rather than named members (ADR 0036): when the lead is in a deck,
	// generation guarantees at least Floor cards matching Match, in any House. Mark
	// the lead with card.PullsMatching rather than constructing this directly.
	FilteredCluster = deckgen.FilteredCluster
)

// ClusterStrategy groups the cluster fill strategies, e.g.
// card.ClusterStrategy.OnePerHouse.
var ClusterStrategy = clusterStrategies{
	WholePool:   deckgen.WholePool,
	RandomCount: deckgen.RandomCount,
	SelfPull:    deckgen.SelfPull,
	PullExact:   deckgen.PullExact,
	Pull:        deckgen.Pull,
	OnePerHouse: deckgen.OnePerHouse,
	PerGigantic: deckgen.PerGigantic,
}

type clusterStrategies struct {
	// WholePool places every member of the cluster (the four Horsemen).
	WholePool deckgen.ClusterStrategy
	// RandomCount places a random number of members in [Cluster.Min, Cluster.Max]
	// (the seven sins).
	RandomCount deckgen.ClusterStrategy
	// SelfPull places a random number of copies of the single triggering member
	// itself — Cluster.Min at least, averaging about Cluster.Mean, with a thin tail
	// to a full pod (Plague Rat pulls more Plague Rats).
	SelfPull deckgen.ClusterStrategy
	// PullExact places one copy of each non-lead member per lead instance in the
	// pod (two Timetravellers pull two Help from Future Self). It is always ByLead.
	PullExact deckgen.ClusterStrategy
	// Pull places a per-partner random count of each non-lead member when the lead
	// rolls in (Troop Call pulls Niffle Apes, less often a Niffle Queen). It is
	// always ByLead; each pulled partner sets its own rate with card.Pulled.
	Pull deckgen.ClusterStrategy
	// OnePerHouse places one member in each of the deck's Houses; it is deck-wide
	// and gated complete-by-construction — every House must have a member (the
	// Shards).
	OnePerHouse deckgen.ClusterStrategy
	// PerGigantic places one member, chosen at random, into each gigantic base's
	// pod, stamped to that pod's House. Its members are Houseless reservoir cards
	// pulled only by a gigantic, never drawn on their own (the tutors).
	PerGigantic deckgen.ClusterStrategy
}

// ClusterTrigger groups the cluster triggers, e.g.
// card.ClusterTrigger.ByAnyMember.
var ClusterTrigger = clusterTriggers{
	ByLead:      deckgen.ByLead,
	ByAnyMember: deckgen.ByAnyMember,
}

type clusterTriggers struct {
	// ByLead fires the cluster only when its lead member is placed (the Horseman
	// that pulls the others).
	ByLead deckgen.ClusterTrigger
	// ByAnyMember fires the cluster when any member is placed (any Shard, any sin).
	ByAnyMember deckgen.ClusterTrigger
}

// InCluster marks the card a member of the given cluster, copying the cluster's
// strategy and trigger onto it so deck generation can resolve the whole family
// from any one member.
func InCluster(c Cluster) Option {
	return func(b *builder) { b.profile.Cluster = c }
}

// LeadsCluster marks the card the lead member of a ByLead cluster — the Horseman
// whose placement pulls the rest — in addition to making it a member.
func LeadsCluster(c Cluster) Option {
	c.Lead = true
	return func(b *builder) { b.profile.Cluster = c }
}

// Pulled returns a copy of a Pull cluster carrying this partner's own pull rate:
// at least minCopies when the lead rolls in, averaging about mean, on a Poisson
// tail (minCopies 0 pulls none most of the time). Hand it to card.InCluster on
// each pulled partner; the lead carries the plain cluster via card.LeadsCluster.
func Pulled(c Cluster, minCopies int, mean float64) Cluster {
	c.Min, c.Mean = minCopies, mean
	return c
}

// PullsMatching marks this card the lead of a deck-wide filtered pull (ADR 0036):
// whenever it is in a generated deck, generation guarantees at least floor cards
// matching match somewhere in the deck, in any House, topping up from the pool.
// Each deck rolls a target of floor + Poisson(mean − floor), so the count is at
// least floor and averages about mean; pass mean equal to floor for a flat floor.
// Cards already in the deck that match count toward the target. Use it for a
// deck-building payoff keyed to a card family the card cares about — Chief Engineer
// Walls guaranteeing the Upgrades and Robots it retrieves. It is independent of
// card.InCluster, so a card can lead a filtered pull and belong to a named cluster.
func PullsMatching(name string, floor int, mean float64, match func(Definition) bool) Option {
	fc := FilteredCluster{
		Name:  name,
		Floor: floor,
		Mean:  mean,
		Match: match,
	}
	return func(b *builder) { b.profile.Leads = &fc }
}

// MaterializeFunc adapts a plain function to a Materializer, so a template can be
// written inline: card.Template(func(ctx card.SlotContext, r *rand.Rand) card.Definition { ... }).
type MaterializeFunc func(SlotContext, *rand.Rand) Definition

// Materialize satisfies the Materializer interface.
func (f MaterializeFunc) Materialize(ctx SlotContext, r *rand.Rand) Definition { return f(ctx, r) }

// Template marks a card as a template: a single pool entry whose concrete cards
// deck generation produces from f. The card.New face (house, type, rarity) buckets
// the template in the pool; f builds the playable card.
func Template(f MaterializeFunc) Option { return func(b *builder) { b.materializer = f } }

// OneCopyPerDeck bars deck generation from placing more than one copy of the card.
func OneCopyPerDeck() Option { return func(b *builder) { b.profile.OneCopyPerDeck = true } }

// Houseless marks a Special card that carries no House of its own until deck
// generation places it, when it adopts the House of the pod it fills (ADR 0004).
// A houseless Special enters a deck through the special slot, so it can appear in
// any deck whatever its houses; author it with card.House.None. Dark Æmber Vault
// is the model.
func Houseless() Option { return func(b *builder) { b.profile.Houseless = true } }

// RarityWeight scales how often deck generation draws this card among its
// house+rarity peers, relative to the default weight of 1 — the card-level
// companion to a Set's Tuning.RarityWeights. A fraction makes a card rarer within
// its rarity: the five Master-of-N variants each carry card.RarityWeight(0.2), so
// the family drafts about as often as one ordinary Rare card.
func RarityWeight(w float64) Option { return func(b *builder) { b.profile.RarityWeight = w } }
