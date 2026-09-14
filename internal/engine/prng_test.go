package engine

import "testing"

// A fixed seed yields a fixed stream, and an independent copy of the state carries
// the RNG position with it (ADR 0039): FastCopy must reproduce the same next draws
// as the original, which is what makes replay and MCTS rollouts deterministic.
func TestPRNGDeterministicAndCopied(t *testing.T) {
	a := PRNG{State: 1}
	b := PRNG{State: 1}
	for range 100 {
		if got, want := a.next(), b.next(); got != want {
			t.Fatalf("same seed diverged: %d != %d", got, want)
		}
	}

	// A snapshot of the state carries the counter, so the copy and the original draw
	// the same continuation.
	s := GameState{PRNG: PRNG{State: 42}}
	s.PRNG.next()
	cp := s.FastCopy()
	if got, want := s.PRNG.next(), cp.PRNG.next(); got != want {
		t.Fatalf("copied RNG diverged: %d != %d", got, want)
	}
}

// Intn stays within range across many draws and rejects a non-positive bound,
// matching math/rand.Intn so a random index into an empty range is a caller bug.
func TestPRNGIntn(t *testing.T) {
	p := PRNG{State: 7}
	for range 1000 {
		if v := p.Intn(6); v < 0 || v >= 6 {
			t.Fatalf("Intn(6) out of range: %d", v)
		}
	}
	if v := p.Intn(1); v != 0 {
		t.Fatalf("Intn(1) = %d, want 0", v)
	}

	defer func() {
		if recover() == nil {
			t.Fatal("Intn(0) did not panic")
		}
	}()
	p.Intn(0)
}
