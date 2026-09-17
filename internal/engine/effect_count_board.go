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
	// Trait counts only creatures with this trait on both sides; the unset zero
	// value counts every creature (Pismire compares Mutant counts).
	Trait Trait
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

// sideCount counts the creatures one player controls, restricted to Trait when set.
func (e ExcessCreatures) sideCount(ctx *EffectContext, player int) int {
	f := e.filter()
	n := 0
	for _, id := range ctx.Resolver.Battleline(player) {
		if f.admits(ctx.Resolver, id) {
			n++
		}
	}
	return n
}

// filter is the identity predicate a counted creature must satisfy — its Trait
// when set, admitting every creature otherwise.
func (e ExcessCreatures) filter() CardFilter {
	return CardFilter{Trait: e.Trait}
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

// InPlay selects the cards a player has in play that match its filters — of a
// given type and house — and serves two roles from one description. As a Count (a
// Per clause) it yields the match count and renders the repeated "for each ..."
// noun; as a Condition it is met when the count reaches Amount, which defaults to
// one. This unifies the several friendly-creature counts and conditions, e.g.
// InPlay{Player: Controller, Type: Creature} or the house-filtered
// InPlay{Player: Controller, Type: Creature, House: namedHouse(Mars)}.
type InPlay struct {
	// Player names whose cards to count (Controller or Opponent).
	Player Player
	// Type filters by card type; the zero value counts any type.
	Type CardType
	// House filters by house; the zero value (MatchAnyHouse) counts any house
	// (ADR 0038).
	House HouseMatcher
	// Trait filters by trait; the unset zero value counts any trait, and a set
	// trait replaces the rendered noun with it ("friendly Shard").
	Trait Trait
	// Ready counts only cards that are ready (not exhausted).
	Ready bool
	// Damaged counts only creatures that have damage on them.
	Damaged bool
	// WithAember counts only creatures that have Æmber on them (Faust the Great
	// counts each friendly creature with Æmber on it).
	WithAember bool
	// MinPower counts only creatures whose power is at least this value; zero (the
	// unset default) applies no power floor (Grump Buggy counts power 5 or higher).
	MinPower int
	// Other leaves the source card itself out of the count, which is how a card
	// counts its companions (Phylyx the Disintegrator).
	Other bool
	// Name filters by printed card name, and replaces the rendered noun with it:
	// "if there are no Ancient Bears in play" (Bear Flute).
	Name string
	// None inverts the Condition role: it is met when nothing matches. It reads as
	// its own word rather than as Amount 0 because "if there are no X" is a
	// different sentence from "if there are N X", not the N = 0 case of it.
	None bool
	// Amount is the minimum the Condition role requires; zero means at least one.
	Amount int
}

// Value counts the matching cards the player has in play.
func (e InPlay) Value(ctx *EffectContext) int {
	f := e.filter()
	n := 0
	for _, id := range e.set(ctx) {
		if !e.House.matches(ctx, id) {
			continue
		}
		if !f.admits(ctx.Resolver, id) {
			continue
		}
		if e.Ready && ctx.Resolver.Exhausted(id) {
			continue
		}
		if e.WithAember && ctx.Resolver.AmberOn(id) == 0 {
			continue
		}
		if e.Damaged && ctx.Resolver.Damage(id) == 0 {
			continue
		}
		if e.MinPower > 0 && ctx.Resolver.Power(id) < e.MinPower {
			continue
		}
		if e.Other && id == ctx.Source {
			continue
		}
		n++
	}
	return n
}

// filter is the identity predicate a counted card must satisfy, conjoining the
// Type, Trait, and Name filters (Chosen's filter has the same shape).
func (e InPlay) filter() CardFilter {
	return CardFilter{Type: e.Type, Trait: e.Trait, Name: e.Name}
}

// Met reports whether at least Amount (default one) matching cards are in play,
// or, under None, that none are.
func (e InPlay) Met(ctx *EffectContext) bool {
	if e.None {
		return e.Value(ctx) == 0
	}
	return e.Value(ctx) >= e.threshold()
}

// set returns every in-play id the count considers — the battleline and the
// artifact row together, for one player or both. The Type filter (via filter)
// narrows creatures or artifacts out of this universe, so the zones need no
// type-specific selection here.
func (e InPlay) set(ctx *EffectContext) []LocalID {
	if e.Player == EachPlayer {
		return append(e.playerSet(ctx, 0), e.playerSet(ctx, 1)...)
	}
	return e.playerSet(ctx, ctx.PlayerFor(e.Player))
}
func (e InPlay) playerSet(ctx *EffectContext, p int) []LocalID {
	return append(ctx.Resolver.Battleline(p), ctx.Resolver.Artifacts(p)...)
}

// threshold is the Condition's required count, defaulting to one.
func (e InPlay) threshold() int {
	if e.Amount < 1 {
		return 1
	}
	return e.Amount
}

// who renders the controlling side as "friendly" or "enemy".
func (e InPlay) who() string {
	switch e.Player {
	case Opponent:
		return "enemy"
	case EachPlayer:
		return ""
	default:
		return "friendly"
	}
}

// typeNoun renders the filtered type as a noun. A trait filter names the trait
// itself ("Shard"), and combines with a type as "Thief creature".
func (e InPlay) typeNoun() string {
	if e.Trait != traitUnset {
		switch e.Type {
		case Creature:
			return e.Trait.String() + " creature"
		case Artifact:
			return e.Trait.String() + " artifact"
		default:
			return e.Trait.String()
		}
	}
	switch e.Type {
	case Creature:
		return "creature"
	case Artifact:
		return "artifact"
	default:
		return "card"
	}
}

// noun renders the "<side> [ready ][house ]<type>" phrase the text roles share.
func (e InPlay) noun() string {
	if e.Name != "" {
		return e.Name
	}
	parts := []string{}
	if e.Other {
		parts = append(parts, "other")
	}
	if who := e.who(); who != "" {
		parts = append(parts, who)
	}
	if e.Ready {
		parts = append(parts, "ready")
	}
	if e.Damaged {
		parts = append(parts, "damaged")
	}
	noun := strings.Join(append(parts, e.House.qualifyNoun(e.typeNoun())), " ")
	noun = e.House.qualifyPhrase(noun)
	if e.MinPower > 0 {
		noun += fmt.Sprintf(" with power %d or higher", e.MinPower)
	}
	if e.WithAember {
		noun += " with \u00c6mber on it"
	}
	return noun
}

// CountText renders the singular noun the "for each" clause repeats. A
// house- or trait-filtered count reads "friendly Mars creature" / "friendly
// Shard"; an unfiltered one adds "in play" to distinguish it from cards in hand.
func (e InPlay) CountText() string {
	if (e.House.filters() || e.Trait != traitUnset || e.MinPower > 0 || e.WithAember) &&
		e.Player != EachPlayer {
		return e.noun()
	}
	return e.noun() + " in play"
}

// cardinalCountText renders the count as "the number of friendly Mars creatures
// you control", the cardinal form for a clause that compares against it.
func (e InPlay) cardinalCountText() string {
	return "the number of " + plural(2, e.noun()) + " " + e.controls()
}

// controls renders which side's board the cardinal count reads from.
func (e InPlay) controls() string {
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
func (e InPlay) CondText() string {
	if e.None {
		return fmt.Sprintf("if there are no %s in play", plural(0, e.noun()))
	}
	if n := e.threshold(); n > 1 {
		return fmt.Sprintf("if there are %d or more %s in play", n, plural(n, e.noun()))
	}
	return fmt.Sprintf("if there is %s in play", indefinite(e.noun()))
}
