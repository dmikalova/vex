package engine

import (
	"fmt"
)

// This file holds the effects that play a card out of one of a player's piles as
// part of resolving another card — the "play a card" clause, as opposed to a
// player taking their own play action.

// PlayFrom has the controller play a card out of their own hand or discard pile
// right now, ignoring the active-house gate — Phase Shift's off-house card,
// Sacrificial Altar's creature back from the discard pile. From names the source
// pile; House and Types narrow which cards may be chosen.
//
// KeyForge prints the from-hand form as a permission held open for the rest of the
// turn ("you may play one non-Logos card this turn"); it is rendered and resolved
// as an immediate play instead, which needs no turn-scoped memory of an unspent
// allowance (see card-wording-rules.md rule 21). With no legal card in the source
// pile it does nothing.
type PlayFrom struct {
	// From is the zone the card is played out of: Hand or Discard. It has no
	// default — an effect names the pile it reaches into.
	From Zone
	// Player names whose pile the card comes from: Controller (the default, an
	// unset zero value) plays out of your own pile; Opponent reaches into the
	// opponent's hand (Lateral Shift) or discard pile (Mimicry). A card played from
	// the opponent's pile is still your play — it counts against your own card-play
	// limit — but is owned by, and returns to, the opponent.
	Player Player
	// House narrows which houses may be chosen; the zero value admits every house
	// (Com. Officer Kirby plays a non-Star Alliance card).
	House HouseMatcher
	// Types narrows the admitted card types — Com. Officer Kirby plays a non-Star
	// Alliance Artifact, Upgrade, or Tactic. The zero value admits every type.
	Types CardTypes
}

// validate rejects a malformed house filter and a source zone a card may not be
// played from.
func (e PlayFrom) validate() error {
	if err := e.House.validate(); err != nil {
		return fmt.Errorf("PlayFrom: %w", err)
	}
	if e.From != Hand && e.From != Discard && e.From != Archives {
		return fmt.Errorf("PlayFrom: From must be Hand, Discard, or Archives, got %v", e.From)
	}
	if e.Player == Opponent && e.From != Discard && e.From != Hand {
		return fmt.Errorf("PlayFrom: only the opponent's hand or discard pile may be played from")
	}
	return nil
}

// Text renders the effect, e.g. "play a non-Logos card" or "play a creature from
// your discard pile". Playing from your own hand is the default a printed card
// leaves unsaid ("Play a non-Logos card"), so only another zone — or the
// opponent's hand — is named.
func (e PlayFrom) Text() string {
	text := "play " + indefinite(e.noun())
	if e.From != Hand || e.Player == Opponent {
		text += " from " + e.pileOwner() + " " + e.From.noun()
	}
	return text
}

// pileOwner names whose pile the card is played from: your own by default, or
// the opponent's when Player is Opponent (Mimicry).
func (e PlayFrom) pileOwner() string {
	if e.Player == Opponent {
		return "your opponent's"
	}
	return "your"
}

// noun names the cards the filters admit, e.g. "card", "non-Logos card", or
// "Mars creature".
func (e PlayFrom) noun() string {
	noun := "card"
	if !e.Types.all() {
		noun = e.Types.list()
	}
	return e.House.qualify(noun)
}

// Resolve has the controller choose a matching card in the source zone and plays
// it.
func (e PlayFrom) Resolve(ctx *EffectContext) { e.resolveGate(ctx) }

// resolveGate has the controller choose a matching card in the source zone and
// plays it, binding the played card in context (ctx.It) so a following effect can
// act on "that creature" — Imperial Road plays a Saurian creature, then stuns it.
// It reports whether a card was played, so a Then hangs its follow-up on the play
// happening: an empty pile or a declined choice reports false.
func (e PlayFrom) resolveGate(ctx *EffectContext) bool {
	candidates := e.candidates(ctx)
	if len(candidates) == 0 {
		return false
	}
	id, ok := ctx.ChooseCreature("Choose "+indefinite(e.noun())+" to play", candidates)
	if !ok {
		return false
	}
	switch {
	case e.From == Hand && e.Player == Opponent:
		ctx.Resolver.PlayFromOpponentHand(ctx.Controller, id)
	case e.From == Discard && e.Player == Opponent:
		ctx.Resolver.PlayFromOpponentDiscard(ctx.Controller, id)
	case e.From == Discard:
		ctx.Resolver.PlayFromDiscard(ctx.Controller, id)
	case e.From == Archives:
		ctx.Resolver.PlayFromArchives(ctx.Controller, id)
	default:
		ctx.Resolver.PlayFromHand(ctx.Controller, id)
	}
	ctx.It, ctx.HasIt = id, true
	return true
}

// candidates are the cards in the source zone the filters admit. The source is
// the controller's own pile, or the opponent's discard pile when Player is
// Opponent (Mimicry).
func (e PlayFrom) candidates(ctx *EffectContext) []LocalID {
	player := ctx.Controller
	if e.Player == Opponent {
		player = ctx.Opponent()
	}
	source := ctx.Resolver.Hand(player)
	switch e.From {
	case Discard:
		source = ctx.Resolver.Discard(player)
	case Archives:
		source = ctx.Resolver.Archives(player)
	}
	var out []LocalID
	for _, id := range source {
		if !e.Types.all() && !e.Types.has(ctx.Resolver.TypeOf(id)) {
			continue
		}
		if !e.House.matches(ctx, id) {
			continue
		}
		out = append(out, id)
	}
	return out
}

// PlayFromOpponent has the controller play a card from a zone of their opponent's
// as their own (the Murkens options): the top card of the opponent's deck (From
// Deck), or a uniformly random card from their facedown archives (From Archives).
// The card is played as the controller's own play — it counts against their
// card-play limit — and, if it stays in play, the controller controls it while its
// owner stays the opponent. An empty source zone does nothing.
type PlayFromOpponent struct {
	// From is the opponent's zone the card comes from: Deck (its top card) or
	// Archives (a uniformly random card, since the archives are facedown).
	From Zone
}

// validate requires a supported source zone.
func (e PlayFromOpponent) validate() error {
	switch e.From {
	case Deck, Archives:
		return nil
	default:
		return fmt.Errorf("PlayFromOpponent: From must be Deck or Archives")
	}
}

// Text renders the effect's printed clause.
func (e PlayFromOpponent) Text() string {
	if e.From == Deck {
		return "play the top card of your opponent's deck"
	}
	return "play a random card from your opponent's archives"
}

// Resolve plays the card from the opponent's zone as the controller's own.
func (e PlayFromOpponent) Resolve(ctx *EffectContext) {
	ctx.Resolver.PlayFromOpponent(ctx.Controller, e.From)
}
