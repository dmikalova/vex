package engine

import (
	"fmt"
	"strings"
	"unicode"
)

// This file holds the low-level, card-agnostic string helpers the text renderers
// share: punctuation, article and plural inflection, small-number words, and
// phrase trimming. They take plain strings and ints — no CardDefinition, Effect,
// or game state — so they read as a self-contained toolkit apart from the
// card-shaped rendering in text.go.

// damageAmount renders a quantity of damage for card text. Damage is lowercase
// everywhere it is dealt, and every renderer that names an amount routes through
// here so the wording cannot drift between effects.
func damageAmount(n int) string { return fmt.Sprintf("%d damage", n) }

// dealDamageTo renders the full clause "deal <n> damage to <target>".
func dealDamageTo(n int, target string) string {
	return "deal " + damageAmount(n) + " to " + target
}

// consideredBlank renders the blanking predicate, which reads the same for one
// text box and for many, so the two renderers that blank cannot drift apart.
const consideredBlank = "considered blank, except for traits"

// punctuate ends an ability body with a period. A body that already ends in a
// period is left alone; one that ends in a closing quote (an embedded ability
// such as Charge!'s granted "Play: ...") takes its period inside the quote, so
// the line reads `... an enemy creature."` rather than doubling or misplacing it.
func punctuate(body string) string {
	if strings.HasSuffix(body, `"`) {
		if inner := body[:len(body)-1]; !strings.HasSuffix(inner, ".") {
			return inner + `."`
		}
		return body
	}
	if strings.HasSuffix(body, ".") {
		return body
	}
	return body + "."
}

// oxfordAnd joins parts into an English list with a serial comma, e.g. "a, b, and
// c"; two parts read "a and b" and one reads as itself.
func oxfordAnd(parts []string) string {
	switch len(parts) {
	case 0:
		return ""
	case 1:
		return parts[0]
	case 2:
		return parts[0] + " and " + parts[1]
	default:
		return strings.Join(parts[:len(parts)-1], ", ") + ", and " + parts[len(parts)-1]
	}
}

// ordinalWord renders a small positive integer as its English ordinal word,
// covering the key ordinals the Key Imps bar (first, second, third).
func ordinalWord(n int) string {
	switch n {
	case 1:
		return "first"
	case 2:
		return "second"
	case 3:
		return "third"
	default:
		return fmt.Sprintf("%dth", n)
	}
}

// countWord renders a small count as an English word ("one") for the common
// single-card grant, falling back to the numeral for larger counts.
func countWord(n int) string {
	if n == 1 {
		return "one"
	}
	return fmt.Sprintf("%d", n)
}

// capitalizeFirst upper-cases the first rune of s.
func capitalizeFirst(s string) string {
	if s == "" {
		return s
	}
	r := []rune(s)
	r[0] = unicode.ToUpper(r[0])
	return string(r)
}

// repeatedSubjects are the noun phrases a printed line may name more than once in
// one sentence: the demonstratives a Target renders for a creature or card the
// sentence already picked out. A card's own name is deliberately absent — it is
// often preceded by another creature noun, so "in a fight with Chonkers, give it
// +1 power counters" would hand "it" to the creature that was destroyed. The
// self-reference collapses only where it prints as "this creature", which a
// granted ability's substitution produces (see grantedAbilityText).
var repeatedSubjects = []string{"this creature", "that creature", "this card", "that card"}

// collapseRepeatedSubject replaces a subject's second and later mention within one
// sentence with "it", so a sentence names a card once and then refers back to it —
// "ready it" (Rocket Boots), "deal 2 damage to it" (Bonerot Venom).
//
// Scope resets at a sentence end and at a quote boundary: an ability quoted inside
// another card's text is its own sentence and may not borrow the outer subject, so
// `This creature gains, "Reap: Ward this creature."` keeps both mentions.
//
// A possessive mention stays spelled out, because "it" carries no case and would
// attach to the nearest preceding noun instead: Pain Reaction's "that creature's
// neighbors" after "this damage" would read as the damage's neighbors.
//
// Pinned by TestCollapsesRepeatedSubject.
func collapseRepeatedSubject(line string) string {
	var out strings.Builder
	start := 0
	for i := 0; i < len(line); i++ {
		if !endsSubjectScope(line, i) {
			continue
		}
		out.WriteString(collapseSubjectsInScope(line[start : i+1]))
		start = i + 1
	}
	out.WriteString(collapseSubjectsInScope(line[start:]))
	return out.String()
}

// endsSubjectScope reports whether the byte at i closes a scope a pronoun may not
// reach across: a quote mark, or the period ending a sentence another follows.
func endsSubjectScope(line string, i int) bool {
	if line[i] == '"' {
		return true
	}
	return line[i] == '.' && i+1 < len(line) && line[i+1] == ' '
}

// collapseSubjectsInScope collapses every repeatable subject within one sentence.
func collapseSubjectsInScope(scope string) string {
	for _, subject := range repeatedSubjects {
		scope = collapseSubjectMentions(scope, subject)
	}
	return scope
}

// collapseSubjectMentions rewrites every mention of one subject after the first,
// matching case-insensitively so a sentence-initial "This creature" still counts
// as the mention the later ones refer back to.
func collapseSubjectMentions(scope, subject string) string {
	var out strings.Builder
	rest, seen := scope, false
	for {
		i := strings.Index(strings.ToLower(rest), subject)
		if i < 0 {
			break
		}
		out.WriteString(rest[:i])
		end := i + len(subject)
		if seen && !strings.HasPrefix(rest[end:], "'s") {
			out.WriteString("it")
		} else {
			out.WriteString(rest[i:end])
			seen = true
		}
		rest = rest[end:]
	}
	out.WriteString(rest)
	return out.String()
}

// lowerFirst lower-cases the first rune of s, for folding an effect's own
// sentence-cased Text() into the middle of a longer sentence.
func lowerFirst(s string) string {
	if s == "" {
		return s
	}
	r := []rune(s)
	r[0] = unicode.ToLower(r[0])
	return string(r)
}

// indefinite prefixes a noun with the indefinite article "a" or "an", choosing
// "an" before a word that starts with a vowel — e.g. "an Urchin", "a Knight". A
// noun already led by "another" carries its own article ("an other") and takes
// none, so "another creature" is left as is rather than "an another creature".
func indefinite(noun string) string {
	if noun == "" {
		return noun
	}
	if noun == "another" || strings.HasPrefix(noun, "another ") {
		return noun
	}
	// "an" + "other …" reads as "another …" — the merged indefinite of a noun the
	// Other flag prefixes ("other friendly creature" → "another friendly creature").
	if strings.HasPrefix(noun, "other ") {
		return "an" + noun
	}
	switch unicode.ToLower([]rune(noun)[0]) {
	case 'a', 'e', 'i', 'o', 'u':
		return "an " + noun
	default:
		return "a " + noun
	}
}

// plural gives a noun the form a count of n calls for: "card" for one, "cards"
// for any other number, including zero.
func plural(n int, noun string) string {
	if n == 1 {
		return noun
	}
	return noun + "s"
}

// countNoun renders a quantity and the noun it counts, e.g. "1 card", "3 cards".
// Reach for it instead of writing "%d cards", which renders "1 cards" the first
// time a card or game passes 1 — the bug that shipped as "1 keys". Mass nouns
// (damage, power, armor, Æmber) take a bare %d; they never pluralize. Guarded by
// TestLogLinesNeverPrintOneAsPlural and the countNoun rule in universalRetired.
func countNoun(n int, noun string) string {
	return fmt.Sprintf("%d %s", n, plural(n, noun))
}

// singularNoun strips the leading article or quantifier from a Target's phrase,
// leaving the bare noun a "ready a <noun>" or "up to 3 <noun>s" clause needs. The
// adjectives stay: "an enemy damaged creature" becomes "enemy damaged creature",
// which pluralizes correctly and keeps the "enemy" the card is scoped to.
func singularNoun(phrase string) string {
	if rest, ok := strings.CutPrefix(phrase, "another "); ok {
		return "other " + rest
	}
	for _, p := range []string{"each other ", "each ", "an ", "a "} {
		if after, ok := strings.CutPrefix(phrase, p); ok {
			return after
		}
	}
	return phrase
}
