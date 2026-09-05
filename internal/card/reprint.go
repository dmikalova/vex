package card

import "github.com/dmikalova/vactrol/internal/cards/provenance"

// ReprintRef claims that a card first implemented in another set is also printed
// in Set at collector number Number. A reprint is a full member of Set's
// deck-generation pool, the same as a card newly implemented in that set. The
// claim carries only the card's name, so the set package that registers it (its
// generated 0set.go) imports no other set — the cards aggregator resolves Name
// against the card registry. See docs/deck-generation.md and ADR 0021.
type ReprintRef struct {
	Set    provenance.SourceSet
	Number int
	Name   string
}

// reprints holds every reprint claim, populated at init by each set's 0set.go.
var reprints []ReprintRef

// Reprint records that Name — implemented in the set that introduced it — is also
// printed in set at collector number number. A set's generated 0set.go calls it
// once per reprint so the card joins that set's pool without the package importing
// another set. It is a deck-generation membership claim, never read by the engine.
func Reprint(set provenance.SourceSet, number int, name string) {
	reprints = append(reprints, ReprintRef{Set: set, Number: number, Name: name})
}

// ReprintRefs returns a copy of every registered reprint claim. The cards
// aggregator uses it to fold reprints into each set's own deck-generation pool.
func ReprintRefs() []ReprintRef {
	out := make([]ReprintRef, len(reprints))
	copy(out, reprints)
	return out
}
