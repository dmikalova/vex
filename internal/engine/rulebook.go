package engine

import (
	"sort"
	"strings"
)

// This file assembles the flat rule-term registry (ruleterm.go) into the ordered
// shape every rulebook surface renders from. It is the shared facet behind the
// Markdown generator, the /rulebook page, and the /glossary page: RuleBook groups
// the terms into ordered sections, and Glossary flattens the same registry into an
// alphabetical Title→Definition list. Because all three surfaces draw from here,
// the ordering, headings, and term set cannot drift between them.

// ruleSpineEntry is one section of the rulebook spine: its registry key and the
// heading it prints under. The spine fixes the section order and their titles in
// one place so no surface can reorder or rename a section on its own.
type ruleSpineEntry struct {
	key   Section
	title string
}

// ruleSpine is the rulebook's sections in render order, each paired with its
// printed heading.
var ruleSpine = []ruleSpineEntry{
	{SectionTurn, "The Turn"},
	{SectionCombat, "Combat"},
	{SectionCardType, "Card Types"},
	{SectionKeyword, "Keywords"},
	{SectionAbility, "Abilities"},
	{SectionEffect, "Effects"},
}

// RuleSection is one assembled section of the rulebook: its printed heading, its
// framing intro, and the term groups filed under it in render order.
type RuleSection struct {
	Key    Section
	Title  string
	Intro  string
	Groups []RuleGroup
}

// RuleGroup gathers every term sharing a Title under one heading. Parts render in
// order: the untitled lead parts first, in the order they were registered, then
// the subheaded parts sorted by Subtitle (the numbered turn steps).
type RuleGroup struct {
	Title string
	Parts []RuleTerm
}

// RuleBook assembles the registered rule terms into ordered sections — the shared,
// rendering-agnostic shape behind every rulebook surface. A section with no terms
// and no intro is omitted.
func RuleBook() []RuleSection {
	return assembleRuleBook(ruleTerms, ruleSectionIntros)
}

// assembleRuleBook is RuleBook's engine over explicit inputs, so its ordering and
// skip rules can be tested without reaching for the package-level registry.
func assembleRuleBook(terms []RuleTerm, intros map[Section]string) []RuleSection {
	byKey := map[Section][]RuleTerm{}
	for _, t := range terms {
		byKey[t.Section] = append(byKey[t.Section], t)
	}
	var out []RuleSection
	for _, sp := range ruleSpine {
		sectionTerms := byKey[sp.key]
		intro := intros[sp.key]
		if len(sectionTerms) == 0 && intro == "" {
			continue
		}
		out = append(out, RuleSection{
			Key:    sp.key,
			Title:  sp.title,
			Intro:  intro,
			Groups: groupTerms(sectionTerms, sp.title),
		})
	}
	return out
}

// groupTerms merges terms sharing a Title into groups, orders the groups
// alphabetically with the section-titled group (a section's own overview entry)
// leading, and orders each group's parts.
func groupTerms(terms []RuleTerm, sectionTitle string) []RuleGroup {
	var groups []RuleGroup
	at := map[string]int{}
	for _, t := range terms {
		i, ok := at[t.Title]
		if !ok {
			i = len(groups)
			at[t.Title] = i
			groups = append(groups, RuleGroup{Title: t.Title})
		}
		groups[i].Parts = append(groups[i].Parts, t)
	}
	for i := range groups {
		groups[i].Parts = orderParts(groups[i].Parts)
	}
	sort.SliceStable(groups, func(i, j int) bool {
		return strings.ToLower(groups[i].Title) < strings.ToLower(groups[j].Title)
	})
	sort.SliceStable(groups, func(i, j int) bool {
		return strings.EqualFold(groups[i].Title, sectionTitle) &&
			!strings.EqualFold(groups[j].Title, sectionTitle)
	})
	return groups
}

// orderParts puts a group's untitled lead parts first, in registration order,
// then its subheaded parts sorted by subtitle.
func orderParts(parts []RuleTerm) []RuleTerm {
	var lead, subbed []RuleTerm
	for _, p := range parts {
		if p.Subtitle == "" {
			lead = append(lead, p)
		} else {
			subbed = append(subbed, p)
		}
	}
	sort.SliceStable(subbed, func(i, j int) bool {
		return strings.ToLower(subbed[i].Subtitle) < strings.ToLower(subbed[j].Subtitle)
	})
	return append(lead, subbed...)
}

// GlossaryEntry is one term as the glossary lists it: its Title and the one-line
// Definition. A term whose Definition is not yet written yields an empty string.
type GlossaryEntry struct {
	Title      string
	Definition string
}

// Glossary returns every rule term as an alphabetical Title→Definition list, one
// entry per title — the glossary facet of the same registry RuleBook draws from.
// Terms sharing a title collapse to one entry, keeping the first Definition given.
func Glossary() []GlossaryEntry {
	return assembleGlossary(ruleTerms)
}

// assembleGlossary is Glossary's engine over explicit input, so its dedupe and
// sort can be tested without the package-level registry.
func assembleGlossary(terms []RuleTerm) []GlossaryEntry {
	at := map[string]int{}
	var out []GlossaryEntry
	for _, t := range terms {
		if i, ok := at[t.Title]; ok {
			if out[i].Definition == "" {
				out[i].Definition = t.Definition
			}
			continue
		}
		at[t.Title] = len(out)
		out = append(out, GlossaryEntry{Title: t.Title, Definition: t.Definition})
	}
	sort.SliceStable(out, func(i, j int) bool {
		return strings.ToLower(out[i].Title) < strings.ToLower(out[j].Title)
	})
	return out
}
