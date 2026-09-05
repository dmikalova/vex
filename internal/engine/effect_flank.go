package engine

// ConsiderFlank makes each creature its Target selects count as a flank creature
// for the remainder of the turn, regardless of where it actually sits in its
// battleline (Spectral Tunneler). A flank is normally a battleline position, so
// this is a lasting override the ready phase lifts, not a move; every flank check
// — combat's FlankOnly damage bonus and the OnFlank/NotOnFlank target filters —
// honors it.
type ConsiderFlank struct {
	Target Target
}

// validate requires an explicit target.
func (e ConsiderFlank) validate() error {
	if !e.Target.valid() {
		return errUnsetTarget("ConsiderFlank")
	}
	return nil
}

// Text renders the effect, e.g. "for the remainder of the turn, it is considered a
// flank creature".
func (e ConsiderFlank) Text() string {
	return "for the remainder of the turn, " + e.Target.Text() + " is considered a flank creature"
}

// Resolve makes each selected creature count as a flank creature for the turn.
func (e ConsiderFlank) Resolve(ctx *EffectContext) {
	for _, id := range e.Target.Select(ctx) {
		ctx.Resolver.ConsiderFlank(id)
	}
}
