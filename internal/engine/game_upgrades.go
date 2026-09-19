package engine

// Upgrades attached to a creature form an intrusive singly-linked list threaded
// through CardCore: FirstUpgradePlus on the host points at the head, each upgrade's
// NextUpgradePlus points at the next sibling, and HostPlus points back at the host.
// This lets a creature carry any number of upgrades with three bytes of flat,
// comparable state instead of a fixed-size array — KeyForge itself sets no limit on
// how many upgrades a creature may hold. Every read of a creature's upgrades walks
// this list; attaching threads onto the tail (preserving attach order) and detaching
// stitches the neighbours back together so a chain stays whole when an upgrade in the
// middle leaves play.

// upgradePlus encodes a LocalID into the +1 form the linkage fields use, where 0
// means "none". This mirrors ControlPlus and keeps the zero value "unattached".
func upgradePlus(id LocalID) uint8 { return uint8(id) + 1 }

// decodeUpgrade reverses upgradePlus, returning ok=false for the zero "none" value.
func decodeUpgrade(v uint8) (LocalID, bool) {
	if v == 0 {
		return 0, false
	}
	return LocalID(v - 1), true
}

// firstUpgrade returns the first upgrade attached to host, or ok=false if it has
// none. Paired with nextUpgrade it walks a host's chain with no allocation:
//
//	for up, ok := g.firstUpgrade(host); ok; up, ok = g.nextUpgrade(up) { ... }
func (g *Game) firstUpgrade(host LocalID) (LocalID, bool) {
	return decodeUpgrade(g.State.Cards[host].FirstUpgradePlus)
}

// nextUpgrade returns the upgrade after up on its host, or ok=false at the tail.
func (g *Game) nextUpgrade(up LocalID) (LocalID, bool) {
	return decodeUpgrade(g.State.Cards[up].NextUpgradePlus)
}

// upgradesOf collects a host's attached upgrades into a fresh slice, in attach
// order. Callers that mutate the chain while iterating (detaching each) use this so
// they read every link before any is rewritten.
func (g *Game) upgradesOf(host LocalID) []LocalID {
	var out []LocalID
	for up, ok := g.firstUpgrade(host); ok; up, ok = g.nextUpgrade(up) {
		out = append(out, up)
	}
	return out
}

// AttachUpgrade threads an upgrade onto the tail of a host creature's chain,
// preserving attach order, and sets the upgrade's back-link to the host.
func (g *Game) AttachUpgrade(host, up LocalID) {
	upCore := &g.State.Cards[up]
	upCore.HostPlus = upgradePlus(host)
	upCore.NextUpgradePlus = 0
	tail, ok := g.firstUpgrade(host)
	if !ok {
		g.State.Cards[host].FirstUpgradePlus = upgradePlus(up)
		return
	}
	for next, ok := g.nextUpgrade(tail); ok; next, ok = g.nextUpgrade(tail) {
		tail = next
	}
	g.State.Cards[tail].NextUpgradePlus = upgradePlus(up)
}

// detachUpgrade unlinks an upgrade from its host's chain, stitching the upgrades on
// either side of it back together so the rest of the chain stays reachable, and
// clears the upgrade's own links. It returns the host, or ok=false when the id was
// not attached. This is what keeps the other upgrades on their host when an upgrade
// in the middle of the chain leaves play (e.g. one of several is destroyed).
func (g *Game) detachUpgrade(up LocalID) (LocalID, bool) {
	upCore := &g.State.Cards[up]
	host, ok := decodeUpgrade(upCore.HostPlus)
	if !ok {
		return 0, false
	}
	next := upCore.NextUpgradePlus
	if head, _ := g.firstUpgrade(host); head == up {
		g.State.Cards[host].FirstUpgradePlus = next
	} else {
		prev := head
		for g.State.Cards[prev].NextUpgradePlus != upgradePlus(up) {
			prev, _ = g.nextUpgrade(prev)
		}
		g.State.Cards[prev].NextUpgradePlus = next
	}
	upCore.HostPlus = 0
	upCore.NextUpgradePlus = 0
	return host, true
}

// hostOf returns the creature an attached upgrade is on, or ok=false when the id is
// not attached to any creature. The back-link makes this O(1). An id outside the
// card table is no card at all, so it is unattached rather than a panic — inPlay
// asks this of ids it has not yet established are real.
func (g *Game) hostOf(upgrade LocalID) (LocalID, bool) {
	if int(upgrade) >= len(g.State.Cards) {
		return 0, false
	}
	return decodeUpgrade(g.State.Cards[upgrade].HostPlus)
}

// HostOf is the exported form of hostOf, letting an effect (a blaster's payoff
// targeting the instance it bound to) read an upgrade's current host.
func (g *Game) HostOf(upgrade LocalID) (LocalID, bool) {
	return g.hostOf(upgrade)
}

// MoveUpgrade relocates an already-attached upgrade from its current host onto
// newHost, keeping the upgrade in play. It is how a movable upgrade — a Star
// Alliance "blaster" — homes to its signature creature. It does nothing when the
// upgrade is not currently attached to any creature, so a stray call cannot graft
// an unattached card onto the battleline.
func (g *Game) MoveUpgrade(upgrade, newHost LocalID) {
	if _, ok := g.detachUpgrade(upgrade); !ok {
		return
	}
	g.AttachUpgrade(newHost, upgrade)
	g.record(UpgradeAttached{
		Player:  g.owner(upgrade),
		Upgrade: upgrade,
		Host:    newHost,
	})
}
