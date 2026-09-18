package engine

import "fmt"

// HouseMatchKind is how a HouseMatcher names the houses it admits. Its zero value
// (MatchAnyHouse) admits every house: a house filter is usually absent, so the
// permissive case is the one an author leaves unwritten (ADR 0038).
type HouseMatchKind uint8

const (
	// MatchAnyHouse admits every house. It is the zero value, so an unset matcher
	// filters nothing.
	MatchAnyHouse HouseMatchKind = iota
	// MatchNamedHouse admits only cards of HouseMatcher.House ("a Mars card").
	MatchNamedHouse
	// MatchExceptHouse admits every card but those of HouseMatcher.House ("a
	// non-Star Alliance card").
	MatchExceptHouse
	// MatchChosenHouse admits cards of the house an enclosing ChooseHouseThen picked.
	MatchChosenHouse
	// MatchActiveHouse admits cards of the active house.
	MatchActiveHouse
	// MatchContextualHouse admits cards sharing the house of the card in context
	// (ctx.It) — the card a preceding effect put in focus.
	MatchContextualHouse
	// MatchEachHouse admits cards of the house an enclosing ForEachHouse is on. It
	// reads the same ctx.ChosenHouse binding as MatchChosenHouse but renders "of
	// that house" — the loop's house, not one the player picked.
	MatchEachHouse
)

// HouseMatcher is the one way a filter or target names which houses it admits: a
// flat, comparable value (a kind plus a house) rendered as a noun qualifier that
// each effect welds into its own phrase (ADR 0038). It carries the per-card house
// notions only; the grant-only "houses you control" and the set-relative "house
// with the most creatures" are deliberately not kinds here. The House field is
// exported, so a SelfHouse sentinel in it resolves by reflection with no per-type
// seam (self_house.go).
type HouseMatcher struct {
	Kind  HouseMatchKind
	House House
}

// filters reports whether the matcher narrows anything (false for the zero value).
func (m HouseMatcher) filters() bool { return m.Kind != MatchAnyHouse }

// matches reports whether the card is admitted by the matcher.
func (m HouseMatcher) matches(ctx *EffectContext, id LocalID) bool {
	switch m.Kind {
	case MatchNamedHouse:
		return ctx.Resolver.House(id) == m.House
	case MatchExceptHouse:
		return ctx.Resolver.House(id) != m.House
	case MatchChosenHouse, MatchEachHouse:
		return ctx.Resolver.House(id) == ctx.ChosenHouse
	case MatchActiveHouse:
		return ctx.Resolver.House(id) == ctx.Resolver.ActiveHouse()
	case MatchContextualHouse:
		return ctx.HasIt && ctx.Resolver.House(id) == ctx.Resolver.House(ctx.It)
	default: // MatchAnyHouse
		return true
	}
}

// qualify weaves the house into a base noun, e.g. "Mars card", "non-Logos card",
// or "creature of the chosen house". The any-house matcher returns the noun
// unchanged. It is qualifyNoun followed by qualifyPhrase, so an effect that
// renders its subject in one piece (the pile selections) gets both the prefix
// kinds and the suffix kinds from one call.
func (m HouseMatcher) qualify(noun string) string {
	return m.qualifyPhrase(m.qualifyNoun(noun))
}

// adjective returns the house's prefix adjective — "Mars", "non-Sanctum" — and
// whether the matcher has one. Only the prefix kinds do; the suffix kinds render
// after the noun and the any-house matcher renders nothing.
func (m HouseMatcher) adjective() (string, bool) {
	switch m.Kind {
	case MatchNamedHouse:
		return m.House.String(), true
	case MatchExceptHouse:
		return "non-" + m.House.String(), true
	}
	return "", false
}

// qualifyNoun prefixes the house onto a base noun for the prefix kinds, e.g.
// "Mars creature" or "non-Sanctum creature". The suffix kinds and any-house leave
// the noun unchanged, to be rendered later by qualifyPhrase — Target needs the two
// halves apart because power and flank qualifiers sit between the noun and the
// suffix kinds' "of the chosen house".
func (m HouseMatcher) qualifyNoun(noun string) string {
	if adj, ok := m.adjective(); ok {
		return adj + " " + noun
	}
	return noun
}

// qualifyPhrase suffixes the house onto a rendered phrase for the suffix kinds,
// e.g. "each creature of the chosen house". The prefix kinds and any-house leave
// the phrase unchanged.
func (m HouseMatcher) qualifyPhrase(phrase string) string {
	switch m.Kind {
	case MatchChosenHouse:
		return phrase + " of the chosen house"
	case MatchActiveHouse, MatchEachHouse:
		return phrase + " of that house"
	case MatchContextualHouse:
		return phrase + " of that card's house"
	default:
		return phrase
	}
}

// validate rejects a matcher whose kind names a house but leaves it unset.
func (m HouseMatcher) validate() error {
	if (m.Kind == MatchNamedHouse || m.Kind == MatchExceptHouse) && m.House == HouseNone {
		return fmt.Errorf("HouseMatcher: a named or non-house matcher needs a house")
	}
	return nil
}
