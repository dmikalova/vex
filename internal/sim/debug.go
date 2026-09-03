package sim

import (
	"encoding/hex"
	"fmt"
	"io"
	"math/rand"

	"github.com/dmikalova/vactrol/internal/engine"
)

// SeedScripts builds the fixed-seed batch of game scripts that TestSimulateSeeds
// plays. The scripts are drawn in sequence from one seeded source, so the whole
// batch has to be generated together for script i to stay the same game between
// the property test and a debug replay of it.
func SeedScripts(count int) [][]byte {
	r := rand.New(rand.NewSource(1))
	scripts := make([][]byte, count)
	for i := range scripts {
		scripts[i] = make([]byte, 8+r.Intn(400))
		r.Read(scripts[i])
	}
	return scripts
}

// Failure is a replayed simulation failure: the script that produced it, the
// violation, and the game log leading up to it. A zero Err means the script
// played out cleanly.
type Failure struct {
	Script []byte
	Log    []string
	Err    error
}

// Debug replays one script with the game log turned on — the seam behind `mage
// debug`, so a soak or fuzz find can be read as the sequence of plays that caused
// it rather than a bare invariant message.
func Debug(script []byte) Failure {
	g, err := simulate(script, true)
	if err == nil {
		return Failure{Script: script}
	}
	return Failure{Script: script, Log: g.LogText(), Err: err}
}

// FirstSeedFailure plays the fixed-seed batch and replays the first script that
// fails, so `mage debug` with no script reproduces whatever `mage test` is
// currently catching.
func FirstSeedFailure(count int) Failure {
	for _, script := range SeedScripts(count) {
		if Simulate(script) != nil {
			return Debug(script)
		}
	}
	return Failure{}
}

// Trace replays the first count games of the fixed-seed property batch with the
// game log on and writes every line of every game to w, each under a header
// naming its batch index and script. It is the seam behind `mage trace`: where
// Debug shows the tail of the one game that broke, a trace is the whole log of
// games that pass, for reading how the engine actually plays a turn.
//
// A failing game is written too, with its violation after the log, so a trace of
// a red batch still reads end to end. Trace returns the first write error.
func Trace(count int, w io.Writer) error {
	for i, script := range SeedScripts(count) {
		g, playErr := simulate(script, true)
		if _, err := fmt.Fprintf(
			w, "=== game %d\nscript: %s\n", i, hex.EncodeToString(script),
		); err != nil {
			return err
		}
		for _, line := range g.LogText() {
			if _, err := fmt.Fprintln(w, line); err != nil {
				return err
			}
		}
		if playErr != nil {
			if _, err := fmt.Fprintf(w, "VIOLATION: %v\n", playErr); err != nil {
				return err
			}
		}
		if _, err := fmt.Fprintf(w, "outcome: %s\n", outcome(g)); err != nil {
			return err
		}
		if _, err := fmt.Fprintln(w); err != nil {
			return err
		}
	}
	return nil
}

// outcome names how a traced game closed out, so a reader can tell a real win from
// a game that simply ran out of script.
func outcome(g *engine.Game) string {
	if w := g.Winner(); w >= 0 {
		return fmt.Sprintf("P%d wins", w)
	}
	return "no winner — the script ran out of decisions"
}
