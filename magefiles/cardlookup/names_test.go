package main

import (
	"strings"
	"testing"
)

// apostropheForms are the apostrophe characters a card name can carry: the ASCII
// ' (U+0027), the typographic right single quote (U+2019), and the modifier
// letter apostrophe (U+02BC, itself a Unicode letter). varName and fileName must
// drop every form so none reaches a Go identifier or file base name.
var apostropheForms = []rune{'\u0027', '\u2019', '\u02bc'}

// TestNamesStripApostrophes asserts the invariant that an apostrophe makes no
// difference: a name carrying any apostrophe form yields the same identifier and
// file name as the same name without it, and neither output keeps an apostrophe.
// It compares two computed results rather than a golden literal, so a config that
// rewrites a string literal (a typos autocorrect) cannot invalidate the test.
func TestNamesStripApostrophes(t *testing.T) {
	for _, ap := range apostropheForms {
		bare := "Frane" + "s Blaster"
		withAp := "Frane" + string(ap) + "s Blaster"
		if got, want := varName(withAp), varName(bare); got != want {
			t.Errorf("varName(%q) = %q, want %q (apostrophe should not matter)", withAp, got, want)
		}
		if got, want := fileName(withAp), fileName(bare); got != want {
			t.Errorf("fileName(%q) = %q, want %q (apostrophe should not matter)", withAp, got, want)
		}
		for _, form := range apostropheForms {
			gotVar, gotFile := varName(withAp), fileName(withAp)
			if strings.ContainsRune(gotVar, form) || strings.ContainsRune(gotFile, form) {
				t.Errorf("apostrophe survived in %q / %q", gotVar, gotFile)
			}
		}
	}
}
