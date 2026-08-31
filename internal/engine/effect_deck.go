package engine

// RevealTopOfDeck reveals the top card of the controller's deck — logging it and
// putting it in context (ctx.It) so a following effect can inspect or play it (Chaos
// Portal plays it when it is of the chosen house). Revealing does not move the card;
// an empty deck reveals nothing.
//
//rulebook:effect Reveal Top of Deck
type RevealTopOfDeck struct{}

// Text renders the effect.
func (RevealTopOfDeck) Text() string { return "reveal the top card of your deck" }

// Resolve reveals the top card, putting it in context.
func (RevealTopOfDeck) Resolve(ctx *EffectContext) {
	id, ok := ctx.Resolver.TopOfDeck(ctx.Controller)
	ctx.It, ctx.HasIt = id, ok
	if ok {
		ctx.Resolver.Logf(
			"%s reveals %s",
			ctx.Resolver.PlayerName(ctx.Controller),
			ctx.Resolver.Name(id),
		)
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
