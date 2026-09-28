package engine

import (
	"fmt"
)

// CardsDiscarded is a Condition met when the specified player has discarded at
// least Amount cards its Filter admits from hand this turn. Amount must be at
// least 1: a check for "discarded 0 or more cards" is always true, so an unset
// threshold is rejected at registration rather than silently treated as one.
type CardsDiscarded struct {
	// Player names whose discards to count. It picks a turn log rather than
	// narrowing a card, so it is a set selector and not a filter axis.
	Player Player
	// Filter narrows which discarded cards count and supplies the noun the clause
	// prints. The log holds real cards, so the filter is enforceable — but a
	// discarded card is in no battleline, so validate rejects an in-play axis.
	Filter Filter
	Amount int
}

// Value counts the cards the filter admits that the player discarded from hand
// this turn.
func (e CardsDiscarded) Value(ctx *EffectContext) int {
	return countMatching(
		ctx, ctx.Resolver.DiscardedThisTurn(ctx.PlayerFor(e.Player)), e.Filter)
}

// Met reports whether at least Amount matching cards were discarded.
func (e CardsDiscarded) Met(ctx *EffectContext) bool { return e.Value(ctx) >= e.Amount }

// validate rejects a non-positive Amount: "discarded 0 or more" is always met, so
// an omitted threshold is an authoring error, not a silent default.
func (e CardsDiscarded) validate() error {
	if e.Amount < 1 {
		return fmt.Errorf("CardsDiscarded: Amount must be at least 1")
	}
	return e.Filter.validateIdentityOnly("CardsDiscarded")
}

// CondText renders the condition text.
func (e CardsDiscarded) CondText() string {
	player, whose := "you have", "your"
	if e.Player == Opponent {
		player, whose = "your opponent has", "their"
	}
	return fmt.Sprintf(
		"if %s discarded %s from %s hand this turn",
		player,
		e.discardPhrase(),
		whose,
	)
}

// discardPhrase renders the required discards: "an Untamed card" for one, or
// "3 Untamed cards" for more. With a zero filter the qualifier is omitted ("a
// card", "3 cards"), so an Earthbind-style "discarded a card" reads naturally.
func (e CardsDiscarded) discardPhrase() string {
	noun := e.Filter.noun("card")
	if e.Amount == 1 {
		return indefinite(noun)
	}
	return fmt.Sprintf("%d %ss", e.Amount, noun)
}
