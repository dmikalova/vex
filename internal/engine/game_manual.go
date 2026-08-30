package engine

// This file holds manual (sandbox) mode: unrestricted operations a UI exposes so
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

var manualZoneNames = [...]string{"hand", "top of deck", "bottom of deck", "discard", "archives", "purge"}

// String returns the printed zone name.
func (z ManualZone) String() string {
	if int(z) < len(manualZoneNames) {
		return manualZoneNames[z]
	}
	return "unknown"
}

// Manual reports whether sandbox mode is on.
func (g *Game) Manual() bool { return g.manual }

// SetManual turns sandbox mode on or off. While on, active-house checks on
// playing and using cards are lifted (see inActiveHouse).
func (g *Game) SetManual(on bool) { g.manual = on }

// ManualMove takes a card from wherever it is — a resting zone or in play — and
// places it into dest for its owner, shedding any in-play state (upgrades to the
// discard, damage/Æmber/counters cleared) if it was on the board.
func (g *Game) ManualMove(id LocalID, dest ManualZone) {
	o := g.owner(id)
	g.removeFromAnyZone(id)
	switch dest {
	case ManualHand:
		g.State.Hand[o].add(id)
	case ManualDeckTop:
		g.State.Deck[o].addFront(id)
	case ManualDeckBottom:
		g.State.Deck[o].add(id)
	case ManualDiscard:
		g.State.Discard[o].add(id)
	case ManualArchives:
		g.State.Archives[o].add(id)
	case ManualPurge:
		g.State.Purge[o].add(id)
	}
	g.logf("%s manually moves %s to %s", g.names[o], g.Name(id), dest)
}

// removeFromAnyZone removes id from whatever holds it. A card in play sheds its
// upgrades (to the discard) and its per-match state; a card in a resting zone is
// simply unlisted.
func (g *Game) removeFromAnyZone(id LocalID) {
	o := g.owner(id)
	if g.inPlay(id) {
		g.removeFromPlay(id)
		g.discardUpgrades(id)
		g.resetCore(id)
		return
	}
	g.State.Hand[o].remove(id)
	g.State.Deck[o].remove(id)
	g.State.Discard[o].remove(id)
	g.State.Archives[o].remove(id)
	g.State.Purge[o].remove(id)
}

// ManualSetExhausted sets or clears a card's exhausted flag — readying an
// exhausted creature, or exhausting a ready one.
func (g *Game) ManualSetExhausted(id LocalID, exhausted bool) {
	g.State.Cards[id].Exhausted = exhausted
	if exhausted {
		g.logf("%s is manually exhausted", g.Name(id))
	} else {
		g.logf("%s is manually readied", g.Name(id))
	}
}

// ManualAddCard registers def as a new card owned by player and places it in
// their hand, returning its id — so a sandbox can pull any card from the pool.
func (g *Game) ManualAddCard(def CardDefinition, player int) LocalID {
	id := g.Register(def, player)
	g.State.Hand[player].add(id)
	g.logf("%s manually adds %s to hand", g.names[player], g.Name(id))
	return id
}

// ManualAddAmber adjusts player's Æmber pool by delta (clamped at zero), so a
// sandbox can dial each player's Æmber up or down.
func (g *Game) ManualAddAmber(player, delta int) {
	n := g.State.Aember[player] + delta
	if n < 0 {
		n = 0
	}
	g.State.Aember[player] = n
	g.logf("%s now has %d Æmber (manual)", g.names[player], n)
}

// ManualAddChains adjusts player's chain count by delta (clamped at zero).
func (g *Game) ManualAddChains(player, delta int) {
	n := g.State.Chains[player] + delta
	if n < 0 {
		n = 0
	}
	g.State.Chains[player] = n
	g.logf("%s now has %d %s (manual)", g.names[player], n, chainNoun(n))
}

// ManualSetActiveHouse sets the active player's active house directly, so a
// sandbox can switch houses mid-turn.
func (g *Game) ManualSetActiveHouse(h House) {
	g.State.ActiveHouse = h
	g.logf("%s manually chooses %s as their active house", g.names[g.State.ActivePlayer], h)
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
	if g.State.Keys[player] >= KeysToWin {
		return
	}
	g.State.KeyColors[player][g.State.Keys[player]] = c
	g.State.Keys[player]++
	g.logf("%s manually forges a %s key (%d/%d)", g.names[player], c, g.State.Keys[player], KeysToWin)
}

// ManualUnforgeKey removes player's most recently forged key, if any.
func (g *Game) ManualUnforgeKey(player int) {
	if g.State.Keys[player] <= 0 {
		return
	}
	g.State.Keys[player]--
	g.State.KeyColors[player][g.State.Keys[player]] = KeyColorNone
	g.logf("%s manually unforges a key (%d/%d)", g.names[player], g.State.Keys[player], KeysToWin)
}
