package web

import "github.com/dmikalova/vactrol/internal/cards/provenance"

// Set-specific CSS classes. Like palette.go for houses, this only maps a set to
// the .set-<slug> class defined in web/app.css, which supplies that set's accent
// colour (--set-accent) so the Go markup carries labelled class names only.

// setSlug is the provenance slug of the deck-generation set with this display
// name, or "" when no source set matches. It is what turns a set the picker lists
// (named for a player) into the stable key its colours are keyed under.
func setSlug(name string) string {
	for _, s := range provenance.Sets() {
		if s.Name == name {
			return s.Slug
		}
	}
	return ""
}

// setAccent returns the class supplying a set's accent colour to its picker
// button, or "" for a set without a known slug (which then falls back to the
// default control colour).
func setAccent(name string) string {
	if slug := setSlug(name); slug != "" {
		return "set-" + slug
	}
	return ""
}
