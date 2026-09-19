//go:build mage

package main

import (
	"bufio"
	jsontext "encoding/json/jsontext"
	"encoding/json/v2"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/magefile/mage/sh"
)

// Fuzz runs the whole-game fuzzer with assertions on. It drives
// internal/sim.FuzzPlay under the assert build tag, explores the game tree with
// coverage-guided mutation, and saves any minimized failing script under
// internal/sim/testdata/fuzz/FuzzPlay.
//
// The budget defaults to 60s. Pass -fuzztime a Go duration, or an "Nx" execution
// count that a duration cannot express:
//
//	mage fuzz
//	mage fuzz -fuzztime=5m
//	mage fuzz -fuzztime=10000x
//
// The generated corpus lives in the Go build cache, not the repo — run
// `mage fuzzClean` to reset it.
func Fuzz(fuzztime *string) error {
	ft := "60s"
	if fuzztime != nil && *fuzztime != "" {
		ft = *fuzztime
	}
	return sh.RunV("go", "test", "-tags", "assert",
		"-run", "^$", "-fuzz", "^FuzzPlay$", "-fuzztime", ft,
		"./internal/sim")
}

// FuzzClean discards the generated fuzz corpus. It lives in the Go build cache.
// Reach for this when a stale corpus keeps steering the fuzzer down the same
// paths. The checked-in seed corpus under internal/sim/testdata is untouched.
func FuzzClean() error {
	return sh.RunV("go", "clean", "-fuzzcache")
}

// CorpusPrune prunes FuzzPlay's seed corpus of fixed bugs. Each bug that still
// reproduces is left as one minimized entry. The soak saves every failing
// script verbatim, so a single bug can leave hundreds of near-identical
// multi-kilobyte entries; this replays them all, drops the ones whose bug is
// fixed, and shrinks what remains.
func CorpusPrune() error {
	return sh.RunV("go", "run", "./magefiles/simcorpus")
}

// Debug replays a failing simulated game with the game log on. It prints the log
// tail next to the invariant violation that ended the game, turning a soak, fuzz,
// or property-test find into a readable sequence of plays.
//
// With no -script it searches the fixed-seed property batch that `mage test` plays
// and replays the first failure; pass -script the hex a failure printed to replay
// that one instead. -tail widens the log tail (default 60 lines):
//
//	mage debug
//	mage debug -script=00ff1a...
//	mage debug -tail=200
func Debug(script *string, tail *int) error {
	s := ""
	if script != nil {
		s = *script
	}
	tailArg := ""
	if tail != nil && *tail > 0 {
		tailArg = strconv.Itoa(*tail)
	}
	return sh.RunV("go", "run", "./magefiles/simdebug", s, tailArg)
}

// Trace writes a full game log to a file, end to end. It plays the fixed-seed
// property games once with the log turned on, so a whole game can be read at once
// instead of one prompt at a time. Unlike `mage debug`, which shows the tail of a
// game that broke, a trace is the full log of games that pass.
//
// -count sets how many of the property batch's games to play (default 1) and -out
// the destination; it defaults to tmp/sim/trace.log, under the repo's gitignored
// scratch directory:
//
//	mage trace
//	mage trace -count=25
//	mage trace -out=tmp/sim/mine.log
func Trace(count *int, out *string) error {
	countArg := ""
	if count != nil && *count > 0 {
		countArg = strconv.Itoa(*count)
	}
	o := ""
	if out != nil {
		o = *out
	}
	return sh.RunV("go", "run", "./magefiles/simtrace", countArg, o)
}

// Soak runs the long-running game soak with assertions on. It
// drives internal/sim.TestSoak under the assert build tag, churning fresh random
// games across GOMAXPROCS workers until the time budget elapses. It does not stop
// at the first failure: every failing script is saved into
// internal/sim/testdata/fuzz/FuzzPlay as a permanent FuzzPlay regression.
//
// The budget defaults to 30s. Pass -duration a Go duration to change it:
//
//	mage soak
//	mage soak -duration=5m
func Soak(duration *time.Duration) error {
	d := 30 * time.Second
	if duration != nil {
		d = *duration
	}
	os.Setenv("SOAK_DURATION", d.String())
	return sh.RunV("go", "test", "-tags", "assert",
		"-timeout", soakTimeout(d),
		"-run", "^TestSoak$", "-count", "1", "-v", "./internal/sim")
}

// soakTimeout gives go test a panic timeout past the soak's own budget. The soak
// stops itself at the budget, so the timeout is only there to catch a genuine
// hang — but go test defaults it to 10m, which kills a soak of 10m or longer
// mid-run and reports the budget elapsing as a timeout panic.
func soakTimeout(budget time.Duration) string {
	return (budget + 5*time.Minute).String()
}

// profileDir is the gitignored scratch directory the profiles and the profiling
// test binary are written to.
const profileDir = "tmp/sim"

// baselineFile is the committed performance baseline `mage profile` diffs against
// and `mage profile -save` re-blesses. It holds only counts, never wall-clock, so
// they are independent of CPU speed and git history reads as change over time. The
// counts carry a small per-run jitter (the engine randomizes map iteration), so the
// baseline is an advisory regression signal, not an exact gate; only FastCopy
// allocs/op is exact (it must stay 0, guarding ADR 0005).
const baselineFile = "internal/sim/testdata/perf_baseline.json"

// benchGamesPerOp is the number of games one BenchmarkPlaySeeds op replays, matching
// benchSeedBatch in internal/sim, so per-op totals divide down to per-game counts.
const benchGamesPerOp = 1000

// perfBaseline holds the counts from the seeded workload plus the snapshot copy's
// allocation count. Only counts live here — no ns/op — so the file is independent of
// CPU speed. The play counts carry a small per-run jitter; FastCopyAllocs is exact.
type perfBaseline struct {
	AllocsPerGame  float64 `json:"allocs_per_game"`
	BytesPerGame   float64 `json:"bytes_per_game"`
	ActionsPerGame float64 `json:"actions_per_game"`
	FastCopyAllocs float64 `json:"fastcopy_allocs_per_op"`
}

// Profile profiles the engine under the sim workload and diffs it to a baseline.
// It runs the sim benchmarks in internal/sim — the closest proxy the repo has for
// the load an MCTS bot would put on the engine — writes CPU and allocation profiles
// under tmp/sim, and prints the deterministic per-game counts next to the committed
// baseline so a regression shows as a delta. It does not block; open the flame graph
// afterwards with `mage profileServer`.
//
// Default runs the 1000 deterministic seeded games (cpu.prof, mem.prof). -random
// instead runs a fresh random batch — the outlier pass, whose counts vary run to run
// so it is never diffed — to cpu-random.prof and mem-random.prof. -save re-blesses
// the baseline from the seeded run and cannot combine with -random:
//
//	mage profile
//	mage profile -random
//	mage profile -save
func Profile(random, save *bool) error {
	r := random != nil && *random
	s := save != nil && *save
	if r && s {
		return fmt.Errorf("-save works only on the seeded workload; drop -random")
	}
	if err := os.MkdirAll(profileDir, 0o755); err != nil {
		return err
	}
	if r {
		return profileRandom()
	}
	return profileSeeded(s)
}

// profileSeeded profiles the deterministic batch, prints per-game counts against the
// baseline, and re-blesses the baseline when save is set.
func profileSeeded(save bool) error {
	cpu := filepath.Join(profileDir, "cpu.prof")
	mem := filepath.Join(profileDir, "mem.prof")
	play, err := profiledBench("PlaySeeds", cpu, mem)
	if err != nil {
		return err
	}
	micro, err := microBench()
	if err != nil {
		return err
	}
	cur := perfBaseline{
		AllocsPerGame:  play["allocs/op"] / benchGamesPerOp,
		BytesPerGame:   play["B/op"] / benchGamesPerOp,
		ActionsPerGame: play["actions/game"],
		FastCopyAllocs: micro["FastCopy"]["allocs/op"],
	}
	base, hasBase := readBaseline()
	printSeeded(play, micro, cur, base, hasBase)
	if save {
		if err := writeBaseline(cur); err != nil {
			return err
		}
		fmt.Printf("\nbaseline saved to %s\n", baselineFile)
	}
	fmt.Println("\nflame graph: `mage profileServer` (add -mem for allocations)")
	return nil
}

// profileRandom profiles a fresh random batch — the outlier pass. Its counts vary
// run to run, so it prints them for inspection but never diffs a baseline.
func profileRandom() error {
	cpu := filepath.Join(profileDir, "cpu-random.prof")
	mem := filepath.Join(profileDir, "mem-random.prof")
	play, err := profiledBench("PlayRandom", cpu, mem)
	if err != nil {
		return err
	}
	fmt.Println("\nrandom outlier workload — per game (not baselined):")
	fmt.Printf("  ns/game       %12.1f\n", play["ns/op"]/benchGamesPerOp)
	fmt.Printf("  allocs/game   %12.1f\n", play["allocs/op"]/benchGamesPerOp)
	fmt.Printf("  bytes/game    %12.1f\n", play["B/op"]/benchGamesPerOp)
	fmt.Printf("  actions/game  %12.1f\n", play["actions/game"])
	fmt.Println("\nflame graph: `mage profileServer -random` (add -mem for allocations)")
	return nil
}

// ProfileServer opens a profile from the last `mage profile` run in the interactive
// pprof web UI (flame graph, call graph, source view) and blocks until stopped. Go's
// flame graph is served, not printed, so this is its own target: `mage profile`
// writes the profiles without blocking, and this serves one on demand.
//
// Default serves the seeded CPU profile. -mem serves the allocation profile, -random
// the outlier run's profiles:
//
//	mage profileServer
//	mage profileServer -mem
//	mage profileServer -random
func ProfileServer(mem, random *bool) error {
	name := "cpu"
	if mem != nil && *mem {
		name = "mem"
	}
	useRandom := random != nil && *random
	if useRandom {
		name += "-random"
	}
	file := filepath.Join(profileDir, name+".prof")
	if _, err := os.Stat(file); err != nil {
		hint := "mage profile"
		if useRandom {
			hint += " -random"
		}
		return fmt.Errorf("no profile at %s — run `%s` first", file, hint)
	}
	addr, err := freeAddr()
	if err != nil {
		return err
	}
	fmt.Printf("serving %s at http://%s — opening a browser, Ctrl-C to stop\n", file, addr)
	return sh.RunV("go", "tool", "pprof", "-http="+addr, file)
}

// freeAddr picks a concrete free localhost address for pprof's web UI. pprof's own
// -http=localhost:0 binds an ephemeral port but then prints and opens the literal
// ":0", which no browser can reach, so the port has to be resolved up front.
func freeAddr() (string, error) {
	l, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		return "", err
	}
	defer l.Close()
	return l.Addr().String(), nil
}

// profiledBench runs one play benchmark with CPU and allocation profiling on and
// returns its parsed metrics. The profiles and the test binary land under tmp/sim;
// absolute paths keep go test from resolving them against the package directory.
func profiledBench(name, cpu, mem string) (benchMetrics, error) {
	absCPU, _ := filepath.Abs(cpu)
	absMem, _ := filepath.Abs(mem)
	absBin, _ := filepath.Abs(filepath.Join(profileDir, "sim.test"))
	out, err := benchOutput(
		"-bench", "^Benchmark"+name+"$",
		"-cpuprofile", absCPU,
		"-memprofile", absMem,
		"-o", absBin,
	)
	if err != nil {
		return nil, err
	}
	m, ok := parseBench(out)[name]
	if !ok {
		return nil, fmt.Errorf("benchmark %s produced no result", name)
	}
	return m, nil
}

// microBench runs the snapshot micro-benchmarks without profiling and at a short
// benchtime, since they only supply deterministic counts (FastCopy allocs/op guards
// ADR 0005) and their nanosecond costs — not a CPU profile, which the fast loops
// would swamp. It returns each benchmark's metrics keyed by name.
func microBench() (map[string]benchMetrics, error) {
	out, err := benchOutput(
		"-bench", "^Benchmark(FastCopy|SnapshotApplyRestore)$",
		"-benchtime", "2000x",
	)
	if err != nil {
		return nil, err
	}
	return parseBench(out), nil
}

// benchOutput runs go test in benchmark mode against internal/sim with the given
// extra flags, echoes the raw output, and returns it for parsing.
func benchOutput(extra ...string) (string, error) {
	args := append([]string{
		"test", "-run", "^$", "-benchmem",
	}, extra...)
	args = append(args, "./internal/sim")
	out, err := sh.Output("go", args...)
	fmt.Println(out)
	if err != nil {
		return "", fmt.Errorf("running benchmarks: %w", err)
	}
	return out, nil
}

// benchMetrics maps each metric unit a `go test -benchmem` line reports — ns/op,
// B/op, allocs/op, and any custom ReportMetric such as actions/game — to its value.
type benchMetrics map[string]float64

// parseBench reads `go test -bench` output and returns one metric set per benchmark,
// keyed by the benchmark's name without its "Benchmark" prefix or "-N" GOMAXPROCS
// suffix. A benchmark line is the name, the iteration count, then value/unit pairs.
func parseBench(out string) map[string]benchMetrics {
	res := map[string]benchMetrics{}
	sc := bufio.NewScanner(strings.NewReader(out))
	for sc.Scan() {
		fields := strings.Fields(sc.Text())
		if len(fields) < 4 || !strings.HasPrefix(fields[0], "Benchmark") {
			continue
		}
		name := strings.TrimPrefix(fields[0], "Benchmark")
		if i := strings.LastIndex(name, "-"); i >= 0 {
			name = name[:i]
		}
		m := benchMetrics{}
		for i := 2; i+1 < len(fields); i += 2 {
			if v, err := strconv.ParseFloat(fields[i], 64); err == nil {
				m[fields[i+1]] = v
			}
		}
		res[name] = m
	}
	return res
}

// printSeeded prints the seeded workload's per-game counts and the snapshot micro-
// benchmarks, each count next to its baseline delta.
func printSeeded(
	play benchMetrics,
	micro map[string]benchMetrics,
	cur, base perfBaseline,
	hasBase bool,
) {
	fmt.Println("\nseeded workload — per game (1000 deterministic games):")
	fmt.Printf(
		"  ns/game       %12.1f   (local only, not baselined)\n",
		play["ns/op"]/benchGamesPerOp,
	)
	fmt.Printf(
		"  allocs/game   %12.1f   %s\n",
		cur.AllocsPerGame,
		delta(cur.AllocsPerGame, base.AllocsPerGame, hasBase),
	)
	fmt.Printf(
		"  bytes/game    %12.1f   %s\n",
		cur.BytesPerGame,
		delta(cur.BytesPerGame, base.BytesPerGame, hasBase),
	)
	fmt.Printf(
		"  actions/game  %12.1f   %s\n",
		cur.ActionsPerGame,
		delta(cur.ActionsPerGame, base.ActionsPerGame, hasBase),
	)
	fmt.Println("\nsnapshot micro-benchmarks:")
	fmt.Printf("  fastcopy ns/op          %12.1f\n", micro["FastCopy"]["ns/op"])
	fmt.Printf(
		"  fastcopy allocs/op      %12.1f   %s\n",
		cur.FastCopyAllocs,
		delta(cur.FastCopyAllocs, base.FastCopyAllocs, hasBase),
	)
	fmt.Printf("  snapshot+apply ns/op    %12.1f\n", micro["SnapshotApplyRestore"]["ns/op"])
}

// delta renders a value's relation to its baseline as an absolute baseline and a
// percentage change, or a plain note when there is no baseline yet.
func delta(cur, base float64, hasBase bool) string {
	if !hasBase {
		return "(no baseline)"
	}
	if base == 0 {
		if cur == 0 {
			return "(baseline 0)"
		}
		return fmt.Sprintf("(baseline 0 -> %.1f)", cur)
	}
	pct := (cur - base) / base * 100
	sign := "+"
	if pct < 0 {
		sign = ""
	}
	return fmt.Sprintf("(baseline %.1f, %s%.1f%%)", base, sign, pct)
}

// readBaseline loads the committed baseline, reporting whether one exists yet.
func readBaseline() (perfBaseline, bool) {
	data, err := os.ReadFile(baselineFile)
	if err != nil {
		return perfBaseline{}, false
	}
	var b perfBaseline
	if json.Unmarshal(data, &b) != nil {
		return perfBaseline{}, false
	}
	return b, true
}

// writeBaseline commits a new baseline as indented JSON with a trailing newline.
func writeBaseline(b perfBaseline) error {
	if err := os.MkdirAll(filepath.Dir(baselineFile), 0o755); err != nil {
		return err
	}
	data, err := json.Marshal(b, jsontext.WithIndent("  "))
	if err != nil {
		return err
	}
	return os.WriteFile(baselineFile, append(data, '\n'), 0o644)
}
