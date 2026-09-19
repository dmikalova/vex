package engine

import "slices"

// This file holds the one internal mechanism the cross-zone movers share:
// gathering a player's cards from a set of source zones and detecting which zone
// each picked card came from. The actual source-by-destination move is the shared
// Destination.moveFrom (ADR 0031); a mover only supplies the origin and delegates.
// The authoring verbs — ReturnNamedToHand, Search,
// ShuffleChosenCreaturesFromZones — are thin
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

// zoneCards reads the Player's cards in one source zone. In play means every card
// they control in play, upgrades included (resolverCardsInPlay), so a mover reaches
// an upgrade as a card in its own right rather than only as cargo on its host.
func (m crossZoneMover) zoneCards(ctx *EffectContext, z Zone) []LocalID {
	switch z {
	case InPlay:
		return resolverCardsInPlay(ctx, m.Player)
	case Hand:
		return ctx.Resolver.Hand(m.Player)
	case Discard:
		return ctx.Resolver.Discard(m.Player)
	case Archives:
		return ctx.Resolver.Archives(m.Player)
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

// originOf reports which source zone holds a picked card. It reports ok false when
// the card has left every source zone since it was picked — a power settle inside
// the batch can destroy one of the picks (Inka the Spider, archived by Biomatrix
// Backup while Song of Spring shuffles) — because moving it out of a zone it is no
// longer in would put it in two places at once
// (TestCrossZoneMoveSkipsCardThatLeftItsSourceZones).
func (m crossZoneMover) originOf(ctx *EffectContext, id LocalID) (Zone, bool) {
	for _, z := range m.Sources {
		if m.inZone(ctx, z, id) {
			return z, true
		}
	}
	return zoneUnset, false
}

// inZone reports whether a card sits in one source zone. It asks zoneCards rather
// than the row-only resolverInPlay so the probe admits exactly what gather offered
// — an upgrade included.
func (m crossZoneMover) inZone(ctx *EffectContext, z Zone, id LocalID) bool {
	return slices.Contains(m.zoneCards(ctx, z), id)
}

// move sends one picked card to Dest through the shared move matrix, keyed on the
// zone the card came from (ADR 0031). A pick that has since left every source zone
// is skipped rather than moved from where it no longer is.
func (m crossZoneMover) move(ctx *EffectContext, id LocalID) {
	if from, ok := m.originOf(ctx, id); ok {
		m.Dest.moveFrom(ctx, from, m.Player, id)
	}
}

// resolverCardsInPlay lists every card a player has in play: their creatures,
// their artifacts, and the upgrades attached to either. A card placed *under*
// another card is not in play — it sits in the out-of-play zone that is "under"
// its host — and is never listed. Each host's upgrades come immediately ahead of
// the host, so a caller that moves the whole list has already dealt with an
// upgrade before its host leaves under it.
//
// This is the one answer to "what is in play" for effects that do not name a card
// type: a Target, Count, or Condition that says "cards in play" reaches upgrades,
// while one that says Creature or Artifact narrows this set by type afterwards.
// Pinned by TestCardsInPlayCountsUpgrades.
func resolverCardsInPlay(ctx *EffectContext, p int) []LocalID {
	var out []LocalID
	for _, host := range append(ctx.Resolver.Battleline(p), ctx.Resolver.Artifacts(p)...) {
		out = append(out, ctx.Resolver.Upgrades(host)...)
		out = append(out, host)
	}
	return out
}

// currentZone reports where a card an effect already selected can still be
// reached and whose copy of that zone holds it: in play under its controller, or
// in its owner's discard pile when the card left play earlier in this same
// resolution. A removal follows its target into the discard rather than fizzling,
// so Yxilo Bolter purges the creature its own damage just killed. A card in
// neither zone has gone somewhere no effect reaches (ADR 0030) and reports false.
// Pinned by TestPurgeCreatureFollowsIntoDiscard.
//
// The rule this implements is about the resolution, not the destination: an
// ability that is **still resolving** may reach its target in any zone whose cards
// are public, because later instructions refer to a departed card as it was
// immediately before it left. The discard pile is public, so a creature killed
// earlier in the same ability is still reachable there; a hand, deck, or archives
// is hidden, so a card that goes there is out of reach for good (Master Rulebook
// line 926, and the Bad Penny ruling — it returns to hand and Yxilo Bolter's purge
// no longer finds it).
//
// A card that has left play and an ability in a LATER window is a different case
// and is not this function's business: Collector Worm does not archive the creature
// it killed because its "after a fight" ability opens a window of its own, by which
// time the target is long gone. That is a fizzle from window separation, not from
// the zone. Do not read the bulk movers' skips as evidence against the rule above;
// theirs is the same window-and-selection question, handled where they select.
func currentZone(ctx *EffectContext, id LocalID) (Zone, int, bool) {
	if resolverInPlay(ctx, id) {
		return InPlay, ctx.Resolver.Controller(id), true
	}
	owner := ctx.Resolver.Owner(id)
	if slices.Contains(ctx.Resolver.Discard(owner), id) {
		return Discard, owner, true
	}
	return zoneUnset, 0, false
}
