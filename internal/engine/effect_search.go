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
}

// Text renders the effect, e.g. "search your deck and discard pile for a
// Timetraveller, reveal it, and put it into your hand".
func (e SearchForName) Text() string {
	return fmt.Sprintf("search your deck and discard pile for %s, reveal it, and put it into your hand", indefinite(e.Name))
}

// Resolve gathers the deck and discard cards with the name, lets the controller
// choose one, reveals it, and moves it to their hand from whichever zone it is in.
func (e SearchForName) Resolve(ctx *EffectContext) {
	inDeck := nameMatches(ctx, ctx.Resolver.Deck(ctx.Controller), e.Name)
	candidates := append(inDeck, nameMatches(ctx, ctx.Resolver.Discard(ctx.Controller), e.Name)...)
	id, ok := ctx.ChooseCreature("Choose "+indefinite(e.Name)+" to put into your hand", candidates)
	if !ok {
		return
	}
	ctx.Resolver.Logf("%s is revealed", ctx.Resolver.Name(id))
	if slices.Contains(inDeck, id) {
		ctx.Resolver.MoveFromDeckToHand(id)
	} else {
		ctx.Resolver.PutFromDiscardIntoHand(id)
	}
}
