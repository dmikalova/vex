package web

import (
	htmlpkg "html"
	"strings"
	"testing"

	"github.com/maxence-charriere/go-app/v11/pkg/app"

	"github.com/dmikalova/vactrol/internal/card"
)

// TestGalleryWindowClamps checks the drawn-face window never points outside the
// shown-card count, so a filter that shrinks the list can never index past its
// end or draw a backwards range.
func TestGalleryWindowClamps(t *testing.T) {
	g := &gallery{winFrom: 10, winTo: 100}
	from, to := g.window(5)
	if from != 5 || to != 5 {
		t.Errorf("window past the end gave [%d,%d), want [5,5)", from, to)
	}

	g = &gallery{winFrom: 1, winTo: 3}
	if from, to := g.window(10); from != 1 || to != 3 {
		t.Errorf("in-range window gave [%d,%d), want [1,3)", from, to)
	}

	g = &gallery{winFrom: -4, winTo: 2}
	if from, to := g.window(10); from != 0 || to != 2 {
		t.Errorf("negative from gave [%d,%d), want [0,2)", from, to)
	}
}

// TestGalleryWindowsFacesButKeepsTextSearchable checks the grid draws a full face
// only for cards inside the window and a lightweight placeholder for the rest,
// while keeping every card's name and rules text in the DOM so the browser's
// find-in-page (Ctrl+F) still matches a card whose heavy face is not mounted.
func TestGalleryWindowsFacesButKeepsTextSearchable(t *testing.T) {
	regs := card.Cards()
	if len(regs) < 3 {
		t.Skipf("need at least 3 registered cards, have %d", len(regs))
	}
	g := &gallery{ready: true, order: "name"}
	for i := 0; i < 3; i++ {
		d := regs[i].Def
		g.cards = append(g.cards, galleryCard{def: &d, nameHay: d.Name, textHay: d.Name})
	}
	// Only the first card in the sorted grid gets a full face; the rest fall back to
	// placeholders.
	g.winFrom, g.winTo = 0, 1

	html := app.HTMLString(g.Render())
	if !strings.Contains(html, "gallery-card--placeholder") {
		t.Fatal("no off-screen card rendered as a text placeholder")
	}
	// find-in-page reads the decoded DOM text, so decode the entity-escaped markup
	// before matching a name that may contain an apostrophe.
	shown := htmlpkg.UnescapeString(html)
	for _, c := range g.cards {
		if !strings.Contains(shown, c.def.Name) {
			t.Errorf("card %q is not in the DOM, so find-in-page could not match it", c.def.Name)
		}
	}
}
