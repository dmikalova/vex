package engine

import "fmt"

// Dealing damage puts that much pending damage on each creature the effect
// targets. Armor prevents pending damage first — each point stops 1, and armor
// spent this way stays spent for the rest of the turn — and whatever is not
// prevented lands as damage tokens. A creature whose total damage reaches or
// exceeds its power is destroyed. When one ability deals damage to several
// creatures they are damaged simultaneously and any that died are destroyed
// together, so no creature's destruction changes another's.
//
//rulebook:effect Deal Damage
type DealDamage struct {
	Amount int
	Per    Count
	Target Target
	// IgnoreArmor makes the damage bypass armor (Qyxxlyx Plague Master's "this
	// damage cannot be prevented by armor").
	IgnoreArmor bool
}

// validate requires an explicit target.
func (e DealDamage) validate() error {
	if !e.Target.valid() {
		return errUnsetTarget("DealDamage")
	}
	return nil
}

// Text renders the effect, e.g. "deal 2 damage to each enemy creature". A "for
// each" count leads the sentence (rule 9), e.g. "for each friendly creature in
// play, deal 1 damage to a creature". Armor-ignoring damage adds a trailing clause.
func (e DealDamage) Text() string {
	body := fmt.Sprintf("deal %d damage to %s", e.Amount, e.Target.Text())
	if e.IgnoreArmor {
		body += ", ignoring armor"
	}
	return forEach(e.Per, body)
}

// Resolve deals the damage to every selected creature simultaneously, resolving
// destruction as part of it. A Per count multiplies the amount dealt.
func (e DealDamage) Resolve(ctx *EffectContext) {
	amount := e.Amount
	if e.Per != nil {
		amount *= e.Per.Value(ctx)
	}
	ids := e.Target.Select(ctx)
	targets := make([]DamageTarget, len(ids))
	for i, id := range ids {
		targets[i] = DamageTarget{ID: id, Amount: amount, IgnoreArmor: e.IgnoreArmor}
	}
	ctx.Resolver.DealDamage(ctx.Controller, targets)
}

// DamageIfDestroyed deals damage to one chosen creature and, only if that damage
// destroys it, resolves a follow-up effect with the destroyed creature in context
// (ctx.It) — Seeker Needle's "deal 1 damage to a creature. If this damage destroys
// that creature, gain 1 Æmber."
type DamageIfDestroyed struct {
	Amount int
	Target Target
	Then   Effect
}

// validate requires a target and a well-formed follow-up effect.
func (e DamageIfDestroyed) validate() error {
	if !e.Target.valid() {
		return errUnsetTarget("DamageIfDestroyed")
	}
	return validateEffect(e.Then)
}

// Text renders the effect, e.g. "deal 2 damage to a creature. If this damage
// destroys that creature, steal 1 Æmber".
func (e DamageIfDestroyed) Text() string {
	return fmt.Sprintf("deal %d damage to %s. If this damage destroys that creature, %s",
		e.Amount, e.Target.Text(), e.Then.Text())
}

// Resolve deals the damage to the chosen creature, then runs Then only if the
// creature has left play. The destroyed creature is placed in context (ctx.It) so
// Then can refer to it ("purge it").
func (e DamageIfDestroyed) Resolve(ctx *EffectContext) {
	ids := e.Target.Select(ctx)
	if len(ids) == 0 {
		return
	}
	id := ids[0]
	ctx.Resolver.DealDamage(ctx.Controller, []DamageTarget{{ID: id, Amount: e.Amount}})
	if !resolverInPlay(ctx, id) {
		ctx.It, ctx.HasIt = id, true
		e.Then.Resolve(ctx)
	}
}

// DamageIfSurvives deals damage to one chosen creature and, only if the creature
// is not destroyed, resolves a follow-up effect with the surviving creature in
// context (ctx.It) — Gongoozle's "deal 3 damage to a creature. If it is not
// destroyed, its owner discards a random card from their hand." This is a plain
// state branch, not a result gate.
type DamageIfSurvives struct {
	Amount int
	Target Target
	Then   Effect
}

// validate requires a target and a well-formed follow-up effect.
func (e DamageIfSurvives) validate() error {
	if !e.Target.valid() {
		return errUnsetTarget("DamageIfSurvives")
	}
	return validateEffect(e.Then)
}

// Text renders the effect, e.g. "deal 3 damage to a creature. If it is not
// destroyed, its owner discards a random card from their hand".
func (e DamageIfSurvives) Text() string {
	return fmt.Sprintf("deal %d damage to %s. If it is not destroyed, %s",
		e.Amount, e.Target.Text(), e.Then.Text())
}

// Resolve deals the damage to the chosen creature, then runs Then only if the
// creature is still in play. The surviving creature is placed in context (ctx.It)
// so Then can refer to it ("its owner").
func (e DamageIfSurvives) Resolve(ctx *EffectContext) {
	ids := e.Target.Select(ctx)
	if len(ids) == 0 {
		return
	}
	id := ids[0]
	ctx.Resolver.DealDamage(ctx.Controller, []DamageTarget{{ID: id, Amount: e.Amount}})
	if resolverInPlay(ctx, id) {
		ctx.It, ctx.HasIt = id, true
		e.Then.Resolve(ctx)
	}
}

// DamageCreatureAndNeighbor deals damage to a chosen creature and damage to one of
// its battleline neighbors (the controller picks when it has two, the sole one
// otherwise), all at once — Mighty Lance's "deal 3 damage to a creature and 3
// damage to a neighbor of that creature."
type DamageCreatureAndNeighbor struct {
	Amount         int
	NeighborAmount int
}

// Text renders the effect.
func (e DamageCreatureAndNeighbor) Text() string {
	return fmt.Sprintf("deal %d damage to a creature and %d damage to a neighbor of that creature",
		e.Amount, e.NeighborAmount)
}

// Resolve chooses a creature, then damages it and one chosen neighbor together.
func (e DamageCreatureAndNeighbor) Resolve(ctx *EffectContext) {
	chosen := (Target{Kind: TargetChosenCreature}).Select(ctx)
	if len(chosen) == 0 {
		return
	}
	targets := []DamageTarget{{ID: chosen[0], Amount: e.Amount}}
	if ns := neighbors(ctx, chosen[0]); len(ns) > 0 {
		n := ns[0]
		if len(ns) > 1 {
			if pick, ok := ctx.ChooseCreature("Choose a neighbor", ns); ok {
				n = pick
			}
		}
		targets = append(targets, DamageTarget{ID: n, Amount: e.NeighborAmount})
	}
	ctx.Resolver.DealDamage(ctx.Controller, targets)
}

// SplashDamage deals damage to one chosen creature that is not on a flank and a
// smaller "splash" amount to each of that creature's neighbors, all at once. The
// not-on-a-flank restriction guarantees the chosen creature has two neighbors.
type SplashDamage struct {
	Amount int
	Splash int
}

// Text renders the effect.
func (e SplashDamage) Text() string {
	return fmt.Sprintf("deal %d damage to a creature that is not on a flank and %d damage to each of its neighbors", e.Amount, e.Splash)
}

// Resolve chooses a non-flank creature, then damages it and its neighbors as one
// simultaneous batch.
func (e SplashDamage) Resolve(ctx *EffectContext) {
	chosen := Target{Kind: TargetChosenCreature}.NotOnFlank().Select(ctx)
	if len(chosen) == 0 {
		return
	}
	targets := []DamageTarget{{ID: chosen[0], Amount: e.Amount}}
	for _, n := range neighbors(ctx, chosen[0]) {
		targets = append(targets, DamageTarget{ID: n, Amount: e.Splash})
	}
	ctx.Resolver.DealDamage(ctx.Controller, targets)
}

// DamageDifferent deals damage to a chosen creature and damage to a second,
// different chosen creature, all at once — Twin Bolt Emission's "deal 2 damage
// to a creature and deal 2 damage to a different creature."
type DamageDifferent struct {
	First  int
	Second int
}

// Text renders the effect.
func (e DamageDifferent) Text() string {
	return fmt.Sprintf("deal %d damage to a creature and deal %d damage to a different creature",
		e.First, e.Second)
}

// Resolve chooses a creature and a different creature, damaging both together.
// When only one creature is in play, the second, different creature cannot be
// chosen and only the first is damaged.
func (e DamageDifferent) Resolve(ctx *EffectContext) {
	chosen := (Target{Kind: TargetChosenCreature}).Select(ctx)
	if len(chosen) == 0 {
		return
	}
	targets := []DamageTarget{{ID: chosen[0], Amount: e.First}}
	if others := creaturesExcept(ctx, chosen[0]); len(others) > 0 {
		if id, ok := ctx.ChooseCreature("Choose a different creature", others); ok {
			targets = append(targets, DamageTarget{ID: id, Amount: e.Second})
		}
	}
	ctx.Resolver.DealDamage(ctx.Controller, targets)
}

// FlankWalkDamage chooses a flank creature and deals decreasing damage inward
// along its battleline: Amounts[0] to the chosen flank creature, Amounts[1] to
// its neighbor, Amounts[2] to that neighbor's other neighbor, and so on. Positron
// Bolt is Amounts{3, 2, 1}. The walk stops at the far flank if the battleline is
// shorter than the list of amounts.
type FlankWalkDamage struct {
	Amounts []int
}

// flankWalkPhrase names the creature at step i of the inward walk.
func flankWalkPhrase(i int) string {
	switch i {
	case 0:
		return "it"
	case 1:
		return "its neighbor"
	default:
		return "the neighbor's other neighbor"
	}
}

// Text renders the effect, e.g. "choose a flank creature. Deal 3 damage to it, 2
// damage to its neighbor, and 1 damage to the neighbor's other neighbor."
func (e FlankWalkDamage) Text() string {
	parts := make([]string, len(e.Amounts))
	for i, a := range e.Amounts {
		parts[i] = fmt.Sprintf("%d damage to %s", a, flankWalkPhrase(i))
	}
	joined := parts[0]
	for i := 1; i < len(parts); i++ {
		sep := ", "
		if i == len(parts)-1 {
			sep = ", and "
		}
		joined += sep + parts[i]
	}
	return "choose a flank creature. Deal " + joined
}

// validate requires at least one amount to deal.
func (e FlankWalkDamage) validate() error {
	if len(e.Amounts) == 0 {
		return fmt.Errorf("FlankWalkDamage: needs at least one amount")
	}
	return nil
}

// Resolve chooses a flank creature and deals the amounts inward, all at once.
func (e FlankWalkDamage) Resolve(ctx *EffectContext) {
	chosen := (Target{Kind: TargetChosenCreature}).OnFlank().Select(ctx)
	if len(chosen) == 0 {
		return
	}
	bl := battlelineContaining(ctx, chosen[0])
	idx := -1
	for i, x := range bl {
		if x == chosen[0] {
			idx = i
			break
		}
	}
	step := 1
	if idx == len(bl)-1 {
		step = -1
	}
	var targets []DamageTarget
	for k, amt := range e.Amounts {
		pos := idx + k*step
		if pos < 0 || pos >= len(bl) {
			break
		}
		targets = append(targets, DamageTarget{ID: bl[pos], Amount: amt})
	}
	ctx.Resolver.DealDamage(ctx.Controller, targets)
}
