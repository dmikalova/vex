package engine

import "fmt"

// This file holds the counts that read what a player did this turn — the cards
// they played and the creatures they used — split out of effect_count.go.

// CardsPlayed counts the cards of a house a player has played this turn. Like
// InPlay it serves two roles from one description: as a Count it yields the tally
// and renders a "for each ... you have played this turn" noun; as a Condition it is
// met once Amount have been played (Epic Quest fires after seven), defaulting to
// one. This replaces a bespoke "played at least N of a house" condition.
type CardsPlayed struct {
	// Player names whose plays to count; House filters by the played card's house.
	Player Player
	House  HouseMatcher
	// Amount is the minimum the Condition role requires; zero means at least one.
	Amount int
}

// Value counts the player's cards of the house played this turn.
func (e CardsPlayed) Value(ctx *EffectContext) int {
	return countOfHouse(ctx, ctx.Resolver.PlayedThisTurn(ctx.PlayerFor(e.Player)), e.House)
}

// countOfHouse counts how many of the ids the matcher admits — the filter a
// turn-log Count applies to the unfiltered record the engine keeps. An any-house
// matcher admits every card, so a Count can ask "how many cards" as well as "how
// many Sanctum cards".
func countOfHouse(ctx *EffectContext, ids []LocalID, house HouseMatcher) int {
	n := 0
	for _, id := range ids {
		if house.matches(ctx, id) {
			n++
		}
	}
	return n
}

// Met reports whether at least Amount (default one) matching cards were played.
func (e CardsPlayed) Met(ctx *EffectContext) bool { return e.Value(ctx) >= e.threshold() }

// threshold is the Condition's required count, defaulting to one.
func (e CardsPlayed) threshold() int {
	if e.Amount < 1 {
		return 1
	}
	return e.Amount
}

// CountText renders the singular noun the "for each" clause repeats.
func (e CardsPlayed) CountText() string {
	return e.cardNoun() + " you have played this turn"
}

// cardNoun is the noun the count repeats, house-qualified when the count filters
// by house and a plain "card" when it counts every card played.
func (e CardsPlayed) cardNoun() string {
	return e.House.qualifyNoun("card")
}

// CondText renders the condition, e.g. "if you have played 7 or more Sanctum cards
// this turn".
func (e CardsPlayed) CondText() string {
	return fmt.Sprintf("if you have played %d or more %ss this turn", e.threshold(), e.cardNoun())
}

// CountClause renders the clause CountIs puts after "if", e.g. "you played
// exactly 1 card this turn".
func (e CardsPlayed) CountClause(quantity string, plural bool) string {
	noun := e.cardNoun()
	if plural {
		noun += "s"
	}
	return fmt.Sprintf("you played %s %s this turn", quantity, noun)
}

// CreaturesUsed counts the creatures a player has used this turn — reaped,
// fought, or fired an "Action:" with. It reads the per-creature use tally the
// engine already keeps, so only creatures still in play are counted.
type CreaturesUsed struct {
	Player Player
}

// Value counts the player's creatures in play that have been used this turn.
func (e CreaturesUsed) Value(ctx *EffectContext) int {
	n := 0
	for _, id := range ctx.Resolver.Battleline(ctx.PlayerFor(e.Player)) {
		if ctx.Resolver.TimesUsedThisTurn(id) > 0 {
			n++
		}
	}
	return n
}

// CountText renders the singular noun the "for each" clause repeats.
func (e CreaturesUsed) CountText() string { return "creature you used this turn" }

// CountClause renders the clause CountIs puts after "if", e.g. "you used 3 or
// more creatures this turn".
func (e CreaturesUsed) CountClause(quantity string, plural bool) string {
	noun := "creature"
	if plural {
		noun += "s"
	}
	return fmt.Sprintf("you used %s %s this turn", quantity, noun)
}
