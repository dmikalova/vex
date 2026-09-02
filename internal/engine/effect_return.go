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
	return e.Destination.clause(e.Target.Text(), false)
}

// validate rejects a destination this effect cannot move a card to; only the hand,
// the top of the deck, and the archives are supported, and the destination must be
// named.
func (e PutFromPlay) validate() error {
	if !e.Target.valid() {
		return errUnsetTarget("PutFromPlay")
	}
	if !e.Destination.movable() {
		return fmt.Errorf("PutFromPlay: unsupported destination %d", e.Destination.zone)
	}
	return nil
}

// Resolve moves each selected card from play to the destination. Cards headed to
// the top of the deck are stacked in an order the controller chooses.
func (e PutFromPlay) Resolve(ctx *EffectContext) {
	ids := e.Target.Select(ctx)
	if e.Destination == ToTopOfDeck {
		ids = ctx.OrderByChoice("Choose the next card to put on top of the deck", ids)
	}
	for _, id := range ids {
		e.Destination.move(ctx, id)
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
	if !e.Destination.movable() {
		return fmt.Errorf("PutUpTo: unsupported destination %d", e.Destination.zone)
	}
	return nil
}

// Text renders the effect, e.g. "put up to 3 artifacts into their owners' hands".
func (e PutUpTo) Text() string {
	noun := singularNoun(e.Target.Text()) + "s"
	return e.Destination.clause(fmt.Sprintf("up to %d %s", e.Max, noun), true)
}

// Resolve moves up to Max cards one at a time. Each step offers the current pool
// as a declinable choice so the controller can stop early; when the pool is empty
// it stops.
func (e PutUpTo) Resolve(ctx *EffectContext) {
	for i := 0; i < e.Max; i++ {
		chosen, ok := ctx.ChooseCardOptional("Choose a card to move", e.Target.Select(ctx))
		if !ok {
			return
		}
		e.Destination.move(ctx, chosen)
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
