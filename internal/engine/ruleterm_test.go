package engine

import (
	"slices"
	"testing"
)

// TestClosedCatalogsAreComplete enforces ADR 0018: each closed catalog the game
// defines is complete in the rulebook registry. A member with no term fails the
// build here rather than going silently undescribed (the gap that once left
// Elusive and Taunt out of the rulebook). Turn phases are covered too. Effects
// are no longer exempt: the node census closes them from the other end, and
// TestCatalogTermsAreRegistered / TestCensusSectionTermsAreClaimed bind its rows
// to the registry in both directions. Combat is deliberately not a step catalog:
// fighting is an action taken during the main phase, described in its own
// section, not a turn step.
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
		for _, ct := range allCardTypes() {
			if name := ct.String(); !have[name] {
				t.Errorf("card type %q has no rulebook term (ADR 0018)", name)
			}
		}
	})

	t.Run("bonus icons", func(t *testing.T) {
		have := titled(SectionBonus)
		for _, b := range allBonusIcons() {
			if name := b.String(); !have[name] {
				t.Errorf("bonus icon %q has no rulebook term (ADR 0018)", name)
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

	t.Run("turn phases", func(t *testing.T) {
		// Phase terms all share the Title "Turn structure" and differ by Subtitle,
		// so completeness is keyed on the subtitle each phase owns (rulebookStep).
		have := map[string]bool{}
		for _, term := range RuleTerms() {
			if term.Section == SectionTurn {
				have[term.Subtitle] = true
			}
		}
		for _, p := range Phases() {
			if !have[p.rulebookStep()] {
				t.Errorf("phase %q has no rulebook term (ADR 0018)", p)
			}
		}
	})
}

// censusTermSections are the rulebook sections the node census claims one for
// one: Card text holds the sentence parts a node leans on (Target, For Each,
// Duration), Effects holds the verbs. A term filed here is owed a row, which is
// what TestCensusSectionTermsAreClaimed enforces. A row may still bind to a term
// in any other section — Ward is a keyword, forging is a turn step — so the
// forward check in TestCatalogTermsAreRegistered matches across the whole
// registry.
var censusTermSections = []Section{SectionCardText, SectionEffect}

// TestCatalogTermsAreRegistered is the forward half of the census's binding to
// the rulebook (ADR 0018): a row that names a Term must name a title some
// registered term carries, so a classification cannot point at prose nobody
// wrote. The match is across the whole registry rather than one section, because
// a node binds to the term a player would actually look up — the ward nodes to
// the Keyword section's Ward, a forge node to the Turn section's forge step.
func TestCatalogTermsAreRegistered(t *testing.T) {
	registered := map[string]bool{}
	for _, term := range RuleTerms() {
		registered[term.Title] = true
	}
	check := func(t *testing.T, rows []FamilyRow) {
		t.Helper()
		for _, row := range rows {
			if title := row.Rules.Term; title != "" && !registered[title] {
				t.Errorf("%s names rulebook term %q, which no term carries; "+
					"write it or correct the row (ADR 0018)", row.Type, title)
			}
		}
	}
	for _, family := range Families() {
		t.Run(family.Name, func(t *testing.T) { check(t, family.Rows) })
	}
	for _, enum := range Enums() {
		t.Run(enum.Type, func(t *testing.T) { check(t, enum.Rows) })
	}
}

// TestCensusSectionTermsAreClaimed is the reverse half: every term in the
// sections the census claims must be named by at least one row. A term left
// behind by a deleted or renamed node fails here instead of rotting in the
// rulebook. A rule that genuinely has no census member — a standing restriction a
// card carries rather than an effect node — belongs in another section, which is
// why "Must Fight When Used" is filed under Combat.
func TestCensusSectionTermsAreClaimed(t *testing.T) {
	claimed := map[string]bool{}
	claim := func(rows []FamilyRow) {
		for _, row := range rows {
			claimed[row.Rules.Term] = true
		}
	}
	for _, family := range Families() {
		claim(family.Rows)
	}
	for _, enum := range Enums() {
		claim(enum.Rows)
	}
	for _, term := range RuleTerms() {
		if !slices.Contains(censusTermSections, term.Section) || claimed[term.Title] {
			continue
		}
		t.Errorf("term %q in section %q is claimed by no census row; "+
			"give a node the term, or file the term in another section (ADR 0018)",
			term.Title, term.Section)
	}
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
		SectionTurn, SectionCardType, SectionKeyword, SectionBonus, SectionAbility,
		SectionCardText, SectionEffect,
	} {
		if RuleSectionIntro(sec) == "" {
			t.Errorf("section %q has no intro", sec)
		}
	}
	if RuleSectionIntro(SectionCombat) != "" {
		t.Error("SectionCombat is not expected to carry an intro")
	}
}
