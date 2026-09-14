package anomalyexpansion

import "github.com/dmikalova/vactrol/internal/card"

// set is Anomaly Expansion's registrar: a reservoir set (ADR 0036). Every card in
// this package registers through set.New, so it is stamped as an Anomaly Expansion
// card and marked undraftable — the set builds no draw pool of its own, and its
// cards (the two Shards that complete the nine-House cycle) reach a deck only
// through the cross-set Shard cluster.
var set = card.ReservoirSet(card.AE)
