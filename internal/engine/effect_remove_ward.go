package engine

// RemoveWard takes the ward off a creature. Its target is any creature, warded or
// not — the KeyForge "remove a ward from a creature" — so removing a ward from a
// creature that carries none does nothing. It is the counterpart to Ward, and the
// two pair in a Sequence for Hunter or Hunted?.
type RemoveWard struct {
	Target Target
}

// validate requires an explicit target.
func (e RemoveWard) validate() error {
	if !e.Target.valid() {
		return errUnsetTarget("RemoveWard")
	}
	return nil
}

// Text renders the effect, e.g. "remove a ward from a creature".
func (e RemoveWard) Text() string {
	return "remove a ward from " + e.Target.Text()
}

// Resolve takes the ward off each selected creature. A creature that carries no
// ward is left unchanged; the choice was still made, so it is worth a log line.
func (e RemoveWard) Resolve(ctx *EffectContext) {
	for _, id := range e.Target.Select(ctx) {
		if !ctx.Resolver.Warded(id) {
			ctx.Resolver.Record(WardRemoved{
				Creature:        id,
				By:              ctx.Source,
				AlreadyUnwarded: true,
			})
			continue
		}
		ctx.Resolver.SetWarded(id, false)
		ctx.Resolver.Record(WardRemoved{
			Creature: id,
			By:       ctx.Source,
		})
	}
}
