package cards

import (
	"regexp"
	"strings"
	"testing"

	"github.com/dmikalova/vactrol/internal/engine"
)

// retiredTerm is a word the controlled Rules voice retired in favor of one
// preferred synonym (ADR 0019). The tests below fail the build when a retired
// term reaches a player-facing surface, so the one-word-per-meaning rule can
// never quietly regress.
type retiredTerm struct {
	pattern *regexp.Regexp
	use     string
}

// universalRetired are retired everywhere: the resource is always Æmber, and each
// generic-gaming verb has one KeyForge word (docs/CONTEXT.md, docs/style-guide.md).
var universalRetired = []retiredTerm{
	{regexp.MustCompile(`(?i)\baember\b`), "Æmber"},
	{regexp.MustCompile(`(?i)\bmana\b`), "Æmber"},
	{regexp.MustCompile(`(?i)\bgets\b`), "gains"},
	{regexp.MustCompile(`(?i)\bsacrifices?\b|\bsacrificed\b`), "destroy"},
	{regexp.MustCompile(`(?i)\btaps?\b|\btapped\b`), "exhaust"},
	{regexp.MustCompile(`(?i)\bexiles?\b|\bexiled\b`), "purge, or put into hand/discard"},
	{regexp.MustCompile(`(?i)\bbounces?\b|\bbounced\b`), "put into hand"},
}

// cardTextRetired adds the conventions that govern only printed card text
// (docs/card-wording-rules.md rules 4 and 19): one movement verb, and the card
// type's name. The rulebook prose keeps "return" for senses card text does not
// use (Æmber returning to the supply), so those stay off the universal list.
var cardTextRetired = append([]retiredTerm{
	{regexp.MustCompile(`(?i)\breturns?\b|\breturned\b`), "put"},
	{regexp.MustCompile(`(?i)\baction cards?\b`), "Tactic"},
}, universalRetired...)

// allowedCardText are the card-text surfaces that still use a retired term, each
// an open STE pre-pass decision tracked in docs/todo.md. The guard exempts these
// exact "<card> | <word>" pairings so it enforces every other surface today; a
// new violation is not on the list and still fails the build. Empty means every
// implemented card conforms.
var allowedCardText = map[string]bool{}

// lint reports every retired term in text that is not on the allow set.
func lint(t *testing.T, where, text string, terms []retiredTerm, allow map[string]bool) {
	t.Helper()
	for _, rt := range terms {
		m := rt.pattern.FindString(text)
		if m == "" || allow[where+" | "+strings.ToLower(m)] {
			continue
		}
		t.Errorf("%s uses retired term %q (use %q)", where, m, rt.use)
	}
}

// TestCardTextUsesControlledVocabulary fails the build if any card's rendered text
// uses a word the Rules voice retired, so the vocabulary is enforced by the
// build rather than by review.
func TestCardTextUsesControlledVocabulary(t *testing.T) {
	all := All()
	for i := range all {
		lint(t, all[i].Name, engine.RenderCardText(&all[i]), cardTextRetired, allowedCardText)
	}
}

// TestRulebookUsesControlledVocabulary applies the universal guard to every
// rulebook surface the registry renders (ADR 0018): the overview, each section
// intro, and every term's prose.
func TestRulebookUsesControlledVocabulary(t *testing.T) {
	lint(t, "overview", engine.RuleOverview(), universalRetired, nil)
	for _, sec := range []engine.Section{
		engine.SectionTurn, engine.SectionCombat, engine.SectionCardType,
		engine.SectionKeyword, engine.SectionAbility, engine.SectionEffect,
	} {
		lint(t, "intro:"+string(sec), engine.RuleSectionIntro(sec), universalRetired, nil)
	}
	for _, term := range engine.RuleTerms() {
		body := term.Title + "\n" + term.Subtitle + "\n" + term.Definition + "\n" + term.Body
		lint(t, "term "+term.Title, body, universalRetired, nil)
	}
}
