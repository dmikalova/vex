package engine

import (
	"fmt"
	"slices"
)

// PutFromPlay takes each card its Target selects out of play and puts it in a
// destination zone — the top of its owner's deck, their hand, or their archives —
// shedding the per-match state the card built up in play (damage, spent armor,
// Aember on it, upgrades). The destination is required. Moving a card out of play
// this way is how a "Destroyed:" ability can save its own creature: the creature
// leaves for the named zone as it is destroyed, so it never reaches the discard
// pile. When several cards move to the top of the deck at once the controller
// chooses the order they stack.
//
//rulebook:effect Put from Play
type PutFromPlay struct {
	Target      Target
	Destination Destination
}

// Text renders the effect, e.g. "put each artifact on top of its owner's deck" or
// "put this creature into its owner's hand".
func (e PutFromPlay) Text() string {
	switch e.Destination {
	case ToTopOfDeck:
		return fmt.Sprintf("put %s on top of its owner's deck", e.Target.Text())
	case ToDeckShuffled:
		return fmt.Sprintf("shuffle %s into its owner's deck", e.Target.Text())
	case ToArchives:
		return fmt.Sprintf("put %s into its owner's archives", e.Target.Text())
	default:
		return fmt.Sprintf("put %s into its owner's hand", e.Target.Text())
	}
}

// validate rejects a destination this effect cannot move a card to; only the hand,
// the top of the deck, and the archives are supported, and the destination must be
// named.
func (e PutFromPlay) validate() error {
	if !e.Target.valid() {
		return errUnsetTarget("PutFromPlay")
	}
	switch e.Destination {
	case ToHand, ToTopOfDeck, ToDeckShuffled, ToArchives:
		return nil
	default:
		return fmt.Errorf("PutFromPlay: unsupported destination %d", e.Destination)
	}
}

// Resolve moves each selected card from play to the destination. Cards headed to
// the top of the deck are stacked in an order the controller chooses.
func (e PutFromPlay) Resolve(ctx *EffectContext) {
	ids := e.Target.Select(ctx)
	if e.Destination == ToTopOfDeck {
		ids = ctx.OrderByChoice("Choose the next card to put on top of the deck", ids)
	}
	for _, id := range ids {
		moveFromPlayTo(ctx, id, e.Destination)
	}
}

// moveFromPlayTo puts one card from play into a destination zone.
func moveFromPlayTo(ctx *EffectContext, id LocalID, d Destination) {
	switch d {
	case ToTopOfDeck:
		ctx.Resolver.PutOnTopOfDeck(id)
	case ToArchives:
		ctx.Resolver.PutIntoArchives(id)
	case ToDeckShuffled:
		ctx.Resolver.PutIntoDeckShuffled(id)
	default:
		ctx.Resolver.PutIntoHand(id)
	}
}

// PutUpTo moves up to Max cards the controller chooses from Target's pool into a
// destination zone, one at a time, stopping early when they choose Done — Grasping
// Vines returns up to 3 artifacts to their owners' hands. It is the bounded-choice
// counterpart to PutFromPlay, which moves every card the Target selects.
type PutUpTo struct {
	Max         int
	Target      Target
	Destination Destination
}

// validate requires a target, a positive maximum, and a supported destination.
func (e PutUpTo) validate() error {
	if !e.Target.valid() {
		return errUnsetTarget("PutUpTo")
	}
	if e.Max <= 0 {
		return fmt.Errorf("PutUpTo: Max must be positive")
	}
	switch e.Destination {
	case ToHand, ToTopOfDeck, ToDeckShuffled, ToArchives:
		return nil
	default:
		return fmt.Errorf("PutUpTo: unsupported destination %d", e.Destination)
	}
}

// Text renders the effect, e.g. "put up to 3 artifacts into their owners' hands".
func (e PutUpTo) Text() string {
	noun := singularNoun(e.Target.Text()) + "s"
	switch e.Destination {
	case ToTopOfDeck:
		return fmt.Sprintf("put up to %d %s on top of their owners' decks", e.Max, noun)
	case ToDeckShuffled:
		return fmt.Sprintf("shuffle up to %d %s into their owners' decks", e.Max, noun)
	case ToArchives:
		return fmt.Sprintf("put up to %d %s into their owners' archives", e.Max, noun)
	default:
		return fmt.Sprintf("put up to %d %s into their owners' hands", e.Max, noun)
	}
}

// Resolve moves up to Max cards one at a time. Each step offers the current pool
// plus a "Done" option to stop early; when the pool is empty it stops.
func (e PutUpTo) Resolve(ctx *EffectContext) {
	const done = "Done"
	for i := 0; i < e.Max; i++ {
		cands := e.Target.Select(ctx)
		if len(cands) == 0 {
			return
		}
		options := make([]string, 0, len(cands)+1)
		for _, id := range cands {
			options = append(options, ctx.Resolver.Name(id))
		}
		options = append(options, done)
		choice := ctx.ChooseOption("Choose a card to move", options)
		if choice >= len(cands) {
			return // "Done" (the last option), or an out-of-range choice
		}
		moveFromPlayTo(ctx, cands[choice], e.Destination)
	}
}

// ReturnNamedToHand puts a card with a specific name that the controller chooses
// into their hand, taken either from a friendly creature in play or from their
// discard pile — Faygin recovering an Urchin. The controller chooses among both
// zones at once; an in-play creature returns to hand (shedding its in-play state)
// and a discard card is recovered.
//
//rulebook:effect Return Named Card to Hand
type ReturnNamedToHand struct {
	Name string
}

// Text renders the effect, e.g. "put an Urchin from play or from your discard pile
// into your hand".
func (e ReturnNamedToHand) Text() string {
	return fmt.Sprintf(
		"put %s from play or from your discard pile into your hand",
		indefinite(e.Name),
	)
}

// Resolve gathers every friendly in-play creature and discard-pile card with the
// name, lets the controller choose one, and moves it to their hand from whichever
// zone it is in.
func (e ReturnNamedToHand) Resolve(ctx *EffectContext) {
	inPlay := nameMatches(ctx, ctx.Resolver.Battleline(ctx.Controller), e.Name)
	candidates := slices.Concat(
		inPlay,
		nameMatches(ctx, ctx.Resolver.Discard(ctx.Controller), e.Name),
	)
	id, ok := ctx.ChooseCreature("Choose "+indefinite(e.Name)+" to put into your hand", candidates)
	if !ok {
		return
	}
	if slices.Contains(inPlay, id) {
		ctx.Resolver.PutIntoHand(id)
	} else {
		ctx.Resolver.PutFromDiscardIntoHand(id)
	}
}

// nameMatches returns the ids whose card has the given name, in order.
func nameMatches(ctx *EffectContext, ids []LocalID, name string) []LocalID {
	var out []LocalID
	for _, id := range ids {
		if ctx.Resolver.Name(id) == name {
			out = append(out, id)
		}
	}
	return out
}
