package sim

import (
	"math/rand"
	"testing"

	"github.com/dmikalova/vactrol/internal/engine"
	"github.com/dmikalova/vactrol/internal/match"
)

// The benchmarks in this file profile the engine under a high-volume driver — the
// closest proxy the repo has for the load an MCTS bot would put on it. Two shapes:
// replaying the deterministic seeded batch for the realistic action and allocation
// mix, and snapshotting then restoring the flat GameState for the copy cost MCTS
// pays at every tree node. They deliberately skip the per-action invariant audits
// Simulate runs: those are the simulation harness's own cost, not the engine's, and
// MCTS never pays them. `mage profile` runs these and diffs them against a committed
// baseline; see docs/testing.md.

// benchSeedBatch is the batch each play benchmark replays per op, matching the
// TestSimulateSeeds workload. One op is the whole batch, so allocs/op and B/op are
// the batch's totals. The scripts are fixed, but the engine randomizes map-
// iteration order per process, so the counts are stable only to a small run-to-run
// jitter (~1-2%) — enough for `mage profile` to track them as an advisory baseline
// where a real regression stands out well beyond the jitter, not an exact gate.
const benchSeedBatch = 1000

// sink keeps a copied state from being optimized away.
var sink engine.GameState

// BenchmarkPlaySeeds replays the 1000 deterministic seeded games per op. It is the
// workload `mage profile` profiles and baselines: the scripts are fixed, so its
// per-game counts are independent of CPU speed and stable to a small per-process
// jitter (Go randomizes map iteration, nudging a few script-indexed choices).
func BenchmarkPlaySeeds(b *testing.B) {
	benchmarkPlay(b, SeedScripts(benchSeedBatch))
}

// BenchmarkPlayRandom replays a fresh random batch per op, for the `mage profile
// -random` outlier pass. Its counts vary run to run, so it is never baselined.
func BenchmarkPlayRandom(b *testing.B) {
	benchmarkPlay(b, RandomScripts(benchSeedBatch))
}

// benchmarkPlay plays every script in the batch per op and reports actions/game
// alongside the standard ns/op, B/op, and allocs/op.
func benchmarkPlay(b *testing.B, scripts [][]byte) {
	b.ReportAllocs()
	b.ResetTimer()
	var actions int64
	for i := 0; i < b.N; i++ {
		for _, s := range scripts {
			actions += int64(playForBench(s))
		}
	}
	b.StopTimer()
	games := float64(b.N) * float64(len(scripts))
	b.ReportMetric(float64(actions)/games, "actions/game")
}

// BenchmarkFastCopy measures one undo/MCTS snapshot: a plain value copy of the flat
// GameState. GameState is fixed-size and pointerless (ADR 0005), so the copy is
// allocation-free and its cost tracks the struct's size — the number that grows
// every time a field is added. allocs/op must stay 0; the baseline guards it.
func BenchmarkFastCopy(b *testing.B) {
	g, _, _ := benchMidgame()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sink = g.State.FastCopy()
	}
}

// BenchmarkSnapshotApplyRestore measures the MCTS inner loop around a real engine
// action: snapshot the state, apply one legal action, restore. Splitting it from
// BenchmarkFastCopy separates "the snapshot got more expensive" (state grew) from
// "the action got more expensive" when a number moves.
func BenchmarkSnapshotApplyRestore(b *testing.B) {
	g, player, ply := benchMidgame()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		snap := g.State.FastCopy()
		d := &decoder{script: ply}
		_, _ = doAction(g, player, d)
		g.State = snap
	}
}

// playForBench plays one script to completion and returns the number of actions
// applied, skipping the per-step invariant checks Simulate runs. The play
// benchmarks measure the engine's own play cost as an MCTS-rollout proxy, and the
// harness's invariant audits are neither the engine's cost nor something MCTS pays.
// It mirrors simulate's turn scaffolding minus those checks and reuses doAction.
func playForBench(script []byte) int {
	g, houses, d := benchGame(script)
	actions := 0
	for turn := 0; turn < maxTurns && g.Winner() < 0 && !d.done(); turn++ {
		player := turn % 2
		if turn == 0 {
			g.StartGame(player)
		} else {
			g.StartTurn(player)
		}
		if g.Winner() >= 0 {
			break
		}
		if len(houses[player]) > 0 {
			_ = g.ChooseHouse(player, houses[player][int(d.byte())%len(houses[player])])
		}
		for step := 0; step < maxDecisionsPerTurn; step++ {
			if g.Winner() >= 0 || d.done() {
				break
			}
			acted, err := doAction(g, player, d)
			if err != nil || !acted {
				break
			}
			actions++
		}
		g.EndPlayPhase(player)
	}
	return actions
}

// benchScript is a fixed, deliberately long script for the snapshot micro-
// benchmarks: long enough to play several full turns onto a populated board and
// still leave bytes to drive one more action. Seed 42 keeps the position identical
// across runs, so the micro-benchmark counts stay reproducible.
func benchScript() []byte {
	r := rand.New(rand.NewSource(42))
	s := make([]byte, 2000)
	_, _ = r.Read(s)
	return s
}

// benchMidgame replays benchScript through several full turns, then starts one more
// turn and chooses its house, leaving the game paused mid-turn on a populated board.
// It returns the game, the active player, and the still-unconsumed script bytes that
// drive one further action — the position the snapshot micro-benchmarks copy from.
func benchMidgame() (*engine.Game, int, []byte) {
	const warmupTurns = 8
	script := benchScript()
	g, houses, d := benchGame(script)
	for turn := 0; turn < warmupTurns && g.Winner() < 0 && !d.done(); turn++ {
		player := turn % 2
		if turn == 0 {
			g.StartGame(player)
		} else {
			g.StartTurn(player)
		}
		if g.Winner() >= 0 {
			break
		}
		_ = playTurn(g, player, houses[player], d)
		g.EndPlayPhase(player)
	}
	player := warmupTurns % 2
	g.StartTurn(player)
	if len(houses[player]) > 0 {
		_ = g.ChooseHouse(player, houses[player][int(d.byte())%len(houses[player])])
	}
	return g, player, append([]byte(nil), script[d.pos:]...)
}

// benchGame decodes a script's seed and builds the match with the script chooser
// wired to both players, the shared setup the bench drivers start from.
func benchGame(script []byte) (*engine.Game, [2][]engine.House, *decoder) {
	d := &decoder{script: script}
	seed := int64(d.uint64())
	g, houses := match.New("P0", "P1", seed)
	ch := &scriptChooser{d: d}
	g.SetChooser(0, ch)
	g.SetChooser(1, ch)
	g.SetPlayerHouses(0, houses[0])
	g.SetPlayerHouses(1, houses[1])
	return g, houses, d
}
