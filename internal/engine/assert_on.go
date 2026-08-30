//go:build assert

package engine

// DebugAssert is true in an -tags assert build: the engine validates its flat
// state at turn boundaries and panics on the first violation, so a corruption
// crashes loudly at its source instead of silently producing a wrong game. Soak
// and fuzz builds carry this tag; the normal build uses the no-op in assert_off.go.
const DebugAssert = true

// assertInvariants panics if the flat state violates an invariant. It runs only
// in an -tags assert build, called at turn boundaries where the state is settled.
func (g *Game) assertInvariants() {
	if err := g.InvariantError(); err != nil {
		panic("engine invariant violated: " + err.Error())
	}
}
