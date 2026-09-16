package engine

// The two halves of a gigantic creature are linked while in play by a single
// symmetric CardCore byte, GiganticPartnerPlus, encoded +1 so its zero value is
// "no partner" — the same idiom as the Upgrade and Under chains (ADR 0001), but
// a plain pair rather than a list because a gigantic has exactly two halves. The
// base half holds the battleline slot and is the creature; the art half is in
// play but slot-less and lends only its bonus icons. See ADR 0042.

// giganticPlus encodes a LocalID into the +1 form the partner link uses, where 0
// means "no partner".
func giganticPlus(id LocalID) uint8 { return uint8(id) + 1 }

// decodeGigantic reverses giganticPlus, returning ok=false for the zero value.
func decodeGigantic(v uint8) (LocalID, bool) {
	if v == 0 {
		return 0, false
	}
	return LocalID(v - 1), true
}

// linkGiganticPartners points the two halves of a gigantic at each other, so
// either half can find the other in O(1).
func (g *Game) linkGiganticPartners(base, art LocalID) {
	g.State.Cards[base].GiganticPartnerPlus = giganticPlus(art)
	g.State.Cards[art].GiganticPartnerPlus = giganticPlus(base)
}

// giganticPartner returns a half's partner, or ok=false when the card is not a
// linked gigantic half.
func (g *Game) giganticPartner(id LocalID) (LocalID, bool) {
	return decodeGigantic(g.State.Cards[id].GiganticPartnerPlus)
}

// unlinkGigantic clears the link on both halves. It is safe to call on a lone
// card, which does nothing.
func (g *Game) unlinkGigantic(id LocalID) {
	partner, ok := g.giganticPartner(id)
	if !ok {
		return
	}
	g.State.Cards[partner].GiganticPartnerPlus = 0
	g.State.Cards[id].GiganticPartnerPlus = 0
}

// giganticHalves returns both halves of the gigantic id belongs to — the card
// itself and its partner — or just id when it is a lone card. Callers moving a
// gigantic out of play use this so both halves travel together (ADR 0042).
func (g *Game) giganticHalves(id LocalID) []LocalID {
	if partner, ok := g.giganticPartner(id); ok {
		return []LocalID{id, partner}
	}
	return []LocalID{id}
}

// oppositeGiganticRole returns the role a half must find in its partner: the base
// pairs with the art half and vice versa.
func oppositeGiganticRole(r GiganticRole) GiganticRole {
	if r == GiganticBase {
		return GiganticArt
	}
	return GiganticBase
}

// giganticBaseArt sorts a matched pair of halves into (base, art) by their roles,
// so the caller always knows which one holds the battleline slot.
func (g *Game) giganticBaseArt(a, b LocalID) (base, art LocalID) {
	if g.cat.def(a).GiganticRole == GiganticBase {
		return a, b
	}
	return b, a
}

// recruitGiganticPartner finds the opposite-role half of the gigantic being
// played — same name, matched one-to-one — in one of the given source zones,
// returning it with a closure that removes it from that zone. It reports ok=false
// when no partner is available, which is why a lone half cannot be played.
func (g *Game) recruitGiganticPartner(
	id LocalID,
	sources []*deckList,
) (LocalID, func(), bool) {
	def := g.cat.def(id)
	want := oppositeGiganticRole(def.GiganticRole)
	for _, src := range sources {
		for _, other := range src.slice() {
			if other == id {
				continue
			}
			od := g.cat.def(other)
			if od.GiganticRole == want && od.Name == def.Name {
				s := src
				return other, func() { s.remove(other) }, true
			}
		}
	}
	return 0, nil, false
}

// bonusIconsOf returns the bonus icons a card contributes when played. A gigantic
// base half prints none of its own — they live on its art half — so it reports
// its partner's icons, which is how the combined creature resolves its icons once.
func (g *Game) bonusIconsOf(id LocalID) []BonusIcon {
	def := g.cat.def(id)
	if def.GiganticRole == GiganticBase {
		if art, ok := g.giganticPartner(id); ok {
			return g.cat.def(art).Bonuses
		}
	}
	return def.Bonuses
}

// playGigantic plays a gigantic creature from hand: it recruits the played half's
// opposite-role partner from the same hand by one-to-one match, brings both
// halves into play together — the base half holds the battleline slot, the art
// half is linked to it slot-less — and runs the single standard creature-play
// sequence on the base. A lone half whose partner is not in hand cannot be played
// (ADR 0042): nothing is removed and the play returns ErrGiganticNoPartner.
func (g *Game) playGigantic(
	player int,
	id LocalID,
	remove func(),
	opts playCardOptions,
) (LocalID, error) {
	if g.cannotPlayCreatures(player) {
		return 0, ErrCannotPlayCreature
	}
	partner, removePartner, ok := g.recruitGiganticPartner(
		id, []*deckList{&g.State.Hand[player]},
	)
	if !ok {
		return 0, ErrGiganticNoPartner
	}
	base, art := g.giganticBaseArt(id, partner)
	g.recordCardPlayed(player, base, opts)
	remove()
	removePartner()
	// Link before playing the base so its bonus icons — read through the link from
	// the art half — resolve once as part of the standard sequence.
	g.linkGiganticPartners(base, art)
	g.playCreatureCard(player, base, opts.flank)
	g.applyTreachery(player, base)
	return base, nil
}
