package engine

import (
	"fmt"
	"slices"
)

// SearchForName lets the controller search their deck and discard pile for a card
// with a specific name, reveal it, and put it into their hand — Help from Future
// Self tutoring a Timetraveller. Nothing happens if no matching card is found.
type SearchForName struct {
	// Name is the card name searched for.
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

// SearchDeck is the KeyForge "search" keyword: the controller searches their deck
// for a card — any card, or one of a given House — puts it into their hand, and
// then shuffles their deck (Orb of Wonder searches for any card, Saurus Rex for a
// Saurian card). A House-restricted search reveals the card it takes. The deck is
// always shuffled, even when nothing was taken.
type SearchDeck struct {
	// House restricts the search to cards of that house; HouseNone searches for any
	// card and does not reveal what it takes.
	House House
}

// Text renders the effect, ending in "then shuffle your deck" so the shuffle is
// always stated (e.g. "search your deck for a Saurian card, reveal it, and put it
// into your hand, then shuffle your deck").
func (e SearchDeck) Text() string {
	if e.House == HouseNone {
		return "search your deck for a card and put it into your hand" +
			", then shuffle your deck"
	}
	return "search your deck for " + indefinite(e.House.String()+" card") +
		", reveal it, and put it into your hand, then shuffle your deck"
}

// Resolve gathers the deck cards matching the House filter, lets the controller
// choose one to put into their hand (revealing it when the search was restricted),
// and then shuffles the deck regardless of whether a card was taken.
func (e SearchDeck) Resolve(ctx *EffectContext) {
	var cands []LocalID
	for _, id := range ctx.Resolver.Deck(ctx.Controller) {
		if e.House == HouseNone || ctx.Resolver.House(id) == e.House {
			cands = append(cands, id)
		}
	}
	if id, ok := ctx.ChooseCard("Choose a card to put into your hand", cands); ok {
		if e.House != HouseNone {
			ctx.Resolver.Record(CardsRevealedToAll{Player: ctx.Controller, Cards: []LocalID{id}})
		}
		ctx.Resolver.MoveFromDeckToHand(id)
	}
	ctx.Resolver.Shuffle(ctx.Controller)
	ctx.Resolver.Record(DeckShuffled{Player: ctx.Controller})
}
