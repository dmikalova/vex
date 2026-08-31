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
		g.logf("%s takes control of %s", g.names[controller], g.Name(id))
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
	g.logf("%s takes control of %s", g.names[controller], g.Name(id))
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
			g.logf("%s returns to %s's control", g.Name(id), g.names[owner])
		}
	}
}
