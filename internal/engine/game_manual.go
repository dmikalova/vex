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

// Manual reports whether manual mode is on.
func (g *Game) Manual() bool { return g.manual }

// SetManual turns manual mode on or off. While on, active-house checks on
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
	g.record(ManualCardMoved{Player: o, Card: id, To: dest})
}

// removeFromAnyZone removes id from whatever holds it. A card in play sheds its
// upgrades (to the discard) and its per-match state; a card in a resting zone is
// simply unlisted.
func (g *Game) removeFromAnyZone(id LocalID) {
	o := g.owner(id)
	if g.inPlay(id) {
		art, hasArt := g.giganticPartner(id)
		g.removeFromPlay(id)
		g.discardUpgrades(id)
		g.resetCore(id)
		if hasArt { // tear down the gigantic art half so it does not dangle
			g.removeFromPlay(art)
			g.resetCore(art)
		}
		return
	}
	g.State.Hand[o].remove(id)
	g.State.Deck[o].remove(id)
	g.State.Discard[o].remove(id)
	g.State.Archives[o].remove(id)
	g.State.Purge[o].remove(id)
}

// ManualAttachUnder removes a card from wherever it rests or sits in play and
// places it under host, face up (graft) or face down (place under). Like every
// manual operation it performs no rule checks; a card taken from play sheds its
// upgrades and per-match state on the way under (removeFromAnyZone).
func (g *Game) ManualAttachUnder(host, id LocalID, faceDown bool) {
	o := g.owner(id)
	g.removeFromAnyZone(id)
	g.AttachUnder(host, id, faceDown)
	g.record(CardPutUnder{Player: o, Card: id, Host: host, FaceDown: faceDown})
}

// ManualDetachToHand sends a selected upgrade or under-card to its owner's hand,
// detaching it from its host first and shedding its per-match state. It is the
// manual counterpart to "return to hand" for an attached card; a card that is
// neither an upgrade nor placed under a host is left where it is.
func (g *Game) ManualDetachToHand(id LocalID) {
	if _, ok := g.detachUpgrade(id); !ok {
		if _, ok := g.detachUnder(id); !ok {
			return
		}
	}
	o := g.owner(id)
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
// artifact row. The card sheds its old zone first (removeFromAnyZone).
func (g *Game) ManualPlaceInPlay(id LocalID, index int) {
	o := g.owner(id)
	g.removeFromAnyZone(id)
	g.State.Cards[id].ArmorRemaining = int16(g.Def(id).Armor)
	if g.Def(id).Type == Creature {
		line := &g.State.Battleline[o]
		if index < 0 {
			index = 0
		}
		if index > int(line.Count) {
			index = int(line.Count)
		}
		line.insertAt(index, id)
	} else {
		g.State.Artifacts[o].add(id)
	}
	g.record(ManualPlacedInPlay{Player: o, Card: id})
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
	n := g.State.Chains[player] + delta
	if n < 0 {
		n = 0
	}
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
