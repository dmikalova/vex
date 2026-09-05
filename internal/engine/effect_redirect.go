package engine

// RedirectFightDamage is a "Before Fight" effect: the controller chooses a
// creature (its Target), and the attacker deals its own fight damage to that
// creature instead of to the creature it is fighting (Gabos Longarms). It only
// redirects the attacker's outgoing fight damage — the attacker still takes
// damage back from the creature it fights. The chosen creature is stored on the
// game state for the fight in progress; the combat step reads and clears it.
type RedirectFightDamage struct {
	Target Target
}

// validate requires an explicit target.
func (e RedirectFightDamage) validate() error {
	if !e.Target.valid() {
		return errUnsetTarget("RedirectFightDamage")
	}
	return nil
}

// Text renders the effect, e.g. "choose a creature - Gabos Longarms deals its
// fight damage to the chosen creature instead of to the creature it is fighting".
func (e RedirectFightDamage) Text() string {
	return "choose " + e.Target.Text() + " - " + SelfName + " deals its fight damage to the chosen creature instead of to the creature it is fighting"
}

// Resolve records the chosen creature as the target of the attacker's fight
// damage for the fight in progress.
func (e RedirectFightDamage) Resolve(ctx *EffectContext) {
	ids := e.Target.Select(ctx)
	if len(ids) == 0 {
		return
	}
	ctx.Resolver.SetFightDamageRedirect(ids[0])
}
