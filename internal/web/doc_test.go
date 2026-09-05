package web

import (
	"strings"
	"testing"

	"github.com/maxence-charriere/go-app/v11/pkg/app"
)

// These tests draw the reference pages and read the markup back. Both pages are
// facets of the engine's rule-term registry, so the tests pin that what the engine
// assembles is what the page puts on screen.

// TestRulebookPage draws /rulebook and checks the frame, the section spine, and a
// known term all reach the markup.
func TestRulebookPage(t *testing.T) {
	h := app.HTMLString(NewRulebook().Render())

	for _, want := range []string{
		`class="doc-page"`,
		"doc-nav-link--active", // this page marked in the nav
		`href="/glossary"`,     // cross-link to the sibling
		"The Turn",             // first section heading
		"Keywords",             // a later section heading
		"Turn structure",       // a term group
	} {
		if !strings.Contains(h, want) {
			t.Errorf("/rulebook does not show %q", want)
		}
	}

	// The page draws the same terms the engine assembles: spot-check a keyword.
	if !strings.Contains(h, "Elusive") {
		t.Error("/rulebook missing the Elusive keyword term")
	}
}

// TestGlossaryPage draws /glossary and checks the frame and that every engine
// term title reaches the list.
func TestGlossaryPage(t *testing.T) {
	h := app.HTMLString(NewGlossary().Render())

	if !strings.Contains(h, `class="glossary"`) {
		t.Error("/glossary missing its list")
	}
	if !strings.Contains(h, `href="/rulebook"`) {
		t.Error("/glossary missing the cross-link to the rulebook")
	}
	// Spot-check a few plain-text keyword titles reach the list. (Titles with
	// apostrophes are HTML-escaped in the markup, so we avoid asserting on those.)
	for _, title := range []string{"Elusive", "Assault", "Taunt"} {
		if !strings.Contains(h, title) {
			t.Errorf("/glossary missing term %q", title)
		}
	}
	// Every term now carries an authored definition, so no definition cell should
	// render as the bare placeholder, and a real definition must reach the list.
	if strings.Contains(h, `class="glossary-def">`+glossaryPending+`</dd>`) {
		t.Errorf(
			"/glossary still shows the %q placeholder; a term is missing its definition",
			glossaryPending,
		)
	}
	if !strings.Contains(h, "no fight damage is dealt by or to it") {
		t.Error("/glossary missing the Elusive definition text")
	}
}
