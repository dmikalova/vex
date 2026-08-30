package card

import "github.com/dmikalova/vactrol/internal/engine"

// Option helpers — thin wrappers over the engine's card options, each returning
// an authoring Option that appends its engine option to the builder. New (in
// register.go) applies them and enrolls the card.
var (
	WithPower            = func(p int) Option { return gameplay(engine.WithPower(p)) }
	WithArmor            = func(a int) Option { return gameplay(engine.WithArmor(a)) }
	WithTraits           = func(t ...Trait) Option { return gameplay(engine.WithTraits(t...)) }
	WithKeywords         = func(k ...engine.Keyword) Option { return gameplay(engine.WithKeywords(k...)) }
	WithAssault          = func(n int) Option { return gameplay(engine.WithAssault(n)) }
	WithHazardous        = func(n int) Option { return gameplay(engine.WithHazardous(n)) }
	WithAttackDamage     = func(ad engine.AttackDamage) Option { return gameplay(engine.WithAttackDamage(ad)) }
	WithFightRestriction = func(t engine.Target) Option { return gameplay(engine.WithFightRestriction(t)) }
	WithEntersPlay       = func(e Effect) Option { return gameplay(engine.WithEntersPlay(e)) }
	WithAemberBonus      = func(n int) Option { return gameplay(engine.WithAemberBonus(n)) }
	WithStatic           = func(m StaticModifier) Option { return gameplay(engine.WithStatic(m)) }
	WithConstantAbility  = func(c ConstantAbility) Option { return gameplay(engine.WithConstantAbility(c)) }
	WithRestrictions     = func(r Restrictions) Option { return gameplay(engine.WithRestrictions(r)) }
	WithKeyCost          = func(kc engine.KeyCostChange) Option { return gameplay(engine.WithKeyCost(kc)) }
	WithPlayPermission   = func(p engine.PlayPermission) Option { return gameplay(engine.WithPlayPermission(p)) }
	WithReplaces         = func(r Instead) Option { return gameplay(engine.WithReplaces(r)) }
	WithAbility          = func(t engine.Trigger, e Effect) Option { return gameplay(engine.WithAbility(t, e)) }
	// WithFightOrReap adds effect as both a Fight and a Reap ability, so it resolves
	// whenever the creature is used to fight or to reap; the two print as one
	// "Fight/Reap:" line.
	WithFightOrReap = func(e Effect) Option {
		return gameplay(func(d *engine.CardDefinition) {
			engine.WithAbility(engine.TriggerAfterReap, e)(d)
			engine.WithAbility(engine.TriggerAfterFight, e)(d)
		})
	}
	// WithPlayReap adds effect as both a Play and a Reap ability; the two print as
	// one "Play/Reap:" line.
	WithPlayReap = func(e Effect) Option {
		return gameplay(func(d *engine.CardDefinition) {
			engine.WithAbility(engine.TriggerAfterPlay, e)(d)
			engine.WithAbility(engine.TriggerAfterReap, e)(d)
		})
	}
	// WithPlayFightReap adds effect as a Play, a Fight, and a Reap ability; the
	// three print as one "Play/Fight/Reap:" line.
	WithPlayFightReap = func(e Effect) Option {
		return gameplay(func(d *engine.CardDefinition) {
			engine.WithAbility(engine.TriggerAfterPlay, e)(d)
			engine.WithAbility(engine.TriggerAfterFight, e)(d)
			engine.WithAbility(engine.TriggerAfterReap, e)(d)
		})
	}
)

// FightOrReap grants effect as both a Fight and a Reap ability, so a creature
// that gains these abilities resolves effect whenever it is used to fight or to
// reap; the pair prints as one "Fight/Reap:" line. It is the granted-ability
// analog of WithFightOrReap (which adds the pair to a card directly), for the
// Granted list of a StaticModifier or ConstantAbility — Rocket Boots grants its
// host "Fight/Reap: ready it".
func FightOrReap(e Effect) []Ability {
	return []Ability{
		{Trigger: Trigger.Reap, Effect: e},
		{Trigger: Trigger.Fight, Effect: e},
	}
}
