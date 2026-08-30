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
		ctx.Resolver.Logf("%s reveals %s", ctx.Resolver.PlayerName(ctx.Controller), ctx.Resolver.Name(id))
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
// controller's first, then the opponent's — and records each discarded card's
// house on the context so a following effect can act on it. An empty deck
// contributes no card and no house. Bonkers Killing Machine pairs it with
// DestroyOfEachDiscardedHouse.
type DiscardTopOfEachDeck struct{}

// Text renders the effect.
func (DiscardTopOfEachDeck) Text() string {
	return "discard the top card of each player's deck"
}

// Resolve discards the controller's top deck card, then the opponent's, recording
// the discarded cards' houses on the context.
func (DiscardTopOfEachDeck) Resolve(ctx *EffectContext) {
	ctx.DiscardedHouses = nil
	for _, player := range []int{ctx.Controller, ctx.Opponent()} {
		if discarded, ok := ctx.Resolver.DiscardTopOfDeck(player); ok {
			ctx.DiscardedHouses = append(ctx.DiscardedHouses, ctx.Resolver.House(discarded))
		}
	}
}

// DestroyOfEachDiscardedHouse destroys, for each house a preceding
// DiscardTopOfEachDeck recorded on the context, one creature or artifact of that
// house that the controller chooses, tallying how many were destroyed on the
// context (read by CardsDestroyedFewerThan). It is the second half of Bonkers
// Killing Machine.
type DestroyOfEachDiscardedHouse struct{}

// Text renders the effect.
func (DestroyOfEachDiscardedHouse) Text() string {
	return "for each card discarded this way, destroy a creature or artifact of that card's house"
}

// Resolve destroys one creature or artifact of each recorded house, counting the
// destructions on the context.
func (DestroyOfEachDiscardedHouse) Resolve(ctx *EffectContext) {
	for _, house := range ctx.DiscardedHouses {
		chosen, ok := ctx.ChooseCreature("Choose a "+house.String()+" creature or artifact", inPlayOfHouse(ctx, house))
		if !ok {
			continue
		}
		if destroyAndReport(ctx, chosen) {
			ctx.Destroyed++
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

// inPlayOfHouse returns every creature and artifact in play of house, from the
// controller's point of view: friendly battleline, enemy battleline, friendly
// artifacts, then enemy artifacts.
func inPlayOfHouse(ctx *EffectContext, house House) []LocalID {
	ids := ctx.Resolver.Battleline(ctx.Controller)
	ids = append(ids, ctx.Resolver.Battleline(ctx.Opponent())...)
	ids = append(ids, ctx.Resolver.Artifacts(ctx.Controller)...)
	ids = append(ids, ctx.Resolver.Artifacts(ctx.Opponent())...)
	var out []LocalID
	for _, id := range ids {
		if ctx.Resolver.House(id) == house {
			out = append(out, id)
		}
	}
	return out
}

// destroyAndReport destroys id and reports whether it actually left play.
func destroyAndReport(ctx *EffectContext, id LocalID) bool {
	before := resolverInPlay(ctx, id)
	ctx.Resolver.DestroyEach(ctx.Controller, []LocalID{id})
	return before && !resolverInPlay(ctx, id)
}

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
