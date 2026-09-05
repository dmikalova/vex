package engine

// This file holds the effects that place a card under another card, and that
// play it back out from there — Masterplan puts a card from hand facedown under
// itself then later plays it; Jargogle and Graft do the same shape of thing (see
// ADR 0016 for the mechanic these effects sit on top of).

// PutUnderFromHand has the controller choose a card from their hand and place it
// under the resolving card, face up or face down. Masterplan and Jargogle place
// theirs facedown; Graft always places its card faceup. It does nothing with an
// empty hand.
type PutUnderFromHand struct {
	// FaceDown places the chosen card hidden from the opponent, viewable only by
	// the controller of the resolving card (Peekable).
	FaceDown bool
}

// Text renders the effect, e.g. "put a card from your hand facedown under
// {self}".
func (e PutUnderFromHand) Text() string {
	face := "faceup"
	if e.FaceDown {
		face = "facedown"
	}
	return "put a card from your hand " + face + " under " + SelfName
}

// Resolve has the controller choose a card from their hand and place it under
// the resolving card.
func (e PutUnderFromHand) Resolve(ctx *EffectContext) {
	candidates := ctx.Resolver.Hand(ctx.Controller)
	if len(candidates) == 0 {
		return
	}
	id, ok := ctx.ChooseCreature("Choose a card to put under "+SelfName, candidates)
	if !ok {
		return
	}
	ctx.Resolver.PutCardUnder(ctx.Controller, id, ctx.Source, e.FaceDown)
}

// PlayCardUnder plays the card placed under the resolving card, putting the one
// played in context (It) — Masterplan's and Jargogle's own "play the card under
// me." With more than one card underneath, the controller chooses which; it does
// nothing with none.
type PlayCardUnder struct{}

// Text renders the effect, e.g. "play the card under {self}".
func (PlayCardUnder) Text() string {
	return "play the card under " + SelfName
}

// Resolve plays the card placed under the resolving card.
func (PlayCardUnder) Resolve(ctx *EffectContext) {
	candidates := ctx.Resolver.Under(ctx.Source)
	if len(candidates) == 0 {
		return
	}
	id := candidates[0]
	if len(candidates) > 1 {
		var ok bool
		id, ok = ctx.ChooseCreature("Choose the card to play", candidates)
		if !ok {
			return
		}
	}
	ctx.Resolver.PlayFromUnder(ctx.Controller, id)
	ctx.It, ctx.HasIt = id, true
}

// Graft moves a target card in play faceup under the resolving card, out of play
// (rulebook: Graft). The grafted card leaves play — firing its Leaves Play
// abilities, not Destroyed — and waits under its new host until that host leaves
// play. Spangler Box grafts a chosen creature onto itself.
type Graft struct {
	// Target chooses the card to graft under the resolving card.
	Target Target
}

// validate requires an explicit target.
func (e Graft) validate() error {
	if !e.Target.valid() {
		return errUnsetTarget("Graft")
	}
	return nil
}

// Text renders the effect, e.g. "graft a creature from play".
func (e Graft) Text() string { return "graft " + e.Target.Text() + " from play" }

// Resolve grafts each selected card faceup under the resolving card.
func (e Graft) Resolve(ctx *EffectContext) {
	for _, id := range e.Target.Select(ctx) {
		ctx.Resolver.GraftUnder(id, ctx.Source)
	}
}

// PutUnderIntoPlay puts every card placed under the resolving card into play
// under its owner's control — Spangler Box's Destroyed ability returns the
// creatures grafted onto it. It does nothing with nothing underneath.
type PutUnderIntoPlay struct{}

// Text renders the effect, e.g. "put each card under {self} into play under its
// owner's control".
func (PutUnderIntoPlay) Text() string {
	return "put each card under " + SelfName + " into play under its owner's control"
}

// Resolve puts each card under the resolving card into play under its owner.
func (PutUnderIntoPlay) Resolve(ctx *EffectContext) {
	ctx.Resolver.PutUnderIntoPlay(ctx.Source)
}
