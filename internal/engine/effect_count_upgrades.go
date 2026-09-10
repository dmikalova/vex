package engine

// UpgradesOn counts the upgrades attached to the creature its Target names, for
// an effect whose magnitude scales with how many upgrades a creature carries —
// Walls' Blaster stuns a creature for each upgrade on Chief Engineer Walls after
// it homes there, itself now among them. It selects the creature once, so a
// Target that names nothing in play counts zero.
type UpgradesOn struct {
	Target Target
}

// Value returns the number of upgrades attached to the creature the Target names,
// or 0 when it selects nothing.
func (e UpgradesOn) Value(ctx *EffectContext) int {
	ids := e.Target.Select(ctx)
	if len(ids) == 0 {
		return 0
	}
	return len(ctx.Resolver.Upgrades(ids[0]))
}

// CountText renders the singular noun the "for each" clause repeats, e.g.
// "upgrade on Chief Engineer Walls".
func (e UpgradesOn) CountText() string { return "upgrade on " + e.Target.Text() }
