package engine

import "strings"

// This file holds the effects that place a card under another card, and that
// play it back out from there — Masterplan puts a card from hand facedown under
// itself then later plays it; Jargogle and Graft do the same shape of thing (see
// ADR 0016 for the mechanic these effects sit on top of).

// PutUnderFromHand has the controller choose a card from their hand and place it
// under the resolving card, face up or face down. Masterplan and Jargogle place
// theirs facedown; Graft always places its card faceup. Type restricts the choice
// to cards of that type; the zero value allows any card, and Tactic is the "action
// card" a graft-from-hand takes (Infomancer, Memolith). It does nothing with an
// empty hand.
type PutUnderFromHand struct {
	// FaceDown places the chosen card hidden from the opponent, viewable only by
	// the controller of the resolving card (Peekable).
	FaceDown bool
	// Type restricts the choice to cards of that type; the zero value (an unset
	// CardType) allows any card.
	Type CardType
}

// noun renders the kind of card the effect places, e.g. "Tactic" for a
// Tactic filter, "card" for none.
func (e PutUnderFromHand) noun() string {
	if e.Type == Tactic {
		return "Tactic"
	}
	if e.Type != TypeUnset {
		return strings.ToLower(e.Type.String()) + " card"
	}
	return "card"
}

// Text renders the effect, e.g. "put an action card from your hand faceup under
// {self}".
func (e PutUnderFromHand) Text() string {
	face := "faceup"
	if e.FaceDown {
		face = "facedown"
	}
	return "put " + indefinite(e.noun()) + " from your hand " + face + " under " + SelfName
}

// Resolve has the controller choose a card of the allowed type from their hand
// and place it under the resolving card.
func (e PutUnderFromHand) Resolve(ctx *EffectContext) {
	candidates := handCardsWhere(ctx, ctx.Controller, func(id LocalID) bool {
		return e.Type == TypeUnset || ctx.Resolver.TypeOf(id) == e.Type
	})
	if len(candidates) == 0 {
		return
	}
	id, ok := ctx.ChooseCard("Choose "+indefinite(e.noun())+" to put under "+SelfName, candidates)
	if !ok {
		return
	}
	ctx.Resolver.PutCardUnder(ctx.Controller, id, ctx.Source, e.FaceDown)
}

// TriggerGraftedPlayEffect triggers the Play effect of an action card grafted
// faceup under the resolving card — Infomancer's Reap and Memolith's Action.
// Unlike PlayCardUnder the grafted card does not move: only its Play abilities
// resolve, in place, and the card stays grafted (a tactic resolves its own Play
// abilities while out of play, since it never enters play at all). With more than
// one grafted action the controller chooses which; it does nothing with none.
type TriggerGraftedPlayEffect struct{}

// Text renders the effect, e.g. "trigger the play effect of a Tactic
// grafted onto {self}".
func (TriggerGraftedPlayEffect) Text() string {
	return "trigger the play effect of a Tactic grafted onto " + SelfName
}

// Resolve triggers the Play abilities of a grafted action the controller chooses,
// leaving it grafted under the source. Only faceup grafted cards that carry a Play
// ability are offered, so the choice is never a wasted one.
func (TriggerGraftedPlayEffect) Resolve(ctx *EffectContext) {
	var candidates []LocalID
	for _, id := range ctx.Resolver.Under(ctx.Source) {
		if !ctx.Resolver.UnderFaceDown(id) && ctx.Resolver.HasTrigger(id, TriggerAfterPlay) {
			candidates = append(candidates, id)
		}
	}
	if len(candidates) == 0 {
		return
	}
	id := candidates[0]
	if len(candidates) > 1 {
		var ok bool
		id, ok = ctx.ChooseCard("Choose a grafted Tactic to trigger", candidates)
		if !ok {
			return
		}
	}
	ctx.It, ctx.HasIt = id, true
	ctx.Resolver.TriggerAbilityOf(ctx.Controller, id, TriggerAfterPlay)
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

// ArchiveCardUnder archives the card placed under the resolving card, moving it
// to its owner's archives — Jargogle's Destroyed ability when it is not its
// controller's turn. It does nothing with nothing underneath.
type ArchiveCardUnder struct{}

// Text renders the effect, e.g. "archive the card under {self}".
func (ArchiveCardUnder) Text() string {
	return "archive the card under " + SelfName
}

// Resolve archives each card under the resolving card.
func (ArchiveCardUnder) Resolve(ctx *EffectContext) {
	ctx.Resolver.ArchiveCardUnder(ctx.Source)
}
