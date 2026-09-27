// Package narrationaudit checks that a card cannot change the game without the
// game log saying so.
//
// It lives beside package engine rather than inside it for two reasons: the
// engine's coverage gate holds internal/engine at 100%, and the ~139 wrappers in
// resolver_gen.go exist only for a test, so gating them would measure the test
// rather than the engine; and the audit needs nothing unexported — Resolver, its
// role interfaces, Game.State and Game.Log are all exported.
//
// The audit has two halves, the same pair the state-field ratchet in
// internal/engine/state_test.go has:
//
//   - methodNarration classifies every mutating Resolver method, and
//     TestEveryMutatingPortMethodDeclaresItsNarration fails when a method has no
//     entry — so a new port method cannot be added without a decision about
//     whether it narrates.
//   - Auditor enforces the classification at runtime: a method that says it
//     narrates and then changes the state without appending a log entry fails the
//     test that drove it.
package narrationaudit

import "github.com/dmikalova/vex/internal/engine"

// The Auditor must be a complete Resolver: a method it failed to wrap would be
// a capability cards use and the audit never sees.
var _ engine.Resolver = (*Auditor)(nil)

// Reporter is the part of *testing.T the audit needs, so the package does not
// import testing.
type Reporter interface {
	Helper()
	Errorf(format string, args ...any)
}

// Auditor is a Resolver that delegates every method to the game it wraps and,
// for each method methodNarration classifies as narrating directly, fails the
// test when the call changed the game state without appending a log entry.
//
// The comparison, rather than a bare "must log on every call", is what gives the
// check teeth: many mutating methods are legitimately a no-op on some calls —
// SetDamage(id, 4) on a creature already at 4, GainAember(p, 0), SetWarded on an
// unwarded creature, ForgeKeyAtExtraCost when the player cannot pay. Requiring a
// log entry from those would force them to be reclassified as not narrating,
// which is how the check loses its teeth. GameState is comparable by design
// (ADR 0005), so "did anything change?" is one struct comparison — it gives no
// field attribution, and naming the method is enough to act on.
//
// Nesting needs no bookkeeping: an outer call passes when an inner one narrated,
// so the innermost frame that mutates silently is the one that fails.
//
// Known blind spot: BeginShuffleBatch and EndShuffleBatch emit one combined entry
// for a whole batch ("Lost in the Woods shuffles A and B into P2's deck"), so
// every move inside the batch mutates without a log entry of its own. The audit
// suspends itself while batch depth is non-zero rather than declassifying the
// methods a batch calls, which would blind it outside batches too.
type Auditor struct {
	*engine.Game

	t Reporter
	// batch is the open shuffle-batch depth; the audit is suspended above zero.
	batch int
}

// Install wraps a game in an Auditor and hands the game back its own decorated
// port, so every EffectContext the game builds from here carries the audit. It
// returns the Auditor for a caller that wants to drive the port directly.
func Install(g *engine.Game, t Reporter) *Auditor {
	a := &Auditor{Game: g, t: t}
	g.SetResolver(a)
	return a
}

// audit takes the before-state of a narrating call and returns the check to run
// after it, so a wrapper is one deferred line whatever it returns.
func (a *Auditor) audit(method string) func() {
	if method == beginBatch {
		a.batch++
	}
	if methodNarration[method] != narratedDirectly || a.batch > 0 {
		return func() {
			if method == endBatch {
				a.batch--
			}
		}
	}
	before := a.State.FastCopy()
	logged := len(a.Log)
	return func() {
		if method == endBatch {
			a.batch--
		}
		after := a.State.FastCopy()
		// The match RNG advances on every shuffle and random pick and is the one
		// field classified as never narrated (ADR 0039), so it is normalised out
		// rather than counted as a change that needed a line.
		before.PRNG, after.PRNG = engine.PRNG{}, engine.PRNG{}
		if after != before && len(a.Log) == logged {
			a.t.Helper()
			a.t.Errorf(
				"%s changed the game state without narrating it: methodNarration "+
					"classifies it as narrating directly, so it must append a log "+
					"entry (ADR 0046) on a call that changes anything",
				method,
			)
		}
	}
}

// The two methods that fence a grouped shuffle narration; see Auditor.
const (
	beginBatch = "BeginShuffleBatch"
	endBatch   = "EndShuffleBatch"
)
