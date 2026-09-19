package engine

import (
	"fmt"
	"slices"
)

// PlayOrUse has the controller immediately either play a matching card from their
// hand or use a matching card they have in play — CXO Taber plays or uses a
// non-Star Alliance card. It is the immediate form of a MayPlayOrUse grant (ADR
// 0037) that frees one out-of-house card: rather than holding the permission open
// for the turn, it makes the one play-or-use happen now (see card-wording-rules.md
// rule 21). Grant is the verb set it exercises — the same HouseGrant axis
// MayPlayOrUse carries — so a card that only plays, or only uses, is this node with
// a narrower Grant, not a new node; the zero value is the full play-or-use. House
// narrows which cards qualify ("a non-Star Alliance card").
type PlayOrUse struct {
	House HouseMatcher
	Grant HouseGrant
}

// grant is the verb set the node exercises, reading the zero value as the full
// play-or-use so an unset Grant renders and resolves as CXO Taber's default.
func (e PlayOrUse) grant() HouseGrant {
	if e.Grant == 0 {
		return GrantPlay | GrantUse
	}
	return e.Grant
}

// validate rejects a malformed house filter or a verb this immediate node cannot
// exercise. It plays and uses only; GrantFight is the lasting per-creature
// narrowing MayPlayOrUse carries, meaningless for a single forced action.
func (e PlayOrUse) validate() error {
	if err := e.House.validate(); err != nil {
		return fmt.Errorf("PlayOrUse: %w", err)
	}
	if e.grant()&^(GrantPlay|GrantUse) != 0 {
		return fmt.Errorf("PlayOrUse: Grant must be play, use, or both")
	}
	return nil
}

// Text renders the effect, e.g. "play or use a non-Star Alliance card", narrowing
// the verb to what Grant frees.
func (e PlayOrUse) Text() string {
	return e.verb() + " " + indefinite(e.noun())
}

// verb renders the granted verb phrase — "play", "use", or "play or use".
func (e PlayOrUse) verb() string {
	switch g := e.grant(); {
	case g&GrantPlay == 0:
		return "use"
	case g&GrantUse == 0:
		return "play"
	default:
		return "play or use"
	}
}

// noun names the cards the filter admits, e.g. "card", "non-Star Alliance card",
// or "Mars card".
func (e PlayOrUse) noun() string {
	return e.House.qualify("card")
}

// Resolve offers, in one prompt, every matching card the controller may play from
// hand or use in play — limited to the verbs Grant frees — and plays or uses the
// chosen one. A card in hand is played; a card in play is used (a creature reaps,
// fights, or fires its Action; an artifact fires its Action). With nothing to offer
// it does nothing.
func (e PlayOrUse) Resolve(ctx *EffectContext) {
	var playable, usable []LocalID
	if e.grant()&GrantPlay != 0 {
		playable = e.playable(ctx)
	}
	if e.grant()&GrantUse != 0 {
		usable = e.usable(ctx)
	}
	cands := append(append([]LocalID{}, playable...), usable...)
	if len(cands) == 0 {
		return
	}
	id, ok := ctx.ChooseCard("Choose a card to "+e.verb(), cands)
	if !ok {
		return
	}
	if slices.Contains(playable, id) {
		ctx.Resolver.PlayFromHand(ctx.Controller, id)
		return
	}
	ctx.It, ctx.HasIt = id, true
	useCard(ctx, id)
}

// playable are the controller's hand cards the house filter admits.
func (e PlayOrUse) playable(ctx *EffectContext) []LocalID {
	var out []LocalID
	for _, id := range ctx.Resolver.Hand(ctx.Controller) {
		if e.admits(ctx, id) {
			out = append(out, id)
		}
	}
	return out
}

// usable are the controller's in-play cards the house filter admits and that can
// be used now — a ready creature, or a ready artifact holding an Action.
func (e PlayOrUse) usable(ctx *EffectContext) []LocalID {
	var inPlay []LocalID
	for _, id := range ctx.Resolver.Battleline(ctx.Controller) {
		if e.admits(ctx, id) {
			inPlay = append(inPlay, id)
		}
	}
	for _, id := range ctx.Resolver.Artifacts(ctx.Controller) {
		if e.admits(ctx, id) {
			inPlay = append(inPlay, id)
		}
	}
	return usableCards(ctx, inPlay, false)
}

// admits reports whether the house filter accepts the card.
func (e PlayOrUse) admits(ctx *EffectContext, id LocalID) bool {
	return e.House.matches(ctx, id)
}
