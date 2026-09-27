package engine

import (
	"testing"

	"github.com/dmikalova/vex/internal/census"
)

// TestFamilyTotality binds every gated family to the set of node types the
// package's own source declares, so a node with no census row — and a row naming
// a type that no longer exists — fails the build (ADR 0018, the shape ADR 0046
// uses for log entries). The families come from Families(), so a family added
// later is checked without editing this test. A family still being catalogued is
// not gated yet: `mage tool:census` reports its gaps until the task that
// finishes it switches the gate on.
func TestFamilyTotality(t *testing.T) {
	for _, family := range Families() {
		if !family.Gated {
			continue
		}
		t.Run(family.Name, func(t *testing.T) {
			declared := familyTypes(t, family)
			catalogued := map[string]bool{}
			for _, row := range family.Rows {
				catalogued[row.Type] = true
			}
			for name := range declared {
				if !catalogued[name] {
					t.Errorf("%s %s has no census row; catalogue it (ADR 0018)", family.Name, name)
				}
			}
			for name := range catalogued {
				if _, ok := declared[name]; !ok {
					t.Errorf("census row %s is not a %s the package declares", name, family.Name)
				}
			}
		})
	}
}

// TestFamilyRowsWellFormed checks each row itself: it classifies its node as
// owing exactly one of a term or a reason for owing none, and the node it
// carries renders text. The render check is what stops a zero-valued literal
// standing in for a node that reads a Target and cataloguing something that
// prints nothing.
func TestFamilyRowsWellFormed(t *testing.T) {
	for _, family := range Families() {
		t.Run(family.Name, func(t *testing.T) {
			for _, row := range family.Rows {
				switch {
				case row.Rules.Term != "" && row.Rules.NoTerm != "":
					t.Errorf("%s names both a term (%q) and a reason for none (%q)",
						row.Type, row.Rules.Term, row.Rules.NoTerm)
				case row.Rules.Term == "" && row.Rules.NoTerm == "":
					t.Errorf("%s is unclassified: name the rulebook term it owes, "+
						"or the reason it owes none", row.Type)
				}
				if text := row.Text(); text == "" {
					t.Errorf("%s renders no text; fill its census literal far enough to print",
						row.Type)
				}
			}
		})
	}
}

// familyTypes returns the name of every type in the package's non-test source
// that declares the family's method, failing the test when the scan finds none —
// a family whose method was renamed would otherwise report every row as an
// orphan.
func familyTypes(t *testing.T, family Family) map[string]string {
	t.Helper()
	types, err := census.Implementations(".", family.Method, census.Params(family.Params...))
	if err != nil {
		t.Fatalf("scanning for %s implementations: %v", family.Name, err)
	}
	if len(types) == 0 {
		t.Fatalf("found no %s implementations; the source scan is broken", family.Name)
	}
	return types
}
