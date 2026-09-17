package engine

import "slices"

// This file holds the one internal mechanism the cross-zone movers share:
// gathering a player's cards from a set of source zones and detecting which zone
// each picked card came from. The actual source-by-destination move is the shared
// Destination.moveFrom (ADR 0031); a mover only supplies the origin and delegates.
// The authoring verbs — ReturnNamedToHand, Search,
// ShuffleNamedFromDiscardIntoDeck, ShuffleChosenCreaturesFromZones — are thin
// wrappers that supply their own filter, prompt, and printed text and delegate
// the zone plumbing here.

// crossZoneMover gathers a player's cards from Sources and moves each picked card
// to Dest. Sources is ordered with the assumed fallback zone last — a picked card
// always came from a source, so originOf checks every earlier source and takes the
// last without a redundant lookup, mirroring the zone probes the movers used
// before they shared this.
type crossZoneMover struct {
	Player  int
	Dest    Destination
	Sources []Zone
}

// zoneCards reads the Player's cards in one source zone.
func (m crossZoneMover) zoneCards(ctx *EffectContext, z Zone) []LocalID {
	switch z {
	case inPlay:
		return ctx.Resolver.Battleline(m.Player)
	case Hand:
		return ctx.Resolver.Hand(m.Player)
	case Discard:
		return ctx.Resolver.Discard(m.Player)
	default:
		return ctx.Resolver.Deck(m.Player)
	}
}

// gather collects, in source order, every card in a source zone that keep admits.
func (m crossZoneMover) gather(ctx *EffectContext, keep func(LocalID) bool) []LocalID {
	var out []LocalID
	for _, z := range m.Sources {
		out = append(out, filterIDs(m.zoneCards(ctx, z), keep)...)
	}
	return out
}

// originOf reports which source zone holds a picked card, checking each source but
// the last and assuming the last when none matched.
func (m crossZoneMover) originOf(ctx *EffectContext, id LocalID) Zone {
	last := len(m.Sources) - 1
	for _, z := range m.Sources[:last] {
		if m.inZone(ctx, z, id) {
			return z
		}
	}
	return m.Sources[last]
}

// inZone reports whether a card sits in one source zone.
func (m crossZoneMover) inZone(ctx *EffectContext, z Zone, id LocalID) bool {
	if z == inPlay {
		return resolverInPlay(ctx, id)
	}
	return slices.Contains(m.zoneCards(ctx, z), id)
}

// move sends one picked card to Dest through the shared move matrix, keyed on the
// zone the card came from (ADR 0031).
func (m crossZoneMover) move(ctx *EffectContext, id LocalID) {
	m.Dest.moveFrom(ctx, m.originOf(ctx, id), m.Player, id)
}
