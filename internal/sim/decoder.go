package sim

import "github.com/dmikalova/vactrol/internal/engine"

// decoder turns the fuzz/soak byte script into a stream of small decisions. When
// the script runs out every read returns a zero value, which makes players stop
// acting and the game wind down to a natural end — so a short script still yields a
// complete, terminating game.
type decoder struct {
	script []byte
	pos    int
}

// done reports whether the script is exhausted.
func (d *decoder) done() bool { return d.pos >= len(d.script) }

// byte consumes and returns the next script byte, or 0 when exhausted.
func (d *decoder) byte() byte {
	if d.done() {
		return 0
	}
	b := d.script[d.pos]
	d.pos++
	return b
}

// bool consumes a byte and returns its low bit as a boolean (e.g. a flank side).
func (d *decoder) bool() bool { return d.byte()&1 == 1 }

// uint64 consumes eight bytes (big-endian) for a seed. Exhausted bytes read as 0.
func (d *decoder) uint64() uint64 {
	var v uint64
	for i := 0; i < 8; i++ {
		v = v<<8 | uint64(d.byte())
	}
	return v
}

// scriptChooser answers the engine's target/option choices from the script,
// matching engine.FirstChooser's determinism when the script is exhausted (it then
// always picks the first candidate). One instance serves both players, since the
// script already encodes whose decision each byte drives.
type scriptChooser struct{ d *decoder }

// ChooseCreature picks a candidate by script byte; declining only when there are no
// candidates. Always committing to a target keeps optional effects legal and simple.
func (c *scriptChooser) ChooseCreature(
	_, _ string,
	candidates []engine.LocalID,
) (engine.LocalID, bool) {
	if len(candidates) == 0 {
		return 0, false
	}
	return candidates[int(c.d.byte())%len(candidates)], true
}

// ChooseOption picks one of a set of option prompts by script byte.
func (c *scriptChooser) ChooseOption(_, _ string, options []string) int {
	if len(options) == 0 {
		return 0
	}
	return int(c.d.byte()) % len(options)
}
