package engine

import "testing"

// TestEnumTotality binds every text-bearing enum's enumerating function to the
// constants the package's own source declares with its type, so a constant that
// is neither enumerated nor named as excluded fails the build (ADR 0018). This is
// the guard the hand-written lists never had: Keywords() and its siblings were
// only as complete as the author who last added a constant remembered to make
// them, which left the rulebook completeness test unable to see a keyword nobody
// listed.
func TestEnumTotality(t *testing.T) {
	for _, enum := range Enums() {
		t.Run(enum.Type, func(t *testing.T) {
			declared := enumConstants(t, enum)
			for _, name := range enum.Excluded {
				if _, ok := declared[name]; !ok {
					t.Errorf("excluded constant %s is not declared with type %s; "+
						"drop it from Excluded", name, enum.Type)
					continue
				}
				delete(declared, name)
			}
			if len(declared) != enum.Enumerated {
				t.Errorf("%s declares %d constants the census expects enumerated, "+
					"but its enumerating function returns %d; add the new constant to it "+
					"or name it in Excluded", enum.Type, len(declared), enum.Enumerated)
			}
			if len(enum.Rows) == 0 {
				return
			}
			catalogued := map[string]bool{}
			for _, row := range enum.Rows {
				catalogued[row.Type] = true
			}
			for name := range declared {
				if !catalogued[name] {
					t.Errorf("%s %s has no census row; catalogue it (ADR 0018)", enum.Type, name)
				}
			}
			for name := range catalogued {
				if _, ok := declared[name]; !ok {
					t.Errorf("census row %s is not a %s constant the package declares",
						name, enum.Type)
				}
			}
		})
	}
}

// TestEnumRowsWellFormed checks each enum row the way TestFamilyRowsWellFormed
// checks a node row: it classifies its value as owing exactly one of a term or a
// reason for owing none, and the value renders text.
func TestEnumRowsWellFormed(t *testing.T) {
	for _, enum := range Enums() {
		t.Run(enum.Type, func(t *testing.T) {
			for _, row := range enum.Rows {
				checkRowClassified(t, row)
			}
		})
	}
}

// enumConstants returns every constant the package's non-test source declares
// with the enum's type, failing the test when the scan finds none — an enum whose
// type was renamed would otherwise report every row as an orphan.
func enumConstants(t *testing.T, enum Enum) map[string]string {
	t.Helper()
	declared, err := enum.Declared(".")
	if err != nil {
		t.Fatalf("scanning for %s constants: %v", enum.Type, err)
	}
	if len(declared) == 0 {
		t.Fatalf("found no %s constants; the source scan is broken", enum.Type)
	}
	return declared
}
