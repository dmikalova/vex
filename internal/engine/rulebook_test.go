package engine

import (
	"strings"
	"testing"
)

// TestRuleBookRealRegistry checks the assembled book over the live registry: the
// spine leads with The Turn, every section carries at least one group, and a known
// term is present under its section.
func TestRuleBookRealRegistry(t *testing.T) {
	book := RuleBook()
	if len(book) == 0 {
		t.Fatal("RuleBook() is empty")
	}
	if book[0].Title != "The Turn" {
		t.Fatalf("first section = %q, want The Turn", book[0].Title)
	}
	var sawTurnStructure bool
	for _, sec := range book {
		if len(sec.Groups) == 0 {
			t.Errorf("section %q has no groups", sec.Title)
		}
		if sec.Key != SectionTurn {
			continue
		}
		for _, gr := range sec.Groups {
			if gr.Title == "Turn structure" {
				sawTurnStructure = true
				assertTurnStepsOrdered(t, gr)
			}
		}
	}
	if !sawTurnStructure {
		t.Error("Turn structure group missing from The Turn section")
	}
}

// assertTurnStepsOrdered checks the numbered turn steps render in subtitle order.
func assertTurnStepsOrdered(t *testing.T, gr RuleGroup) {
	t.Helper()
	var subs []string
	for _, p := range gr.Parts {
		if p.Subtitle != "" {
			subs = append(subs, p.Subtitle)
		}
	}
	for i := 1; i < len(subs); i++ {
		if subs[i-1] > subs[i] {
			t.Errorf("turn steps out of order: %q before %q", subs[i-1], subs[i])
		}
	}
}

// TestAssembleRuleBook drives the ordering and skip rules over synthetic input:
// an intro-only section is kept, a section with neither terms nor intro is
// dropped, groups sort with the section-titled group leading, and a group's
// untitled parts precede its subtitle-sorted parts.
func TestAssembleRuleBook(t *testing.T) {
	terms := []RuleTerm{
		{Section: SectionKeyword, Title: "Elusive", Body: "b1"},
		{Section: SectionKeyword, Title: "Keywords", Body: "lead"},
		{Section: SectionKeyword, Title: "Assault", Body: "b2"},
		{Section: SectionTurn, Title: "Turn structure", Subtitle: "2. b", Body: "s2"},
		{Section: SectionTurn, Title: "Turn structure", Subtitle: "1. a", Body: "s1"},
		{Section: SectionTurn, Title: "Turn structure", Body: "intro-part"},
	}
	intros := map[Section]string{
		SectionKeyword: "kw intro",
		SectionCombat:  "combat intro", // intro-only section, no terms
	}

	book := assembleRuleBook(terms, intros)

	// The Turn, Combat (intro only), and Keywords survive; cardtype/ability/effect
	// have neither terms nor intro and are dropped.
	if got := len(book); got != 3 {
		t.Fatalf("assembled %d sections, want 3", got)
	}
	if book[0].Key != SectionTurn || book[1].Key != SectionCombat || book[2].Key != SectionKeyword {
		t.Fatalf("section order = %v", []Section{book[0].Key, book[1].Key, book[2].Key})
	}
	if book[1].Intro != "combat intro" || len(book[1].Groups) != 0 {
		t.Errorf(
			"combat should be intro-only, got intro=%q groups=%d",
			book[1].Intro,
			len(book[1].Groups),
		)
	}

	// Keywords: the section-titled "Keywords" group leads, then Assault, Elusive.
	kw := book[2].Groups
	if len(kw) != 3 || kw[0].Title != "Keywords" || kw[1].Title != "Assault" ||
		kw[2].Title != "Elusive" {
		t.Fatalf("keyword groups = %v", groupTitles(kw))
	}

	// Turn structure: untitled lead part first, then subtitles in order.
	parts := book[0].Groups[0].Parts
	if len(parts) != 3 || parts[0].Subtitle != "" || parts[1].Subtitle != "1. a" ||
		parts[2].Subtitle != "2. b" {
		t.Fatalf("turn parts = %v", partSubs(parts))
	}
}

func groupTitles(gs []RuleGroup) []string {
	var out []string
	for _, g := range gs {
		out = append(out, g.Title)
	}
	return out
}

func partSubs(ps []RuleTerm) []string {
	var out []string
	for _, p := range ps {
		out = append(out, p.Subtitle)
	}
	return out
}

// TestAssembleGlossary drives dedupe and sort: a title collapses to one entry
// keeping the first Definition, a later Definition fills an earlier empty one, and
// entries come out alphabetical regardless of registration order.
func TestAssembleGlossary(t *testing.T) {
	terms := []RuleTerm{
		{Title: "Taunt", Definition: ""},
		{Title: "Elusive", Definition: "dodge the first fight"},
		{Title: "Taunt", Definition: "protects its neighbors"},
		{Title: "Assault", Definition: "damage on the way in"},
		{Title: "Elusive", Definition: "ignored second gloss"},
	}

	gl := assembleGlossary(terms)

	if len(gl) != 3 {
		t.Fatalf("glossary has %d entries, want 3", len(gl))
	}
	if gl[0].Title != "Assault" || gl[1].Title != "Elusive" || gl[2].Title != "Taunt" {
		t.Fatalf("glossary order = %v", []string{gl[0].Title, gl[1].Title, gl[2].Title})
	}
	if gl[1].Definition != "dodge the first fight" {
		t.Errorf("Elusive kept later gloss: %q", gl[1].Definition)
	}
	if gl[2].Definition != "protects its neighbors" {
		t.Errorf("Taunt empty-then-filled failed: %q", gl[2].Definition)
	}
}

// TestGlossaryRealRegistry checks the live glossary is non-empty, sorted, and
// carries a known term.
func TestGlossaryRealRegistry(t *testing.T) {
	gl := Glossary()
	if len(gl) == 0 {
		t.Fatal("Glossary() is empty")
	}
	var sawElusive bool
	for i, e := range gl {
		if e.Title == "Elusive" {
			sawElusive = true
		}
		if i > 0 && strings.ToLower(gl[i-1].Title) > strings.ToLower(e.Title) {
			t.Errorf("glossary out of order: %q before %q", gl[i-1].Title, e.Title)
		}
	}
	if !sawElusive {
		t.Error("Elusive missing from the glossary")
	}
}
