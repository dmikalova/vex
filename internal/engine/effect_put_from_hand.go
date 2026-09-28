package engine

// PutFromHand puts a card the controller chooses from their own hand directly
// into play — Swap Widget swapping in a replacement creature. Filter restricts the
// choice; its zero value allows any card. ExceptSameName excludes a card sharing
// the name of the card currently in context (ctx.It), the "with a different name"
// clause — meant to follow a gate that left a card in context, e.g.
// Then{PutFromPlay, PutFromHand}.
//
// It stays separate from PutIntoPlay because PutIntoPlay selects with a Target
// and a Target only reaches cards in play. Folding the two would mean teaching
// Target a hand domain — a hidden zone only its own controller may be prompted
// over — for one source, so this node walks the hand directly through the shared
// handCardsWhere gather instead, the way
// EachPlayerPutsHandCreaturesIntoPlay does.
type PutFromHand struct {
	// Filter narrows which cards in hand may be chosen and supplies the noun the
	// clause prints. A card in hand is in no battleline, so a filter that narrows on
	// an in-play axis is a definition error rather than a narrowing that matches
	// nothing; validate rejects it.
	Filter         Filter
	ExceptSameName bool
}

// validate rejects a filter that reads the board, which a card in hand is not on.
func (e PutFromHand) validate() error {
	return e.Filter.validateIdentityOnly("PutFromHand")
}

// noun renders the kind of card the effect puts into play, e.g. "Mars creature
// with a different name".
func (e PutFromHand) noun() string {
	base := e.Filter.noun("card")
	if e.ExceptSameName {
		base += " with a different name"
	}
	return base
}

// Text renders the effect, e.g. "put a Mars creature with a different name from
// your hand into play".
func (e PutFromHand) Text() string {
	return "put " + indefinite(e.noun()) + " from your hand into play"
}

// Resolve puts a card the controller chooses from their hand into play. Nothing
// happens if there is no candidate or the choice is declined. The chosen card is
// left in context (ctx.It) so a following effect can act on "it" (Swap Widget
// readying the creature it just put into play).
func (e PutFromHand) Resolve(ctx *EffectContext) {
	candidates := handCardsWhere(ctx, ctx.Controller, func(id LocalID) bool {
		if !e.Filter.matches(ctx, id) {
			return false
		}
		return !e.ExceptSameName || !ctx.HasIt ||
			ctx.Resolver.Name(id) != ctx.Resolver.Name(ctx.It)
	})
	id, ok := ctx.ChooseCard("Choose a "+e.noun()+" from your hand to put into play", candidates)
	if !ok {
		return
	}
	ctx.Resolver.PutIntoPlay(id, ctx.Controller)
	ctx.It, ctx.HasIt = id, true
}
