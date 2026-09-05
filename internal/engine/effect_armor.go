package engine

// This file holds the effects that act on armor — the points that stop damage
// before it lands. Armor is spent absorbing damage during a turn and refreshes
// when its controller readies; an effect here takes it away outright instead.

// LoseArmor takes all the remaining armor off each creature its Target selects,
// and tallies what it took so a following effect can scale with it (Red-Hot Armor
// strips armor, then deals damage for each point stripped). The armor comes back
// when its controller readies, the same way armor spent absorbing damage does.
type LoseArmor struct {
	Target Target
}

// validate requires an explicit target.
func (e LoseArmor) validate() error {
	if !e.Target.valid() {
		return errUnsetTarget("LoseArmor")
	}
	return nil
}

// Text renders the effect, e.g. "each enemy creature with armor loses all of its
// armor".
func (e LoseArmor) Text() string {
	return e.Target.Text() + " loses all of its armor"
}

// Resolve empties each selected creature's armor.
func (e LoseArmor) Resolve(ctx *EffectContext) {
	for _, id := range e.Target.Select(ctx) {
		ctx.Resolver.StripArmor(id)
	}
}

// ArmorLostThisWay scales a per-target amount by how much armor an effect has
// taken off that creature this turn — the "for each point of armor it lost this
// way" clause. It reads a strip, not armor spent absorbing damage, so it measures
// only what an effect took.
var ArmorLostThisWay PerTarget = armorLostThisWay{}

type armorLostThisWay struct{}

// perTargetValue reads how much armor was stripped from this target.
func (armorLostThisWay) perTargetValue(ctx *EffectContext, id LocalID) int {
	return ctx.Resolver.ArmorStripped(id)
}

// perTargetText renders the singular unit the "for each" clause repeats.
func (armorLostThisWay) perTargetText() string { return "point of armor it lost this way" }
