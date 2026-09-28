package web

import (
	"strings"
	"testing"
)

// TestGlyphSectionShowsEveryCatalogedNode holds the section's listing equal to
// the census the Iconography pass is measured against. The section is the review
// surface for drawing glyphs, so a node it omits is a drawing nobody reviews —
// which is exactly the failure the card-walk population used to hide.
func TestGlyphSectionShowsEveryCatalogedNode(t *testing.T) {
	shown := map[string]map[string]bool{}
	for _, family := range glyphCatalogViews() {
		shown[family.name] = map[string]bool{}
		for _, row := range family.rows {
			shown[family.name][row.name] = true
		}
	}
	for _, family := range iconCatalog() {
		rows, ok := shown[family.name]
		if !ok {
			t.Errorf("the glyph section shows no %s family", family.name)
			continue
		}
		for _, row := range family.rows {
			if !rows[row.name] {
				t.Errorf("the glyph section omits %s node %s", family.name, row.name)
			}
		}
	}
}

// TestGlyphSectionExplainsEveryRow checks every row carries a classification: a
// rulebook term, or the reason it owes none. A blank strip is shown rather than
// omitted, so without the note beside it a plumbing row that folds into its
// children would read as a mechanic nobody drew.
func TestGlyphSectionExplainsEveryRow(t *testing.T) {
	for _, family := range glyphCatalogViews() {
		for _, row := range family.rows {
			if row.note == "" {
				t.Errorf("%s row %s shows no term and no plumbing reason", family.name, row.name)
			}
		}
	}
}

// TestGlyphSectionRenders checks the section builds without panicking over the
// whole census, ungated families included. It asserts nothing about markup or
// wording (ADR 0014).
func TestGlyphSectionRenders(t *testing.T) {
	s := &style{}
	if ui := s.glyphVocabSection(); ui == nil {
		t.Fatal("glyphVocabSection rendered nil")
	}
	// The second call must come off the cache rather than redrawing the census.
	if ui := s.glyphVocabSection(); ui == nil {
		t.Fatal("glyphVocabSection rendered nil on re-render")
	}
	if len(s.glyphCatalog) == 0 {
		t.Error("the section did not cache its census rows")
	}
}

// TestGlyphSectionClassesAreDefinedInCSS checks the classes the census listing
// draws with exist in web/app.css. The blank modifier is the one that matters: a
// row drawing nothing is meant to show as an empty bordered cell, and an
// undefined class would make it show as nothing at all — the omission the
// section exists to prevent.
func TestGlyphSectionClassesAreDefinedInCSS(t *testing.T) {
	css := repoFile(t, "web/app.css")
	for _, cls := range []string{
		"style-h4",
		"style-glyph-count",
		"style-glyph-rows",
		"style-glyph-row",
		"style-glyph-strip",
		"style-glyph-strip--blank",
		"style-glyph-meta",
		"style-glyph-text",
		"style-glyph-note",
		"style-glyph-note--plumbing",
	} {
		if !strings.Contains(css, "."+cls) {
			t.Errorf("web/app.css defines no .%s rule", cls)
		}
	}
}
