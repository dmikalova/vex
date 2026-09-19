package web

import "strings"

// This file is the /cards gallery's text-search mini-language (ADR 0023). A query
// string filters cards by their name and rules text with a small, hand-rolled
// grammar — no search library, to keep the WASM bundle small:
//
//   term term   all-of (AND): every term must match
//   a|b         either (OR): matches when any alternative matches
//   -x          exclude: matches only when x does NOT appear
//   "two words" phrase: matches the exact run, spaces and all
//   \"          a literal quote to search for (\\ is a literal backslash)
//
// A term matches as a case-insensitive substring of the haystack (a card's name
// and rendered rules text). Terms are ANDead; alternatives within one term are
// ORed; an excluded term negates. An empty query matches everything.

// queryTerm is one space-separated unit of a parsed query: a set of OR
// alternatives that all share the term's negation. A plain term has one
// alternative; `a|b` has two; `-x` negates.
type queryTerm struct {
	negate bool
	alts   []string // lower-cased; the term matches if any alternative matches
}

// parseQuery turns a raw query string into its terms. Whitespace separates terms
// except inside a quoted run; `|` separates alternatives within a term except
// inside quotes; a leading `-` on an unquoted term negates it; `\` escapes the
// next character (so `\"` searches for a quote). A term that parses to nothing
// (a lone `-`, or empty quotes) is dropped.
func parseQuery(q string) []queryTerm {
	var terms []queryTerm
	var alts []string       // completed alternatives of the current term
	var cur strings.Builder // the alternative being built
	negate, started := false, false
	inQuote, escaped := false, false

	flushAlt := func() {
		alts = append(alts, cur.String())
		cur.Reset()
	}
	flushTerm := func() {
		if started {
			flushAlt()
		}
		var nonEmpty []string
		for _, a := range alts {
			if a != "" {
				nonEmpty = append(nonEmpty, normalizeSearch(a))
			}
		}
		if len(nonEmpty) > 0 {
			terms = append(terms, queryTerm{
				negate: negate,
				alts:   nonEmpty,
			})
		}
		alts, negate, started = nil, false, false
	}

	for _, r := range q {
		switch {
		case escaped:
			cur.WriteRune(r)
			escaped, started = false, true
		case r == '\\':
			escaped = true
		case r == '"':
			inQuote = !inQuote
			started = true
		case inQuote:
			cur.WriteRune(r)
			started = true
		case r == ' ' || r == '\t' || r == '\n':
			flushTerm()
		case r == '|':
			flushAlt()
			started = true
		case r == '-' && !started && len(alts) == 0 && !negate:
			negate = true // a leading '-' negates; it is not part of the term text
		default:
			cur.WriteRune(r)
			started = true
		}
	}
	flushTerm()
	return terms
}

// normalizeSearch lower-cases text and folds the Æmber ligature to "ae" so a
// player can type "aember" to find "Æmber". Both the query and the haystack pass
// through it, so the two are compared on the same footing.
func normalizeSearch(s string) string {
	return strings.ReplaceAll(strings.ToLower(s), "æ", "ae")
}

// matchesQuery reports whether the haystack satisfies every term of the parsed
// query: each non-negated term has an alternative that appears, and no negated
// term's alternative appears. An empty query (no terms) matches.
func matchesQuery(terms []queryTerm, haystack string) bool {
	h := normalizeSearch(haystack)
	for _, t := range terms {
		hit := false
		for _, a := range t.alts {
			if strings.Contains(h, a) {
				hit = true
				break
			}
		}
		if hit == t.negate {
			return false
		}
	}
	return true
}
