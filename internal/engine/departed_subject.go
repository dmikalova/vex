package engine

// departedSubject holds the mutable scalars a card carried the instant before it
// left play, captured so an ability that is still resolving reads its last-known
// values rather than the zeroed core it leaves behind (resetCore blanks CardCore).
// Only the scalars a resolving reader can ask for after the subject is gone are
// kept — power, Æmber-on-card, and damage; a departed subject's position and
// controller have their own snapshots (Produced.Neighbors and ctx.ItController).
// See ADR 0030.
type departedSubject struct {
	power  int
	amber  int
	damage int
}

// captureDepartingSubject records id's mutable scalars the instant before it leaves
// play, so a following effect in the same resolution that keeps referencing it (as
// It or Source) reads its last-known power, Æmber-on-card, and damage instead of
// its zeroed core. Call it just before the effect removes the card — Destroy
// captures every creature it is about to destroy, DealDamage{IfDestroyed} captures
// the creature it is about to deal lethal damage to.
func captureDepartingSubject(ctx *EffectContext, id LocalID) {
	if ctx.Departed == nil {
		ctx.Departed = map[LocalID]departedSubject{}
	}
	ctx.Departed[id] = departedSubject{
		power:  ctx.Resolver.Power(id),
		amber:  ctx.Resolver.AmberOn(id),
		damage: ctx.Resolver.Damage(id),
	}
}

// powerOf returns id's live power, or its last-known power when it has left play
// mid-resolution and was captured departing.
func (ctx *EffectContext) powerOf(id LocalID) int {
	if d, ok := ctx.departed(id); ok {
		return d.power
	}
	return ctx.Resolver.Power(id)
}

// amberOn returns the Æmber on id, or the last-known amount when it has left play
// mid-resolution and was captured departing.
func (ctx *EffectContext) amberOn(id LocalID) int {
	if d, ok := ctx.departed(id); ok {
		return d.amber
	}
	return ctx.Resolver.AmberOn(id)
}

// damageOn returns the damage on id, or the last-known amount when it has left play
// mid-resolution and was captured departing.
func (ctx *EffectContext) damageOn(id LocalID) int {
	if d, ok := ctx.departed(id); ok {
		return d.damage
	}
	return ctx.Resolver.Damage(id)
}

// departed returns id's captured last-known scalars, but only once id has actually
// left play — while it is still in play its live values remain authoritative.
func (ctx *EffectContext) departed(id LocalID) (departedSubject, bool) {
	d, ok := ctx.Departed[id]
	if !ok || resolverInPlay(ctx, id) {
		return departedSubject{}, false
	}
	return d, true
}
