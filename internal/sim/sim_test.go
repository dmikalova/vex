package sim

import (
	"math/rand"
	"os"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// FuzzPlay drives a whole random legal game decoded from the fuzz input and fails
// on any invariant violation or engine panic. Go's coverage-guided mutator turns
// the byte script into a smart explorer of the game tree; a failing input is
// minimized and saved under testdata/fuzz/FuzzPlay as a permanent regression.
//
// Run it with the assert build tag so the engine's in-game checks fire too:
//
//	mage fuzz            # go test -tags assert -fuzz=FuzzPlay ./internal/sim
func FuzzPlay(f *testing.F) {
	f.Add([]byte{})
	f.Add([]byte{0, 0, 0, 0, 0, 0, 0, 1})                   // seed 1, then wind down
	f.Add([]byte{0, 0, 0, 0, 0, 0, 0, 7, 1, 2, 3, 4, 5, 6}) // seed 7, a few decisions
	f.Fuzz(func(t *testing.T, script []byte) {
		if err := Simulate(script); err != nil {
			t.Fatalf("script %x failed: %v", script, err)
		}
	})
}

// TestSimulateSeeds is the fixed-seed property test wired into `mage test`: it
// plays a deterministic batch of random games so every run of the suite shakes the
// engine, while staying fast enough for the gate. Because the seed is fixed, a
// regression here is reproducible from the printed script. It then plays a second
// batch from a fresh, non-deterministic source, so each run also explores games
// the fixed batch never touches — a failure there is still replayable from the
// printed hex under `mage debug`.
func TestSimulateSeeds(t *testing.T) {
	for i, script := range SeedScripts(1000) {
		if err := Simulate(script); err != nil {
			t.Fatalf("fixed batch %d, script %x failed: %v", i, script, err)
		}
	}
	for i, script := range RandomScripts(1000) {
		if err := Simulate(script); err != nil {
			t.Fatalf("random batch %d, script %x failed: %v", i, script, err)
		}
	}
}

// TestSoak is the long-running soak, skipped unless SOAK_DURATION is set (see
// `mage soak`). It churns fresh random games across GOMAXPROCS workers until the
// time budget runs out, and — unlike the fixed-seed property test — it does not
// stop at the first failure: every failing script is saved into the FuzzPlay
// corpus (see SaveCorpus) so a soak find becomes a permanent regression, then the
// soak keeps hunting. Build with -tags assert to run the in-engine checks too.
func TestSoak(t *testing.T) {
	budget := os.Getenv("SOAK_DURATION")
	if budget == "" {
		t.Skip("set SOAK_DURATION (e.g. 30s) to run the soak")
	}
	dur, err := time.ParseDuration(budget)
	if err != nil {
		t.Fatalf("invalid SOAK_DURATION %q: %v", budget, err)
	}

	workers := runtime.GOMAXPROCS(0)
	deadline := time.Now().Add(dur)
	var games, failures atomic.Int64
	var mu sync.Mutex // serialize corpus writes and test logging across workers
	var wg sync.WaitGroup
	for w := 0; w < workers; w++ {
		wg.Add(1)
		go func(seed int64) {
			defer wg.Done()
			r := rand.New(rand.NewSource(seed))
			for time.Now().Before(deadline) {
				script := make([]byte, 8+r.Intn(2000))
				r.Read(script)
				games.Add(1)
				if err := Simulate(script); err != nil {
					failures.Add(1)
					mu.Lock()
					if name, saveErr := SaveCorpus(CorpusDir, script, err); saveErr != nil {
						t.Errorf(
							"script %x failed: %v (could not save corpus: %v)",
							script,
							err,
							saveErr,
						)
					} else {
						t.Errorf("script failed: %v\n  saved to %s", err, name)
					}
					mu.Unlock()
				}
			}
		}(time.Now().UnixNano() + int64(w))
	}
	wg.Wait()
	t.Logf("soak completed %d games across %d workers in %s (%d failures)",
		games.Load(), workers, dur, failures.Load())
}
