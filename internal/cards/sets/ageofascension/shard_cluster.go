package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// shardCluster is the Shards' cluster: one Shard per House, pulled in whenever any
// Shard is drawn (ADR 0036). Each Shard joins it with card.InCluster, so the
// generator resolves the whole cluster from any one Shard. The cluster is complete
// by construction — every House Age of Ascension can deck has a Shard — so the
// deckgen build gates on it.
var shardCluster = card.Cluster{
	Name:     "Shard",
	Strategy: card.ClusterStrategy.OnePerHouse,
	Trigger:  card.ClusterTrigger.ByAnyMember,
}
