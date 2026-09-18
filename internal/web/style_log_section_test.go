package web

import (
	"reflect"
	"testing"

	"github.com/dmikalova/vactrol/internal/engine"
	"github.com/dmikalova/vactrol/internal/sim"
)

// sampledStyle builds a style component whose log gallery is filled from a small
// deterministic batch, so the section can be rendered in a test.
func sampledStyle(t *testing.T) *style {
	t.Helper()
	s := &style{}
	s.logCov = sampleLog(sim.SeedScripts(40))
	s.logSampled = true
	return s
}

// TestLogSyntheticCoversUnobserved checks the synthetic fallback is complete:
// every kind no game produced has a catalog instance behind it and renders,
// so a reader sees every kind either observed or synthetic — never missing.
func TestLogSyntheticCoversUnobserved(t *testing.T) {
	s := sampledStyle(t)
	byKind := map[logKind]engine.LogEntry{}
	for _, e := range engine.LogEntrySamples() {
		byKind[reflect.TypeOf(e)] = e
	}
	for _, k := range s.logCov.unobserved {
		e, ok := byKind[k]
		if !ok {
			t.Errorf("unobserved kind %s has no catalog instance to draw", k)
			continue
		}
		if ui := syntheticLine(e); ui == nil {
			t.Errorf("synthetic line for %s rendered nil", k)
		}
	}
}

// TestLogSectionRenders checks the section draws in each state — before sampling,
// after sampling, and with a drill-in open — without panicking.
func TestLogSectionRenders(t *testing.T) {
	empty := &style{}
	if ui := empty.logSection(); ui == nil {
		t.Fatal("unsampled logSection rendered nil")
	}

	s := sampledStyle(t)
	if ui := s.logSection(); ui == nil {
		t.Fatal("sampled logSection rendered nil")
	}
	if ui := s.logHeroGallery(); ui == nil {
		t.Fatal("logHeroGallery rendered nil")
	}
	if len(s.logCov.cover) == 0 {
		t.Fatal("sample produced no cover bubbles to drill into")
	}
	s.drill, s.drillOpen = s.logCov.cover[0], true
	if ui := s.logDrill(); ui == nil {
		t.Fatal("open logDrill rendered nil")
	}
}
