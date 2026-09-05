package engine

import (
	"fmt"
	"strings"
)

// This file holds the effects that play a card out of one of a player's piles as
// part of resolving another card — the "play a card" clause, as opposed to a
// player taking their own play action.

// PlayFrom has the controller play a card out of their own hand or discard pile
// right now, ignoring the active-house gate — Phase Shift's off-house card,
// Sacrificial Altar's creature back from the discard pile. From names the source
// pile; House and Type narrow which cards may be chosen; Except inverts the house
// filter, so House names the house that may *not* be played ("a non-Logos card").
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
	// opponent's discard pile (Mimicry). A card played from the opponent's discard
	// is still your play — it counts against your own card-play limit — but returns
	// to its owner's discard pile.
	Player Player
	House  House
	Except bool
	Type   CardType
}

// validate rejects an inverted filter that names no house to exclude, and a
// source zone a card may not be played from.
func (e PlayFrom) validate() error {
	if e.Except && e.House == HouseNone {
		return fmt.Errorf("PlayFrom: Except needs a house to exclude")
	}
	if e.From != Hand && e.From != Discard {
		return fmt.Errorf("PlayFrom: From must be Hand or Discard, got %v", e.From)
	}
	if e.Player == Opponent && e.From != Discard {
		return fmt.Errorf("PlayFrom: only the opponent's discard pile may be played from")
	}
	return nil
}

// Text renders the effect, e.g. "play a non-Logos card" or "play a creature from
// your discard pile". Playing from hand is the default a printed card leaves
// unsaid ("Play a non-Logos card"), so only another zone is named.
func (e PlayFrom) Text() string {
	text := "play " + indefinite(e.noun())
	if e.From != Hand {
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
	if e.Type != TypeUnset && e.Type != AnyType {
		noun = strings.ToLower(e.Type.String())
	}
	switch {
	case e.Except:
		return "non-" + e.House.String() + " " + noun
	case e.House != HouseNone:
		return e.House.String() + " " + noun
	default:
		return noun
	}
}

// Resolve has the controller choose a matching card in the source zone and plays
// it.
func (e PlayFrom) Resolve(ctx *EffectContext) {
	candidates := e.candidates(ctx)
	if len(candidates) == 0 {
		return
	}
	id, ok := ctx.ChooseCreature("Choose "+indefinite(e.noun())+" to play", candidates)
	if !ok {
		return
	}
	switch {
	case e.From == Discard && e.Player == Opponent:
		ctx.Resolver.PlayFromOpponentDiscard(ctx.Controller, id)
	case e.From == Discard:
		ctx.Resolver.PlayFromDiscard(ctx.Controller, id)
	default:
		ctx.Resolver.PlayFromHand(ctx.Controller, id)
	}
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
	if e.From == Discard {
		source = ctx.Resolver.Discard(player)
	}
	var out []LocalID
	for _, id := range source {
		if e.Type != TypeUnset && e.Type != AnyType && ctx.Resolver.TypeOf(id) != e.Type {
			continue
		}
		if e.House != HouseNone && (ctx.Resolver.House(id) == e.House) == e.Except {
			continue
		}
		out = append(out, id)
	}
	return out
}
