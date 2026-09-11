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
	// Connection is the set of connected cards a puller card brings into its pod;
	// build one with card.Connects.
	Connection = deckgen.Connection
	// ConnectedCard is one card a connection pulls; build one with card.Pull or
	// card.PullSometimes.
	ConnectedCard = deckgen.ConnectedCard
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
// at least min copies when the lead rolls in, averaging about mean, on a Poisson
// tail (min 0 pulls none most of the time). Hand it to card.InCluster on each
// pulled partner; the lead carries the plain cluster via card.LeadsCluster.
func Pulled(c Cluster, min int, mean float64) Cluster {
	c.Min, c.Mean = min, mean
	return c
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

// RarityWeight scales how often deck generation draws this card among its
// house+rarity peers, relative to the default weight of 1 — the card-level
// companion to a Set's Tuning.RarityWeights. A fraction makes a card rarer within
// its rarity: the five Master-of-N variants each carry card.RarityWeight(0.2), so
// the family drafts about as often as one ordinary Rare card.
func RarityWeight(w float64) Option { return func(b *builder) { b.profile.RarityWeight = w } }

// Connects marks the card as a connection puller: when it is placed in a pod,
// deck generation pulls the given cards into that pod. Each pull is built with
// card.Pull or card.PullSometimes, which name the connected card by its own
// definition symbol, so a connection to a card that does not exist is a compile
// error rather than a silently dropped link.
//
//	card.Connects(
//	  card.Pull(NiffleApe, 2),
//	  card.PullSometimes(NiffleQueen, 0.15),
//	)
//
// A card that should never roll without its puller is additionally authored with
// card.Rarity.Connected, which keeps it out of the pool; an ordinary card can be
// pulled too, and still rolls on its own.
func Connects(cards ...ConnectedCard) Option {
	return func(b *builder) { b.profile.Connection = deckgen.Connection{Cards: cards} }
}

// Pull is one connected card brought into the pod every time, in the given number
// of copies.
func Pull(c Definition, copies int) ConnectedCard {
	return ConnectedCard{Name: c.Name, Copies: copies, Chance: 1}
}

// PullExact is one connected card brought into the pod once per copy of the
// puller, so N pullers guarantee N partners — e.g. two Timetravellers pull two
// Help From Future Self. Use it where an ordinary Pull would collapse several
// pullers onto a single shared partner.
func PullExact(c Definition, copies int) ConnectedCard {
	return ConnectedCard{Name: c.Name, Copies: copies, Chance: 1, Exact: true}
}

// PullSometimes is one connected card brought into the pod with the given
// probability, rolled once per pod. It is how a card guarantees a flourish
// without guaranteeing it every deck.
func PullSometimes(c Definition, chance float64) ConnectedCard {
	return ConnectedCard{Name: c.Name, Copies: 1, Chance: chance}
}
