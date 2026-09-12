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

// HouseSelectorKind is how a MayPlayOrUse names the houses its grant frees.
type HouseSelectorKind uint8

const (
	houseSelectorUnset HouseSelectorKind = iota
	// SelectHouse frees one house: HouseSelector.House when set, or the house an
	// enclosing ChooseHouseThen picked when HouseNone.
	SelectHouse
	// SelectAny frees every house (Follow the Leader, Scientifical Hack).
	SelectAny
	// SelectExcept frees every house but HouseSelector.House (the Star Alliance
	// "non-Star Alliance card" cycle).
	SelectExcept
	// SelectControlled frees every house the controller has a card in play for
	// (United Action).
	SelectControlled
)

// HouseSelector is the Houses axis of a MayPlayOrUse: it selects whose cards the
// grant frees. It is flat, comparable state (a small enum plus a House), so it
// lives in the snapshotable GameState (ADR 0005).
type HouseSelector struct {
	Kind  HouseSelectorKind
	House House
}

// MayPlayOrUse lets the controller act with cards outside their active house for
// the remainder of the turn — the one node for every out-of-house permission
// grant. Its four axes fold what were five wordings: Houses selects whose cards
// (a named or chosen house, any house, every house but one, or every house you
// control), Grant selects the verbs it frees (play, use, or fight), Types narrows
// the card types (the zero value frees all), and Count bounds how many cards the
// grant frees (zero is unlimited). The grant lasts only the current turn (the
// ready phase clears it).
type MayPlayOrUse struct {
	Houses HouseSelector
	Grant  HouseGrant
	Types  CardTypes
	Count  int
}

// validate rejects a grant that names no houses, frees no verb, or bounds a
// negative count.
func (e MayPlayOrUse) validate() error {
	if e.Houses.Kind == houseSelectorUnset {
		return fmt.Errorf("MayPlayOrUse: houses must be set")
	}
	if e.Grant == 0 {
		return fmt.Errorf("MayPlayOrUse: at least one grant must be set")
	}
	if e.Count < 0 {
		return fmt.Errorf("MayPlayOrUse: count must not be negative")
	}
	return nil
}

// Text renders the grant as its KeyForge clause, narrowing to the shortest wording
// the axes select — "may fight", "may use", "may play or use" — over the houses,
// types, and count the grant reaches.
func (e MayPlayOrUse) Text() string {
	switch e.Houses.Kind {
	case SelectControlled:
		return "for the remainder of the turn, you may play cards from any house for which you have a card in play"
	case SelectExcept:
		verb := "play"
		if e.Grant&GrantUse != 0 {
			verb = "play or use"
		}
		return "you may " + verb + " " + e.exceptObject() + " this turn"
	default: // SelectHouse, SelectAny
		if e.Grant == GrantFight {
			return "for the remainder of the turn, " + e.fightSubject() + " may fight"
		}
		if e.Houses.Kind == SelectAny {
			return "for the remainder of the turn, you may use friendly artifacts as if they belonged to the active house"
		}
		return "for the remainder of the turn, you may " + e.namedVerbObject()
	}
}

// fightSubject renders who a fight grant frees: every friendly creature, a named
// house's, or the chosen house's.
func (e MayPlayOrUse) fightSubject() string {
	if e.Houses.Kind == SelectAny {
		return "each friendly creature"
	}
	if e.Houses.House == HouseNone {
		return "each friendly creature of the chosen house"
	}
	return "each friendly " + e.Houses.House.String() + " creature"
}

// namedVerbObject renders the verb-and-object of a named-house use/play grant, e.g.
// "use friendly Sanctum creatures" or "play or use a Mars card".
func (e MayPlayOrUse) namedVerbObject() string {
	if e.Grant&GrantPlay != 0 {
		verb := "play"
		if e.Grant&GrantUse != 0 {
			verb = "play or use"
		}
		return verb + " a " + e.Houses.House.String() + " card"
	}
	return "use friendly " + e.Houses.House.String() + " creatures"
}

// exceptObject renders the card an exclusion grant frees, e.g. "a non-Star Alliance
// artifact, upgrade, or Tactic" (Types narrowed) or "one non-Star Alliance card"
// (all types).
func (e MayPlayOrUse) exceptObject() string {
	house := ""
	if e.Houses.House != HouseNone {
		house = "non-" + e.Houses.House.String() + " "
	}
	if e.Types.all() {
		return "one " + house + "card"
	}
	return "a " + house + e.Types.list()
}

// Resolve records the this-turn grant for the controller, resolving a chosen-house
// selector against the house an enclosing ChooseHouseThen picked.
func (e MayPlayOrUse) Resolve(ctx *EffectContext) {
	houses := e.Houses
	if houses.Kind == SelectHouse && houses.House == HouseNone {
		houses.House = ctx.ChosenHouse
	}
	ctx.Resolver.GrantMayPlayOrUse(ctx.Controller, houses, e.Grant, e.Types, e.Count)
}
