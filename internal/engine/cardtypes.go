package engine

import "strings"

// CardTypes is a set of card types a grant reaches — the Types axis of MayPlayOrUse.
// The zero value is the empty set, which every reader treats as "all card types":
// a grant names types only to narrow itself (Scientifical Hack to artifacts, Com.
// Officer Kirby to everything but creatures). It is flat, comparable state (a
// bitset), so it lives in the snapshotable GameState (ADR 0005).
type CardTypes uint8

// CardTypesOf builds a set from the given card types.
func CardTypesOf(types ...CardType) CardTypes {
	var s CardTypes
	for _, t := range types {
		s |= 1 << t
	}
	return s
}

// all reports that the set names no type, so it admits every card type.
func (s CardTypes) all() bool { return s == 0 }

// has reports whether the set admits a card of the given type; the empty set
// admits every type.
func (s CardTypes) has(t CardType) bool {
	return s.all() || s&(1<<t) != 0
}

// list renders the admitted types as a printed noun phrase in rulebook order,
// e.g. "Artifact, Upgrade, or Tactic" or "Artifact". The empty set renders "card".
func (s CardTypes) list() string {
	if s.all() {
		return "card"
	}
	var words []string
	for _, t := range []CardType{Creature, Artifact, Upgrade, Tactic} {
		if s.has(t) {
			words = append(words, typeWord(t))
		}
	}
	return joinOr(words)
}

// typeWord is a card type's printed word in a filter or listing: the lowercase
// type name (creature, artifact, upgrade, tactic). Card types read as lowercase
// common nouns in card text, like every other rendered noun.
func typeWord(t CardType) string {
	return strings.ToLower(t.String())
}

// playablePlural renders the admitted types as a lowercase plural noun phrase,
// e.g. "upgrades" or "artifacts or upgrades", for a play-permission clause. The
// empty set, which admits everything, renders "cards".
func (s CardTypes) playablePlural() string {
	if s.all() {
		return "cards"
	}
	var words []string
	for _, t := range []CardType{Creature, Artifact, Upgrade, Tactic} {
		if s.has(t) {
			words = append(words, typeWord(t)+"s")
		}
	}
	return joinOr(words)
}
