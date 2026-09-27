package engine

import (
	"reflect"
	"testing"

	"github.com/dmikalova/vex/internal/census"
)

// TestLogEntrySamplesTotality binds the sample catalog to the set of LogEntry
// variants the source actually declares, so a new variant fails the build until
// it is catalogued (ADR 0046). It reads the package's own source rather than a
// hand-kept list: a variant is any type with a Text(Namer) method, which is the
// LogEntry interface. Record — the frame-plus-entry wrapper — implements Text the
// same way but is not itself a narrated variant, so it is excluded.
func TestLogEntrySamplesTotality(t *testing.T) {
	variants := logEntryVariantTypes(t)
	sampled := map[string]bool{}
	for _, e := range LogEntrySamples() {
		sampled[reflect.TypeOf(e).Name()] = true
	}
	for name := range variants {
		if !sampled[name] {
			t.Errorf("log entry %s has no sample in LogEntrySamples; add one (ADR 0046)", name)
		}
	}
	for name := range sampled {
		if _, ok := variants[name]; !ok {
			t.Errorf("LogEntrySamples has %s, which is not a LogEntry variant", name)
		}
	}
}

// logEntryVariantTypes returns the name of every type that implements LogEntry —
// a type with a Text(Namer) method, the signature that distinguishes it from
// Effect.Text, which takes none — with Record excluded. It runs the census
// scanner the catalog totality tests use, so the package has one source scan
// rather than two.
func logEntryVariantTypes(t *testing.T) map[string]string {
	t.Helper()
	types, err := census.Implementations(".", "Text", census.Params("Namer"))
	if err != nil {
		t.Fatalf("scanning for LogEntry variants: %v", err)
	}
	delete(types, "Record")
	if len(types) == 0 {
		t.Fatal("found no LogEntry variants; the source scan is broken")
	}
	return types
}
