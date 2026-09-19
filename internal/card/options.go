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
	// WithSplashAttack gives a creature Splash-attack N — when it fights it also deals N damage to each neighbor of the creature it fights.
	WithSplashAttack = func(n int) Option { return gameplay(engine.WithSplashAttack(n)) }
	// WithAttackDamage overrides how much fight damage a creature deals.
	WithAttackDamage = func(ad engine.AttackDamage) Option { return gameplay(engine.WithAttackDamage(ad)) }
	// WithAttackKeywords makes a creature gain keywords while it is attacking (Spyyyder gains poison against a flank creature).
	WithAttackKeywords = func(ak engine.AttackKeywords) Option { return gameplay(engine.WithAttackKeywords(ak)) }
	// WithNoDamageWhenAttacked makes a creature deal no retaliation damage when attacked.
	WithNoDamageWhenAttacked = func() Option { return gameplay(engine.WithNoDamageWhenAttacked()) }
	// WithStealsInsteadOfDamageWhenAttacked makes a creature's controller steal instead of it dealing retaliation damage (Shoulder Id).
	WithStealsInsteadOfDamageWhenAttacked = func(n int) Option {
		return gameplay(engine.WithStealsInsteadOfDamageWhenAttacked(n))
	}
	// WithFriendlyEntersPlayReady makes friendly cards enter play ready while this card is in play, per the grant (Duskwitch, The Curator, Fandangle).
	WithFriendlyEntersPlayReady = func(g engine.EntersReadyGrant) Option { return gameplay(engine.WithFriendlyEntersPlayReady(g)) }
	// WithFightRestriction restricts which creatures this creature may fight.
	WithFightRestriction = func(t engine.Target) Option { return gameplay(engine.WithFightRestriction(t)) }
	// WithCannotBeUsedTo bars a card from named ways of being used (reap, fight, action).
	WithCannotBeUsedTo = func(k ...engine.UseKind) Option {
		return gameplay(engine.WithCannotBeUsedTo(k...))
	}
	// WithCannotBeUsedWhile bars a card from being used at all while a condition holds.
	WithCannotBeUsedWhile = func(c Condition) Option {
		return gameplay(engine.WithCannotBeUsedWhile(c))
	}
	// WithDestroyedWhen destroys a creature for as long as a board condition holds.
	WithDestroyedWhen = func(c Condition) Option { return gameplay(engine.WithDestroyedWhen(c)) }
	// WithTakesDamageFor makes this card take the damage dealt to other creatures.
	WithTakesDamageFor = func(t engine.Target) Option { return gameplay(engine.WithTakesDamageFor(t)) }
	// WithAlsoTakesNeighborFightDamage makes this creature take an equal share of the
	// damage dealt to its neighbors during a fight, on top of the neighbor's own
	// damage (Drecker).
	WithAlsoTakesNeighborFightDamage = func() Option {
		return gameplay(engine.WithAlsoTakesNeighborFightDamage())
	}
	// WithTauntReachingNeighborsNeighbors extends this creature's taunt one step
	// further, so it shields its neighbors' neighbors as well as its neighbors (Lady
	// Loreena).
	WithTauntReachingNeighborsNeighbors = func() Option {
		return gameplay(engine.WithTauntReachingNeighborsNeighbors())
	}
	// WithPowerX gives a creature a variable "X" power: a live Count added to its
	// base power while its text is not blanked (Picaroon's combined-neighbor power).
	WithPowerX = func(c engine.Count) Option { return gameplay(engine.WithPowerX(c)) }
	// WithCannotBeDealtDamageBy makes the card refuse damage dealt to it by the
	// creatures the matcher names (Ardent Hero refuses Mutant creatures or creatures
	// with power 5 or higher).
	WithCannotBeDealtDamageBy = func(m DamageSource) Option {
		return gameplay(engine.WithCannotBeDealtDamageBy(m))
	}
	// WithTriggersFromDiscard keeps a card's triggered abilities live while it sits
	// in its owner's discard pile, so an "after you choose <house>" ability fires
	// from the discard (Relentless Creeper returns itself to hand).
	WithTriggersFromDiscard = func() Option { return gameplay(engine.WithTriggersFromDiscard()) }
	// WithAttackIgnores makes a creature ignore defensive keywords while attacking.
	WithAttackIgnores = func(kws ...engine.Keyword) Option { return gameplay(engine.WithAttackIgnores(kws...)) }
	// WithEntersPlay adds an effect that resolves as the creature enters play.
	WithEntersPlay = func(e Effect) Option { return gameplay(engine.WithEntersPlay(e)) }
	// WithBonus sets the bonus icons printed on a card, in top-to-bottom order —
	// card.WithBonus(card.Bonus.Aember, card.Bonus.Aember, card.Bonus.Draw).
	WithBonus = func(icons ...BonusIcon) Option { return gameplay(engine.WithBonus(icons...)) }
	// WithEnhance makes the card an Enhance source contributing the given bonus icons
	// to the deck at generation time; they have no effect on the card itself.
	WithEnhance = func(icons ...BonusIcon) Option { return gameplay(engine.WithEnhance(icons...)) }
	// WithoutEnhancement bars the given bonus-icon kinds from landing on this card
	// via Enhance (a Vactrol divergence, for a bonus that would only weaken it).
	WithoutEnhancement = func(icons ...BonusIcon) Option { return gameplay(engine.WithoutEnhancement(icons...)) }
	// WithStatic adds a static modifier (an upgrade's granted stats and abilities).
	WithStatic = func(m StaticModifier) Option { return gameplay(engine.WithStatic(m)) }
	// WithPlayableAsUpgrade lets a creature be played as an upgrade instead of a
	// creature, granting its host the card's WithStatic modifier.
	WithPlayableAsUpgrade = func() Option { return gameplay(engine.WithPlayableAsUpgrade()) }
	// WithConstant adds an ability that applies to the board while the card is in play.
	WithConstant = func(c ConstantAbility) Option { return gameplay(engine.WithConstantAbility(c)) }
	// WithRestrictions adds constant restrictions (cannot fight, cannot reap, ...).
	WithRestrictions = func(r Restrictions) Option { return gameplay(engine.WithRestrictions(r)) }
	// WithCannotPlayWhile adds a symmetric play bar: any player who meets the condition cannot play that type.
	WithCannotPlayWhile = func(b ConditionalPlayBar) Option { return gameplay(engine.WithCannotPlayWhile(b)) }
	// WithHouseLock constrains a player's active-house choice while this card is in play.
	WithHouseLock = func(l HouseLock) Option { return gameplay(engine.WithHouseLock(l)) }
	// WithKeyCost adds a change to the cost of forging a key.
	WithKeyCost = func(kc engine.KeyCostChange) Option { return gameplay(engine.WithKeyCost(kc)) }
	// WithPlayPermission sets the conditions under which the card may be played.
	WithPlayPermission = func(p engine.PlayPermission) Option { return gameplay(engine.WithPlayPermission(p)) }
	// WithReplaces adds a replacement effect (Instead) the card applies while in play.
	WithReplaces = func(r Instead) Option { return gameplay(engine.WithReplaces(r)) }
	// WithBonusInstead lets the card substitute one of its controller's bonus icons
	// while in play — resolving an icon as a different icon, or as an effect, instead.
	WithBonusInstead = func(r BonusInstead) Option { return gameplay(engine.WithBonusInstead(r)) }
	// WithDrawModifier changes how many cards a player draws.
	WithDrawModifier = func(p Player, amount int) Option { return gameplay(engine.WithDrawModifier(p, amount)) }
	// WithDrawModifierOffFlank changes how many cards a player draws, but only while
	// the source card is not on a flank (Streke).
	WithDrawModifierOffFlank = func(p Player, amount int) Option {
		return gameplay(engine.WithDrawModifierOffFlank(p, amount))
	}
	// WithDrawModifierInCenter changes how many cards a player draws, but only while
	// the source sits in the center of its battleline (Zenzizenzizenzic).
	WithDrawModifierInCenter = func(p Player, amount int) Option {
		return gameplay(engine.WithDrawModifierInCenter(p, amount))
	}
	// WithDrawModifierPer changes how many cards a player draws, scaled by a running
	// count (Greed refills 1 extra card for each friendly Sin creature).
	WithDrawModifierPer = func(p Player, amount int, per engine.Count) Option {
		return gameplay(engine.WithDrawModifierPer(p, amount, per))
	}
	// WithAemberCannotBeStolen keeps the controller's Æmber from being stolen — with
	// no argument unconditionally, or only while the given condition holds
	// (card.HasAember{Subject: card.Subject.This} while the card has Æmber, a card.PoolAember threshold
	// while the pool is deep enough).
	WithAemberCannotBeStolen = func(cond ...engine.Condition) Option {
		return gameplay(engine.WithAemberCannotBeStolen(cond...))
	}
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
	// WithAbility adds an ability that resolves an effect on a trigger. A composite
	// trigger (Trigger.PlayFightReap, Trigger.FightReap, Trigger.PlayReap,
	// Trigger.PlayFight) fans out into its atomic Play/Fight/Reap abilities, which
	// text rendering merges back into one "Play/Fight/Reap:" line.
	WithAbility = func(t engine.Trigger, e Effect) Option {
		return gameplay(func(d *engine.CardDefinition) {
			for _, at := range fanOutTrigger(t) {
				engine.WithAbility(at, e)(d)
			}
		})
	}
	// WithEachPlayerAbility adds an ability whose turn-scoped trigger — choosing a
	// house, the start of a turn, or the end of a turn — fires for either player's
	// turn or choice, not only its controller's, resolving as the player whose turn
	// or choice it was (Snag's Mirror, Gambling Den, Pincerator). It is WithAbility
	// with the ability's EachPlayer scope set.
	WithEachPlayerAbility = func(t engine.Trigger, e Effect) Option {
		return gameplay(engine.WithEachPlayerAbility(t, e))
	}
)

// fanOutTrigger expands a composite trigger into the atomic engine triggers it
// stands for, preserving the order the printed line reads; a plain trigger is
// returned unchanged.
func fanOutTrigger(t engine.Trigger) []engine.Trigger {
	switch t {
	case triggerPlayFightReap:
		return []engine.Trigger{
			engine.TriggerAfterPlay,
			engine.TriggerAfterFight,
			engine.TriggerAfterReap,
		}
	case triggerFightReap:
		return []engine.Trigger{engine.TriggerAfterReap, engine.TriggerAfterFight}
	case triggerPlayReap:
		return []engine.Trigger{engine.TriggerAfterPlay, engine.TriggerAfterReap}
	case triggerPlayFight:
		return []engine.Trigger{engine.TriggerAfterPlay, engine.TriggerAfterFight}
	default:
		return []engine.Trigger{t}
	}
}

// FightReap grants effect as both a Fight and a Reap ability, so a creature
// that gains these abilities resolves effect whenever it is used to fight or to
// reap; the pair prints as one "Fight/Reap:" line. It is the granted-ability
// analog of the Trigger.FightReap composite (which adds the pair to a card
// directly), for the Granted list of a StaticModifier or ConstantAbility — Rocket
// Boots grants its host "Fight/Reap: ready it".
func FightReap(e Effect) []Ability {
	return []Ability{
		{Trigger: Trigger.Reap, Effect: e},
		{Trigger: Trigger.Fight, Effect: e},
	}
}
