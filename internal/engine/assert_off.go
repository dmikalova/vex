//go:build !assert

package engine

// DebugAssert reports whether the engine was built with the "assert" build tag.
// It is false in normal and production builds, so the in-engine invariant checks
// (see assertInvariants) compile away to nothing and MCTS rollouts pay no cost.
const DebugAssert = false

// assertInvariants is a no-op in the normal build. The -tags assert build (see
// assert_on.go) replaces it with a check that panics on a corrupt state, so
// soak/fuzz runs validate the engine while production stays lean.
func (g *Game) assertInvariants() {}
