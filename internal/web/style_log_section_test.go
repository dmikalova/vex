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

// TestSampleLogGamesResolveCardNames checks each retained game carries a card
// lookup, so a hover over a log mention can resolve the printed card to preview
// it — without it, onLogCardHover would find nothing and no preview would open.
func TestSampleLogGamesResolveCardNames(t *testing.T) {
	s := sampledStyle(t)
	if len(s.logCov.games) == 0 {
		t.Fatal("sample retained no games")
	}
	for i, sg := range s.logCov.games {
		if len(sg.game.defByName) == 0 {
			t.Errorf("game %d has no card lookup, so its log mentions cannot preview", i)
		}
	}
}

// TestLogPreviewFollowsHover checks the style page reflects a sampled game's
// hover: with nothing hovered there is no preview, and once a game holds a hovered
// card the page previews that card.
func TestLogPreviewFollowsHover(t *testing.T) {
	s := sampledStyle(t)
	if def := s.hoveredPreviewDef(); def != nil {
		t.Fatalf("nothing hovered but preview shows %q", def.Name)
	}
	gw := s.logCov.games[0].game
	for _, def := range gw.defByName {
		gw.hoverDef = def
		break
	}
	if gw.hoverDef == nil {
		t.Fatal("could not seed a hovered card")
	}
	if got := s.hoveredPreviewDef(); got != gw.hoverDef {
		t.Fatalf("hoveredPreviewDef = %v, want the hovered card %q", got, gw.hoverDef.Name)
	}
	if ui := s.logPreviewOverlay(); ui == nil {
		t.Fatal("logPreviewOverlay rendered nil with a card hovered")
	}
}
