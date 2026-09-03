package engine

// takeControl moves a creature into controller's battleline without changing its
// owner, recording source as the card whose lasting effect holds the control so it
// can be reverted when that source leaves play. Ownership is immutable in KeyForge
// and still decides the out-of-play zone the card returns to; control is the
// battleline/current-player relationship.
func (g *Game) takeControl(id LocalID, controller int, source LocalID) {
	if g.inPlay(id) && g.cat.def(id).Type == Creature {
		g.removeFromPlay(id)
		controlPlus := uint8(0)
		if controller != g.owner(id) {
			controlPlus = uint8(controller + 1)
		}
		g.State.Cards[id].ControlPlus = controlPlus
		g.State.Cards[id].ControlSource = source
		g.State.Battleline[controller].add(id)
		g.record(ControlTaken{Player: controller, Card: id})
		// Switching sides re-aims every constant ability that reads "friendly" or
		// "enemy", so the creature can land already dead: a damaged creature that was
		// only alive on the +power its old side gave it, or one whose new side is
		// under a power-reducing constant ability. removeFromPlay settled the board
		// it left, not this one.
		g.settleDestroyed(g.State.ActivePlayer)
	}
}

// takeControlOfArtifact moves an artifact into controller's artifact row without
// changing its owner, giving that player control of it for good. Unlike creature
// control there is no reverting source: releaseControlHeldBy only scans
// battlelines, so a controlled artifact is never handed back until it leaves
// play. Ownership stays fixed and still decides the zone it returns to.
func (g *Game) takeControlOfArtifact(id LocalID, controller int) {
	if !g.inPlay(id) || g.cat.def(id).Type != Artifact {
		return
	}
	g.removeFromPlay(id)
	controlPlus := uint8(0)
	if controller != g.owner(id) {
		controlPlus = uint8(controller + 1)
	}
	g.State.Cards[id].ControlPlus = controlPlus
	g.State.Cards[id].ControlSource = 0
	g.State.Artifacts[controller].add(id)
	g.record(ControlTaken{Player: controller, Card: id})
	// An artifact's constant abilities change sides with it, so creatures on its new
	// enemy side can drop to lethal the moment it arrives.
	g.settleDestroyed(g.State.ActivePlayer)
}

// releaseControlHeldBy reverts every creature whose control was taken "until source
// leaves play" back to its owner's battleline, called when source leaves play. It is
// the leave-play half of the UntilThisLeavesPlay duration: the control lasts exactly
// as long as its source card stays in play.
func (g *Game) releaseControlHeldBy(source LocalID) {
	for p := 0; p < 2; p++ {
		for _, id := range g.battlelineCopy(p) {
			core := &g.State.Cards[id]
			if core.ControlPlus == 0 || core.ControlSource != source {
				continue
			}
			owner := g.owner(id)
			g.removeFromPlay(id)
			core.ControlPlus = 0
			core.ControlSource = 0
			g.State.Battleline[owner].add(id)
			g.record(ControlReturned{Card: id, Owner: owner})
		}
	}
	// Handing creatures back swaps which side's constant abilities reach them, so
	// settle once the whole batch has moved.
	g.settleDestroyed(g.State.ActivePlayer)
}
