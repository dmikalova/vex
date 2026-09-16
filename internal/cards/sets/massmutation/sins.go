package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// sinsCluster binds the seven deadly Sins into one RandomCount cluster triggered by
// any member: whenever any Sin is drawn into a deck, three to seven distinct Sins
// ride in together (ADR 0036), so a deck that reaches for one temptation gets a
// connected family of them — an average of five, never fewer than three, never more
// than the seven that exist. Each Sin is OneCopyPerDeck, so a deck holds at most one
// of each, and the whole set reinforces itself: every Sin scales with the friendly
// Sin creatures beside it.
var sinsCluster = card.Cluster{
	Name:     "Sins",
	Strategy: card.ClusterStrategy.RandomCount,
	Trigger:  card.ClusterTrigger.ByAnyMember,
	Min:      3,
	Max:      7,
}
