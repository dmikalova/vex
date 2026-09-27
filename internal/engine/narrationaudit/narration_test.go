package narrationaudit

import (
	"reflect"
	"testing"

	"github.com/dmikalova/vex/internal/engine"
)

// mutatingRoles are the Resolver role interfaces through which a card changes
// the game (ADR 0008). The read-only roles — StateReader and its mirrors, plus
// ChoiceResolver and Logger — are deliberately absent: a read changes nothing,
// so there is nothing for the log to have missed.
func mutatingRoles() map[string]reflect.Type {
	return map[string]reflect.Type{
		"EconomyResolver":  reflect.TypeFor[engine.EconomyResolver](),
		"CreatureResolver": reflect.TypeFor[engine.CreatureResolver](),
		"BoardResolver":    reflect.TypeFor[engine.BoardResolver](),
		"CombatResolver":   reflect.TypeFor[engine.CombatResolver](),
		"PlayResolver":     reflect.TypeFor[engine.PlayResolver](),
		"ZoneResolver":     reflect.TypeFor[engine.ZoneResolver](),
		"TurnResolver":     reflect.TypeFor[engine.TurnResolver](),
	}
}

// TestEveryMutatingPortMethodDeclaresItsNarration is the port-side twin of
// internal/engine's TestEveryStateFieldDeclaresItsNarration: it cannot prove the
// log is complete, but it stops a gap being introduced silently. A method added
// to a mutating role fails the build until methodNarration says how the log
// covers it — and saying "narrates directly" is not free, because Auditor then
// holds the method to it on every call a test makes.
func TestEveryMutatingPortMethodDeclaresItsNarration(t *testing.T) {
	declared := map[string]bool{}
	for role, typ := range mutatingRoles() {
		for method := range typ.Methods() {
			name := method.Name
			declared[name] = true
			if _, ok := methodNarration[name]; !ok {
				t.Errorf(
					"%s.%s has no methodNarration entry; decide whether the method "+
						"narrates its own outcome", role, name,
				)
			}
		}
	}
	for name := range methodNarration {
		if !declared[name] {
			t.Errorf(
				"methodNarration names %q, which no mutating Resolver role declares",
				name,
			)
		}
	}
}
