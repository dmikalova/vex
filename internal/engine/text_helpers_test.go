package engine

import (
	"regexp"
	"testing"
)

// oneAsPlural matches a count of one followed by a plural noun, e.g. "1 cards".
// Mass nouns (damage, power, armor, Æmber) are absent because they never take an
// "s" and so cannot carry this bug.
var oneAsPlural = regexp.MustCompile(
	`\b1 (cards|creatures|artifacts|upgrades|keys|chains|counters|times|entries)\b`,
)

// TestLogLinesNeverPrintOneAsPlural renders every log sample and fails on a line
// that counts one of something with a plural noun. This is the standing guard
// behind countNoun: writing "%d cards" directly renders "1 cards" the first time
// a game passes 1, which is how "1 keys" shipped. It asserts on the rendered line
// rather than on the source so that the legitimate hand-written branches — "the
// top card", "a card revealed this way", "1 of 3 keys" — stay off the report
// instead of needing a growing exemption list.
//
// TestCardTextUsesControlledVocabulary applies the same rule to printed card text
// and the rulebook, through universalRetired in internal/cards.
func TestLogLinesNeverPrintOneAsPlural(t *testing.T) {
	var n stubNamer
	for i, e := range LogEntrySamples() {
		if m := oneAsPlural.FindString(e.Text(n)); m != "" {
			t.Errorf("sample %d (%T) renders %q; use countNoun", i, e, m)
		}
	}
}
