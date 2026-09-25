package web

import (
	"fmt"
	"reflect"
	"strconv"
	"time"

	"github.com/maxence-charriere/go-app/v11/pkg/app"

	"github.com/dmikalova/vex/internal/engine"
	"github.com/dmikalova/vex/internal/sim"
)

// This file is the Game log section of the Style gallery: it shows one real
// bubble for every kind of log line a batch of sampled games produced, a
// drill-in to the whole game behind a bubble, and a synthetic fallback for the
// kinds no game produced. See
// docs/adr/0046-log-entry-catalog-and-sampled-gallery.md and the sampler in
// style_log_sample.go, which supplies the coverage this section draws.

// styleLogSampleGames is how many games a sampling run plays. The loop stops early
// once every observable kind has appeared, but the kinds sim never produces
// (manual edits, concede) keep it running the whole batch, so this is the cost of
// a run — kept modest because it plays synchronously on the wasm main thread.
const styleLogSampleGames = 200

// styleSampleDelay is how long autoSampleGames waits before it starts sampling.
// The sampling loop is CPU-bound and blocks the wasm main thread, so it is
// scheduled through a timer that first yields to the browser — letting the page
// paint and its load indicator settle — rather than run during the initial mount.
const styleSampleDelay = 200 * time.Millisecond

// sampleNamer names cards and players from their ids alone, so a catalogued log
// entry can render its own shape in the synthetic fallback with no game behind it.
type sampleNamer struct{}

func (sampleNamer) Name(id engine.LocalID) string { return fmt.Sprintf("Card%d", id) }
func (sampleNamer) PlayerName(player int) string  { return fmt.Sprintf("P%d", player) }

// logSection draws the Game log gallery, or a loading note while it builds:
// sampling plays games synchronously, so autoSampleGames runs it off the first
// paint (OnMount) and this section fills in when that finishes.
func (s *style) logSection() app.UI {
	if !s.logSampled {
		return app.Div().Class("style-log").Body(
			app.P().Class("style-note").Text(
				"Playing a batch of seeded games to show every kind of log line they "+
					"produce — one real bubble per kind. Kinds no game produced fall "+
					"back to a constructed sample."),
			app.P().Class("style-note").Text("Sampling games…"),
		)
	}
	return app.Div().Class("style-log").Body(
		s.logSummary(),
		app.H3().Class("style-h3").Text("Observed in real games"),
		s.logHeroGallery(),
		s.logDrill(),
		app.H3().Class("style-h3").Text("Synthetic — not produced by a sampled game"),
		s.logSyntheticGallery(),
		s.logPreviewOverlay(),
	)
}

// logPreviewOverlay draws the enlarged card face for whichever sampled bubble is
// hovered. The bubbles are rendered by the game's own logBlockView, so hovering a
// card mention runs that game's onLogCardHover and sets its hoverDef; the style
// page owns no per-mention state, so it finds the one game holding a preview and
// draws it. Nothing hovered draws an empty div.
func (s *style) logPreviewOverlay() app.UI {
	if def := s.hoveredPreviewDef(); def != nil {
		return app.Div().Class("style-log-preview").Body(printedFace(def))
	}
	return app.Div()
}

// hoveredPreviewDef returns the card a sampled bubble's log-mention hover is
// previewing, or nil if none is. The hover handler lives on the game wrapper, so
// the previewed card is read back off whichever wrapper is holding it.
func (s *style) hoveredPreviewDef() *engine.CardDefinition {
	for _, sg := range s.logCov.games {
		if sg.game.hoverDef != nil {
			return sg.game.hoverDef
		}
	}
	return nil
}

// autoSampleGames plays the fixed-seed batch after the page has loaded, so the
// Game log section fills itself in on a fresh visit instead of waiting for a
// button. The batch is heavy and blocks the wasm main thread, so it is scheduled
// through a short timer (ctx.After) that lets the page paint first; the sampling
// then runs and its result re-renders the section.
func (s *style) autoSampleGames(ctx app.Context) {
	if s.logSampled {
		return
	}
	ctx.After(styleSampleDelay, func(app.Context) {
		s.logCov = sampleLog(sim.SeedScripts(styleLogSampleGames))
		s.logSampled, s.drillOpen = true, false
	})
}

// reshuffleGames replays a fresh, non-deterministic batch, so a reader can look
// for a different real example of a given kind.
func (s *style) reshuffleGames(app.Context, app.Event) {
	s.logCov = sampleLog(sim.RandomScripts(styleLogSampleGames))
	s.drillOpen = false
}

// logSummary reports the three-state split — observed, synthetic, over how many
// games — and carries the reshuffle control.
func (s *style) logSummary() app.UI {
	cov := s.logCov
	observed := len(catalogKinds()) - len(cov.unobserved)
	return app.Div().Class("style-log-summary").Body(
		app.Span().Text(fmt.Sprintf(
			"%d kinds observed across %d games · %d synthetic",
			observed, cov.played, len(cov.unobserved))),
		app.Button().Class("style-btn").Text("Reshuffle").OnClick(s.reshuffleGames),
	)
}

// logHeroGallery draws the set-cover: one real bubble per group of kinds, each
// rendered through the production logBlockView and clickable to open the whole
// game it came from.
func (s *style) logHeroGallery() app.UI {
	cov := s.logCov
	out := make([]app.UI, 0, len(cov.cover))
	for _, cb := range cov.cover {
		gw := cov.games[cb.game].game
		blocks := cov.games[cb.game].blocks
		if cb.block >= len(blocks) {
			continue
		}
		open := s.drillOpen && s.drill == cb
		out = append(out, app.Div().
			Class(cx("style-log-bubble", ifCls(open, "style-log-bubble--open"))).
			OnClick(func(app.Context, app.Event) { s.toggleDrill(cb) }).
			Body(app.Div().Class("log").Body(gw.logBlockView(blocks[cb.block]))))
	}
	return app.Div().Class("style-log-gallery").Body(out...)
}

// toggleDrill opens the whole-game log behind a bubble, or closes it if that
// bubble is already open.
func (s *style) toggleDrill(cb coverBubble) {
	if s.drillOpen && s.drill == cb {
		s.drillOpen = false
		return
	}
	s.drill, s.drillOpen = cb, true
}

// closeDrill collapses the drill-in.
func (s *style) closeDrill(app.Context, app.Event) { s.drillOpen = false }

// logDrill draws the whole game log behind the open bubble, with that bubble
// highlighted, so a kind can be read in the context that produced it.
func (s *style) logDrill() app.UI {
	if !s.drillOpen {
		return app.Div()
	}
	gw := s.logCov.games[s.drill.game].game
	blocks := s.logCov.games[s.drill.game].blocks
	body := make([]app.UI, 0, len(blocks))
	for i, b := range blocks {
		view := gw.logBlockView(b)
		if i == s.drill.block {
			view = app.Div().Class("style-log-highlight").Body(view)
		}
		body = append(body, view)
	}
	return app.Div().Class("style-log-drill").Body(
		app.Div().Class("style-log-drill-head").Body(
			app.Span().Text("Full game log"),
			app.Button().Class("style-btn").Text("Close").OnClick(s.closeDrill),
		),
		app.Div().Class("log").Body(app.Div().Class("log-list").Body(body...)),
	)
}

// logSyntheticGallery draws the kinds no sampled game produced, each from its
// constructed catalog instance and badged so it is never mistaken for a real
// example. These are the manual-edit entries, concede, and effects of cards not
// yet reachable in a random game.
func (s *style) logSyntheticGallery() app.UI {
	byKind := map[logKind]engine.LogEntry{}
	for _, e := range engine.LogEntrySamples() {
		if k := reflect.TypeOf(e); byKind[k] == nil {
			byKind[k] = e
		}
	}
	out := make([]app.UI, 0, len(s.logCov.unobserved))
	for _, k := range s.logCov.unobserved {
		out = append(out, app.Div().Class("style-log-synthetic").Body(
			app.Div().Class("log").Body(
				app.Div().Class("log-group").Body(syntheticLine(byKind[k])),
			),
			app.Span().Class("style-log-badge").Text(fmt.Sprintf(
				"synthetic — not observed in %d games", s.logCov.played)),
		))
	}
	return app.Div().Class("style-log-gallery").Body(out...)
}

// syntheticLine renders one catalogued entry as a log line from its own shape,
// naming cards and players by id. It mirrors the live logSegments drawing but
// without the game-bound hover handlers, since a synthetic entry has no game to
// preview a card against.
func syntheticLine(e engine.LogEntry) app.UI {
	var out []app.UI
	for _, seg := range engine.RenderEntry(e, sampleNamer{}) {
		switch {
		case seg.HasCard:
			out = append(out, app.Span().Class("log-card").Text(seg.Text))
		case seg.HasPlayer:
			out = append(out, app.Span().
				Class("log-player log-player--p"+strconv.Itoa(seg.Player)).
				Text(seg.Text))
		case seg.Icon != "":
			out = append(out, logIcon(seg.Icon), app.Text(seg.Text))
		default:
			out = append(out, app.Text(seg.Text))
		}
	}
	return app.Div().Class(cx("log-line", logToneClass(e))).Body(out...)
}
