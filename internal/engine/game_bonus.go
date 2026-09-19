package engine

import "slices"

// resolveBonusIcons resolves a played card's printed bonus icons, one at a time,
// top to bottom, before its "Play:" abilities. Each icon's source is the card
// itself.
//
// An armed one-shot boost (Wild Bounty) is consumed here, on the next card its
// owner plays, and makes each icon resolve an additional time — interleaved, so an
// icon resolves and then immediately resolves again before the next icon, rather
// than in a second full pass. Each resolution independently rechecks the
// leaves-play and prevention gates, so a boost re-resolution that would land after
// the creature left play, or after a bar arrives, is skipped like a first one.
//
// Two deliberate divergences from KeyForge: the source of each icon's effect is
// the card carrying it (KeyForge treats the game as the source), and if the
// creature leaves play before or during resolution the remaining icons do not
// resolve (KeyForge resolves them all once the card is played). Only a creature
// is in play as its icons resolve — an artifact deals its bonus damage to other
// creatures, an upgrade resolves before it attaches, and an action is never in
// play — so the leaves-play check applies to creatures alone.
func (g *Game) resolveBonusIcons(player int, id LocalID) {
	boosted := g.consumeBonusIconBoost(player)
	icons := g.bonusIconsOf(id)
	if len(icons) == 0 {
		return
	}
	gated := g.TypeOf(id) == Creature
	for _, ic := range icons {
		times := 1
		if boosted {
			times = 2
		}
		for t := 0; t < times; t++ {
			if gated && !g.inPlay(id) {
				return
			}
			// A standing bar (Master of the Grey) stops the player resolving icons.
			// Checked per resolution so a bar arriving mid-resolution stops the rest,
			// and so Wild Bounty's extra resolution is checked independently of the
			// first.
			if g.cannotResolveBonusIcons(player) {
				return
			}
			g.resolveBonusIcon(player, id, ic)
		}
	}
}

// ResolveBonusIconsOn resolves the bonus icons printed on a card as a foreign
// resolution — player resolves the icons on id (read from its definition, so a
// card in hand, discard, or purge still exposes them) without the card being
// played. No leaves-play gate applies, since the card is not in play as a creature;
// the prevention bar (Master of the Grey) and per-icon substitution still apply,
// checked per icon so a bar arriving mid-resolution stops the rest.
func (g *Game) ResolveBonusIconsOn(player int, id LocalID) {
	for _, ic := range g.cat.def(id).Bonuses {
		if g.cannotResolveBonusIcons(player) {
			return
		}
		g.resolveBonusIcon(player, id, ic)
	}
}

// cannotResolveBonusIcons reports whether an in-play card bars this player from
// resolving the bonus icons on cards they play (Master of the Grey).
func (g *Game) cannotResolveBonusIcons(player int) bool {
	for owner := range 2 {
		for _, id := range g.allInPlay(owner) {
			switch g.cat.def(id).Restricts.BonusIcons {
			case Controller:
				if player == owner {
					return true
				}
			case Opponent:
				if player != owner {
					return true
				}
			case EachPlayer:
				return true
			}
		}
	}
	return false
}

// resolveBonusIcon resolves a single bonus icon, crediting the card as its source.
// Before it resolves as an icon, each in-play card the player controls may
// substitute it — resolving the icon as a different icon (Amphora Captura) or as a
// concrete effect (Scrivener Favian) instead. A substitution marked May asks the
// controller first; otherwise it always applies. Substitutions chain: a swapped
// icon is offered in turn to the player's other cards, so Amphora's Capture can be
// caught by Scrivener and turned into a steal. Each source is offered at most once
// and a vacuous swap (an icon to the same icon) is never offered, so the chain
// terminates. The moment a substitution resolves a concrete effect the result is no
// longer a bonus icon, so no further substitution applies and the icon's own effect
// never runs.
func (g *Game) resolveBonusIcon(player int, id LocalID, ic BonusIcon) {
	var offered []LocalID
	for {
		rm, src, ok := g.bonusInsteadFor(player, ic, offered)
		if !ok {
			break
		}
		offered = append(offered, src)
		if rm.May && !g.chooseBonusInstead(player, src, rm) {
			continue
		}
		if rm.Instead != nil {
			rm.Instead.Resolve(&EffectContext{
				Resolver:   g,
				Source:     src,
				Controller: player,
			})
			return
		}
		ic = rm.As
	}
	g.resolveBonusIconEffect(player, id, ic)
}

// bonusInsteadFor returns the first in-play card the player controls, not already in
// offered, whose substitution applies to icon ic: its From naming ic (or any icon)
// and, for an icon swap, an As different from ic — a swap to the same icon is
// vacuous and never offered, which (with the once-per-source rule) is what
// terminates a chain of swaps.
func (g *Game) bonusInsteadFor(
	player int,
	ic BonusIcon,
	offered []LocalID,
) (BonusInstead, LocalID, bool) {
	for _, id := range g.allInPlay(player) {
		if slices.Contains(offered, id) {
			continue
		}
		rm := g.cat.def(id).BonusInstead
		if !rm.set() {
			continue
		}
		if rm.From != bonusUnset && rm.From != ic {
			continue
		}
		if rm.Instead == nil && rm.As == ic {
			continue
		}
		return rm, id, true
	}
	return BonusInstead{}, 0, false
}

// chooseBonusInstead asks the player whether to apply an optional substitution,
// attributing the prompt to the offering card.
func (g *Game) chooseBonusInstead(player int, src LocalID, rm BonusInstead) bool {
	prompt := capitalizeFirst(bonusInsteadClause(rm)) + "?"
	return g.chooseOption(player, g.sourceName(src), prompt, []string{"Yes", "No"}) == 0
}

// resolveBonusIconEffect resolves a single bonus icon's own effect, crediting the
// card as its source.
func (g *Game) resolveBonusIconEffect(player int, id LocalID, ic BonusIcon) {
	switch ic {
	case BonusAember:
		if capturer, ok := g.gainAember(player, 1); ok {
			g.record(BonusAemberCaptured{
				Creature: capturer,
				Card:     id,
				Player:   player,
				Amount:   1,
			})
			return
		}
		g.record(BonusAemberGained{
			Player: player,
			Card:   id,
			Amount: 1,
		})
	case BonusDraw:
		if g.draw(player, 1) > 0 {
			g.record(BonusCardDrawn{
				Player: player,
				Card:   id,
				Amount: 1,
			})
			g.afterBonusReaction(player, TriggerAfterBonusDraw, 0, false)
		}
	case BonusDamage:
		if target, ok := g.resolveBonusDamage(player, id); ok {
			g.afterBonusReaction(player, TriggerAfterBonusDamage, target, true)
		}
	case BonusCapture:
		g.resolveBonusCapture(player, id)
	}
}

// afterBonusReaction fires a reaction window over the resolving player's own
// in-play cards after one of their bonus icons resolves — Chronus (Draw) and
// Maleficorn (Damage, with the creature it hit bound as "it").
func (g *Game) afterBonusReaction(player int, tr Trigger, it LocalID, hasIt bool) {
	w := g.window()
	for _, src := range g.allInPlay(player) {
		w.add(src, tr, it, hasIt)
	}
	g.resolveWindow(g.orderTriggered(player, w.pending))
}

// resolveBonusDamage deals 1 damage from a Damage bonus icon to a creature in play
// the player chooses, returning that creature so a reaction can act on it. The
// damage may hit any creature, friendly or enemy — and with no enemy creature in
// play it must land on a friendly one.
func (g *Game) resolveBonusDamage(player int, id LocalID) (LocalID, bool) {
	cands := g.creaturesInPlay(0)
	cands = append(cands, g.creaturesInPlay(1)...)
	if len(cands) == 0 {
		return 0, false
	}
	target, ok := g.ChooseCreature(player, id, "Choose a creature to deal 1 bonus damage to", cands)
	if !ok {
		return 0, false
	}
	g.dealDamage(player, DamageTarget{
		ID:            target,
		Amount:        1,
		Source:        id,
		SourceKeyword: bonusDamage,
	})
	return target, true
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
	g.record(BonusCaptured{
		Creature: captor,
		Card:     id,
		Amount:   1,
	})
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
