package engine

import "fmt"

// HouseGrant is a bitset of what a MayPlayOrUse frees this turn: playing the
// house's cards from hand, using its creatures in play, or fighting with them.
// Fight is the narrow case of use — a grant that only lets creatures fight, not
// reap or fire "Action:" abilities.
type HouseGrant uint8

const (
	// GrantPlay lets the controller play the freed cards from hand.
	GrantPlay HouseGrant = 1 << iota
	// GrantUse lets the controller use (fight, reap, or Action:) the freed creatures.
	GrantUse
	// GrantFight lets the controller's freed creatures fight only.
	GrantFight
)

// HouseSelector is the Houses axis of a MayPlayOrUse: it selects whose cards the
// grant frees. It is a grant-only superset of the filter-side HouseMatcher —
// Match names the houses the way any filter would (a named or chosen house, any
// house, or every house but one), while Controlled adds the one selection only a
// grant can make: every house the controller has a card in play for (United
// Action). Keeping Controlled off HouseMatcher means a filter field can never
// express it. It is flat, comparable state, so it lives in the snapshotable
// GameState (ADR 0005).
type HouseSelector struct {
	Match      HouseMatcher
	Controlled bool
}

// MayPlayOrUse lets the controller act with cards outside their active house for
// the remainder of the turn — the one node for every out-of-house permission
// grant. Its axes fold what were several wordings: Houses selects whose cards
// (a named or chosen house, any house, every house but one, or every house you
// control), Trait scopes a use grant to a creature trait instead of a house
// (Mutagenic Serum's "use friendly Mutant creatures"), Grant selects the verbs it
// frees (play, use, or fight), Types narrows the card types (the zero value frees
// all), and Cards bounds how many cards the grant frees (zero is unlimited). The
// grant lasts only the current turn (the ready phase clears it).
type MayPlayOrUse struct {
	Houses HouseSelector
	Trait  Trait
	Grant  HouseGrant
	Types  CardTypes
	Cards  int
}

// validate rejects a grant that frees no verb or bounds a negative count.
func (e MayPlayOrUse) validate() error {
	if e.Grant == 0 {
		return fmt.Errorf("MayPlayOrUse: at least one grant must be set")
	}
	if e.Cards < 0 {
		return fmt.Errorf("MayPlayOrUse: Cards must not be negative")
	}
	return nil
}

// Text renders the grant as its KeyForge clause, narrowing to the shortest wording
// the axes select — "may fight", "may use", "may play or use" — over the houses,
// types, and count the grant reaches.
func (e MayPlayOrUse) Text() string {
	remainder := durationClause(RemainderOfPlayerTurn, "")
	if e.Trait != traitUnset {
		return remainder + ", you may use friendly " +
			e.Trait.String() + " creatures"
	}
	if e.Houses.Controlled {
		return remainder + ", you may play cards from any house for which you have a card in play"
	}
	switch e.Houses.Match.Kind {
	case MatchExceptHouse:
		// An except-house grant keeps KeyForge's shorter permission phrasing ("you
		// may play a non-Logos card this turn"), the printed wording every such card
		// carries — the same remainder-of-turn window as the branches below, said in
		// the permission voice rather than the leading duration clause.
		verb := "play"
		if e.Grant&GrantUse != 0 {
			verb = "play or use"
		}
		return "you may " + verb + " " + e.exceptObject() + " this turn"
	default: // MatchNamedHouse, MatchChosenHouse, MatchAnyHouse
		if e.Grant == GrantFight {
			return remainder + ", " + e.fightSubject() + " may fight"
		}
		if e.Houses.Match.Kind == MatchAnyHouse {
			return remainder + ", you may use friendly artifacts as if they belonged to the active house"
		}
		return remainder + ", you may " + e.namedVerbObject()
	}
}

// fightSubject renders who a fight grant frees: every friendly creature, a named
// house's, or the chosen house's.
func (e MayPlayOrUse) fightSubject() string {
	switch e.Houses.Match.Kind {
	case MatchAnyHouse:
		return "each friendly creature"
	case MatchChosenHouse:
		return "each friendly creature of the chosen house"
	default:
		return "each friendly " + e.Houses.Match.House.String() + " creature"
	}
}

// namedVerbObject renders the verb-and-object of a named-house use/play grant, e.g.
// "use friendly Sanctum creatures" or "play or use a Mars card".
func (e MayPlayOrUse) namedVerbObject() string {
	if e.Grant&GrantPlay != 0 {
		verb := "play"
		if e.Grant&GrantUse != 0 {
			verb = "play or use"
		}
		return verb + " a " + e.Houses.Match.House.String() + " card"
	}
	return "use friendly " + e.Houses.Match.House.String() + " creatures"
}

// exceptObject renders the card an exclusion grant frees, e.g. "a non-Star Alliance
// artifact, upgrade, or Tactic" (Types narrowed) or "one non-Star Alliance card"
// (all types).
func (e MayPlayOrUse) exceptObject() string {
	house := ""
	if e.Houses.Match.House != HouseNone {
		house = "non-" + e.Houses.Match.House.String() + " "
	}
	if e.Types.all() {
		return "one " + house + "card"
	}
	return "a " + house + e.Types.list()
}

// Resolve records the this-turn grant for the controller, resolving a chosen-house
// selector against the house an enclosing ChooseHouseThen picked.
func (e MayPlayOrUse) Resolve(ctx *EffectContext) {
	if e.Trait != traitUnset {
		ctx.Resolver.GrantMayUseTrait(ctx.Controller, e.Trait)
		return
	}
	houses := e.Houses
	if houses.Match.Kind == MatchChosenHouse && houses.Match.House == HouseNone {
		houses.Match.House = ctx.ChosenHouse
	}
	ctx.Resolver.GrantMayPlayOrUse(ctx.Controller, houses, e.Grant, e.Types, e.Cards)
}
