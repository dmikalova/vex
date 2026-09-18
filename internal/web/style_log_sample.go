package web

import (
	"reflect"
	"sort"

	"github.com/dmikalova/vactrol/internal/engine"
	"github.com/dmikalova/vactrol/internal/sim"
)

// This file samples real played games against the engine's log-entry catalog, so
// the style gallery can show every kind of log line as it appears in an actual
// game rather than as a constructed specimen. See
// docs/adr/0046-log-entry-catalog-and-sampled-gallery.md: the catalog
// (engine.LogEntrySamples) is the target set, a sampled game supplies real
// bubbles, and a greedy set-cover picks the fewest bubbles that show every kind a
// run produced. Kinds no game produced are left for the synthetic fallback the
// gallery draws from the catalog instance itself.

// logKind identifies a catalogued log-entry variant by its Go type — the same
// grain the catalog's totality test enumerates, so a kind covered here is a kind
// the gallery would otherwise have to draw synthetically.
type logKind = reflect.Type

// catalogKinds is the set of log-entry types engine.LogEntrySamples enumerates,
// the kinds a sampled run is covered against.
func catalogKinds() map[logKind]bool {
	kinds := map[logKind]bool{}
	for _, e := range engine.LogEntrySamples() {
		kinds[reflect.TypeOf(e)] = true
	}
	return kinds
}

// blockKinds returns the catalogued kinds one drawn log block contains: its
// header if it rules a line, else the types of its bubble's records. It
// intersects with want so an entry outside the catalog (there are none today)
// never widens the target set.
func blockKinds(b logBlock, want map[logKind]bool) map[logKind]bool {
	got := map[logKind]bool{}
	add := func(e engine.LogEntry) {
		if t := reflect.TypeOf(e); want[t] {
			got[t] = true
		}
	}
	if b.header != nil {
		add(b.header)
	}
	for _, rec := range b.lines {
		add(rec.Entry)
	}
	return got
}

// sampledGame is a played-out game retained for the log gallery: the script that
// reproduces it and the web game wrapping it, so its blocks render through the
// production logBlockView and a drill-in shows its whole log.
type sampledGame struct {
	script []byte
	game   *game
}

// coverBubble points at one bubble in the set-cover: which retained game it came
// from and which of that game's blocks it is, so the gallery renders the block and
// a click opens the game's full log with the block highlighted.
type coverBubble struct {
	game  int
	block int
}

// logCoverage is the result of sampling games against the catalog.
type logCoverage struct {
	// games are the retained games, each contributing a bubble the cover uses;
	// cover and drill-in index into this slice.
	games []sampledGame
	// cover is the greedy set-cover: the fewest retained bubbles whose union shows
	// every observed kind.
	cover []coverBubble
	// unobserved are the catalogued kinds no sampled game produced, left for the
	// synthetic fallback, sorted for a stable gallery.
	unobserved []logKind
	// played is how many games were played, for the "not observed in N games" badge.
	played int
}

// sampleLog plays the given scripts against the catalog, stopping as soon as every
// catalogued kind has been observed or the scripts run out. It retains only games
// that introduce a kind not yet seen, computes a greedy set-cover of their
// bubbles, then prunes to the games that cover actually uses. It is a pure
// function of the scripts, so the same batch reproduces the same gallery.
func sampleLog(scripts [][]byte) logCoverage {
	want := catalogKinds()
	observed := map[logKind]bool{}
	cov := logCoverage{}
	for _, script := range scripts {
		if len(observed) == len(want) {
			break
		}
		g, err := sim.Play(script)
		cov.played++
		if err != nil || g == nil {
			continue
		}
		gw := &game{selHand: -1, zonesPlayer: -1, forgingKey: -1, handSlot: -1}
		gw.g = g
		if !foldsNewKind(gw, want, observed) {
			continue
		}
		cov.games = append(cov.games, sampledGame{script: script, game: gw})
	}
	cov.cover = coverBubbles(cov.games, observed, want)
	for t := range want {
		if !observed[t] {
			cov.unobserved = append(cov.unobserved, t)
		}
	}
	sort.Slice(cov.unobserved, func(i, j int) bool {
		return cov.unobserved[i].String() < cov.unobserved[j].String()
	})
	return cov.pruned()
}

// foldsNewKind reports whether a game's blocks contain a catalogued kind not yet
// in observed, folding every kind the game produced into observed as it checks.
func foldsNewKind(gw *game, want, observed map[logKind]bool) bool {
	fresh := false
	for _, b := range gw.logBlocks() {
		for t := range blockKinds(b, want) {
			if !observed[t] {
				fresh = true
			}
			observed[t] = true
		}
	}
	return fresh
}

// coverBubbles greedily picks the fewest blocks across the retained games whose
// kinds together cover every observed kind. At each step it takes the block adding
// the most still-uncovered kinds, breaking ties by game then block order so the
// cover is deterministic.
func coverBubbles(games []sampledGame, observed, want map[logKind]bool) []coverBubble {
	type candidate struct {
		bubble coverBubble
		kinds  map[logKind]bool
	}
	var cands []candidate
	for gi, sg := range games {
		for bi, b := range sg.game.logBlocks() {
			if ks := blockKinds(b, want); len(ks) > 0 {
				cands = append(cands, candidate{coverBubble{gi, bi}, ks})
			}
		}
	}
	uncovered := map[logKind]bool{}
	for t := range observed {
		uncovered[t] = true
	}
	var cover []coverBubble
	for len(uncovered) > 0 {
		best, bestGain := -1, 0
		for i, c := range cands {
			gain := 0
			for t := range c.kinds {
				if uncovered[t] {
					gain++
				}
			}
			if gain > bestGain {
				best, bestGain = i, gain
			}
		}
		if best < 0 {
			break // no block advances coverage; the rest is genuinely unreachable
		}
		for t := range cands[best].kinds {
			delete(uncovered, t)
		}
		cover = append(cover, cands[best].bubble)
	}
	return cover
}

// pruned drops retained games no cover bubble references and reindexes the cover
// onto the games kept, so the coverage holds only the games it actually shows.
func (cov logCoverage) pruned() logCoverage {
	used := map[int]bool{}
	for _, cb := range cov.cover {
		used[cb.game] = true
	}
	remap := map[int]int{}
	var kept []sampledGame
	for gi := range cov.games {
		if used[gi] {
			remap[gi] = len(kept)
			kept = append(kept, cov.games[gi])
		}
	}
	cover := make([]coverBubble, len(cov.cover))
	for i, cb := range cov.cover {
		cover[i] = coverBubble{remap[cb.game], cb.block}
	}
	cov.games, cov.cover = kept, cover
	return cov
}
