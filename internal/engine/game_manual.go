package engine

// This file holds manual mode: unrestricted operations a UI exposes so
// the player can rearrange the game outside the normal rules — moving cards
// between zones, readying or exhausting cards, adding an arbitrary card, and
// (via the manual flag) playing and using cards regardless of house. These
// perform no rule checks; the frontend gates them behind a mode the player turns
// on, and the engine merely carries them out.

// ManualZone names a resting zone a manual move can send a card to.
type ManualZone uint8

// The manual resting zones a card can be moved to.
const (
	ManualHand ManualZone = iota
	ManualDeckTop
	ManualDeckBottom
	ManualDiscard
	ManualArchives
	ManualPurge
)

var manualZoneNames = [...]string{
	"hand",
	"top of deck",
	"bottom of deck",
	"discard",
	"archives",
	"purge",
}

// String returns the printed zone name.
func (z ManualZone) String() string {
	if int(z) < len(manualZoneNames) {
		return manualZoneNames[z]
	}
	return "unknown"
}

// zone returns the resting zone a manual destination lands in. ManualZone stays
// its own type rather than becoming a Destination because its ordinals are
// persisted in saved games (web/replay.go reads one back as an int).
func (z ManualZone) zone() Zone {
	switch z {
	case ManualDeckTop, ManualDeckBottom:
		return Deck
	case ManualDiscard:
		return Discard
	case ManualArchives:
		return Archives
	case ManualPurge:
		return Purged
	default:
		return Hand
	}
}

// Manual reports whether manual mode is on.
func (g *Game) Manual() bool { return g.manual }

// SetManual turns manual mode on or off. While on, active-house checks on
// playing and using cards are lifted (see inActiveHouse).
func (g *Game) SetManual(on bool) { g.manual = on }

// ManualMove takes a card from wherever it is — a resting zone or in play — and
// places it into dest for its owner. A card leaving the board gets the standard
// leave-play teardown, so its upgrades and the cards under it are discarded and
// the Æmber on it is released just as an effect-driven move would release it
// (TestManualMoveReleasesAember). Manual mode grants permission to take an action
// that would normally need a card to authorize it; it does not change the action.
func (g *Game) ManualMove(id LocalID, dest ManualZone) {
	g.manualRelocate(id, func(card LocalID, o int) {
		if dest == ManualDeckTop {
			g.State.Deck[o].addFront(card)
		} else {
			g.pile(zoneRef{Player: o, Zone: dest.zone()}).add(card)
		}
		g.record(ManualCardMoved{Player: o, Card: card, To: dest})
	})
}

// manualRelocate takes id out of whatever holds it and hands it, with its owner,
// to file. A card in play leaves through the ordinary leavePlayInto funnel, so a
// manual move is a removal attempt like any other and a **ward absorbs it**, the
// card staying put and file never running. That is deliberate: manual mode exists
// to let a playtester exercise the real engine paths, and it has its own buttons
// for clearing a ward when the ward is what is in the way. A gigantic is relocated
// as a whole, so file runs once per half (ADR 0042). A card in a resting zone is
// simply unlisted, with no teardown and nothing to absorb it.
// Pinned by TestManualMoveIsAbsorbedByWard.
func (g *Game) manualRelocate(id LocalID, file func(card LocalID, owner int)) {
	if g.inPlay(id) {
		g.leavePlayInto(id, file)
		return
	}
	file(id, g.removeFromRestingZones(id))
}

// ManualAttachUnder removes a card from wherever it rests or sits in play and
// places it under host, face up (graft) or face down (place under). A card taken
// from play sheds its upgrades and per-match state on the way under, and a ward
// absorbs the move as it would any other removal (manualRelocate).
func (g *Game) ManualAttachUnder(host, id LocalID, faceDown bool) {
	g.manualRelocate(id, func(card LocalID, o int) {
		g.AttachUnder(host, card, faceDown)
		g.record(CardPutUnder{Player: o, Card: card, Host: host, FaceDown: faceDown})
	})
}

// ManualDetachToHand sends a selected upgrade or under-card to its owner's hand,
// detaching it from its host first and shedding its per-match state. It is the
// manual counterpart to "return to hand" for an attached card; a card that is
// neither an upgrade nor placed under a host is left where it is.
//
// Releasing the card's Æmber is not a rule check manual mode may skip — it is
// what stops the Æmber being destroyed outright, which no zone edit should do.
// Pinned by TestUpgradeReleasesAemberOnLeavingHost.
func (g *Game) ManualDetachToHand(id LocalID) {
	if _, ok := g.detachUpgrade(id); !ok {
		if _, ok := g.detachUnder(id); !ok {
			return
		}
	}
	o := g.owner(id)
	g.releaseAemberOnLeavePlay(id)
	g.resetCore(id)
	g.State.Hand[o].add(id)
	g.record(ManualCardMoved{Player: o, Card: id, To: ManualHand})
}

// ManualSetExhausted sets or clears a card's exhausted flag — readying an
// exhausted creature, or exhausting a ready one.
func (g *Game) ManualSetExhausted(id LocalID, exhausted bool) {
	g.State.Cards[id].Exhausted = exhausted
	g.record(ManualExhaustSet{Card: id, Exhausted: exhausted})
}

// ManualPlaceInPlay drops a card straight into play for its owner outside the
// normal play flow — no play effects, no bonus Æmber. A creature enters its
// owner's battleline before index (0 the left flank, the line length the right
// flank), so a playtester can deploy it anywhere; anything else enters the
// artifact row. The card sheds its old zone first (manualRelocate), which for a
// card already in play means a ward absorbs the reposition.
func (g *Game) ManualPlaceInPlay(id LocalID, index int) {
	g.manualRelocate(id, func(card LocalID, o int) {
		g.State.Cards[card].ArmorRemaining = int16(g.Def(card).Armor)
		if g.Def(card).Type == Creature {
			line := &g.State.Battleline[o]
			at := min(max(index, 0), int(line.Count))
			line.insertAt(at, card)
		} else {
			g.State.Artifacts[o].add(card)
		}
		g.record(ManualPlacedInPlay{Player: o, Card: card})
	})
}

// ManualAddCard registers def as a new card owned by player and places it in
// their hand, returning its id — so manual mode can pull any card from the pool.
// A match's LocalID space is finite, so it reports false and adds nothing once it
// is exhausted rather than panicking mid-game.
func (g *Game) ManualAddCard(def CardDefinition, player int) (LocalID, bool) {
	if !g.cat.hasRoom() {
		g.record(ManualMatchFull{Player: player})
		return 0, false
	}
	id := g.Register(def, player)
	g.State.Hand[player].add(id)
	g.record(ManualCardAdded{Player: player, Card: id})
	return id, true
}

// ManualAddAmber adjusts player's Æmber pool by delta (clamped at zero), so
// manual mode can dial each player's Æmber up or down.
func (g *Game) ManualAddAmber(player, delta int) {
	n := max(g.Aember(player)+delta, 0)
	g.SetAember(player, n)
	g.record(ManualAemberSet{Player: player, Amount: n})
}

// ManualAddChains adjusts player's chain count by delta (clamped at zero).
func (g *Game) ManualAddChains(player, delta int) {
	n := max(g.State.Chains[player]+delta, 0)
	g.State.Chains[player] = n
	g.record(ManualChainsSet{Player: player, Amount: n})
}

// ManualSetActiveHouse sets the active player's active house directly, so
// manual mode can switch houses mid-turn.
func (g *Game) ManualSetActiveHouse(h House) {
	g.State.ActiveHouse = h
	g.record(ManualHouseChosen{Player: g.State.ActivePlayer, House: h})
}

// ManualForgeKey forges one more key for player using the next unused colour.
func (g *Game) ManualForgeKey(player int) {
	if remaining := g.remainingKeyColors(player); len(remaining) > 0 {
		g.ManualForgeKeyColor(player, remaining[0])
	}
}

// ManualForgeKeyColor forges one more key of colour c for player, up to
// KeysToWin — no cost and no forge triggers.
func (g *Game) ManualForgeKeyColor(player int, c KeyColor) {
	if g.Keys(player) >= KeysToWin {
		return
	}
	g.State.KeyColors[player][g.Keys(player)] = c
	g.record(ManualKeyForged{
		Player: player,
		Color:  c,
		Keys:   g.Keys(player),
		Needed: KeysToWin,
	})
}

// ManualUnforgeKey removes player's most recently forged key, if any.
func (g *Game) ManualUnforgeKey(player int) {
	if g.Keys(player) <= 0 {
		return
	}
	g.State.KeyColors[player][g.Keys(player)-1] = KeyColorNone
	g.record(ManualKeyUnforged{Player: player, Keys: g.Keys(player), Needed: KeysToWin})
}
