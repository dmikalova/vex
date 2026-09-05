package engine

import "testing"

// TestClosedCatalogsAreComplete enforces ADR 0018: each closed catalog the game
// defines is complete in the rulebook registry. A member with no term fails the
// build here rather than going silently undescribed (the gap that once left
// Elusive and Taunt out of the rulebook). Effects and turn/combat steps are not
// yet closed catalogs and are exempt; see docs/todo.md.
func TestClosedCatalogsAreComplete(t *testing.T) {
	titled := func(section Section) map[string]bool {
		have := map[string]bool{}
		for _, term := range RuleTerms() {
			if term.Section == section {
				have[term.Title] = true
			}
		}
		return have
	}

	t.Run("keywords", func(t *testing.T) {
		have := titled(SectionKeyword)
		for _, k := range Keywords() {
			if name := k.String(); !have[name] {
				t.Errorf("keyword %q has no rulebook term (ADR 0018)", name)
			}
		}
	})

	t.Run("card types", func(t *testing.T) {
		have := titled(SectionCardType)
		for _, ct := range CardTypes() {
			if name := ct.String(); !have[name] {
				t.Errorf("card type %q has no rulebook term (ADR 0018)", name)
			}
		}
	})

	t.Run("triggers", func(t *testing.T) {
		have := titled(SectionAbility)
		for _, tr := range Triggers() {
			if !tr.Printed() {
				continue
			}
			name := tr.String()
			if name == "" {
				t.Errorf("trigger %d is printed on cards but has no String name (ADR 0018)", tr)
				continue
			}
			if !have[name] {
				t.Errorf("trigger %q has no rulebook term (ADR 0018)", name)
			}
		}
	})
}

// TestRuleTermsWellFormed checks the registry itself: it is non-empty and every
// term carries a section, a title, and a body, so a half-filled term cannot slip
// into the rulebook.
func TestRuleTermsWellFormed(t *testing.T) {
	terms := RuleTerms()
	if len(terms) == 0 {
		t.Fatal("RuleTerms() is empty")
	}
	for _, term := range terms {
		if term.Section == "" {
			t.Errorf("term %q has no section", term.Title)
		}
		if term.Title == "" {
			t.Errorf("term in section %q has no title", term.Section)
		}
		if term.Body == "" {
			t.Errorf("term %q/%q has no body", term.Section, term.Title)
		}
	}
}

// TestGlossaryComplete enforces that the glossary column of the registry (ADR
// 0018) has no blank rows: every glossary entry carries a Definition. The
// definitions were authored once the rulebook page and glossary landed, so this
// keeps a newly added term from silently reintroducing an empty entry.
func TestGlossaryComplete(t *testing.T) {
	for _, e := range Glossary() {
		if e.Definition == "" {
			t.Errorf("glossary entry %q has no definition (ADR 0018)", e.Title)
		}
	}
}

// TestRuleFramingRegistered checks the framing prose that moved out of
// docs/rulebook/*.md and into the registry (ADR 0018): the document overview is
// present, each section that carries an intro still has one, and a section with
// none (Combat) reports empty.
func TestRuleFramingRegistered(t *testing.T) {
	if RuleOverview() == "" {
		t.Error("RuleOverview() is empty")
	}
	for _, sec := range []Section{
		SectionTurn, SectionCardType, SectionKeyword, SectionAbility, SectionEffect,
	} {
		if RuleSectionIntro(sec) == "" {
			t.Errorf("section %q has no intro", sec)
		}
	}
	if RuleSectionIntro(SectionCombat) != "" {
		t.Error("SectionCombat is not expected to carry an intro")
	}
}
