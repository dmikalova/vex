package engine

// RuleTerm is one entry in the rulebook: a self-describing rule co-located with
// the code it governs. The engine owns the terms as typed data (ADR 0018) so the
// rulebook generator renders them instead of re-parsing Go source, and the
// completeness test can fail the build when a closed catalog — every Keyword,
// card type, or trigger — has a member no term describes.
type RuleTerm struct {
	// Section files the term under one part of the rulebook spine; its value is
	// the section key the renderer groups by (see the Section constants).
	Section Section
	// Title is the term's heading, e.g. "Capture Æmber". Terms that share a Title
	// render under one heading, so several code sites can feed one entry.
	Title string
	// Subtitle, when set, groups the term beneath its Title as a subheading, e.g.
	// the numbered steps under "Turn structure".
	Subtitle string
	// Definition is the one-line, Rules-voice gloss that doubles as the glossary
	// entry. It is optional for now; the glossary that surfaces it is future work.
	Definition string
	// Body is the term's Markdown prose, rendered verbatim under its heading.
	Body string
}

// Section identifies which part of the rulebook spine a term belongs to. The
// constant values are the section keys the renderer buckets by.
type Section string

// The rulebook sections, in spine order.
const (
	SectionTurn     Section = "turn"
	SectionCombat   Section = "combat"
	SectionCardType Section = "cardtype"
	SectionKeyword  Section = "keyword"
	SectionAbility  Section = "ability"
	SectionEffect   Section = "effect"
)

// ruleTerms is the accumulated registry. Each ruleterms_<section>.go file
// registers its section's terms at init time; RuleTerms hands the whole table to
// the renderer and the completeness test.
var ruleTerms []RuleTerm

// registerRuleTerms adds terms to the registry. Each section's file calls it once
// from an init function, so a section's terms keep their listed order regardless
// of the order Go runs the files' init functions.
func registerRuleTerms(terms []RuleTerm) {
	ruleTerms = append(ruleTerms, terms...)
}

// RuleTerms returns every registered rulebook term. The rulebook generator
// renders it and the completeness test checks the closed catalogs against it.
func RuleTerms() []RuleTerm {
	return ruleTerms
}

// ruleOverview is the rulebook's document preamble — the title and core-concepts
// front matter that frames every section. It belongs to no one part of the spine,
// so unlike a section intro it registers on its own.
var ruleOverview string

// registerRuleOverview records the rulebook's preamble Markdown. ruleterms_overview.go
// calls it once from init.
func registerRuleOverview(md string) {
	ruleOverview = md
}

// RuleOverview returns the rulebook's document preamble, rendered above the first
// section.
func RuleOverview() string {
	return ruleOverview
}

// ruleSectionIntros holds each section's intro Markdown, keyed by section. A
// section with no intro (Combat) is simply absent.
var ruleSectionIntros = map[Section]string{}

// registerRuleSectionIntro records a section's intro Markdown. Each
// ruleterms_<section>.go calls it once from init, beside that section's terms, so
// the framing prose lives with the code it frames (ADR 0018) rather than in a
// loose Markdown file.
func registerRuleSectionIntro(sec Section, md string) {
	ruleSectionIntros[sec] = md
}

// RuleSectionIntro returns a section's intro Markdown, or "" when the section has
// none.
func RuleSectionIntro(sec Section) string {
	return ruleSectionIntros[sec]
}
