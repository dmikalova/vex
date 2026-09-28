package engine

import (
	"fmt"
	"strings"
)

// ExcessCreatures counts how many more creatures one player controls than the
// other (never below zero). Player names whose excess is counted: Opponent for
// "each creature your opponent controls in excess of you" (Glorious Few),
// Controller for "each creature you have in excess of your opponent"
// (Unguarded Camp).
type ExcessCreatures struct {
	Player Player
	// NotCountingSelf excludes the source creature from its controller's side of
	// the comparison — Dr. Milli's "in excess of you, not counting Dr. Milli".
	NotCountingSelf bool
	// Filter narrows which creatures count, on both sides alike; the zero value
	// counts every creature (Spare Arm Carmine compares Mutant counts). The scan
	// walks battlelines, so the creature noun is implied and a card names only what
	// narrows it further.
	Filter Filter
}

// Value returns the named player's creature count minus the other's, floored at 0.
func (e ExcessCreatures) Value(ctx *EffectContext) int {
	more := ctx.PlayerFor(e.Player)
	moreCount := e.sideCount(ctx, more)
	lessCount := e.sideCount(ctx, 1-more)
	if e.NotCountingSelf {
		if ctx.Controller == more {
			moreCount = max(0, moreCount-1)
		} else {
			lessCount = max(0, lessCount-1)
		}
	}
	return max(0, moreCount-lessCount)
}

// sideCount counts the creatures one player controls that the filter admits.
func (e ExcessCreatures) sideCount(ctx *EffectContext, player int) int {
	n := 0
	for _, id := range ctx.Resolver.Battleline(player) {
		if e.Filter.matches(ctx, id) {
			n++
		}
	}
	return n
}

// noun renders the counted creature as a singular noun phrase, over "creature" as
// the base the scan implies by walking battlelines: "creature", "Mutant creature".
func (e ExcessCreatures) noun() string {
	return e.Filter.noun("creature")
}

// CountText renders the singular noun the "for each" clause repeats.
func (e ExcessCreatures) CountText() string {
	base := "creature you have in excess of your opponent"
	if e.Player == Opponent {
		base = "creature your opponent controls in excess of you"
	}
	if e.NotCountingSelf {
		base += ", not counting " + SelfName
	}
	return base
}

// CardsInPlay selects the cards a player has in play that its Filter admits and
// serves two roles from one description. As a Count (a Per clause) it yields the
// match count and renders the repeated "for each ..." noun; as a Condition it is
// met when the count reaches Amount, which defaults to one. This unifies the
// several friendly-creature counts and conditions, e.g.
// CardsInPlay{Player: Controller, Filter: Filter{Type: Creature}} or the
// house-filtered CardsInPlay{Player: Controller,
// Filter: Filter{Type: Creature, House: namedHouse(Mars)}}.
type CardsInPlay struct {
	// Player names whose cards to count (Controller, Opponent, or EachPlayer). It
	// picks whose battleline to scan rather than narrowing what a card must be, so
	// it is a set selector and not a filter axis.
	Player Player
	// Filter narrows which of those cards count — their type, house, trait, name,
	// power, and per-card state — and supplies the noun the two text roles print.
	// The zero value counts every card in play.
	Filter Filter
	// None inverts the Condition role: it is met when nothing matches. It reads as
	// its own word rather than as Amount 0 because "if there are no X" is a
	// different sentence from "if there are N X", not the N = 0 case of it.
	None bool
	// Amount is the minimum the Condition role requires; zero means at least one.
	Amount int
}

// Value counts the matching cards the player has in play.
func (e CardsInPlay) Value(ctx *EffectContext) int {
	n := 0
	for _, id := range e.set(ctx) {
		if e.Filter.matches(ctx, id) {
			n++
		}
	}
	return n
}

// Met reports whether at least Amount (default one) matching cards are in play,
// or, under None, that none are.
func (e CardsInPlay) Met(ctx *EffectContext) bool {
	if e.None {
		return e.Value(ctx) == 0
	}
	return e.Value(ctx) >= e.threshold()
}

// set returns every in-play id the count considers, for one player or both. The
// Filter's type axis narrows creatures, artifacts, or upgrades out of this
// universe, so the zones need no type-specific selection here.
func (e CardsInPlay) set(ctx *EffectContext) []LocalID {
	if e.Player == EachPlayer {
		return append(e.playerSet(ctx, 0), e.playerSet(ctx, 1)...)
	}
	return e.playerSet(ctx, ctx.PlayerFor(e.Player))
}
func (e CardsInPlay) playerSet(ctx *EffectContext, p int) []LocalID {
	return resolverCardsInPlay(ctx, p)
}

// threshold is the Condition's required count, defaulting to one.
func (e CardsInPlay) threshold() int {
	if e.Amount < 1 {
		return 1
	}
	return e.Amount
}

// who renders the controlling side as "friendly" or "enemy".
func (e CardsInPlay) who() string {
	switch e.Player {
	case Opponent:
		return "enemy"
	case EachPlayer:
		return ""
	default:
		return "friendly"
	}
}

// base is the noun the filter qualifies. A count that names no type counts
// "cards", except that a trait names a kind of card by itself, so a trait-only
// count passes no base and reads "friendly Shard" rather than "friendly Shard
// card" (Shard of Life).
func (e CardsInPlay) base() string {
	if e.Filter.Type == TypeUnset && e.Filter.Trait != traitUnset {
		return ""
	}
	return "card"
}

// noun renders the "<side> [ready ][house ]<type>" phrase the text roles share.
// The side is the count's own word, since Player picks a battleline rather than
// narrowing a card; everything else comes from the filter.
func (e CardsInPlay) noun() string {
	if e.Filter.Name != "" {
		// A proper name identifies one specific card, so it needs no side: "if there
		// are no Ancient Bears in play" (Bear Flute).
		return e.Filter.Name
	}
	parts := []string{}
	if e.Filter.Except.filters() {
		parts = append(parts, "other")
	}
	if who := e.who(); who != "" {
		parts = append(parts, who)
	}
	parts = append(parts, e.Filter.qualifyNoun(e.base()))
	return e.Filter.clauses(strings.Join(parts, " "))
}

// CountText renders the singular noun the "for each" clause repeats. A
// house- or trait-filtered count reads "friendly Mars creature" / "friendly
// Shard"; an unfiltered one adds "in play" to distinguish it from cards in hand.
func (e CardsInPlay) CountText() string {
	f := e.Filter
	if (f.House.filters() || f.Trait != traitUnset ||
		f.Power.filters() || f.Aember.filters()) && e.Player != EachPlayer {
		return e.noun()
	}
	return e.noun() + " in play"
}

// cardinalCountText renders the count as "the number of friendly Mars creatures
// you control", the cardinal form for a clause that compares against it.
func (e CardsInPlay) cardinalCountText() string {
	return "the number of " + plural(2, e.noun()) + " " + e.controls()
}

// controls renders which side's board the cardinal count reads from.
func (e CardsInPlay) controls() string {
	switch e.Player {
	case Opponent:
		return "your opponent controls"
	case EachPlayer:
		return "in play"
	default:
		return "you control"
	}
}

// CondText renders the condition, e.g. "if there is a friendly creature in play"
// or "if there are 2 or more friendly creatures in play".
func (e CardsInPlay) CondText() string {
	if e.None {
		return fmt.Sprintf("if there are no %s in play", plural(0, e.noun()))
	}
	if n := e.threshold(); n > 1 {
		return fmt.Sprintf("if there are %d or more %s in play", n, plural(n, e.noun()))
	}
	return fmt.Sprintf("if there is %s in play", indefinite(e.noun()))
}
