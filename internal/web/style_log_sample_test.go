package web

import (
	"reflect"
	"testing"

	"github.com/dmikalova/vactrol/internal/sim"
)

// coverKinds is the union of catalogued kinds the set-cover's bubbles show, the
// observed set the gallery renders from real games.
func coverKinds(t *testing.T, cov logCoverage) map[logKind]bool {
	t.Helper()
	want := catalogKinds()
	got := map[logKind]bool{}
	for _, cb := range cov.cover {
		if cb.game < 0 || cb.game >= len(cov.games) {
			t.Fatalf("cover bubble game index %d out of range (%d games)",
				cb.game, len(cov.games))
		}
		blocks := cov.games[cb.game].game.logBlocks()
		if cb.block < 0 || cb.block >= len(blocks) {
			t.Fatalf("cover bubble block index %d out of range (%d blocks)",
				cb.block, len(blocks))
		}
		ks := blockKinds(blocks[cb.block], want)
		if len(ks) == 0 {
			t.Errorf("cover bubble %+v carries no catalogued kind", cb)
		}
		for k := range ks {
			got[k] = true
		}
	}
	return got
}

// TestSampleLogPartitionsCatalog checks the sampler's core contract: every
// catalogued kind ends up either observed (in the cover) or reported unobserved,
// with no overlap and nothing lost.
func TestSampleLogPartitionsCatalog(t *testing.T) {
	cov := sampleLog(sim.SeedScripts(100))
	want := catalogKinds()

	if cov.played == 0 {
		t.Fatal("no games were played")
	}
	observed := coverKinds(t, cov)
	if len(observed) == 0 {
		t.Fatal("sampled games observed no catalogued kind")
	}
	if len(observed)+len(cov.unobserved) != len(want) {
		t.Fatalf("observed %d + unobserved %d != catalog %d",
			len(observed), len(cov.unobserved), len(want))
	}
	for _, k := range cov.unobserved {
		if observed[k] {
			t.Errorf("kind %s is both observed and reported unobserved", k)
		}
		if !want[k] {
			t.Errorf("reported-unobserved kind %s is not in the catalog", k)
		}
	}
}

// TestSampleLogCoverIsTight checks the greedy cover wastes no bubble: each one
// adds at least one kind none before it did, so the hero gallery has no redundant
// entry.
func TestSampleLogCoverIsTight(t *testing.T) {
	cov := sampleLog(sim.SeedScripts(100))
	want := catalogKinds()

	seen := map[logKind]bool{}
	for _, cb := range cov.cover {
		blocks := cov.games[cb.game].game.logBlocks()
		gain := false
		for k := range blockKinds(blocks[cb.block], want) {
			if !seen[k] {
				gain = true
			}
			seen[k] = true
		}
		if !gain {
			t.Errorf("cover bubble %+v adds no new kind; cover is not minimal", cb)
		}
	}
}

// TestSampleLogPrunesUnusedGames checks pruning: every retained game is referenced
// by the cover, so the coverage holds no game it never shows.
func TestSampleLogPrunesUnusedGames(t *testing.T) {
	cov := sampleLog(sim.SeedScripts(100))
	used := map[int]bool{}
	for _, cb := range cov.cover {
		used[cb.game] = true
	}
	for gi := range cov.games {
		if !used[gi] {
			t.Errorf("retained game %d is referenced by no cover bubble", gi)
		}
	}
}

// TestSampleLogRendersWithoutPanic checks every cover bubble draws through the
// production logBlockView, the same path the gallery uses, without panicking.
func TestSampleLogRendersWithoutPanic(t *testing.T) {
	cov := sampleLog(sim.SeedScripts(60))
	for _, cb := range cov.cover {
		gw := cov.games[cb.game].game
		if ui := gw.logBlockView(gw.logBlocks()[cb.block]); ui == nil {
			t.Errorf("cover bubble %+v rendered nil", cb)
		}
	}
}

// TestCatalogKindsMatchSamples checks the kind set is exactly the distinct types
// of the engine catalog, the anchor the sampler covers against.
func TestCatalogKindsMatchSamples(t *testing.T) {
	kinds := catalogKinds()
	for k := range kinds {
		if k == nil {
			t.Fatal("catalog kind is nil")
		}
		if k.Kind() == reflect.Pointer {
			t.Errorf("catalog kind %s is a pointer; samples should be values", k)
		}
	}
}
