package engine

import (
	"fmt"
	"slices"
)

// SearchForName lets the controller search their deck and discard pile for a card
// with a specific name, reveal it, and put it into their hand — Help from Future
// Self tutoring a Timetraveller. Nothing happens if no matching card is found.
//
//rulebook:effect Search for Named Card
type SearchForName struct {
	Name string
	// All takes every copy found instead of one the controller chooses, which
	// leaves nothing to choose and so asks nothing (Bear Flute).
	All bool
}

// Text renders the effect, e.g. "search your deck and discard pile for a
// Timetraveller, reveal it, and put it into your hand".
func (e SearchForName) Text() string {
	if e.All {
		return fmt.Sprintf(
			"search your deck and discard pile and put each %s from them into your hand",
			e.Name,
		)
	}
	return fmt.Sprintf(
		"search your deck and discard pile for %s, reveal it, and put it into your hand",
		indefinite(e.Name),
	)
}

// Resolve gathers the deck and discard cards with the name, lets the controller
// choose one, reveals it, and moves it to their hand from whichever zone it is in.
func (e SearchForName) Resolve(ctx *EffectContext) { e.resolveGate(ctx) }

// resolveGate searches and reports whether it found anything, so a Then can hang
// a follow-up off the search succeeding (Bear Flute reshuffles only if it did).
func (e SearchForName) resolveGate(ctx *EffectContext) bool {
	inDeck := nameMatches(ctx, ctx.Resolver.Deck(ctx.Controller), e.Name)
	candidates := slices.Concat(
		inDeck,
		nameMatches(ctx, ctx.Resolver.Discard(ctx.Controller), e.Name),
	)
	if e.All {
		for _, id := range candidates {
			e.take(ctx, inDeck, id)
		}
		return len(candidates) > 0
	}
	id, ok := ctx.ChooseCreature("Choose "+indefinite(e.Name)+" to put into your hand", candidates)
	if !ok {
		return false
	}
	e.take(ctx, inDeck, id)
	return true
}

// take reveals one found card and moves it to hand from whichever zone holds it.
func (e SearchForName) take(ctx *EffectContext, inDeck []LocalID, id LocalID) {
	ctx.Resolver.Record(CardsRevealedToAll{Player: ctx.Controller, Cards: []LocalID{id}})
	if slices.Contains(inDeck, id) {
		ctx.Resolver.MoveFromDeckToHand(id)
	} else {
		ctx.Resolver.PutFromDiscardIntoHand(id)
	}
}
