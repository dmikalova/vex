package engine

import "strings"

// Revealing cards shows them from a hand to both players and records them in the
// log, turning hidden information public. A card reveals cards so that what
// follows can be trusted — you reveal the Mars cards you are drawing for, or an
// opponent's whole hand before discarding from it — which is why the printed text
// is careful about which cards are shown.
//
// A House narrows the reveal to cards of that house (the wording "reveal any
// number of Mars cards"): a player would only ever reveal cards that help them,
// so every matching card is revealed. An unset House reveals the whole hand.
type Reveal struct {
	Player Player
	House  House
}

// validate rejects a Reveal whose player was left unset.
func (e Reveal) validate() error {
	if !e.Player.valid() {
		return errUnsetPlayer("Reveal")
	}
	return nil
}

// Text renders the effect, e.g. "reveal any number of Mars cards from your hand"
// or "reveal your opponent's hand".
func (e Reveal) Text() string {
	whose := "your"
	if e.Player == Opponent {
		whose = "your opponent's"
	}
	if e.House == HouseNone {
		return "reveal " + whose + " hand"
	}
	return "reveal any number of " + e.House.String() + " cards from " + whose + " hand"
}

// Resolve shows the matching cards, logs them, and records how many were revealed.
func (e Reveal) Resolve(ctx *EffectContext) {
	owner := ctx.PlayerFor(e.Player)
	var names []string
	for _, id := range ctx.Resolver.Hand(owner) {
		if e.House == HouseNone || ctx.Resolver.House(id) == e.House {
			names = append(names, ctx.Resolver.Name(id))
		}
	}
	ctx.Revealed = len(names)
	if len(names) > 0 {
		ctx.Resolver.Logf("%s reveals %s", ctx.Resolver.PlayerName(owner), strings.Join(names, ", "))
	}
}

// CardsRevealed counts the cards the most recent Reveal showed — the "for each
// card revealed this way" clause. Reveal records the tally on the context, so
// pairing it after a Reveal lets an effect scale with the reveal.
type CardsRevealed struct{}

// Value returns how many cards the preceding Reveal showed.
func (CardsRevealed) Value(ctx *EffectContext) int { return ctx.Revealed }

// CountText renders the singular noun the "for each" clause repeats.
func (CardsRevealed) CountText() string { return "card revealed this way" }

// DiscardTopOfDeckAndRevealHandForAember ties the house of a discarded deck card
// to a hand reveal. The player discards the top card of their deck and reveals
// their hand; then the gainer gains 1 Æmber for each revealed card sharing the
// discarded card's house. If the deck is empty, no card establishes a house, so
// the reveal and gain are skipped.
type DiscardTopOfDeckAndRevealHandForAember struct {
	Player Player
	Gainer Player
}

// validate rejects the effect when either relative player was left unset.
func (e DiscardTopOfDeckAndRevealHandForAember) validate() error {
	if !e.Player.valid() {
		return errUnsetPlayer("DiscardTopOfDeckAndRevealHandForAember")
	}
	if !e.Gainer.valid() {
		return errUnsetPlayer("DiscardTopOfDeckAndRevealHandForAember")
	}
	return nil
}

// Text renders the linked discard, reveal, and reward. Every half spells its own
// effect out from the same template rather than pointing back at a preceding one,
// so A Fair Game's second half prints as "Discard the top card of your deck and
// reveal your hand. Your opponent gains 1 Æmber for each card of the discarded
// card's house revealed this way."
func (e DiscardTopOfDeckAndRevealHandForAember) Text() string {
	deck, hand := "your deck", "your hand"
	if e.Player == Opponent {
		deck, hand = "your opponent's deck", "their hand"
	}
	gain := "You gain"
	if e.Gainer == Opponent {
		gain = "Your opponent gains"
	}
	return "discard the top card of " + deck + " and reveal " + hand + ". " +
		gain + " 1 Æmber for each card of the discarded card's house revealed this way."
}

// Resolve discards the top deck card, reveals the same player's hand, and pays the
// gainer for each revealed card of the discarded card's house.
func (e DiscardTopOfDeckAndRevealHandForAember) Resolve(ctx *EffectContext) {
	player := ctx.PlayerFor(e.Player)
	discarded, ok := ctx.Resolver.DiscardTopOfDeck(player)
	if !ok {
		ctx.Revealed = 0
		return
	}
	discardedHouse := ctx.Resolver.House(discarded)
	revealed, matching := e.revealAndCount(ctx, player, discardedHouse)
	ctx.Revealed = revealed
	if matching > 0 {
		GainAember{Player: e.Gainer, Amount: matching}.Resolve(ctx)
	}
}

// revealAndCount reveals the player's whole hand and counts cards of house among
// the revealed cards.
func (e DiscardTopOfDeckAndRevealHandForAember) revealAndCount(ctx *EffectContext, player int, house House) (revealed, matching int) {
	var names []string
	for _, id := range ctx.Resolver.Hand(player) {
		names = append(names, ctx.Resolver.Name(id))
		if ctx.Resolver.House(id) == house {
			matching++
		}
	}
	if len(names) > 0 {
		ctx.Resolver.Logf("%s reveals %s", ctx.Resolver.PlayerName(player), strings.Join(names, ", "))
	}
	return len(names), matching
}
