package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// horsemenCluster binds the four Horsemen into one WholePool cluster led by the
// Horseman of Pestilence: whenever Pestilence rolls into a pod, the other three
// ride in with it (ADR 0036). Death, Famine, and War are Rarity.Connected, so the
// only way they reach a deck is on Pestilence's coattails.
var horsemenCluster = card.Cluster{
	Name:     "Horseman",
	Strategy: card.ClusterStrategy.WholePool,
	Trigger:  card.ClusterTrigger.ByLead,
}
