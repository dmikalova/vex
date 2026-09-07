package engine

// Ward applies the ward status. It is a simple "verb the target" effect, so a
// ward that runs beside another status change on the same target folds into one
// phrase in a Sequence (see combinable).

// A ward is a one-shot shield placed on a creature. The first time a warded
// creature would be dealt damage or would leave play, that damage or removal is
// absorbed and the ward is spent instead. Ward intercepts every removal, even the
// controller's own; it covers only damage and leaving play. Warding applies this
// status to each creature the effect targets.
type Ward struct {
	Target Target
}

// validate requires an explicit target.
func (e Ward) validate() error {
	if !e.Target.valid() {
		return errUnsetTarget("Ward")
	}
	return nil
}

func (e Ward) verb() string       { return "ward" }
func (e Ward) targetText() string { return e.Target.Text() }

// Text renders the effect, e.g. "ward each friendly creature".
func (e Ward) Text() string { return e.verb() + " " + e.targetText() }

// Resolve wards each selected creature. A creature already warded still gets a log
// line — the source still had to choose it — just without a state change.
func (e Ward) Resolve(ctx *EffectContext) {
	for _, id := range e.Target.Select(ctx) {
		if ctx.Resolver.Warded(id) {
			ctx.Resolver.Record(CreatureWarded{Creature: id, By: ctx.Source, AlreadyWarded: true})
			continue
		}
		ctx.Resolver.SetWarded(id, true)
		ctx.Resolver.Record(CreatureWarded{Creature: id, By: ctx.Source})
	}
}
