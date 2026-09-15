package engine

// resolveBonusIcons resolves a played card's printed bonus icons, one at a time,
// top to bottom, before its "Play:" abilities. Each icon's source is the card
// itself.
//
// Two deliberate divergences from KeyForge: the source of each icon's effect is
// the card carrying it (KeyForge treats the game as the source), and if the
// creature leaves play before or during resolution the remaining icons do not
// resolve (KeyForge resolves them all once the card is played). Only a creature
// is in play as its icons resolve — an artifact deals its bonus damage to other
// creatures, an upgrade resolves before it attaches, and an action is never in
// play — so the leaves-play check applies to creatures alone.
func (g *Game) resolveBonusIcons(player int, id LocalID) {
	icons := g.cat.def(id).Bonuses
	if len(icons) == 0 {
		return
	}
	gated := g.TypeOf(id) == Creature
	for _, ic := range icons {
		if gated && !g.inPlay(id) {
			return
		}
		g.resolveBonusIcon(player, id, ic)
	}
}

// resolveBonusIcon resolves a single bonus icon, crediting the card as its source.
func (g *Game) resolveBonusIcon(player int, id LocalID, ic BonusIcon) {
	switch ic {
	case BonusAember:
		if capturer, ok := g.gainAember(player, 1); ok {
			g.record(BonusAemberCaptured{Creature: capturer, Card: id, Amount: 1})
			return
		}
		g.record(BonusAemberGained{Player: player, Card: id, Amount: 1})
	case BonusDraw:
		if g.draw(player, 1) > 0 {
			g.record(BonusCardDrawn{Player: player, Card: id, Amount: 1})
		}
	case BonusDamage:
		g.resolveBonusDamage(player, id)
	case BonusCapture:
		g.resolveBonusCapture(player, id)
	}
}

// resolveBonusDamage deals 1 damage from a Damage bonus icon to a creature in play
// the player chooses. The damage may hit any creature, friendly or enemy — and
// with no enemy creature in play it must land on a friendly one.
func (g *Game) resolveBonusDamage(player int, id LocalID) {
	cands := g.creaturesInPlay(0)
	cands = append(cands, g.creaturesInPlay(1)...)
	if len(cands) == 0 {
		return
	}
	target, ok := g.ChooseCreature(player, id, "Choose a creature to deal 1 bonus damage to", cands)
	if !ok {
		return
	}
	g.dealDamage(player, DamageTarget{
		ID:            target,
		Amount:        1,
		Source:        id,
		SourceKeyword: bonusDamage,
	})
}

// resolveBonusCapture has a friendly creature the player chooses capture 1 Æmber
// from the opponent. It does nothing when the opponent's pool is empty or the
// player controls no creature to hold the Æmber.
func (g *Game) resolveBonusCapture(player int, id LocalID) {
	opp := 1 - player
	if g.State.Aember[opp] == 0 {
		return
	}
	cands := g.creaturesInPlay(player)
	if len(cands) == 0 {
		return
	}
	captor, ok := g.ChooseCreature(player, id, "Choose a creature to capture 1 bonus Æmber", cands)
	if !ok {
		return
	}
	g.State.Aember[opp]--
	g.addAmberOn(captor, 1)
	g.record(BonusCaptured{Creature: captor, Card: id, Amount: 1})
}

// creaturesInPlay lists a player's in-play creatures.
func (g *Game) creaturesInPlay(player int) []LocalID {
	var out []LocalID
	for _, c := range g.allInPlay(player) {
		if g.TypeOf(c) == Creature {
			out = append(out, c)
		}
	}
	return out
}
