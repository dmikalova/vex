// Package clusters holds cluster declarations shared across set packages, so a
// cross-set card family has one declaration whose strategy and trigger cannot
// drift between the members printed in different sets (ADR 0036).
package clusters

import "github.com/dmikalova/vactrol/internal/card"

// Shard is the Shards' cluster: one Shard per House, pulled into every House of a
// deck whenever any Shard is drawn. Age of Ascension prints seven (one per its
// Houses) and Anomaly Expansion the last two (Saurian, Star Alliance), so the
// cycle completes across every House a deck can reach — including a House an
// errant pod brings in from another set.
var Shard = card.Cluster{
	Name:     "Shard",
	Strategy: card.ClusterStrategy.OnePerHouse,
	Trigger:  card.ClusterTrigger.ByAnyMember,
}

// Ludo pairs Monument to Ludo (Mass Mutation) with Praefectus Ludo (Worlds
// Collide), printed in different sets, so the Monument's discard-pile bonus has
// its namesake to feed it (ADR 0036). The Monument leads and pulls the Praefectus.
var Ludo = card.Cluster{
	Name:     "Monument to Ludo",
	Strategy: card.ClusterStrategy.Pull,
	Trigger:  card.ClusterTrigger.ByLead,
}
