package web

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/dmikalova/vactrol/internal/cards"
	"github.com/dmikalova/vactrol/internal/engine"
)

// The UI is mostly markup, which a test can only assert by restating it. What a
// test can catch — and what actually breaks in a browser — is the seam between
// the Go markup and the files it names: an icon whose asset does not exist and a
// house whose CSS class was never defined both render as silent breakage. These
// tests guard that seam, plus the few pure helpers that carry real logic.

// repoFile reads a file relative to the repository root.
func repoFile(t *testing.T, rel string) string {
	t.Helper()
	b, err := os.ReadFile(filepath.Join("..", "..", rel))
	if err != nil {
		t.Fatalf("read %s: %v", rel, err)
	}
	return string(b)
}

func assetExists(t *testing.T, stem string) bool {
	t.Helper()
	_, err := os.Stat(filepath.Join("..", "..", "web", "assets", stem+".svg"))
	return err == nil
}

// iconCall matches an icon("stem", …) call so the test can check every stem the
// package names, not just the ones a helper computes.
var iconCall = regexp.MustCompile(`\bicon\("([a-z0-9-]+)"`)

func TestIconNamesHaveAssets(t *testing.T) {
	files, err := filepath.Glob("*.go")
	if err != nil {
		t.Fatalf("glob: %v", err)
	}
	for _, f := range files {
		if strings.HasSuffix(f, "_test.go") {
			continue
		}
		b, err := os.ReadFile(f)
		if err != nil {
			t.Fatalf("read %s: %v", f, err)
		}
		for _, m := range iconCall.FindAllStringSubmatch(string(b), -1) {
			if !assetExists(t, m[1]) {
				t.Errorf("%s: icon %q has no web/assets/%s.svg", f, m[1], m[1])
			}
		}
	}
}

// stemLiteral matches a lowercase-hyphen string literal ("zone-hand") and the
// same stem inside a .svg path ("/web/assets/favicon.svg"), which is how the code
// names an asset.
var stemLiteral = regexp.MustCompile(`([a-z][a-z0-9-]+)(?:\.svg)?"`)

// stripGalleryIcons removes the galleryIcons declaration from style.go's source.
// That list names every asset by design, so counting it as a reference would make
// every asset trivially "used" and defeat TestNoDeadAssets.
func stripGalleryIcons(src string) string {
	const marker = "var galleryIcons = []string{"
	i := strings.Index(src, marker)
	if i < 0 {
		return src
	}
	if j := strings.Index(src[i:], "}"); j >= 0 {
		return src[:i] + src[i+j+1:]
	}
	return src
}

// TestNoDeadAssets is the reverse of TestIconNamesHaveAssets: every SVG in
// web/assets is drawn by some functional code path, not just parked in the Style
// gallery's showcase. The gallery lists every asset on purpose (galleryIcons), so
// its declaration is excluded from the scan — otherwise a leftover asset from a
// removed feature would count as "used" and never be caught. An asset is
// considered referenced if the code names its stem: as a string literal or .svg
// path in the package (icon("stem"), glyph{asset: "stem"}, every resolver that
// returns a literal) or cmd/web (the favicon and PWA icons), plus the one computed
// family — house emblems, whose stem is "house-"+slug.
func TestNoDeadAssets(t *testing.T) {
	referenced := map[string]bool{}
	addStems := func(src string) {
		for _, m := range stemLiteral.FindAllStringSubmatch(src, -1) {
			referenced[m[1]] = true
		}
	}

	// House emblems are the only stems built from a slug rather than a literal.
	referenced["house-none"] = true
	for h := engine.HouseNone + 1; int(h) < engine.NumHouses; h++ {
		if name := houseIconName(h); name != "" {
			referenced[name] = true
		}
	}

	files, err := filepath.Glob("*.go")
	if err != nil {
		t.Fatalf("glob: %v", err)
	}
	for _, f := range files {
		if strings.HasSuffix(f, "_test.go") {
			continue
		}
		b, err := os.ReadFile(f)
		if err != nil {
			t.Fatalf("read %s: %v", f, err)
		}
		addStems(stripGalleryIcons(string(b)))
	}
	// The favicon and PWA icons are named in the server entrypoint, not the board.
	addStems(repoFile(t, "cmd/web/main.go"))

	assets, err := filepath.Glob(filepath.Join("..", "..", "web", "assets", "*.svg"))
	if err != nil {
		t.Fatalf("glob assets: %v", err)
	}
	for _, a := range assets {
		stem := strings.TrimSuffix(filepath.Base(a), ".svg")
		if !referenced[stem] {
			t.Errorf(
				"web/assets/%s.svg is a dead asset: no code outside galleryIcons draws it",
				stem,
			)
		}
	}
}

func TestHouseIconNamesHaveAssets(t *testing.T) {
	if got := houseIconName(engine.HouseNone); got != "" {
		t.Errorf("houseIconName(HouseNone) = %q, want empty", got)
	}
	// houseIconName itself stays "" for HouseNone (a card's own emblem hides when
	// its house is unset), but houseIcon falls back to a real asset for it, since
	// the Style gallery and house pickers draw HouseNone as its own labelled row.
	if !assetExists(t, "house-none") {
		t.Error("houseIcon's HouseNone fallback \"house-none\" has no web/assets/house-none.svg")
	}
	for h := engine.HouseNone + 1; int(h) < engine.NumHouses; h++ {
		name := houseIconName(h)
		if name == "" {
			t.Errorf("house %v has no icon name", h)
			continue
		}
		if !assetExists(t, name) {
			t.Errorf("house %v: no web/assets/%s.svg", h, name)
		}
	}
}

// TestLogHouseIconsResolve guards the seam between the engine's log-icon keys and
// the web assets they name. The engine keys a house emblem by its lowercased
// printed name, so a house whose name carries a space — Star Alliance — arrives as
// "house-star alliance", which has no matching asset stem and would render blank
// unless the web side normalises it. This fails loudly if any house the log can
// reference resolves to a missing asset, rather than shipping a silent gap.
func TestLogHouseIconsResolve(t *testing.T) {
	for h := engine.HouseNone + 1; int(h) < engine.NumHouses; h++ {
		key := "house-" + strings.ToLower(h.String())
		stem := logIconStem(key)
		if !assetExists(t, stem) {
			t.Errorf(
				"house %v: log icon key %q resolves to stem %q, which has no web/assets/%s.svg",
				h,
				key,
				stem,
				stem,
			)
		}
	}
}

func TestTypeAndKeyIconNamesHaveAssets(t *testing.T) {
	for _, ct := range []engine.CardType{
		engine.Creature, engine.Artifact, engine.Tactic, engine.Upgrade,
	} {
		name := typeIconName(ct)
		if name == "" || !assetExists(t, name) {
			t.Errorf("card type %v: icon %q missing", ct, name)
		}
	}
	for _, tc := range []struct {
		color engine.KeyColor
		label string
	}{
		{engine.KeyColorRed, "Red"},
		{engine.KeyColorBlue, "Blue"},
		{engine.KeyColorYellow, "Yellow"},
	} {
		name := keyColorIconName(tc.color)
		if name == "" || !assetExists(t, name) {
			t.Errorf("key colour %v: icon %q missing", tc.color, name)
		}
		if got := keyColorByName(tc.label); got != tc.color {
			t.Errorf("keyColorByName(%q) = %v, want %v", tc.label, got, tc.color)
		}
	}
	if got := keyColorIconName(engine.KeyColorNone); got != "" {
		t.Errorf("keyColorIconName(none) = %q, want empty", got)
	}
	if got := keyColorByName("Chartreuse"); got != engine.KeyColorNone {
		t.Errorf("keyColorByName(unknown) = %v, want none", got)
	}
}

// TestCounterIconNamesHaveAssets guards the generic-counter icon seam: every
// counter kind the engine defines resolves to a unique icon asset that exists on
// disk. A new CounterKind with no counterAsset entry (or one reusing another
// kind's icon) fails here, so the web strip can never draw two counters the same
// or fall back to a missing SVG.
func TestCounterIconNamesHaveAssets(t *testing.T) {
	if got := counterAsset(engine.CounterNone); got != "" {
		t.Errorf("counterAsset(CounterNone) = %q, want empty", got)
	}
	seen := map[string]engine.CounterKind{}
	for k := engine.CounterNone + 1; int(k) < engine.NumCounterKinds; k++ {
		name := counterAsset(k)
		if name == "" {
			t.Errorf("counter kind %d has no icon asset", k)
			continue
		}
		if !assetExists(t, name) {
			t.Errorf("counter kind %d: no web/assets/%s.svg", k, name)
		}
		if prev, dup := seen[name]; dup {
			t.Errorf(
				"counter kinds %d and %d share icon %q; each needs a unique one",
				prev,
				k,
				name,
			)
		}
		seen[name] = k
	}
}

// cssRuleDeclares reports whether the CSS rule for selector defines every one of
// props. It isolates the rule body (from "selector {" to the next "}") so a
// property defined in some other rule cannot satisfy the check — catching a house
// or set whose class exists but supplies no colour, which would render blank.
func cssRuleDeclares(css, selector string, props ...string) bool {
	i := strings.Index(css, selector+" {")
	if i < 0 {
		return false
	}
	body := css[i+len(selector)+2:]
	if end := strings.Index(body, "}"); end >= 0 {
		body = body[:end]
	}
	for _, p := range props {
		if !strings.Contains(body, p) {
			return false
		}
	}
	return true
}

func TestHouseClassesAreDefinedInCSS(t *testing.T) {
	css := repoFile(t, "web/app.css")
	for h := engine.HouseNone + 1; int(h) < engine.NumHouses; h++ {
		cls := houseClasses(h)
		if !cssRuleDeclares(css, "."+cls, "--nm:", "--edge:") {
			t.Errorf("house %v: web/app.css .%s rule must define house colours "+
				"(--nm and --edge); a missing colour must fail here, not render blank", h, cls)
		}
	}
}

// TestSetAccentClassesAreDefinedInCSS guards the set-picker seam: every set the
// picker lists resolves to a .set-<slug> accent class defined in app.css, and the
// shared emblem the buttons mask lives on disk. A set whose accent or emblem is
// missing renders as an uncoloured, unmarked button rather than a visible break.
func TestSetAccentClassesAreDefinedInCSS(t *testing.T) {
	css := repoFile(t, "web/app.css")
	if !assetExists(t, "set") {
		t.Error("set emblem has no web/assets/set.svg")
	}
	for _, name := range cards.DeckSetNames() {
		cls := setAccent(name)
		if cls == "" {
			t.Errorf("set %q has no accent class", name)
			continue
		}
		if !cssRuleDeclares(css, "."+cls, "--set-accent:") {
			t.Errorf("set %q: web/app.css .%s rule must define --set-accent; a "+
				"missing colour must fail here, not render blank", name, cls)
		}
	}
}

func TestPulseClassesAreDefinedInCSS(t *testing.T) {
	css := repoFile(t, "web/app.css")
	if got := pulseClass(false, false, "dmg"); got != "" {
		t.Errorf("pulseClass(off) = %q, want empty", got)
	}
	for _, kind := range []string{"dmg", "pow", "gain"} {
		a, b := pulseClass(true, false, kind), pulseClass(true, true, kind)
		if a == b {
			t.Errorf("kind %q: both parities give %q, so a repeat would not replay", kind, a)
		}
		for _, cls := range []string{a, b} {
			if !strings.Contains(css, "."+cls) {
				t.Errorf("web/app.css defines no .%s rule", cls)
			}
		}
	}
}

func TestCardDOMIDsAreDistinct(t *testing.T) {
	// The same card can be measured in hand and on the board in one animation, so
	// the two ids must never collide.
	if boardCardID(7) == handCardID(7) {
		t.Fatalf("board and hand ids collide: %q", boardCardID(7))
	}
	if boardCardID(7) == boardCardID(8) {
		t.Fatalf("distinct cards share the id %q", boardCardID(7))
	}
}

func TestCx(t *testing.T) {
	if got := cx("card", "", "card--selected", ""); got != "card card--selected" {
		t.Errorf("cx = %q", got)
	}
	if got := cx(ifCls(true, "on"), ifCls(false, "off")); got != "on" {
		t.Errorf("ifCls composition = %q", got)
	}
}

func TestContainsAndIndexOfID(t *testing.T) {
	ids := []engine.LocalID{3, 9, 4}
	if !containsID(ids, 9) || containsID(ids, 5) {
		t.Error("containsID is wrong")
	}
	if got := indexOfID(ids, 4); got != 2 {
		t.Errorf("indexOfID = %d, want 2", got)
	}
	if got := indexOfID(ids, 5); got != -1 {
		t.Errorf("indexOfID(absent) = %d, want -1", got)
	}
}

func TestRarityMarks(t *testing.T) {
	for _, tc := range []struct {
		rarity engine.Rarity
		mark   rarityMark
		icon   string
	}{
		{engine.Common, rarityCommon, "rarity-triangle"},
		{engine.Uncommon, rarityUncommon, "rarity-square"},
		{engine.Rare, rarityRare, "rarity-pentagon"},
		{engine.Special, raritySpecial, "rarity-hexagon"},
		{engine.Connected, rarityConnected, "rarity-connected"},
	} {
		got := rarityMarkOf(tc.rarity)
		if got != tc.mark {
			t.Errorf("rarityMarkOf(%v) = %v, want %v", tc.rarity, got, tc.mark)
		}
		if name := got.iconName(); name != tc.icon {
			t.Errorf("%v iconName = %q, want %q", tc.rarity, name, tc.icon)
		}
	}
	if name := rarityNone.iconName(); name != "" {
		t.Errorf("rarityNone iconName = %q, want empty", name)
	}
}

func TestKindAndTraitLabels(t *testing.T) {
	def := &engine.CardDefinition{
		Type:   engine.Creature,
		Traits: []engine.Trait{engine.Human, engine.Knight},
	}
	if got := kindLabel(def); got == "" {
		t.Error("kindLabel(creature) is empty")
	}
	if got := traitLabel(def); !strings.Contains(got, "Human") || !strings.Contains(got, "Knight") {
		t.Errorf("traitLabel = %q", got)
	}
	if got := traitLabel(&engine.CardDefinition{Type: engine.Creature}); got != "" {
		t.Errorf("traitLabel(no traits) = %q, want empty", got)
	}
}
