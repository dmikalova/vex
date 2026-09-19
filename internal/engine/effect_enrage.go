package engine

// Enrage applies the enrage status. It is a simple "verb the target" effect, so
// an enrage that runs beside another status change on the same target folds into
// one phrase in a Sequence (see combinable).

// An enraged creature must be used to fight on its controller's turn if it is
// able to — it cannot reap or use an "Action:" ability while there is an enemy
// creature it can fight. Enrage is removed once the creature is used to fight
// (clearEnrageOnFight); otherwise it persists until an effect removes it.
// Enraging applies this status to each creature the effect targets.
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
			ctx.Resolver.Record(CreatureEnraged{
				Creature:       id,
				By:             ctx.Source,
				AlreadyEnraged: true,
			})
			continue
		}
		ctx.Resolver.SetEnraged(id, true)
		ctx.Resolver.Record(CreatureEnraged{
			Creature: id,
			By:       ctx.Source,
		})
	}
}

// clearEnrageOnFight removes the enrage a creature carried into a fight, once it
// has been used to fight. The fight counts even when Elusive turned its damage
// aside — the creature was still used to fight — so the enrage goes either way
// (KeyForge removes all enrage counters after a creature is used to fight). Only
// the attacker is used, so a defender keeps any enrage of its own. Runs before the
// post-fight window so a "Fight:" ability that re-enrages the attacker
// (Gladiodontus) is not undone. Pinned by TestEnrageClearedByFightingElusive and
// TestEnrageClearedBeforeFightAbilityReenrages.
func (g *Game) clearEnrageOnFight(attacker LocalID) {
	if !g.inPlay(attacker) || !g.Enraged(attacker) {
		return
	}
	g.State.Cards[attacker].Enraged = false
	g.record(CreatureEnrageRemoved{Creature: attacker})
}
