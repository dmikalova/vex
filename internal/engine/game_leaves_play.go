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
// "Destroyed:" abilities resolve (see destroyTogether). A gigantic's two halves
// leave together (ADR 0042).
//
// It deliberately does NOT consult ward. Ward prevents a creature from being
// tagged for destruction, but a ward counter placed on a creature that is already
// tagged does not remove the tag — so once the "Destroyed:" abilities finish, the
// tagged card goes to the discard regardless (Master Rulebook, Ward). A ward
// gained mid-window still stops the creature leaving play any OTHER way, which is
// why purgeFromPlay and the bounce and archive paths do check it. Pinned by
// TestDiscardDestroyedIgnoresWard and TestWardAbsorbsDestructionAndPurgeSeparately.
func (g *Game) discardDestroyed(id LocalID) {
	g.fileFromPlay(id, func(half LocalID, o int) { g.State.Discard[o].add(half) })
}

// fileFromPlay tears down each half of the gigantic id belongs to (ADR 0042) and
// hands each to file with its owner, so a destination only says which pile the
// card lands in and what the log calls it. It is the FILING step, downstream of
// ward: a caller whose move is a removal attempt owes a ward check and goes
// through leavePlayInto instead. Only discardDestroyed files directly, because the
// destruction it completes was already attempted.
func (g *Game) fileFromPlay(id LocalID, file func(half LocalID, owner int)) {
	for _, half := range g.giganticHalves(id) {
		file(half, g.leavePlayTeardown(half))
	}
}

// leavePlayInto is fileFromPlay preceded by the one ward check a removal attempt
// owes. The check is per attempt, not per half, so a warded gigantic absorbs the
// whole removal and neither half leaves.
func (g *Game) leavePlayInto(id LocalID, file func(half LocalID, owner int)) {
	if g.absorbedByWard(id, wardLeavePlay, 0) {
		return
	}
	g.fileFromPlay(id, file)
}

// purgeFromPlay moves a card from play to its owner's purge pile (set aside out of
// the game), shedding its upgrades, Æmber, and per-match state on the way — the
// "purge this creature" a Destroyed ability can do (Annihilation Ritual). A card
// purged as it is destroyed leaves play, so destroyTogether then skips discarding it.
//
// Unlike discardDestroyed, this is a removal ATTEMPT, so it consults ward — even on
// a creature already tagged for destruction. A creature whose ward absorbed the
// destruction and which a later "Destroyed:" ability then wards again absorbs this
// purge with the fresh ward.
func (g *Game) purgeFromPlay(id LocalID) {
	g.leavePlayInto(id, func(half LocalID, o int) {
		g.State.Purge[o].add(half)
		g.record(CardPurged{Card: half})
	})
}

// absorbedByWard reports whether a warded creature's ward absorbs a removal from
// play or an instance of damage. A warded creature instead loses its ward and
// stays; the effect that tried to remove or damage it still resolves, just with no
// effect on the creature. Ward covers only leaving play and damage — not stun,
// enrage, capture, control, or power loss — so only those funnels ask it.
// Each removal attempt asks separately, so a creature can absorb a destruction and
// then absorb a later purge with a ward granted in between; filing an
// already-tagged card into the discard is not an attempt and never asks
// (see discardDestroyed).
// prevented names what the ward stopped (and amount the damage it refused) so the
// log can say what it saved the creature from.
func (g *Game) absorbedByWard(id LocalID, prevented wardPrevented, amount int) bool {
	c := g.stateOf(id)
	if c == nil || !c.Warded {
		return false
	}
	c.Warded = false
	g.record(WardAbsorbed{Creature: id, Prevented: prevented, Amount: amount})
	return true
}

// leavePlayTeardown performs the shared teardown every card leaving play gets,
// whether it was destroyed or merely relocated: it removes the card from the
// battle line and artifact row, discards its upgrades and the cards under it,
// releases any Æmber on it (releaseAemberOnLeavePlay), and resets its per-match
// state. It returns the owner so the caller can file the card in the right zone.
// Every exit routes through here — a caller that hand-inlines these five steps
// will drift from the ones that do not.
//
// The upgrades are shed here and now, rather than by a later sweep over upgrades
// whose host has gone. That eager discard was investigated as a refactor and
// rejected: it runs with `settling` held, so the power drop it causes is noticed
// at the same resolution boundary a sweep would use, and a dangling upgrade is
// already an invariant violation (see checkStaleInPlayState) that -tags assert
// runs never trip. A sweep would be cleanup for a state that cannot occur. If a
// card ever does orphan an upgrade, that invariant fires first.
// An attached upgrade leaves play through here too, not through a chain-only
// counterpart. It is a card in play in its own right, so shedding it releases the
// Æmber sitting on it whichever pile it lands in — discarded with its host,
// returned to hand, or archived. removeFromPlay already unlinks it from its host's
// chain (detachUpgrade), so the one path serves a card in a row and a card on a
// chain alike, and emitLeavesPlay is a no-op for it because an upgrade's abilities
// are granted to its host rather than held by itself
// (TestUpgradeAbilitiesBelongToItsHost). A hand-inlined counterpart is what once
// let an upgrade carry its Æmber away (TestUpgradeReleasesAemberOnLeavingHost).
func (g *Game) leavePlayTeardown(id LocalID) int {
	o := g.owner(id)
	g.removeFromPlay(id)
	g.discardUpgrades(id)
	g.discardUnder(id)
	g.releaseAemberOnLeavePlay(id)
	g.resetCore(id)
	return o
}

// releaseAemberOnLeavePlay hands the Æmber sitting on a card that is leaving play
// to its destination, before resetCore zeroes it. A creature's Æmber goes to the
// pool of its controller's opponent — the side it was captured from, so a creature
// taken into the opponent's control gives its Æmber back to its owner — while a
// non-creature card's Æmber returns to the common supply (Master Rulebook line
// 927). Every exit funnels through here, so bouncing, archiving, or shuffling a
// laden creature away releases its Æmber just as destroying it does.
func (g *Game) releaseAemberOnLeavePlay(id LocalID) {
	amt := int(g.State.Cards[id].Amber)
	if amt <= 0 {
		return
	}
	if g.TypeOf(id) == Creature {
		to := 1 - g.controller(id)
		g.SetAember(to, g.Aember(to)+amt)
		g.record(AemberOnCardReleased{Card: id, Amount: amt, To: to})
		return
	}
	g.record(AemberMovedToCommonSupply{Card: id, Amount: amt})
}

// removeFromPlay takes a card out of play for good: it fires the card's Leaves
// Play abilities, sheds its counters, reverts any cards it was controlling and
// sheds its own control entries, then unlists it from every in-play zone. Every
// real exit — destroyed, purged, returned to hand, archived, put on or shuffled
// into a deck, grafted — funnels through here; the resolution boundary that drove
// the exit settles the board it left (ADR 0029).
// A change of control is NOT an exit: the card stays in play on the other side, so
// control uses unlistFromPlay directly and never fires Leaves Play or sheds
// counters.
func (g *Game) removeFromPlay(id LocalID) {
	g.emitLeavesPlay(id)
	g.clearCounters(id)
	g.releaseControlHeldBy(id)
	g.unlistFromPlay(id)
	// An attached upgrade is listed in its host's upgrade chain, not a battleline
	// or artifact row, so unlink it there too — otherwise a leave-play mover
	// (ArchiveFromPlay on "each upgrade on <self>", Away Team) leaves it dangling
	// on the host, which then re-sheds it to the discard pile. A no-op for anything
	// not attached.
	g.detachUpgrade(id)
	g.clearControls(id)
}

// unlistFromPlay removes id from both players' battlelines and artifact rows,
// without the Leaves Play teardown. It is the shared zone move under both a real
// exit (removeFromPlay) and a change of control, which pulls a creature from one
// side to re-add it on the other while it stays in play. A card normally appears
// only under its owner, but control moves it into another player's rows while
// ownership stays fixed, so the move must scan both rows. The card may have been
// buffing what it leaves behind, but the resolution boundary settles that, not
// this low-level move (ADR 0029).
func (g *Game) unlistFromPlay(id LocalID) {
	for p := 0; p < 2; p++ {
		g.State.Battleline[p].remove(id)
		g.State.Artifacts[p].remove(id)
	}
}

// emitLeavesPlay fires a card's "Leaves Play:" abilities while it is still on the
// board. Every exit — destroyed, purged, returned to hand, archived, shuffled away
// — funnels through removeFromPlay, so this one call covers them all. It also fires
// any effect the card armed to resolve when it leaves play (Turnkey's forced
// forge), which likewise wants the card still on the board so its controller reads.
//
// Inside a simultaneous batch the abilities are gathered here — still while the
// card is on the board, so its controller and ability list read correctly — but
// held for the end of the batch, so every card has moved before any of them
// reacts (flushDeferredLeaves). Pinned by TestLeavesPlayWaitsForTheWholeBatch.
func (g *Game) emitLeavesPlay(id LocalID) {
	if g.deferringLeaves {
		g.deferredLeaves = append(g.deferredLeaves, g.leavesPlayWindow(id)...)
		g.fireScheduledOnLeave(id)
		return
	}
	g.triggerAbilities(id, TriggerLeavesPlay, 0, false)
	g.fireScheduledOnLeave(id)
}

// leavesPlayWindow gathers a card's "Leaves Play:" abilities without resolving
// them, stamped with the controller it had while still on the board.
func (g *Game) leavesPlayWindow(id LocalID) []triggeredAbility {
	pending := g.triggeredBy(id, TriggerLeavesPlay)
	actor := g.controller(id)
	for i := range pending {
		pending[i].actor = int8(actor)
		pending[i].fromLeave = true
	}
	return pending
}

// discardUpgrades moves a card's attached upgrades to their owner's discard pile.
// A card that leaves play — destroyed or relocated — sheds its upgrades this way;
// they do not follow it to hand, deck, or archives.
func (g *Game) discardUpgrades(id LocalID) {
	for _, up := range g.upgradesOf(id) {
		g.State.Discard[g.leavePlayTeardown(up)].add(up)
		g.record(UpgradeDiscarded{Upgrade: up, Host: id})
	}
}

// archiveUpgrade detaches an attached upgrade from its host and puts it into its
// owner's archives — Ghostform grants its host "Fight/Reap: Archive Ghostform",
// which sends the upgrade itself to the archives.
func (g *Game) archiveUpgrade(up LocalID) {
	o := g.leavePlayTeardown(up)
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

// applyDestructionReplacement runs the replacement that stands in for a creature's
// destruction, if one applies — the creature's own definition-level Replace
// (Reassembling Automaton "instead of destroying it, ... move it to a flank") or an
// attached Upgrade's (Armageddon Cloak). The replacement resolves with the creature
// as "it" and its carrier as the source, so a Sequence of "fully heal it" and
// "destroy this Upgrade" saves the creature and consumes the Upgrade. It reports
// whether the destruction was replaced.
func (g *Game) applyDestructionReplacement(id LocalID) bool {
	// A destruction replacement is bounded by the Rule of Six: it counts as a usage
	// of the creature's name (Reassembling Automaton "instead move it to a flank"),
	// so it stands in for at most six of that name's destructions in one turn across
	// a player's copies. A creature stuck at 0 power — which the replacement's heal
	// does not lift — would otherwise be re-detected and replaced on every sweep pass,
	// hanging the game; once the name's pool is spent the next destruction resolves
	// for real instead.
	if g.atRuleOfSix(id) {
		return false
	}
	src, r, ok := g.destructionReplacement(id)
	if !ok {
		return false
	}
	g.record(DestructionReplaced{Card: id, By: src})
	g.recordUsage(id)
	r.With.Resolve(&EffectContext{
		Resolver:   g,
		Source:     src,
		It:         id,
		HasIt:      true,
		Controller: g.controller(id),
	})
	return true
}

// destructionReplacement finds the replacement that stands in for this creature's
// destruction, returning its carrier and the Replace. The creature's own
// definition-level replacement takes precedence over an attached Upgrade's; a
// replacement whose condition is not met (Reassembling Automaton with no other
// friendly creature) does not apply.
func (g *Game) destructionReplacement(id LocalID) (LocalID, Replace, bool) {
	if r := g.cat.def(id).Static.Replaces; g.replacesDestruction(id, r) {
		return id, r, true
	}
	for up, ok := g.firstUpgrade(id); ok; up, ok = g.nextUpgrade(up) {
		if r := g.cat.def(up).Static.Replaces; g.replacesDestruction(id, r) {
			return up, r, true
		}
	}
	return 0, Replace{}, false
}

// replacesDestruction reports whether r replaces the destruction of creature id: it
// must target destruction and, when it carries a condition, that condition must
// hold read from the creature's own controller's perspective (so "there is another
// friendly creature" counts the creature owner's board even when an enemy destroys
// it).
func (g *Game) replacesDestruction(id LocalID, r Replace) bool {
	if !r.valid() || r.When != EventCreatureDestroyed {
		return false
	}
	if r.Cond != nil && !r.Cond.Met(&EffectContext{
		Resolver:   g,
		Source:     id,
		Controller: g.controller(id),
	}) {
		return false
	}
	return true
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
// never counts as destroyed. A ward absorbs the destruction (spent first), or a
// replacement stands in for it — the creature's own (Reassembling Automaton) or an
// attached Upgrade's (Armageddon Cloak).
func (g *Game) filterUndestroyed(ids []LocalID) []LocalID {
	var out []LocalID
	for _, id := range ids {
		if g.absorbedByWard(id, wardDestruction, 0) || g.hasKeyword(id, Invulnerable) ||
			g.applyDestructionReplacement(id) {
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
// they choose — one ability at a time, and for a creature carrying more than one
// Destroyed ability, choosing which resolves next. A Destroyed ability that
// destroys more creatures folds their Destroyed abilities into this same window, so
// a whole chain of deaths resolves together and only then reaches the discard pile
// (point 1; ADR 0013). A creature that leaves play mid-window (Annihilation Ritual
// purges it) cannot resolve its remaining abilities. Once every enrolled creature
// is in the discard pile — and so counts as destroyed — the "after ... destroyed"
// reactions fire (Neffru, Pile of Skulls, Loot the Bodies).
//
// The first destruction to run opens the window and owns it; a destruction that
// runs while it is open (a Destroyed ability's own destruction) enrolls its
// creatures and returns, leaving the owner to resolve, discard, and fire the
// after-destruction reactions for the whole window.
func (g *Game) destroyTogether(controller int, ids []LocalID) {
	ids = g.filterUndestroyed(ids)
	// A creature caught by a nested destruction while its Destroyed window is still
	// open (Harbinger of Doom's "Destroyed: destroy each creature" re-selects
	// Harbinger) is already enrolled; drop it so its Destroyed abilities do not fire
	// twice and it is discarded once.
	ids = g.excludeOpenWindows(ids)

	owner := !g.destroyWindowOpen
	if owner {
		g.destroyWindowOpen = true
		g.destroyWindowController = controller
	}
	g.enrollDestroyed(ids)
	if !owner {
		return
	}
	g.resolveDestroyedWindow()
	members := append([]LocalID(nil), g.destroyingWindow...)
	g.discardDestroyWindow()
	// Close the window before the after-destruction reactions: they are a fresh
	// event, so a destruction one of them causes opens its own new window.
	g.destroyWindowOpen = false
	g.destroyingWindow = nil
	g.destroyPending = nil
	g.afterDestroyedWindow(members)
}

// enrollDestroyed adds a batch of creatures to the open Destroyed window: it
// records their destruction, raises the enemy-destroyed tally, adds them to the
// window's membership, and queues their Destroyed abilities. The window owner
// resolves the queued abilities, discards every member, and fires the
// after-destruction reactions once for the whole window.
func (g *Game) enrollDestroyed(ids []LocalID) {
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
	// "Each time an enemy creature is destroyed": the enemy-destroyed tally rises now,
	// in the destruction window, so a Destroyed ability resolving mid-window reads the
	// deaths that have already happened. The Æmber reactions that respond to the
	// destruction (Loot the Bodies) wait for the after-destruction window, which fires
	// once the batch reaches the discard pile.
	for _, id := range ids {
		if g.TypeOf(id) == Creature {
			g.State.TurnHistory[1-g.controller(id)][EnemyCreaturesDestroyed]++
			g.State.TurnHistory[g.controller(id)][FriendlyCreaturesDestroyed]++
		}
	}
	g.destroyingWindow = append(g.destroyingWindow, ids...)
	g.destroyPending = append(g.destroyPending, g.destroyedAbilities(ids)...)
}

// resolveDestroyedWindow resolves the open window's Destroyed abilities one at a
// time, re-gathering the queue after each so a Destroyed ability that destroys more
// creatures folds their abilities into the same window (enrollDestroyed queued
// them). The window's controller picks the order (ADR 0013), and each picked entry
// resolves through the shared resolveTriggered step every trigger window uses, so a
// creature taken out of play mid-window (Annihilation Ritual purges it) drops its
// remaining abilities by the same source guard (RAW §190, ADR 0030). Nothing is
// discarded until the queue drains, so every creature stays in play — and can be
// seen by the abilities that fire — for the whole event; this window therefore does
// not settle per entry (a lethal creature must not be swept away mid-window), unlike
// the up-front-ordered windows.
func (g *Game) resolveDestroyedWindow() {
	for len(g.destroyPending) > 0 {
		i := g.pickNextReaction(g.destroyWindowController, orderDestroyedPrompt, g.destroyPending)
		t := g.destroyPending[i]
		g.destroyPending = append(g.destroyPending[:i], g.destroyPending[i+1:]...)
		// Resolving a Destroyed: ability is a usage of its creature's name, so it is
		// bounded by the Rule of Six: once that name's pool is spent this turn the
		// remaining Destroyed: abilities of that name do not resolve, though the
		// creature still leaves play.
		if g.atRuleOfSix(t.source) {
			continue
		}
		g.recordUsage(t.source)
		g.resolveTriggered(t)
	}
}

// discardDestroyWindow moves every enrolled creature still in play to its discard
// pile, skipping one that already left play mid-window (Annihilation Ritual purged
// it). A creature whose destruction was replaced was dropped before enrollment, so
// it never reaches this window.
func (g *Game) discardDestroyWindow() {
	for _, id := range g.destroyingWindow {
		if g.inPlay(id) {
			g.discardDestroyed(id)
		}
	}
}

// afterDestroyedWindow fires, as one ordered window, every reaction to the window's
// creatures now reaching their discard piles — so a card destroyed in the same
// window is out of play and cannot be chosen or react to the deaths alongside it
// (e.g. Pile of Skulls cannot capture onto a friendly creature that died in the
// same combat).
func (g *Game) afterDestroyedWindow(members []LocalID) {
	pending := g.afterDestroyedReactions(members)
	g.resolveWindow(g.orderTriggered(g.State.ActivePlayer, pending))
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
// deck, clearing the per-match state it accrued while in play. A gigantic's two
// halves leave together (ADR 0042).
func (g *Game) putOnTopOfDeck(id LocalID) {
	g.leavePlayInto(id, func(half LocalID, o int) {
		g.State.Deck[o].addFront(half)
		g.record(CardPutOnTopOfDeck{Card: half, Owner: o})
	})
}

// putIntoHand removes a card from play and places it into its owner's hand,
// clearing the per-match state it accrued while in play. A gigantic's two halves
// leave together (ADR 0042).
func (g *Game) putIntoHand(id LocalID) {
	g.leavePlayInto(id, func(half LocalID, o int) {
		g.State.Hand[o].add(half)
		g.record(CardPutIntoHand{Card: half, Owner: o})
	})
}

// putIntoArchives removes a card from play and places it into its owner's
// archives, clearing the per-match state it accrued while in play. A gigantic's
// two halves leave together (ADR 0042).
func (g *Game) putIntoArchives(id LocalID) {
	g.leavePlayInto(id, func(half LocalID, o int) {
		g.State.Archives[o].add(half)
		g.record(CardPutIntoArchives{Card: half, Owner: o})
	})
}

// putIntoArchivesEach archives a snapshot of in-play cards simultaneously, then
// settles once. Holding the settling flag stops one card's leave-play from
// destroying another still in the batch (Epic Quest archives "Lion" Bautrem and
// the neighbor it was buffing at the same time, so the neighbor is archived, not
// destroyed for the power it just lost).
func (g *Game) putIntoArchivesEach(controller int, ids []LocalID) {
	g.simultaneously(controller, func() {
		for _, id := range ids {
			g.putIntoArchives(id)
		}
	})
}

// putIntoDeckShuffled removes a card from play and shuffles it into its owner's
// deck, clearing the per-match state it accrued while in play. A gigantic's two
// halves leave together (ADR 0042).
func (g *Game) putIntoDeckShuffled(id LocalID) {
	g.leavePlayInto(id, func(half LocalID, o int) {
		g.State.Deck[o].add(half)
		g.Shuffle(o)
		if g.batchingShuffle {
			g.shuffleBatch = append(g.shuffleBatch, half)
			return
		}
		g.record(CardShuffledIntoDeck{Card: half, Owner: o})
	})
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
	g.leavePlayInto(id, func(half LocalID, o int) {
		g.State.Archives[player].add(half)
		g.record(CardAbducted{Player: player, Card: half, Owner: o})
	})
}
