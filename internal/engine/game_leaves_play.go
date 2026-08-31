package engine

// This file holds how a card LEAVES PLAY: destruction (with the KeyForge
// simultaneous "Destroyed:" timing) and relocation to the deck, hand, or
// archives. All of these shed the card's per-match state and its attached
// upgrades on the way out.

// resetCore returns a card to its fresh, out-of-play state by zeroing its entire
// CardCore. Every field there is per-match, in-play state (damage, armor, Æmber,
// exhaustion, attached upgrades), so one zeroing covers them all and any field
// added later is reset automatically — no leaves-play path has to remember to
// clear it. Callers apply field-specific side effects (moving upgrades to the
// discard, handing Æmber to the opponent) BEFORE resetting.
func (g *Game) resetCore(id LocalID) { g.State.Cards[id] = CardCore{} }

// discardDestroyed moves a destroyed card from play to its owner's discard: it
// discards attached upgrades, hands any Æmber on the card to the owner's
// opponent, resets the card's per-match state, and adds it to the discard. It is
// the final step of destruction, run only for creatures still in play after their
// "Destroyed:" abilities resolve (see destroyTogether).
func (g *Game) discardDestroyed(id LocalID) {
	o := g.leavePlayDestroyed(id)
	g.State.Discard[o].add(id)
}

// purgeFromPlay moves a card from play to its owner's purge pile (set aside out of
// the game), shedding its upgrades, Æmber, and per-match state on the way — the
// "purge this creature" a Destroyed ability can do (Annihilation Ritual). A card
// purged as it is destroyed leaves play, so destroyTogether then skips discarding it.
func (g *Game) purgeFromPlay(id LocalID) {
	o := g.leavePlayDestroyed(id)
	g.State.Purge[o].add(id)
	g.logf("%s is purged", g.Name(id))
}

// leavePlayDestroyed performs the shared teardown when a destroyed card leaves
// play: it removes the card from the battle line and artifact row, discards its
// upgrades, hands any Æmber on it to the owner's opponent, and resets its per-match
// state. It returns the owner so the caller can file the card in the right zone.
func (g *Game) leavePlayDestroyed(id LocalID) int {
	o := g.owner(id)
	g.removeFromPlay(id)
	g.discardUpgrades(id)
	core := &g.State.Cards[id]
	if core.Amber > 0 {
		g.State.Aember[1-o] += int(core.Amber)
		g.logf("%d Æmber on %s goes to %s's pool", core.Amber, g.Name(id), g.names[1-o])
	}
	g.resetCore(id)
	return o
}

// removeFromPlay removes id from every in-play zone. A card normally appears only
// under its owner, but control-changing effects move creatures into another
// player's battleline while ownership stays fixed, so leave-play teardown must
// scan both players' rows.
func (g *Game) removeFromPlay(id LocalID) {
	for p := 0; p < 2; p++ {
		g.State.Battleline[p].remove(id)
		g.State.Artifacts[p].remove(id)
	}
}

// discardUpgrades moves a card's attached upgrades to their owner's discard pile.
// A card that leaves play — destroyed or relocated — sheds its upgrades this way;
// they do not follow it to hand, deck, or archives.
func (g *Game) discardUpgrades(id LocalID) {
	for _, up := range g.upgradesOf(id) {
		g.detachUpgrade(up)
		g.releaseControlHeldBy(up)
		g.State.Discard[g.owner(up)].add(up)
	}
}

// applyDestructionReplacement runs the first attached Upgrade that replaces its
// host's destruction (Armageddon Cloak). The Upgrade's replacement effect resolves
// with the host as "it" and the Upgrade as its source, so a Sequence of "fully heal
// it" and "destroy this Upgrade" saves the host and consumes the Upgrade. It reports
// whether the destruction was replaced.
func (g *Game) applyDestructionReplacement(controller int, id LocalID) bool {
	up, r, ok := g.destructionReplacement(id)
	if !ok {
		return false
	}
	g.logf("%s would be destroyed, so %s replaces its destruction", g.Name(id), g.Name(up))
	r.With.Resolve(&EffectContext{
		Resolver:   g,
		Source:     up,
		It:         id,
		HasIt:      true,
		Controller: controller,
	})
	return true
}

// destructionReplacement finds the first attached Upgrade whose StaticModifier
// replaces this creature's destruction, returning the Upgrade and its replacement.
func (g *Game) destructionReplacement(id LocalID) (LocalID, Replace, bool) {
	for up, ok := g.firstUpgrade(id); ok; up, ok = g.nextUpgrade(up) {
		if r := g.cat.def(up).Static.Replaces; r.valid() && r.When == EventCreatureDestroyed {
			return up, r, true
		}
	}
	return 0, Replace{}, false
}

// destroyAttachedUpgrade detaches an Upgrade from its host and discards it — an
// Upgrade destroying itself through its own effect (Armageddon Cloak). Detaching
// stitches the host's remaining upgrades back together, so destroying one in the
// middle of a chain leaves the others attached.
func (g *Game) destroyAttachedUpgrade(upgrade LocalID) {
	if _, ok := g.detachUpgrade(upgrade); !ok {
		return
	}
	g.logf("%s is destroyed", g.Name(upgrade))
	g.releaseControlHeldBy(upgrade)
	g.resetCore(upgrade)
	g.State.Discard[g.owner(upgrade)].add(upgrade)
}

// applyDestructionReplacements removes creatures whose destruction was replaced
// from the pending destruction set before any Destroyed abilities are collected.
func (g *Game) applyDestructionReplacements(controller int, ids []LocalID) []LocalID {
	var out []LocalID
	for _, id := range ids {
		if g.applyDestructionReplacement(controller, id) {
			continue
		}
		out = append(out, id)
	}
	return out
}

// destroyTogether destroys several creatures as one simultaneous event, matching
// KeyForge timing. Every creature remains in play while all of their Destroyed
// abilities are collected; the active player resolves those abilities in an order
// they choose — one creature at a time (its creatures highlight for selection),
// and, for a creature carrying more than one Destroyed ability, choosing which of
// its abilities resolves next. A creature that leaves play (e.g. Annihilation
// Ritual purges it) cannot resolve any of its remaining abilities. Then every
// creature still in play goes to its discard pile.
func (g *Game) destroyTogether(controller int, ids []LocalID) {
	ids = g.applyDestructionReplacements(controller, ids)
	for _, id := range ids {
		g.logf("%s is destroyed", g.Name(id))
	}
	// "Each time an enemy creature is destroyed": the destroyed creature's controller
	// is the enemy of whoever watches, so the reaction fires for that opponent.
	for _, id := range ids {
		g.emitLasting(EventEnemyCreatureDestroyed, 1-g.controller(id), id)
		g.emitEnemyDestroyed(id)
	}
	pending := g.destroyedAbilities(ids)
	for {
		// Drop abilities whose source has left play (e.g. an earlier Destroyed
		// ability purged the creature); they can no longer resolve.
		pending = g.inPlayAbilities(pending)
		if len(pending) == 0 {
			break
		}
		// Choose the next creature whose Destroyed ability to resolve. With several
		// destroyed creatures still holding abilities the controller picks one (they
		// highlight on the board); a lone creature is forced.
		sources := distinctSources(pending)
		src := sources[0]
		if len(sources) > 1 {
			if chosen, ok := g.pickCreature(controller, "", "Choose the next creature whose Destroyed ability resolves", sources); ok {
				src = chosen
			}
		}
		pick := g.pickDestroyedAbility(controller, pending, src)
		ab := pending[pick]
		g.logf("%s: %s", g.Name(ab.source), renderAbilityLine(g.cat.def(ab.source), ab.ability))
		ab.ability.Effect.Resolve(&EffectContext{Resolver: g, Source: ab.source, Controller: g.controller(ab.source)})
		pending = append(pending[:pick], pending[pick+1:]...)
	}
	for _, id := range ids {
		if g.inPlay(id) {
			g.discardDestroyed(id)
		}
	}
}

// inPlayAbilities keeps only the pending abilities whose source is still in play.
func (g *Game) inPlayAbilities(pending []triggeredAbility) []triggeredAbility {
	var kept []triggeredAbility
	for _, ab := range pending {
		if g.inPlay(ab.source) {
			kept = append(kept, ab)
		}
	}
	return kept
}

// distinctSources lists the distinct source creatures in pending, first-seen order.
func distinctSources(pending []triggeredAbility) []LocalID {
	var out []LocalID
	for _, ab := range pending {
		seen := false
		for _, id := range out {
			if id == ab.source {
				seen = true
				break
			}
		}
		if !seen {
			out = append(out, ab.source)
		}
	}
	return out
}

// pickDestroyedAbility returns the pending index of the Destroyed ability to
// resolve next for src: forced when the creature carries one, otherwise the
// controller picks which of its abilities via a labeled prompt.
func (g *Game) pickDestroyedAbility(controller int, pending []triggeredAbility, src LocalID) int {
	var idxs []int
	for i, ab := range pending {
		if ab.source == src {
			idxs = append(idxs, i)
		}
	}
	if len(idxs) == 1 {
		return idxs[0]
	}
	labels := make([]string, len(idxs))
	for i, j := range idxs {
		labels[i] = renderAbilityLine(g.cat.def(src), pending[j].ability)
	}
	return idxs[g.chooseOption(controller, "", "Choose which Destroyed ability resolves", labels)]
}

// destroyEach destroys each id simultaneously (KeyForge's shared Destroyed
// timing), letting the controller order how their "Destroyed:" abilities resolve.
// An id that is an attached Upgrade is detached and discarded instead — that is how
// "destroy this Upgrade" (Armageddon Cloak destroying itself) resolves through the
// ordinary Destroy effect. Callers pass a snapshot of distinct ids.
func (g *Game) destroyEach(controller int, ids []LocalID) {
	var creatures []LocalID
	for _, id := range ids {
		if _, ok := g.hostOf(id); ok {
			g.destroyAttachedUpgrade(id)
			continue
		}
		creatures = append(creatures, id)
	}
	g.destroyTogether(controller, creatures)
}

// putOnTopOfDeck removes a card from play and places it on top of its owner's
// deck, clearing the per-match state it accrued while in play.
func (g *Game) putOnTopOfDeck(id LocalID) {
	o := g.owner(id)
	g.removeFromPlay(id)
	g.discardUpgrades(id)
	g.resetCore(id)
	g.State.Deck[o].addFront(id)
	g.logf("%s is put on top of %s's deck", g.Name(id), g.names[o])
}

// putIntoHand removes a card from play and places it into its owner's hand,
// clearing the per-match state it accrued while in play.
func (g *Game) putIntoHand(id LocalID) {
	o := g.owner(id)
	g.removeFromPlay(id)
	g.discardUpgrades(id)
	g.resetCore(id)
	g.State.Hand[o].add(id)
	g.logf("%s is returned to %s's hand", g.Name(id), g.names[o])
}

// putIntoArchives removes a card from play and places it into its owner's
// archives, clearing the per-match state it accrued while in play.
func (g *Game) putIntoArchives(id LocalID) {
	o := g.owner(id)
	g.removeFromPlay(id)
	g.discardUpgrades(id)
	g.resetCore(id)
	g.State.Archives[o].add(id)
	g.logf("%s is put into %s's archives", g.Name(id), g.names[o])
}

// putIntoDeckShuffled removes a card from play and shuffles it into its owner's
// deck, clearing the per-match state it accrued while in play.
func (g *Game) putIntoDeckShuffled(id LocalID) {
	o := g.owner(id)
	g.removeFromPlay(id)
	g.discardUpgrades(id)
	g.resetCore(id)
	g.State.Deck[o].add(id)
	g.Shuffle(o)
	g.logf("%s is shuffled into %s's deck", g.Name(id), g.names[o])
}
