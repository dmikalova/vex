package engine

import "fmt"

// Healing takes damage tokens off a creature — a fixed amount, or all of them at
// once. It can never remove more damage than is on the creature (a creature with
// no damage is unaffected), and it never changes a creature's power.
type Heal struct {
	Amount int    // damage to remove; must be zero when Fully is set
	Fully  bool   // remove all damage instead of a fixed Amount
	Target Target // the creatures to heal
}

// Text renders the effect, e.g. "heal 1 damage from each creature" or
// "fully heal each other friendly creature".
func (e Heal) Text() string {
	if e.Fully {
		return "fully heal " + e.Target.Text()
	}
	return fmt.Sprintf("heal %d damage from %s", e.Amount, e.Target.Text())
}

// Resolve removes the damage from each selected creature (all of it when Fully),
// and records how many were actually healed on the context so a following effect
// can scale with it (see CreaturesHealed). A creature with no damage is unaffected
// and does not count as healed.
func (e Heal) Resolve(ctx *EffectContext) { e.resolveGate(ctx) }

// resolveGate heals the target and reports whether any creature was actually
// healed, so Heal can gate a Then — Protectrix protects a creature only if it
// healed damage. The last creature healed is left in context (ctx.It) so a
// following effect can act on "that creature".
//
// Undamaged creatures are dropped before the prompt, not after it: healing one
// does nothing, so offering it is a vacuous choice (Guardian Demon asked for a
// creature to heal even with no damage anywhere on the board). Neighbours pulled
// in after the choice can still be undamaged, so the loop skips them too.
func (e Heal) resolveGate(ctx *EffectContext) bool {
	return e.heal(ctx, e.Target.selectWith(ctx, false, e.damagedOnly(ctx)))
}

// declinable reports that the healing is a single clickable creature.
func (e Heal) declinable() bool { return e.Target.isChosen() }

// resolveOptional is resolveGate under a May: the creature to heal is asked
// declinably, with a Done to decline, instead of a separate Yes/No before the
// pick (Protectrix).
func (e Heal) resolveOptional(ctx *EffectContext) bool {
	return e.heal(ctx, e.Target.selectWith(ctx, true, e.damagedOnly(ctx)))
}

// damagedOnly drops undamaged creatures from a Heal's candidates: healing one
// does nothing, so offering it is a vacuous choice.
func (e Heal) damagedOnly(ctx *EffectContext) func(LocalID) bool {
	return func(id LocalID) bool { return ctx.Resolver.Damage(id) > 0 }
}

// heal removes damage from an already-selected set of creatures (all of it when
// Fully), and records how many were actually healed.
func (e Heal) heal(ctx *EffectContext, ids []LocalID) bool {
	healed := 0
	damageHealed := 0
	for _, id := range ids {
		before := ctx.Resolver.Damage(id)
		if before == 0 {
			continue
		}
		healed++
		ctx.It, ctx.HasIt = id, true
		removed := before
		if !e.Fully && e.Amount < before {
			removed = e.Amount
		}
		ctx.Resolver.SetDamage(id, before-removed)
		damageHealed += removed
	}
	ctx.Produced.Healed = healed
	ctx.Produced.DamageHealed = damageHealed
	return healed > 0
}

// validate rejects a Heal that sets both a fixed Amount and Fully, since the two
// are different ways to say how much to heal and combining them is ambiguous.
func (e Heal) validate() error {
	if !e.Target.valid() {
		return errUnsetTarget("Heal")
	}
	return errAmountOr("Heal", "Fully", e.Amount, e.Fully)
}

// CreaturesHealed counts the creatures the most recent Heal actually healed — the
// "for each creature healed this way" clause. Heal records the tally on the
// context, so pairing it after a Heal in a Sequence lets any effect (gain Æmber,
// deal damage, ...) scale with the number healed without a bespoke combined
// effect.
type CreaturesHealed struct{}

// Value returns how many creatures the preceding Heal healed.
func (CreaturesHealed) Value(ctx *EffectContext) int { return ctx.Produced.Healed }

// CountText renders the singular noun the "for each" clause repeats.
func (CreaturesHealed) CountText() string { return "creature healed this way" }

// DamageHealed is the amount of damage the most recent Heal removed — the "that
// amount of damage" a following DealDamage.AmountFrom deals on (Guardian Demon).
type DamageHealed struct{}

// Value returns how much damage the preceding Heal removed.
func (DamageHealed) Value(ctx *EffectContext) int { return ctx.Produced.DamageHealed }

// CountText renders the singular noun a "for each" clause would repeat.
func (DamageHealed) CountText() string { return "damage healed this way" }

// CountClause renders the clause CountIs puts after "if", e.g. "you healed
// exactly 3 damage" — Vigor's bonus turns on having healed the full amount.
// Damage is a mass noun, so the plural flag does not change it.
func (DamageHealed) CountClause(quantity string, _ bool) string {
	return "you healed " + quantity + " damage"
}
