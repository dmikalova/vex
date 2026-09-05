package engine

// Cards placed under a host form an intrusive singly-linked list threaded
// through CardCore, exactly like Upgrades (game_upgrades.go, ADR 0001) — but
// unlike an upgrade, a card placed under a host is OUT OF PLAY: Masterplan and
// Jargogle set a card aside this way rather than leaving it in play. See ADR
// 0016 for why this is a second instance of the intrusive-list decision rather
// than a new mechanism, and for the visibility and leave-play rules that come
// with it (Peekable in game_read.go, discardUnder in game_leaves_play.go).
//
// FirstUnderPlus on the host points at the head, each card's NextUnderPlus
// points at the next sibling, and UnderHostPlus points back at the host.
// UnderFaceDown records whether that particular card is hidden: Graft always
// places its card faceup; Masterplan and Jargogle place theirs facedown.

// underPlus encodes a LocalID into the +1 form the linkage fields use, where 0
// means "none". This mirrors upgradePlus and keeps the zero value "unattached".
func underPlus(id LocalID) uint8 { return uint8(id) + 1 }

// decodeUnder reverses underPlus, returning ok=false for the zero "none" value.
func decodeUnder(v uint8) (LocalID, bool) {
	if v == 0 {
		return 0, false
	}
	return LocalID(v - 1), true
}

// firstUnder returns the first card placed under host, or ok=false if it has
// none. Paired with nextUnder it walks a host's chain with no allocation:
//
//	for id, ok := g.firstUnder(host); ok; id, ok = g.nextUnder(id) { ... }
func (g *Game) firstUnder(host LocalID) (LocalID, bool) {
	return decodeUnder(g.State.Cards[host].FirstUnderPlus)
}

// nextUnder returns the card placed after id under the same host, or ok=false
// at the tail.
func (g *Game) nextUnder(id LocalID) (LocalID, bool) {
	return decodeUnder(g.State.Cards[id].NextUnderPlus)
}

// underOf collects the cards placed under host into a fresh slice, in the order
// they were placed. Callers that mutate the chain while iterating (detaching
// each) use this so they read every link before any is rewritten.
func (g *Game) underOf(host LocalID) []LocalID {
	var out []LocalID
	for id, ok := g.firstUnder(host); ok; id, ok = g.nextUnder(id) {
		out = append(out, id)
	}
	return out
}

// AttachUnder threads a card onto the tail of a host's under-chain, preserving
// placement order, records whether it went facedown or faceup, and sets its
// back-link to the host.
func (g *Game) AttachUnder(host, id LocalID, faceDown bool) {
	core := &g.State.Cards[id]
	core.UnderHostPlus = underPlus(host)
	core.NextUnderPlus = 0
	core.UnderFaceDown = faceDown
	tail, ok := g.firstUnder(host)
	if !ok {
		g.State.Cards[host].FirstUnderPlus = underPlus(id)
		return
	}
	for next, ok := g.nextUnder(tail); ok; next, ok = g.nextUnder(tail) {
		tail = next
	}
	g.State.Cards[tail].NextUnderPlus = underPlus(id)
}

// detachUnder unlinks a card from its host's under-chain, stitching the cards on
// either side of it back together, and clears the card's own links and facedown
// flag. It returns the host, or ok=false when id was not placed under anything.
func (g *Game) detachUnder(id LocalID) (LocalID, bool) {
	core := &g.State.Cards[id]
	host, ok := decodeUnder(core.UnderHostPlus)
	if !ok {
		return 0, false
	}
	next := core.NextUnderPlus
	if head, _ := g.firstUnder(host); head == id {
		g.State.Cards[host].FirstUnderPlus = next
	} else {
		prev := head
		for g.State.Cards[prev].NextUnderPlus != underPlus(id) {
			prev, _ = g.nextUnder(prev)
		}
		g.State.Cards[prev].NextUnderPlus = next
	}
	core.UnderHostPlus = 0
	core.NextUnderPlus = 0
	core.UnderFaceDown = false
	return host, true
}

// underHostOf returns the card a card is placed under, or ok=false when it is
// not placed under anything. The back-link makes this O(1).
func (g *Game) underHostOf(id LocalID) (LocalID, bool) {
	return decodeUnder(g.State.Cards[id].UnderHostPlus)
}

// PutCardUnder removes a card from owner's hand and places it under host, face
// up or face down — Masterplan's and Jargogle's own "put a card from your hand
// under me." It does nothing if the card is not in that hand.
func (g *Game) PutCardUnder(owner int, id, host LocalID, faceDown bool) {
	hand := &g.State.Hand[owner]
	i := hand.indexOf(id)
	if i < 0 {
		return
	}
	hand.removeAt(i)
	g.AttachUnder(host, id, faceDown)
	g.record(CardPutUnder{Player: owner, Card: id, Host: host, FaceDown: faceDown})
}

// GraftUnder moves a card from play to faceup under host, out of play (rulebook:
// Graft). Like the leave-play relocations in game_leaves_play.go it sheds the
// card's upgrades, the cards under it, and its per-match state on the way out;
// unlike them it lands in the host's under-chain rather than a resting zone.
// AttachUnder runs after resetCore because resetCore zeroes the very link fields
// it sets. Spangler Box grafts a creature onto itself.
func (g *Game) GraftUnder(id, host LocalID) {
	g.removeFromPlay(id)
	g.discardUpgrades(id)
	g.discardUnder(id)
	g.resetCore(id)
	g.AttachUnder(host, id, false)
	g.record(CardGrafted{Card: id, Host: host})
}

// PutUnderIntoPlay puts every card placed under host into play under its owner's
// control, detaching each from the under-chain first — Spangler Box's Destroyed
// ability returns the creatures grafted onto it. Reading the chain into a slice
// up front (underOf) lets each card detach as it enters play without disturbing
// the walk.
func (g *Game) PutUnderIntoPlay(host LocalID) {
	for _, u := range g.underOf(host) {
		g.detachUnder(u)
		g.putIntoPlay(u, g.owner(u))
	}
}
