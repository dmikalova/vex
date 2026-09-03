package sim

import "github.com/dmikalova/vactrol/internal/engine"

// decoder turns the fuzz/soak byte script into a stream of small decisions. When
// the script runs out every read returns a zero value, which makes players stop
// acting; simulate then ends the game at the next turn boundary, so a short script
// still yields a complete, terminating game.
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

// stopOneIn weights the answers that shorten a game — ending a turn, discarding a
// card away, declining an optional target — at roughly one script byte in this
// many, instead of letting them share the options' uniform odds. As one uniform
// slot among a handful of legal moves, stopping ends most turns after an action or
// two, so neither player ever banks enough Æmber to forge and every game runs out
// the turn limit without a winner.
const stopOneIn = 8

// pick returns a script-driven index into n options, or -1 for the rare "stop"
// answer (see stopOneIn). n of 0 always stops.
func (d *decoder) pick(n int) int {
	b := int(d.byte())
	if n == 0 || b%stopOneIn == 0 {
		return -1
	}
	// Index by the quotient, not the remainder: the stop test already claimed every
	// byte divisible by stopOneIn, so a remainder would leave option 0 unreachable
	// whenever n divides stopOneIn.
	return b / stopOneIn % n
}

// scriptChooser answers the engine's target/option choices from the script,
// matching engine.FirstChooser's determinism when the script is exhausted (it then
// always picks the first candidate). One instance serves both players, since the
// script already encodes whose decision each byte drives.
type scriptChooser struct{ d *decoder }

// ChooseCreature picks a candidate by script byte, declining only when there are
// no candidates. This is the mandatory-choice path, so it never stops early.
func (c *scriptChooser) ChooseCreature(
	_, _ string,
	candidates []engine.LocalID,
) (engine.LocalID, bool) {
	if len(candidates) == 0 {
		return 0, false
	}
	return candidates[int(c.d.byte())%len(candidates)], true
}

// ChooseCardOrDecline answers a "you may" choice, declining on the rare stop odds
// so both branches of an optional effect get exercised without most of them being
// thrown away.
func (c *scriptChooser) ChooseCardOrDecline(
	_, _ string,
	candidates []engine.LocalID,
) (engine.LocalID, bool) {
	i := c.d.pick(len(candidates))
	if i < 0 {
		return 0, false
	}
	return candidates[i], true
}

// ChooseOption picks one of a set of option prompts by script byte.
func (c *scriptChooser) ChooseOption(_, _ string, options []string) int {
	if len(options) == 0 {
		return 0
	}
	return int(c.d.byte()) % len(options)
}
