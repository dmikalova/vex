package engine

import "fmt"

// RevealTopOfDeck reveals the top card of the controller's deck — logging it and
// putting it in context (ctx.It) so a following effect can inspect or play it (Chaos
// Portal plays it when it is of the chosen house). Revealing does not move the card;
// an empty deck reveals nothing.
type RevealTopOfDeck struct{}

// Text renders the effect.
func (RevealTopOfDeck) Text() string { return "reveal the top card of your deck" }

// Resolve reveals the top card, putting it in context.
func (RevealTopOfDeck) Resolve(ctx *EffectContext) {
	id, ok := ctx.Resolver.TopOfDeck(ctx.Controller)
	ctx.It, ctx.HasIt = id, ok
	if ok {
		ctx.Resolver.Record(CardsRevealedToAll{
			Player: ctx.Controller,
			Cards:  []LocalID{id},
		})
	}
}

// PlayRevealedCard plays the card in context (put there by a preceding
// RevealTopOfDeck) from the controller's deck — Chaos Portal. It does nothing when
// no card is in context.
type PlayRevealedCard struct{}

// Text renders the effect.
func (PlayRevealedCard) Text() string { return "play it" }

// Resolve plays the context card from the deck.
func (PlayRevealedCard) Resolve(ctx *EffectContext) {
	if ctx.HasIt {
		ctx.Resolver.PlayFromDeck(ctx.Controller, ctx.It)
	}
}

// PlayTopOfDeck plays the top card of the controller's deck outright (Wild
// Wormhole), resolving that card's own play effect. It does nothing when the deck
// is empty.
type PlayTopOfDeck struct{}

// Text renders the effect.
func (PlayTopOfDeck) Text() string { return "play the top card of your deck" }

// Resolve plays the top card of the controller's deck, if any.
func (PlayTopOfDeck) Resolve(ctx *EffectContext) {
	if id, ok := ctx.Resolver.TopOfDeck(ctx.Controller); ok {
		ctx.Resolver.Record(PlayedFromTopOfDeck{
			Source: ctx.Source,
			Card:   id,
			Player: ctx.Controller,
		})
		ctx.Resolver.PlayFromDeck(ctx.Controller, id)
	}
}

// DiscardTopOfDeck discards the top card of a deck and puts it in context (ctx.It)
// so a following effect can react to it — Evasion Sigil cancels the fight when it
// is of the active house; A Fair Game gains Æmber for hand cards of its house.
//
// Player picks whose deck, and its rendering follows the two ways a card names a
// player. Left unset it speaks from a granted ability's card-neutral perspective
// ("its controller's deck"), resolving against the creature's controller — the
// only sensible default for an ability every creature gains (Evasion Sigil).
// Controller and Opponent are the direct first/second-person perspectives ("your
// deck", "your opponent's deck") a card played from hand uses (A Fair Game). An
// empty deck discards nothing.
type DiscardTopOfDeck struct {
	Player Player
}

// Text renders the effect from the chosen player's perspective.
func (e DiscardTopOfDeck) Text() string {
	switch e.Player {
	case Controller:
		return "discard the top card of your deck"
	case Opponent:
		return "discard the top card of your opponent's deck"
	default:
		return "discard the top card of its controller's deck"
	}
}

// Resolve discards the top card of the chosen deck (the controller's when Player
// is unset), putting it in context.
func (e DiscardTopOfDeck) Resolve(ctx *EffectContext) {
	player := ctx.Controller
	if e.Player.valid() {
		player = ctx.PlayerFor(e.Player)
	}
	id, ok := ctx.Resolver.DiscardTopOfDeck(player)
	ctx.It, ctx.HasIt = id, ok
}

// DiscardTopOfEachDeck discards the top card of each player's deck — the
// controller's first, then the opponent's — and records each discarded card on
// the context so a following ForEachDiscarded can act on it. An empty deck
// contributes no card. Bonkers Killing Machine pairs it with ForEachDiscarded.
type DiscardTopOfEachDeck struct{}

// Text renders the effect.
func (DiscardTopOfEachDeck) Text() string {
	return "discard the top card of each player's deck"
}

// Resolve discards the controller's top deck card, then the opponent's, recording
// the discarded cards on the context.
func (DiscardTopOfEachDeck) Resolve(ctx *EffectContext) {
	ctx.Produced.Discarded = nil
	for _, player := range []int{ctx.Controller, ctx.Opponent()} {
		if discarded, ok := ctx.Resolver.DiscardTopOfDeck(player); ok {
			ctx.Produced.Discarded = append(ctx.Produced.Discarded, discarded)
		}
	}
}

// ForEachDiscarded resolves Do once for each card a preceding DiscardTopOfEachDeck
// discarded, putting that card in context (ctx.It) so Do can refer to it — Bonkers
// Killing Machine destroys a creature or artifact of each discarded card's house
// (Do targets Target.OfContextualHouse).
type ForEachDiscarded struct {
	Do Effect
}

// validate surfaces a configuration error from Do.
func (e ForEachDiscarded) validate() error { return validateEffect(e.Do) }

// Text renders the effect, leading with the iteration clause.
func (e ForEachDiscarded) Text() string {
	return "for each card discarded this way, " + e.Do.Text()
}

// Resolve runs Do for each discarded card, in context as ctx.It.
func (e ForEachDiscarded) Resolve(ctx *EffectContext) {
	for _, id := range ctx.Produced.Discarded {
		ctx.It, ctx.HasIt = id, true
		e.Do.Resolve(ctx)
	}
}

// LookAtTop looks at the top Count cards of the controller's deck, puts one the
// controller chooses into their hand, and discards the others — Eyegor. It looks
// at as many as remain when the deck holds fewer than Count, and does nothing on
// an empty deck.
type LookAtTop struct {
	Count int
}

// validate rejects a non-positive Count: looking at zero cards is meaningless, so
// an omitted count is an authoring error, not a silent default.
func (e LookAtTop) validate() error {
	if e.Count < 1 {
		return fmt.Errorf("LookAtTop: Count must be at least 1")
	}
	return nil
}

// Text renders the effect.
func (e LookAtTop) Text() string {
	return fmt.Sprintf(
		"look at the top %d cards of your deck, put 1 into your hand, and discard the others",
		e.Count,
	)
}

// Resolve looks at the top Count cards, moves the one the controller chooses to
// their hand, and discards the rest.
func (e LookAtTop) Resolve(ctx *EffectContext) {
	deck := ctx.Resolver.Deck(ctx.Controller)
	if len(deck) == 0 {
		return
	}
	top := deck[:min(e.Count, len(deck))]
	keep, ok := ctx.ChooseCard("Choose a card to put into your hand", top)
	if !ok {
		return
	}
	ctx.Resolver.MoveFromDeckToHand(keep)
	for _, id := range top {
		if id != keep {
			ctx.Resolver.MoveFromDeckToDiscard(id)
		}
	}
}

// CancelFight makes the fight in progress not occur — a "Before Fight" effect
// (Evasion Sigil, gated on the discarded card's house). The attacker was still used
// to fight, so it stays exhausted; combat reads the cancellation and skips Assault,
// Hazardous, fight damage, and Fight: abilities.
type CancelFight struct{}

// Text renders the effect.
func (CancelFight) Text() string { return "the fight does not occur" }

// Resolve cancels the current fight.
func (CancelFight) Resolve(ctx *EffectContext) { ctx.Resolver.CancelCurrentFight() }

// resolverInPlay reports whether id appears in either player's battleline or
// artifact row using only Resolver reads.
func resolverInPlay(ctx *EffectContext, id LocalID) bool {
	for p := 0; p < 2; p++ {
		for _, candidate := range ctx.Resolver.Battleline(p) {
			if candidate == id {
				return true
			}
		}
		for _, candidate := range ctx.Resolver.Artifacts(p) {
			if candidate == id {
				return true
			}
		}
	}
	return false
}

// DiscardDeckUntil digs through the top of your deck, discarding as it goes,
// until it turns up a card the filters admit or the deck runs out. The card it
// finds stays in the discard pile and goes into context (ctx.It), so what happens
// to it is a separate effect gated on the dig succeeding — Sound the Horns and
// Invasion Portal both pair it with PutDiscardedIntoHand.
type DiscardDeckUntil struct {
	// Type filters what ends the dig; the zero value stops at any card.
	Type CardType
	// House filters what ends the dig; HouseNone stops at any house.
	House House
}

// Text renders the dig and names both ways it can end, as the cards do.
func (e DiscardDeckUntil) Text() string {
	return "discard cards from the top of your deck until you discard " +
		indefinite(e.noun()) + " or run out of cards"
}

// noun names the cards the filters admit, e.g. "card" or "Brobnar creature".
func (e DiscardDeckUntil) noun() string {
	noun := "card"
	switch e.Type {
	case Creature:
		noun = "creature"
	case Artifact:
		noun = "artifact"
	}
	if e.House != HouseNone {
		noun = e.House.String() + " " + noun
	}
	return noun
}

// matches reports whether a discarded card is the one the dig was looking for.
func (e DiscardDeckUntil) matches(ctx *EffectContext, id LocalID) bool {
	if e.Type != TypeUnset && ctx.Resolver.TypeOf(id) != e.Type {
		return false
	}
	return e.House == HouseNone || ctx.Resolver.House(id) == e.House
}

// Resolve digs, leaving the found card in context.
func (e DiscardDeckUntil) Resolve(ctx *EffectContext) { e.resolveGate(ctx) }

// resolveGate digs and reports whether it found a matching card, so a Then can
// hang "put it into your hand" off the dig succeeding.
func (e DiscardDeckUntil) resolveGate(ctx *EffectContext) bool {
	ctx.It, ctx.HasIt = 0, false
	for {
		id, ok := ctx.Resolver.DiscardTopOfDeck(ctx.Controller)
		if !ok {
			return false
		}
		if e.matches(ctx, id) {
			ctx.It, ctx.HasIt = id, true
			return true
		}
	}
}

// PutDiscardedIntoHand takes the card in context out of the discard pile and
// into its owner's hand. It is the tail of a dig through the deck (DiscardDeckUntil)
// that just discarded the card. Type names what the dig stopped on so the tail
// reads "put the discarded creature into your hand" rather than a bare "it"; the
// zero value stays the generic "card".
type PutDiscardedIntoHand struct {
	// Type names the discarded card the dig stopped on; the zero value is "card".
	Type CardType
}

// Text renders the effect, naming the discarded card the dig stopped on.
func (e PutDiscardedIntoHand) Text() string {
	return "put the discarded " + discardedNoun(e.Type) + " into your hand"
}

// discardedNoun names a discarded card by type for the "put the discarded …"
// tail: a creature, an artifact, or a bare card when the type is unset.
func discardedNoun(t CardType) string {
	switch t {
	case Creature:
		return "creature"
	case Artifact:
		return "artifact"
	default:
		return "card"
	}
}

// Resolve moves the contextual card from the discard pile to hand.
func (e PutDiscardedIntoHand) Resolve(ctx *EffectContext) {
	if ctx.HasIt {
		ctx.Resolver.PutFromDiscardIntoHand(ctx.It)
	}
}
