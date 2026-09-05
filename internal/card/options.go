package card

import "github.com/dmikalova/vactrol/internal/engine"

// Option helpers — thin wrappers over the engine's card options, each returning
// an authoring Option that appends its engine option to the builder. New (in
// register.go) applies them and enrolls the card.
var (
	// WithPower sets a creature's power.
	WithPower = func(p int) Option { return gameplay(engine.WithPower(p)) }
	// WithArmor sets a creature's armor.
	WithArmor = func(a int) Option { return gameplay(engine.WithArmor(a)) }
	// WithTraits sets a card's traits (e.g. card.Traits.Beast, card.Traits.Item).
	WithTraits = func(t ...Trait) Option { return gameplay(engine.WithTraits(t...)) }
	// WithKeywords gives a card keywords (Skirmish, Elusive, ...).
	WithKeywords = func(k ...engine.Keyword) Option { return gameplay(engine.WithKeywords(k...)) }
	// WithAssault gives a creature Assault N — it deals N damage to a creature it fights, first.
	WithAssault = func(n int) Option { return gameplay(engine.WithAssault(n)) }
	// WithHazardous gives a creature Hazardous N — it deals N damage to a creature that fights it, first.
	WithHazardous = func(n int) Option { return gameplay(engine.WithHazardous(n)) }
	// WithAttackDamage overrides how much fight damage a creature deals.
	WithAttackDamage = func(ad engine.AttackDamage) Option { return gameplay(engine.WithAttackDamage(ad)) }
	// WithFightRestriction restricts which creatures this creature may fight.
	WithFightRestriction = func(t engine.Target) Option { return gameplay(engine.WithFightRestriction(t)) }
	// WithCannotBeUsedTo bars a card from named ways of being used (reap, fight, action).
	WithCannotBeUsedTo = func(k ...engine.UseKind) Option {
		return gameplay(engine.WithCannotBeUsedTo(k...))
	}
	// WithDestroyedWhen destroys a creature for as long as a board condition holds.
	WithDestroyedWhen = func(c Condition) Option { return gameplay(engine.WithDestroyedWhen(c)) }
	// WithTakesDamageFor makes this card take the damage dealt to other creatures.
	WithTakesDamageFor = func(t engine.Target) Option { return gameplay(engine.WithTakesDamageFor(t)) }
	// WithAttackIgnores makes a creature ignore defensive keywords while attacking.
	WithAttackIgnores = func(kws ...engine.Keyword) Option { return gameplay(engine.WithAttackIgnores(kws...)) }
	// WithEntersPlay adds an effect that resolves as the creature enters play.
	WithEntersPlay = func(e Effect) Option { return gameplay(engine.WithEntersPlay(e)) }
	// WithAemberBonus sets a card's Æmber bonus (the pips gained when it is played).
	WithAemberBonus = func(n int) Option { return gameplay(engine.WithAemberBonus(n)) }
	// WithStatic adds a static modifier (an upgrade's granted stats and abilities).
	WithStatic = func(m StaticModifier) Option { return gameplay(engine.WithStatic(m)) }
	// WithConstant adds an ability that applies to the board while the card is in play.
	WithConstant = func(c ConstantAbility) Option { return gameplay(engine.WithConstantAbility(c)) }
	// WithRestrictions adds constant restrictions (cannot fight, cannot reap, ...).
	WithRestrictions = func(r Restrictions) Option { return gameplay(engine.WithRestrictions(r)) }
	// WithHouseLock constrains a player's active-house choice while this card is in play.
	WithHouseLock = func(l HouseLock) Option { return gameplay(engine.WithHouseLock(l)) }
	// WithKeyCost adds a change to the cost of forging a key.
	WithKeyCost = func(kc engine.KeyCostChange) Option { return gameplay(engine.WithKeyCost(kc)) }
	// WithPlayPermission sets the conditions under which the card may be played.
	WithPlayPermission = func(p engine.PlayPermission) Option { return gameplay(engine.WithPlayPermission(p)) }
	// WithReplaces adds a replacement effect (Instead) the card applies while in play.
	WithReplaces = func(r Instead) Option { return gameplay(engine.WithReplaces(r)) }
	// WithDrawModifier changes how many cards a player draws.
	WithDrawModifier = func(p Player, amount int) Option { return gameplay(engine.WithDrawModifier(p, amount)) }
	// WithAemberTheftImmunity makes Æmber on this card immune to theft.
	WithAemberTheftImmunity = func() Option { return gameplay(engine.WithAemberTheftImmunity()) }
	// WithSpendableAember lets Æmber banked on this card be spent when forging.
	WithSpendableAember = func() Option { return gameplay(engine.WithSpendableAember()) }
	// WithGainsForgeAember gives this card's controller all the Æmber their
	// opponent spends forging a key, for as long as it stays in play.
	WithGainsForgeAember = func() Option { return gameplay(engine.WithGainsForgeAember()) }
	// WithAemberThreshold requires a pool of at least n to play this card.
	WithAemberThreshold = func(n int) Option {
		return gameplay(engine.WithPlayRequirement(engine.AemberThreshold(n)))
	}
	// WithAemberCost requires — and spends — n Æmber to play this card.
	WithAemberCost = func(n int) Option {
		return gameplay(engine.WithPlayRequirement(engine.AemberCost(n)))
	}
	// WithAbility adds an ability that resolves an effect on a trigger.
	WithAbility = func(t engine.Trigger, e Effect) Option { return gameplay(engine.WithAbility(t, e)) }
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
