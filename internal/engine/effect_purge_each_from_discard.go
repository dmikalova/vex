package engine

import (
	"strings"
)

// PurgeEachFromDiscard purges every matching card from both players' discard
// piles at once — no choice — and gives each purged card's owner 1 Æmber when
// GainOwnerAember is set. It is Soldiers to Flowers: "Purge each Untamed creature
// from each player's discard pile. For each card purged this way, its owner gains
// 1 Æmber."
type PurgeEachFromDiscard struct {
	// House restricts the purge to cards of this house; HouseNone allows any.
	House House
	// Type restricts the purge to cards of this type; the zero value allows any.
	Type CardType
	// GainOwnerAember gives each purged card's owner 1 Æmber per purged card.
	GainOwnerAember bool
}

// validate has no required fields; an unfiltered purge is legal.
func (e PurgeEachFromDiscard) validate() error { return nil }

// noun renders the kind of card purged, e.g. "Untamed creature".
func (e PurgeEachFromDiscard) noun() string {
	noun := "card"
	if e.Type != TypeUnset {
		noun = strings.ToLower(e.Type.String())
	}
	if e.House != HouseNone {
		noun = e.House.String() + " " + noun
	}
	return noun
}

// Text renders the effect, e.g. "purge each Untamed creature from each player's
// discard pile. For each card purged this way, its owner gains 1A".
func (e PurgeEachFromDiscard) Text() string {
	text := "purge each " + e.noun() + " from each player's discard pile"
	if e.GainOwnerAember {
		text += ". For each card purged this way, its owner gains 1 Æmber"
	}
	return text
}

// Resolve purges every matching card from both discard piles and, when
// GainOwnerAember is set, gives each purged card's owner 1 Æmber.
func (e PurgeEachFromDiscard) Resolve(ctx *EffectContext) {
	matches := func(id LocalID) bool {
		if e.Type != TypeUnset && ctx.Resolver.TypeOf(id) != e.Type {
			return false
		}
		if e.House != HouseNone && ctx.Resolver.House(id) != e.House {
			return false
		}
		return true
	}
	for _, owner := range []int{ctx.Controller, ctx.Opponent()} {
		doomed := discardCardsWhere(ctx, owner, matches)
		for _, id := range doomed {
			ctx.Resolver.PurgeFromDiscard(owner, id)
			if e.GainOwnerAember {
				ctx.Resolver.GainAember(owner, 1)
			}
		}
	}
}
