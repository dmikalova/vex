package web

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

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

func TestHouseClassesAreDefinedInCSS(t *testing.T) {
	css := repoFile(t, "web/app.css")
	for h := engine.HouseNone + 1; int(h) < engine.NumHouses; h++ {
		cls := houseClasses(h)
		if !strings.Contains(css, "."+cls) {
			t.Errorf("house %v: web/app.css defines no .%s rule", h, cls)
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
		rarity   engine.Rarity
		mark     rarityMark
		diamonds int
	}{
		{engine.Common, rarityCommon, 1},
		{engine.Uncommon, rarityUncommon, 2},
		{engine.Rare, rarityRare, 3},
		{engine.Special, raritySpecial, 4},
		{engine.Connected, rarityConnected, 0},
	} {
		got := rarityMarkOf(tc.rarity)
		if got != tc.mark {
			t.Errorf("rarityMarkOf(%v) = %v, want %v", tc.rarity, got, tc.mark)
		}
		if n := got.diamonds(); n != tc.diamonds {
			t.Errorf("%v diamonds = %d, want %d", tc.rarity, n, tc.diamonds)
		}
	}
	if !rarityConnected.isConnected() || rarityCommon.isConnected() {
		t.Error("isConnected is wrong")
	}
	if n := len(rarityDiamonds(3)); n != 3 {
		t.Errorf("rarityDiamonds(3) rendered %d icons", n)
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
