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
// discard, handing Æmber to the opponent) BEFORE resetting. Generic counters
// live off CardCore in the global table, so removeFromPlay sheds them separately
// (ADR 0024).
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

// SaveFromDestruction marks a creature whose own "Destroyed:" ability replaced its
// destruction, so the discard step of the current batch leaves it in play
// (Reassembling Automaton). It records the replacement so the log narrates the
// save.
func (g *Game) SaveFromDestruction(id LocalID) {
	if g.savedFromDestruction == nil {
		g.savedFromDestruction = map[LocalID]bool{}
	}
	g.savedFromDestruction[id] = true
	g.record(DestructionReplaced{Card: id, By: id})
}

// purgeFromPlay moves a card from play to its owner's purge pile (set aside out of
// the game), shedding its upgrades, Æmber, and per-match state on the way — the
// "purge this creature" a Destroyed ability can do (Annihilation Ritual). A card
// purged as it is destroyed leaves play, so destroyTogether then skips discarding it.
func (g *Game) purgeFromPlay(id LocalID) {
	if g.absorbedByWard(id) {
		return
	}
	o := g.leavePlayDestroyed(id)
	g.State.Purge[o].add(id)
	g.record(CardPurged{Card: id})
}

// absorbedByWard reports whether a warded creature's ward absorbs a removal from
// play or an instance of damage. A warded creature instead loses its ward and
// stays; the effect that tried to remove or damage it still resolves, just with no
// effect on the creature. Ward covers only leaving play and damage — not stun,
// enrage, capture, control, or power loss — so only those funnels ask it.
func (g *Game) absorbedByWard(id LocalID) bool {
	c := g.stateOf(id)
	if c == nil || !c.Warded {
		return false
	}
	c.Warded = false
	g.record(WardAbsorbed{Creature: id})
	return true
}

// leavePlayDestroyed performs the shared teardown when a destroyed card leaves
// play: it removes the card from the battle line and artifact row, discards its
// upgrades, hands any Æmber on it to the controller's opponent, and resets its
// per-match state. It returns the owner so the caller can file the card in the
// right zone. The Æmber goes to the opponent of whoever controlled the creature
// at the moment it died, not the owner's opponent — a creature taken by the
// opponent (exiled into their control) gives its Æmber back to its owner when it
// dies under that control.
func (g *Game) leavePlayDestroyed(id LocalID) int {
	o := g.owner(id)
	to := 1 - g.controller(id)
	g.removeFromPlay(id)
	g.discardUpgrades(id)
	g.discardUnder(id)
	core := &g.State.Cards[id]
	if core.Amber > 0 {
		g.State.Aember[to] += int(core.Amber)
		g.record(AemberOnCardReleased{Card: id, Amount: int(core.Amber), To: to})
	}
	g.resetCore(id)
	return o
}

// removeFromPlay takes a card out of play for good: it fires the card's Leaves
// Play abilities, sheds its counters, reverts any cards it was controlling and
// sheds its own control entries, then unlists it from every in-play zone and
// settles the board it left. Every real exit — destroyed, purged, returned to
// hand, archived, put on or shuffled into a deck, grafted — funnels through here.
// A change of control is NOT an exit: the card stays in play on the other side, so
// control uses unlistFromPlay directly and never fires Leaves Play or sheds
// counters.
func (g *Game) removeFromPlay(id LocalID) {
	g.emitLeavesPlay(id)
	g.clearCounters(id)
	g.releaseControlHeldBy(id)
	g.unlistFromPlay(id)
	g.clearControls(id)
}

// unlistFromPlay removes id from both players' battlelines and artifact rows and
// settles the board it left, without the Leaves Play teardown. It is the shared
// zone move under both a real exit (removeFromPlay) and a change of control, which
// pulls a creature from one side to re-add it on the other while it stays in play.
// A card normally appears only under its owner, but control moves it into another
// player's rows while ownership stays fixed, so the move must scan both rows.
func (g *Game) unlistFromPlay(id LocalID) {
	for p := 0; p < 2; p++ {
		g.State.Battleline[p].remove(id)
		g.State.Artifacts[p].remove(id)
	}
	// The card may have been buffing the power of what it leaves behind.
	g.settleDestroyed(g.State.ActivePlayer)
}

// emitLeavesPlay fires a card's "Leaves Play:" abilities while it is still on the
// board. Every exit — destroyed, purged, returned to hand, archived, shuffled away
// — funnels through removeFromPlay, so this one call covers them all.
func (g *Game) emitLeavesPlay(id LocalID) {
	g.triggerAbilities(id, TriggerLeavesPlay, 0, false)
}

// discardUpgrades moves a card's attached upgrades to their owner's discard pile.
// A card that leaves play — destroyed or relocated — sheds its upgrades this way;
// they do not follow it to hand, deck, or archives.
func (g *Game) discardUpgrades(id LocalID) {
	for _, up := range g.upgradesOf(id) {
		g.detachUpgrade(up)
		g.releaseControlHeldBy(up)
		g.clearCounters(up)
		g.State.Discard[g.owner(up)].add(up)
	}
}

// returnUpgradesToHand detaches each upgrade attached to a host still in play and
// puts it into its owner's hand instead of the discard pile — Transporter Platform
// returns a creature and its upgrades together. Call it before the host leaves
// play, so its upgrades are already gone when the host's own move would shed them.
func (g *Game) returnUpgradesToHand(host LocalID) {
	for _, up := range g.upgradesOf(host) {
		g.detachUpgrade(up)
		g.releaseControlHeldBy(up)
		g.clearCounters(up)
		g.resetCore(up)
		o := g.owner(up)
		g.State.Hand[o].add(up)
		g.record(CardReturnedToHand{Card: up, Owner: o})
	}
}

// archiveUpgrade detaches an attached upgrade from its host and puts it into its
// owner's archives — Ghostform grants its host "Fight/Reap: Archive Ghostform",
// which sends the upgrade itself to the archives. An attached upgrade is not
// listed in a battleline or artifact row, so it must be unlinked from its host's
// chain rather than routed through the ordinary leave-play path.
func (g *Game) archiveUpgrade(up LocalID) {
	g.detachUpgrade(up)
	g.releaseControlHeldBy(up)
	g.clearCounters(up)
	g.resetCore(up)
	o := g.owner(up)
	g.State.Archives[o].add(up)
	g.record(CardPutIntoArchives{Card: up, Owner: o})
}

// discardUnder moves the cards placed under a host to their owners' discard
// piles when the host leaves play — generalizing Graft's own rule (if the card
// onto which it is grafted leaves play, the grafted card is placed in its
// owner's discard pile) to every card placed under a card, not only a grafted
// one (see ADR 0016). A card that leaves play — destroyed or relocated — sheds
// what is placed under it this way; it does not follow the host to hand, deck,
// or archives.
func (g *Game) discardUnder(id LocalID) {
	for _, u := range g.underOf(id) {
		g.detachUnder(u)
		g.State.Discard[g.owner(u)].add(u)
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
	g.record(DestructionReplaced{Card: id, By: up})
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
	g.record(CardDestroyed{Card: upgrade})
	g.releaseControlHeldBy(upgrade)
	g.resetCore(upgrade)
	g.State.Discard[g.owner(upgrade)].add(upgrade)
}

// filterUndestroyed drops creatures whose destruction does not happen from the
// pending set before any Destroyed abilities are collected, so a saved creature
// never counts as destroyed. A ward absorbs the destruction (spent first), or an
// attached Upgrade replaces it (Armageddon Cloak).
func (g *Game) filterUndestroyed(controller int, ids []LocalID) []LocalID {
	var out []LocalID
	for _, id := range ids {
		if g.absorbedByWard(id) || g.hasKeyword(id, Invulnerable) ||
			g.applyDestructionReplacement(controller, id) {
			continue
		}
		out = append(out, id)
	}
	return out
}

// excludeOpenWindows drops creatures whose Destroyed window is already open — the
// ones an enclosing destruction batch is still resolving. A nested destruction
// that re-selects them (Harbinger of Doom's "destroy each creature" re-selects
// Harbinger) leaves them to the enclosing batch, so they neither re-fire their
// Destroyed abilities nor get discarded twice.
func (g *Game) excludeOpenWindows(ids []LocalID) []LocalID {
	if len(g.destroyingWindow) == 0 {
		return ids
	}
	var out []LocalID
	for _, id := range ids {
		if !g.inDestroyingWindow(id) {
			out = append(out, id)
		}
	}
	return out
}

// inDestroyingWindow reports whether a creature's Destroyed window is already open.
func (g *Game) inDestroyingWindow(id LocalID) bool {
	for _, w := range g.destroyingWindow {
		if w == id {
			return true
		}
	}
	return false
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
	ids = g.filterUndestroyed(controller, ids)
	// A creature caught by a nested destruction while its own Destroyed window is
	// still open (Harbinger of Doom's "Destroyed: destroy each creature" re-selects
	// Harbinger) is already being destroyed by the enclosing batch; drop it here so
	// its Destroyed abilities do not fire twice and it is discarded once.
	ids = g.excludeOpenWindows(ids)
	g.destroyingWindow = append(g.destroyingWindow, ids...)
	defer func() {
		g.destroyingWindow = g.destroyingWindow[:len(g.destroyingWindow)-len(ids)]
	}()
	// The source is consumed by the batch it directly targets; any state-based
	// deaths that follow (a creature that lost a buff) narrate passively.
	source, hasSource := g.destroyingSource, g.hasDestroyingSource
	g.hasDestroyingSource = false
	if hasSource && len(ids) > 0 {
		g.record(CardsDestroyedBy{Source: source, Cards: append([]LocalID(nil), ids...)})
	} else {
		for _, id := range ids {
			g.record(CardDestroyed{Card: id})
		}
	}
	// "Each time an enemy creature is destroyed": the destroyed creature's controller
	// is the enemy of whoever watches, so the reaction fires for that opponent. The
	// count and the lasting "each time destroyed" event fire here, in the destruction
	// window; the "after ... destroyed" reactions wait until the batch reaches the
	// discard pile (below), because a card is not destroyed until it lands there.
	for _, id := range ids {
		if g.TypeOf(id) == Creature {
			g.State.TurnHistory[1-g.controller(id)][EnemyCreaturesDestroyed]++
		}
		g.emitLasting(EventEnemyCreatureDestroyed, 1-g.controller(id), id)
	}
	// The whole window is ordered once, up front, by the active player (ADR 0013).
	// A creature that leaves play mid-window (Annihilation Ritual purges it) simply
	// drops its remaining abilities as they come up.
	for _, t := range g.orderTriggered(controller, TriggerDestroyed, g.destroyedAbilities(ids)) {
		if !g.inPlay(t.source) {
			continue
		}
		closeFrame := g.openFrame(Frame{
			Actor:      g.controller(t.source),
			Source:     t.source,
			HasSource:  true,
			Trigger:    TriggerDestroyed,
			Grantor:    t.grantor,
			HasGrantor: t.grantor != t.source,
		})
		t.ability.Effect.Resolve(
			&EffectContext{Resolver: g, Source: t.source, Controller: g.controller(t.source)},
		)
		closeFrame()
	}
	for _, id := range ids {
		if g.inPlay(id) && !g.savedFromDestruction[id] {
			g.discardDestroyed(id)
		}
	}
	for _, id := range ids {
		delete(g.savedFromDestruction, id)
	}
	// Only now, with the batch in the discard pile, do the "after ... destroyed"
	// reactions fire — Neffru's "after a creature is destroyed" and Pile of Skulls'
	// "after an enemy creature is destroyed" — so a card destroyed in this same batch
	// is out of play and cannot be chosen or react to the deaths alongside it (e.g.
	// Pile of Skulls cannot capture onto a friendly creature that died in the same
	// combat).
	for _, id := range ids {
		if g.TypeOf(id) == Creature {
			g.emitAfterEnemyDestroyed(id)
			g.emitAfterFriendlyDestroyed(id)
			g.emitAfterCreatureDestroyed(id)
		}
	}
}

// destroyEach destroys each id simultaneously (KeyForge's shared Destroyed
// timing), letting the controller order how their "Destroyed:" abilities resolve.
// An id that is an attached Upgrade is detached and discarded instead — that is how
// "destroy this Upgrade" (Armageddon Cloak destroying itself) resolves through the
// ordinary Destroy effect. Callers pass a snapshot of distinct ids. Once the batch
// is done the board is settled, since the dead may have been buffing the living.
func (g *Game) destroyEach(controller int, ids []LocalID) {
	g.destroyBatch(controller, ids)
	g.settleDestroyed(controller)
}

// destroyBatch is one simultaneous destruction, with no state-based sweep after
// it. It holds the settling flag for its duration so the leave-play of each card
// in the batch does not trigger a sweep mid-batch.
func (g *Game) destroyBatch(controller int, ids []LocalID) {
	was := g.settling
	g.settling = true
	defer func() { g.settling = was }()
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
	if g.absorbedByWard(id) {
		return
	}
	o := g.owner(id)
	g.removeFromPlay(id)
	g.discardUpgrades(id)
	g.discardUnder(id)
	g.resetCore(id)
	g.State.Deck[o].addFront(id)
	g.record(CardPutOnTopOfDeck{Card: id, Owner: o})
}

// putIntoHand removes a card from play and places it into its owner's hand,
// clearing the per-match state it accrued while in play.
func (g *Game) putIntoHand(id LocalID) {
	if g.absorbedByWard(id) {
		return
	}
	o := g.owner(id)
	g.removeFromPlay(id)
	g.discardUpgrades(id)
	g.discardUnder(id)
	g.resetCore(id)
	g.State.Hand[o].add(id)
	g.record(CardReturnedToHand{Card: id, Owner: o})
}

// putIntoArchives removes a card from play and places it into its owner's
// archives, clearing the per-match state it accrued while in play.
func (g *Game) putIntoArchives(id LocalID) {
	if g.absorbedByWard(id) {
		return
	}
	o := g.owner(id)
	g.removeFromPlay(id)
	g.discardUpgrades(id)
	g.discardUnder(id)
	g.resetCore(id)
	g.State.Archives[o].add(id)
	g.record(CardPutIntoArchives{Card: id, Owner: o})
}

// putIntoArchivesEach archives a snapshot of in-play cards simultaneously, then
// settles once. Holding the settling flag stops one card's leave-play from
// destroying another still in the batch (Epic Quest archives "Lion" Bautrem and
// the neighbor it was buffing at the same time, so the neighbor is archived, not
// destroyed for the power it just lost).
func (g *Game) putIntoArchivesEach(controller int, ids []LocalID) {
	was := g.settling
	g.settling = true
	for _, id := range ids {
		g.putIntoArchives(id)
	}
	g.settling = was
	g.settleDestroyed(controller)
}

// putIntoDeckShuffled removes a card from play and shuffles it into its owner's
// deck, clearing the per-match state it accrued while in play.
func (g *Game) putIntoDeckShuffled(id LocalID) {
	if g.absorbedByWard(id) {
		return
	}
	o := g.owner(id)
	g.removeFromPlay(id)
	g.discardUpgrades(id)
	g.discardUnder(id)
	g.resetCore(id)
	g.State.Deck[o].add(id)
	g.Shuffle(o)
	if g.batchingShuffle {
		g.shuffleBatch = append(g.shuffleBatch, id)
		return
	}
	g.record(CardShuffledIntoDeck{Card: id, Owner: o})
}

// shuffleFriendlyInPlayIntoDeck shuffles every card a player controls in play —
// each creature and artifact, and every upgrade attached to them — into their
// deck, and returns how many cards were shuffled. An upgrade is detached first so
// it is shuffled back into the deck rather than shed to the discard pile. The
// caller opens a shuffle batch around it, so the whole sweep narrates as one
// grouped line.
func (g *Game) shuffleFriendlyInPlayIntoDeck(player int) int {
	count := 0
	for _, host := range append(g.Battleline(player), g.Artifacts(player)...) {
		for _, up := range g.upgradesOf(host) {
			g.detachUpgrade(up)
			g.putIntoDeckShuffled(up)
			count++
		}
		g.putIntoDeckShuffled(host)
		count++
	}
	return count
}

// Only three zones of yours may hold a card your opponent owns: your battleline,
// your artifact line, and your archives. A card that would move to any other zone
// of yours — your hand, your discard pile, your deck — goes to its owner's
// matching zone instead. So an enemy creature abducted into your archives goes to
// your opponent's hand when you take your archives up, and to their discard pile
// if those archives are discarded.
// PutIntoYourArchives removes a creature from play into the archives of the player
// who took it rather than its owner's — Mass Abduction, Sample Collection, and
// Uxlyx the Zookeeper all abduct this way. Nothing is marked on the card: the
// ownership rule above sends it home the moment it leaves those archives.
func (g *Game) PutIntoYourArchives(id LocalID, player int) {
	if g.absorbedByWard(id) {
		return
	}
	o := g.owner(id)
	g.removeFromPlay(id)
	g.discardUpgrades(id)
	g.discardUnder(id)
	g.resetCore(id)
	g.State.Archives[player].add(id)
	g.record(CardAbducted{Player: player, Card: id, Owner: o})
}
