package engine

import "fmt"

// GainStats gives each creature its Target selects power and/or armor for the
// remainder of the turn — Abond the Armorsmith's Action grants other friendly
// creatures +1 armor until end of turn. The ready phase clears the bonus. A
// constant "+N armor" that lasts while a card stays in play is a ConstantAbility
// instead; this node is the one-shot, end-of-turn grant.
type GainStats struct {
	Target Target
	Power  int
	Armor  int
}

// validate requires an explicit target and at least one stat to grant.
func (e GainStats) validate() error {
	if !e.Target.valid() {
		return errUnsetTarget("GainStats")
	}
	if e.Power == 0 && e.Armor == 0 {
		return fmt.Errorf("GainStats: grants no power or armor")
	}
	return nil
}

// Text renders the effect, e.g. "for the remainder of the turn, each other
// friendly creature gains +1 armor".
func (e GainStats) Text() string {
	return fmt.Sprintf("for the remainder of the turn, %s gains %s",
		e.Target.Text(), staticBonuses(StaticModifier{PowerBonus: e.Power, ArmorBonus: e.Armor}))
}

// Resolve grants each selected creature the bonus for the remainder of the turn.
func (e GainStats) Resolve(ctx *EffectContext) {
	for _, id := range e.Target.Select(ctx) {
		ctx.Resolver.GainStats(id, e.Power, e.Armor)
	}
}
