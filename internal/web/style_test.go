package web

import (
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"testing"

	"github.com/maxence-charriere/go-app/v11/pkg/app"

	"github.com/dmikalova/vactrol/internal/engine"
)

// The Style gallery's job is to show everything, so what its tests check is
// coverage, not appearance: each one holds two enumerations against each other
// and reports what one has that the other does not. None of them pins markup,
// class names or wording, so restyling the page never breaks them — only
// forgetting to include something does.

// TestGalleryShowsEveryIcon holds the gallery's icon list equal to the assets on
// disk. An icon added to web/assets and drawn on the board but missing here
// would be a piece of the vocabulary the gallery silently omits.
func TestGalleryShowsEveryIcon(t *testing.T) {
	files, err := filepath.Glob(filepath.Join("..", "..", "web", "assets", "*.svg"))
	if err != nil {
		t.Fatalf("glob assets: %v", err)
	}
	var want []string
	for _, f := range files {
		stem := strings.TrimSuffix(filepath.Base(f), ".svg")
		// The favicon is browser chrome, not part of the board's vocabulary.
		if stem == "favicon" {
			continue
		}
		want = append(want, stem)
	}
	shown := map[string]bool{}
	for _, name := range galleryIcons {
		shown[name] = true
	}
	sort.Strings(want)
	for _, stem := range want {
		if !shown[stem] {
			t.Errorf("web/assets/%s.svg is not in galleryIcons", stem)
		}
	}
	for _, name := range galleryIcons {
		if !assetExists(t, name) {
			t.Errorf("galleryIcons has %q, which has no web/assets/%s.svg", name, name)
		}
	}
}

// TestEveryHouseHasAFontToken checks the per-house font override exists for
// every house palette. The gallery sets --font-house to try a face out, so a
// house whose palette never declares it would silently ignore the selection.
func TestEveryHouseHasAFontToken(t *testing.T) {
	css := repoFile(t, "web/app.css")
	for h := engine.HouseNone; int(h) < engine.NumHouses; h++ {
		block := "." + houseSlug(h)
		at := strings.Index(css, "."+houseClasses(h)+" {")
		if at < 0 {
			t.Errorf("%s: no .card-%s palette in web/app.css", block, houseSlug(h))
			continue
		}
		end := strings.Index(css[at:], "}")
		if end < 0 || !strings.Contains(css[at:at+end], "--font-house") {
			t.Errorf("house %s has a palette but no --font-house token", h)
		}
	}
}

// TestSpecimensCoverTheCatalog checks that the queries driving the page answer
// for every card type and rarity. A gap here is a real gap in the loaded sets,
// which the gallery draws — this test says so out loud rather than letting it
// pass unnoticed.
func TestSpecimensCoverTheCatalog(t *testing.T) {
	_, types, cells := houseGrid()
	byType := map[engine.CardType]bool{}
	for _, c := range cells {
		if c.found() {
			byType[c.Def.Type] = true
		}
	}
	for _, ct := range types {
		if !byType[ct] {
			t.Errorf("no specimen anywhere for card type %s", ct)
		}
	}
	// A rarity no loaded set prints is a gap the gallery draws, not a failure —
	// but the question must still be asked, or the mark would go unshown the day a
	// set does print one.
	asked := map[engine.Rarity]bool{}
	for _, sp := range raritySpecimens() {
		asked[engine.Rarity(sp.Caption)] = true
	}
	for _, r := range []engine.Rarity{
		engine.Common, engine.Uncommon, engine.Rare, engine.Special, engine.Connected,
	} {
		if !asked[r] {
			t.Errorf("the gallery never asks for rarity %s", r)
		}
	}
}

// TestEveryKeywordHasASpecimenRow checks the feature section asks about every
// keyword the engine defines. Whether a card answers is up to the loaded sets —
// an unanswered query is drawn as a gap — but the question must always be asked,
// which is what makes a keyword added to the engine show up here on its own.
func TestEveryKeywordHasASpecimenRow(t *testing.T) {
	asked := map[string]bool{}
	for _, sp := range featureSpecimens() {
		asked[sp.Caption] = true
	}
	for _, kw := range engine.Keywords() {
		if !asked[kw.String()] {
			t.Errorf("the feature section never asks for %s", kw)
		}
	}
}

// TestGalleryRenders builds the whole page and walks every section, which is the
// one thing a markup test can usefully assert: that nothing on it panics. The
// harness is a synthetic game, so a surface that assumes a state only a played
// match reaches would fail here rather than in a browser.
func TestGalleryRenders(t *testing.T) {
	s := &style{harness: styleHarness(), attachHost: attachHarness()}
	s.attachments = buildAttachments(s.attachHost.g)
	if s.Render() == nil {
		t.Fatal("the gallery rendered nothing")
	}
	// A loaded font must reach both the compare strip and a house override.
	s.fonts = []string{"Inter"}
	s.uiFont = "Inter"
	s.houseFont[engine.Brobnar] = "Inter"
	if s.Render() == nil {
		t.Fatal("the gallery rendered nothing with fonts selected")
	}
}

// TestGalleryShowsEveryAttachmentCombination holds the attachment specimens to
// the counts and faceup/facedown combinations the board must draw: 1/2/3
// upgrades, every faceup/facedown mask of 1/2/3 under-cards, the peeked case,
// and the two composed. A combination the board can reach but the gallery omits
// would be an unshown state.
func TestGalleryShowsEveryAttachmentCombination(t *testing.T) {
	h := attachHarness()
	specs := buildAttachments(h.g)
	captions := map[string]bool{}
	for _, s := range specs {
		captions[s.caption] = true
	}
	want := []string{
		"1 upgrade", "2 upgrades", "3 upgrades",
		"1 under: up", "1 under: down",
		"1 under, facedown (you peek)",
		"2 upgrades + 2 under", "3 upgrades + 3 under",
	}
	for n := 1; n <= 3; n++ {
		for mask := range 1 << n {
			want = append(want, strconv.Itoa(n)+" under: "+underMaskLabel(n, mask))
		}
	}
	for _, w := range want {
		if !captions[w] {
			t.Errorf("the gallery has no attachment specimen %q", w)
		}
	}
}

// TestGalleryHidesAFacedownUnderFromTheOpponent renders the specimens and checks
// that a facedown under-card on the opponent's card draws as a card-back, while
// the specimen for a card you control draws its face — the visible difference
// between the two the section exists to show.
func TestGalleryHidesAFacedownUnderFromTheOpponent(t *testing.T) {
	h := attachHarness()
	specs := buildAttachments(h.g)
	var hidden, peeked string
	for _, s := range specs {
		html := app.HTMLString(h.hostWithTabs(s.host, h.printedCard(s.host)))
		switch s.caption {
		case "1 under: down":
			hidden = html
		case "1 under, facedown (you peek)":
			peeked = html
		}
	}
	if !strings.Contains(hidden, "card-tab--back") {
		t.Error("a facedown under-card on the opponent's card is not drawn as a card-back")
	}
	if strings.Contains(peeked, "card-tab--back") {
		t.Error("your own facedown under-card is drawn as a card-back instead of its face")
	}
}

// TestGalleryHarnessFillsEveryZone checks the synthetic position leaves no zone
// count at zero and no two the same, since equal counts would hide a Player bar
// that read the wrong zone.
func TestGalleryHarnessFillsEveryZone(t *testing.T) {
	g := styleHarness()
	for player := range 2 {
		counts := map[string]int{
			"hand":     len(g.g.Hand(player)),
			"deck":     len(g.g.Deck(player)),
			"discard":  len(g.g.Discard(player)),
			"archives": len(g.g.Archives(player)),
		}
		seen := map[int]string{}
		for zone, n := range counts {
			if n == 0 {
				t.Errorf("player %d: %s is empty", player, zone)
			}
			if other, dup := seen[n]; dup {
				t.Errorf("player %d: %s and %s both hold %d cards", player, zone, other, n)
			}
			seen[n] = zone
		}
	}
}

// TestFontURLsAreRestrictedToProviders pins the allowlist that keeps a pasted
// URL from becoming a stylesheet this origin loads. It is the gallery's only
// user input, so it is the only place a bad string can get in.
func TestFontURLsAreRestrictedToProviders(t *testing.T) {
	good := []struct{ url, want string }{
		{"https://fonts.googleapis.com/css2?family=Inter:wght@400;700", "Inter"},
		{"https://fonts.googleapis.com/css2?family=Cinzel+Decorative", "Cinzel Decorative"},
		{"https://fonts.bunny.net/css?family=lato", "lato"},
	}
	for _, c := range good {
		got, err := fontFamily(c.url)
		if err != nil || got != c.want {
			t.Errorf("fontFamily(%q) = %q, %v; want %q, nil", c.url, got, err, c.want)
		}
	}
	bad := []string{
		"http://fonts.googleapis.com/css2?family=Inter", // not https
		"https://evil.example/css2?family=Inter",        // not a font provider
		"https://fonts.googleapis.com/css2",             // names no family
		"javascript:alert(1)",                           // not a stylesheet at all
		"",
	}
	for _, u := range bad {
		if got, err := fontFamily(u); err == nil {
			t.Errorf("fontFamily(%q) = %q, nil; want an error", u, got)
		}
	}
}

// TestFamilyParamKeepsTheSeparator pins the one encoding subtlety in the font
// URL: the "+" between words is the provider's separator, so escaping the whole
// name turns it into %2B and the request comes back 400.
func TestFamilyParamKeepsTheSeparator(t *testing.T) {
	cases := map[string]string{
		"Inter":             "Inter",
		"Cinzel Decorative": "Cinzel+Decorative",
		"IBM Plex Mono":     "IBM+Plex+Mono",
	}
	for family, want := range cases {
		if got := familyParam(family); got != want {
			t.Errorf("familyParam(%q) = %q, want %q", family, got, want)
		}
	}
}
