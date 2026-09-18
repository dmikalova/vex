package web

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/engine"
)

// TestTriggerSpecimensCoverPrintedTriggers checks the trigger gallery ranges
// over the engine's registry, not a written list: every printed trigger yields
// exactly one specimen and no unprinted trigger does, so a trigger added to the
// engine cannot silently miss the gallery.
func TestTriggerSpecimensCoverPrintedTriggers(t *testing.T) {
	want := map[string]bool{}
	for _, tr := range engine.Triggers() {
		if tr.Printed() {
			want[tr.String()] = true
		}
	}
	specs := triggerSpecimens()
	if len(specs) != len(want) {
		t.Fatalf("got %d trigger specimens, want %d printed triggers", len(specs), len(want))
	}
	for _, sp := range specs {
		if !want[sp.Caption] {
			t.Errorf("specimen %q is not a printed trigger", sp.Caption)
		}
		delete(want, sp.Caption)
	}
	for name := range want {
		t.Errorf("printed trigger %q has no specimen", name)
	}
}

// TestContinuousSpecimensCoverFeatures checks one specimen per continuous text
// feature, captioned by the feature it queries.
func TestContinuousSpecimensCoverFeatures(t *testing.T) {
	specs := continuousSpecimens()
	if len(specs) != len(cardTextFeatures) {
		t.Fatalf("got %d continuous specimens, want %d features", len(specs), len(cardTextFeatures))
	}
	for i, sp := range specs {
		if sp.Caption != cardTextFeatures[i].caption {
			t.Errorf(
				"specimen %d caption = %q, want %q",
				i,
				sp.Caption,
				cardTextFeatures[i].caption,
			)
		}
	}
}

// TestCardTextSectionRenders checks the whole section and each specimen it draws
// build without panicking, whether the query found a card or fell to a gap.
// It asserts nothing about markup or wording (ADR 0014).
func TestCardTextSectionRenders(t *testing.T) {
	s := &style{}
	if ui := s.cardTextSection(); ui == nil {
		t.Fatal("cardTextSection rendered nil")
	}
	specs := append(triggerSpecimens(), continuousSpecimens()...)
	specs = append(specs, targetShapeSpecimens()...)
	for _, sp := range specs {
		if ui := s.styleCard(sp); ui == nil {
			t.Errorf("styleCard for %q rendered nil", sp.Caption)
		}
	}
	if ui := combinerList(); ui == nil {
		t.Fatal("combinerList rendered nil")
	}
}

// TestTargetShapeSpecimensAreRealCards checks the target-shape gallery is
// data-driven: every phrase it shows was read off a real card, so each specimen
// resolves to a card rather than a gap and the card really targets that phrase.
func TestTargetShapeSpecimensAreRealCards(t *testing.T) {
	specs := targetShapeSpecimens()
	if len(specs) == 0 {
		t.Fatal("no target-shape specimens; expected phrases from the card pool")
	}
	for _, sp := range specs {
		if !sp.found() {
			t.Errorf("target phrase %q has no card, but was gathered from one", sp.Caption)
			continue
		}
		if !defTargetsPhrase(sp.Def, sp.Caption) {
			t.Errorf(
				"card %q does not target the phrase %q it was matched for",
				sp.Def.Name,
				sp.Caption,
			)
		}
	}
}

// TestCombinerRowsAreDistinctPairs checks the combiner demo builds a phrase for
// every row and the rows are distinct, so the minimal pairs it isolates do not
// collapse into one. It asserts the structure of the pairs, not their wording
// (ADR 0014).
func TestCombinerRowsAreDistinctPairs(t *testing.T) {
	rows := combinerRows()
	if len(rows) == 0 {
		t.Fatal("no combiner rows")
	}
	seen := map[string]bool{}
	for _, r := range rows {
		if r.Phrase == "" {
			t.Errorf("combiner %q rendered an empty phrase", r.Caption)
		}
		if seen[r.Caption] {
			t.Errorf("duplicate combiner caption %q", r.Caption)
		}
		seen[r.Caption] = true
	}
}
