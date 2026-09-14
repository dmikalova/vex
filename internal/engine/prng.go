package engine

// This file holds the match's random source. It is a counter-based PRNG kept flat
// inside GameState (ADR 0005, 0039): a single word of state that a snapshot
// captures and that replaying the command log reproduces bit-for-bit, replacing
// the *rand.Rand that used to hang off Game. Every shuffle and random pick during
// play draws from here, so a saved seed plus the command log regenerates the exact
// game.

// PRNG is the match's random state: a SplitMix64 counter whose single word both
// seeds and advances the stream. Each draw adds the golden-ratio increment and
// finalizes the word, so the generator is a pure function of its counter — flat,
// comparable, and copied by value with GameState.FastCopy.
type PRNG struct {
	// State is the running counter; NewGame seeds it and every draw advances it.
	State uint64
}

// splitMix64Gamma is the golden-ratio odd increment SplitMix64 adds per draw.
const splitMix64Gamma = 0x9e3779b97f4a7c15

// next advances the counter one step and returns the finalized 64-bit output.
func (p *PRNG) next() uint64 {
	p.State += splitMix64Gamma
	z := p.State
	z = (z ^ (z >> 30)) * 0xbf58476d1ce4e5b9
	z = (z ^ (z >> 27)) * 0x94d049bb133111eb
	return z ^ (z >> 31)
}

// uint64n returns a uniform value in [0, n) using rejection sampling, so the draw
// carries no modulo bias. n must be positive.
func (p *PRNG) uint64n(n uint64) uint64 {
	if n&(n-1) == 0 { // power of two (including 1): mask is exact and unbiased
		return p.next() & (n - 1)
	}
	// threshold is 2^64 mod n; values at or above it leave a whole number of n-sized
	// buckets below, so accepting only those and taking mod n is unbiased.
	threshold := (-n) % n
	for {
		if v := p.next(); v >= threshold {
			return v % n
		}
	}
}

// Intn returns a uniform int in [0, n). It panics on a non-positive n, matching
// math/rand.Intn, since a random index into an empty range is a caller bug.
func (p *PRNG) Intn(n int) int {
	if n <= 0 {
		panic("engine: PRNG.Intn: n must be positive")
	}
	return int(p.uint64n(uint64(n)))
}

// Shuffle permutes n elements with a Fisher-Yates pass, calling swap(i, j) to
// exchange the elements at those indices — the same shape as math/rand.Shuffle.
func (p *PRNG) Shuffle(n int, swap func(i, j int)) {
	for i := n - 1; i > 0; i-- {
		swap(i, p.Intn(i+1))
	}
}
