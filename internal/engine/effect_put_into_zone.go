package engine

import (
	"fmt"
)

// PutFromPlay takes each card its Target selects out of play and puts it in a
// destination zone — the top of its owner's deck, their hand, or their archives —
// shedding the per-match state the card built up in play (damage, spent armor,
// Æmber on it, upgrades). The destination is required. Moving a card out of play
// this way is how a "Destroyed:" ability can save its own creature: the creature
// leaves for the named zone as it is destroyed, so it never reaches the discard
// pile. When several cards move to the top of the deck at once the controller
// chooses the order they stack.
type PutFromPlay struct {
	Target      Target
	Destination Destination
	// WithUpgrades also returns each upgrade attached to a moved creature to its
	// owner's hand, rather than shedding it to the discard pile (Transporter
	// Platform). It is supported only for the hand.
	WithUpgrades bool
}

// Text renders the effect, e.g. "put each artifact on top of its owner's deck" or
// "put this creature into its owner's hand".
func (e PutFromPlay) Text() string {
	subject := e.Target.Text()
	if e.WithUpgrades {
		subject += " and each upgrade attached to it"
	}
	return e.Destination.clause(subject, false)
}

// validate rejects a destination this effect cannot move a card to; only the hand,
// the top of the deck, and the archives are supported, and the destination must be
// named. Returning attached upgrades along with the creature is supported only for
// the hand.
func (e PutFromPlay) validate() error {
	if !e.Target.valid() {
		return errUnsetTarget("PutFromPlay")
	}
	if !e.Destination.movable() {
		return fmt.Errorf("PutFromPlay: unsupported destination %d", e.Destination.zone)
	}
	if e.WithUpgrades && e.Destination != ToHand {
		return fmt.Errorf("PutFromPlay: WithUpgrades is supported only for the hand")
	}
	return nil
}

// Resolve moves each selected card from play to the destination. Cards headed to
// the top of the deck are stacked in an order the controller chooses. A card an
// earlier move already took out of play — a "Leaves Play:" ability can destroy
// one still on the list — is skipped rather than moved and counted twice.
func (e PutFromPlay) Resolve(ctx *EffectContext) { e.resolveGate(ctx) }

// resolveGate moves the target and reports whether any card actually moved, so
// PutFromPlay can gate a Then — Swap Widget only swaps in a new creature if it
// actually returned one. The last card moved is left in context (ctx.It) so a
// following effect can act on "it" or exclude cards sharing its name.
func (e PutFromPlay) resolveGate(ctx *EffectContext) bool {
	return e.put(ctx, e.Target.Select(ctx))
}

// declinable reports that the move is a single clickable card.
func (e PutFromPlay) declinable() bool { return e.Target.isChosen() }

// vacuous reports that there is no card to move, so a "you may" wrapping it need
// not ask.
func (e PutFromPlay) vacuous(ctx *EffectContext) bool { return e.Target.empty(ctx) }

// resolveOptional asks for the card declinably, so "you may put a creature into
// its owner's hand" is answered by clicking that creature rather than by a
// separate Yes/No.
func (e PutFromPlay) resolveOptional(ctx *EffectContext) bool {
	return e.put(ctx, e.Target.SelectOptional(ctx))
}

// put moves an already-selected set and reports whether any card actually moved.
// The set moves as one moment: the board does not settle between two cards of one
// selection, and no card's "Leaves Play:" ability resolves until every card has
// moved, so the outcome cannot depend on the order the selection is visited
// (TestPutFromPlayMovesTheWholeSelectionTogether).
//
// The per-card in-play recheck survives that batching because one move can still
// carry a second card out with it: a gigantic is two cards and moving either half
// moves both, so the loop reaches the other half after it has already left
// (TestPutFromPlaySkipsTheSecondGiganticHalf).
func (e PutFromPlay) put(ctx *EffectContext, ids []LocalID) bool {
	if e.Destination == ToTopOfDeck {
		ids = ctx.OrderByChoice("Choose the next card to put on top of the deck", ids)
	}
	moved := false
	ctx.Resolver.Simultaneously(ctx.Controller, func() {
		for _, id := range ids {
			if !resolverInPlay(ctx, id) {
				continue
			}
			if e.WithUpgrades {
				ctx.Resolver.ReturnUpgradesToHand(id)
			}
			controller := ctx.Resolver.Controller(id)
			e.Destination.move(ctx, id)
			ctx.Produced.Moved[controller]++
			ctx.It, ctx.HasIt = id, true
			moved = true
		}
	})
	return moved
}

// PutChosen moves Amount cards the controller chooses from Target's pool into a
// destination zone, one at a time — Lost in the Woods shuffles 2 friendly and 2
// enemy creatures into their owners' decks. UpTo makes the choice declinable, the
// "up to 3 artifacts" of Grasping Vines; without it the controller must choose as
// many as the pool allows. It is the bounded-choice counterpart to PutFromPlay,
// which moves every card the Target selects.
type PutChosen struct {
	Amount      int
	UpTo        bool
	Target      Target
	Destination Destination
}

// validate requires a target, a positive count, and a supported destination.
func (e PutChosen) validate() error {
	if !e.Target.valid() {
		return errUnsetTarget("PutChosen")
	}
	if e.Amount <= 0 {
		return fmt.Errorf("PutChosen: Count must be positive")
	}
	if !e.Destination.movable() {
		return fmt.Errorf("PutChosen: unsupported destination %d", e.Destination.zone)
	}
	return nil
}

// Text renders the effect, e.g. "put up to 3 artifacts into their owners' hands"
// or "shuffle 2 friendly creatures into their owners' decks".
func (e PutChosen) Text() string {
	noun := singularNoun(e.Target.Text())
	if e.Amount == 1 {
		return e.Destination.clause(indefinite(noun), false)
	}
	quantity := fmt.Sprintf("%d %ss", e.Amount, noun)
	if e.UpTo {
		quantity = "up to " + quantity
	}
	return e.Destination.clause(quantity, true)
}

// Resolve moves Amount cards one at a time. An UpTo choice is declinable so the
// controller can stop early; either way it stops when the pool runs out. Shuffles
// into a deck are batched so several creatures moved at once narrate as one
// grouped line per owner attributed to this ability's source.
func (e PutChosen) Resolve(ctx *EffectContext) {
	if e.Destination == ToDeckShuffled {
		ctx.Resolver.BeginShuffleBatch()
		e.resolveMoves(ctx)
		ctx.Resolver.EndShuffleBatch()
		return
	}
	e.resolveMoves(ctx)
}

// resolveMoves runs the choose-and-move loop shared by the batched and unbatched
// paths.
func (e PutChosen) resolveMoves(ctx *EffectContext) {
	choose := ctx.ChooseCard
	if e.UpTo {
		choose = ctx.ChooseCardOptional
	}
	for i := 0; i < e.Amount; i++ {
		chosen, ok := choose("Choose a card to move", e.Target.Select(ctx))
		if !ok {
			return
		}
		// Settling before the choice (ADR 0029) can destroy a creature that was in
		// the pool when it was gathered — the buff it relied on left with an earlier
		// pick — so a chosen card no longer in play is skipped rather than moved into
		// a second zone (ADR 0030), exactly as PutFromPlay skips one an earlier move
		// took out.
		if !resolverInPlay(ctx, chosen) {
			continue
		}
		e.Destination.move(ctx, chosen)
	}
}

// PutNamedIntoHand puts a card with a specific name that the controller chooses
// into their hand, taken either from a friendly card in play or from their
// discard pile — Faygin recovering an Urchin. The controller chooses among both
// zones at once; an in-play card returns to hand (shedding its in-play state)
// and a discard card is recovered.
type PutNamedIntoHand struct {
	Name string
}

// Text renders the effect, e.g. "put an Urchin from play or from your discard pile
// into your hand".
func (e PutNamedIntoHand) Text() string {
	return fmt.Sprintf(
		"put %s from play or from your discard pile into your hand",
		indefinite(e.Name),
	)
}

// Resolve gathers every friendly in-play creature and discard-pile card, then
// lets the controller pick the named one through the shared Selection vocabulary
// (Chosen{Name}) and moves it to their hand from whichever zone it is in.
func (e PutNamedIntoHand) Resolve(ctx *EffectContext) {
	mover := crossZoneMover{
		Player:  ctx.Controller,
		Dest:    ToHand,
		Sources: []Zone{InPlay, Discard},
	}
	pool := mover.gather(ctx, func(LocalID) bool { return true })
	for _, id := range (Chosen{Name: e.Name}).pick(ctx, pool) {
		mover.move(ctx, id)
	}
}

// PutItIntoHand puts the creature in context ("it") into its owner's hand,
// recovering it from the discard pile when it has already been destroyed. It is
// the reaction counterpart to PutFromPlay, which only reaches a creature still in
// play: Nizak, The Forgotten returns an enemy destroyed fighting it, and that
// creature already sits in the discard by the time the reaction resolves.
type PutItIntoHand struct{}

// Text renders the effect, e.g. "put it into its owner's hand".
func (e PutItIntoHand) Text() string {
	return ToHand.clause(Target{Kind: TargetTriggeringCreature}.Text(), false)
}

// Resolve moves the contextual creature to its owner's hand — from play if it is
// still there, or recovered from the discard pile if it was already destroyed. It
// does nothing when no card is in context.
func (e PutItIntoHand) Resolve(ctx *EffectContext) {
	if !ctx.HasIt {
		return
	}
	if resolverInPlay(ctx, ctx.It) {
		ctx.Resolver.PutIntoHand(ctx.It)
		return
	}
	ctx.Resolver.PutFromDiscardIntoHand(ctx.It)
}
