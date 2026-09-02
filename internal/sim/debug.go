package sim

import "math/rand"

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
	return Failure{Script: script, Log: g.Log, Err: err}
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
