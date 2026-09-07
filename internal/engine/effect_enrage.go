package engine

// Enrage applies the enrage status. It is a simple "verb the target" effect, so
// an enrage that runs beside another status change on the same target folds into
// one phrase in a Sequence (see combinable).

// An enraged creature must be used to fight on its controller's turn if it is
// able to — it cannot reap or use an "Action:" ability while there is an enemy
// creature it can fight. Enrage persists until an effect removes it. Enraging
// applies this status to each creature the effect targets.
type Enrage struct {
	Target Target
}

// validate requires an explicit target.
func (e Enrage) validate() error {
	if !e.Target.valid() {
		return errUnsetTarget("Enrage")
	}
	return nil
}

func (e Enrage) verb() string       { return "enrage" }
func (e Enrage) targetText() string { return e.Target.Text() }

// Text renders the effect, e.g. "enrage each enemy creature".
func (e Enrage) Text() string { return e.verb() + " " + e.targetText() }

// Resolve enrages each selected creature. A creature already enraged still gets a
// log line — the source still had to choose it — just without a state change.
func (e Enrage) Resolve(ctx *EffectContext) {
	for _, id := range e.Target.Select(ctx) {
		if ctx.Resolver.Enraged(id) {
			ctx.Resolver.Record(CreatureEnraged{Creature: id, By: ctx.Source, AlreadyEnraged: true})
			continue
		}
		ctx.Resolver.SetEnraged(id, true)
		ctx.Resolver.Record(CreatureEnraged{Creature: id, By: ctx.Source})
	}
}
