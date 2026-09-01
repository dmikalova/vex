package engine

import (
	"fmt"
)

// CardsDiscarded is a Condition met when the specified player has discarded at
// least Amount cards of the given House from hand this turn. Amount must be at
// least 1: a check for "discarded 0 or more cards" is always true, so an unset
// threshold is rejected at registration rather than silently treated as one.
type CardsDiscarded struct {
	Player Player
	House  House
	Amount int
}

// Value counts the matching House cards the player discarded from hand this turn.
func (e CardsDiscarded) Value(ctx *EffectContext) int {
	return countOfHouse(ctx, ctx.Resolver.DiscardedThisTurn(ctx.PlayerFor(e.Player)), e.House)
}

// Met reports whether at least Amount matching cards were discarded.
func (e CardsDiscarded) Met(ctx *EffectContext) bool { return e.Value(ctx) >= e.Amount }

// validate rejects a non-positive Amount: "discarded 0 or more" is always met, so
// an omitted threshold is an authoring error, not a silent default.
func (e CardsDiscarded) validate() error {
	if e.Amount < 1 {
		return fmt.Errorf("CardsDiscarded: Amount must be at least 1")
	}
	return nil
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
// "3 Untamed cards" for more.
func (e CardsDiscarded) discardPhrase() string {
	if e.Amount == 1 {
		return indefinite(e.House.String() + " card")
	}
	return fmt.Sprintf("%d %s cards", e.Amount, e.House.String())
}
