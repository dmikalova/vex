package engine

import "fmt"

// Dealing damage puts that much pending damage on each creature the effect
// targets. Armor prevents pending damage first — each point stops 1, and armor
// spent this way stays spent for the rest of the turn — and whatever is not
// prevented lands as damage tokens. A creature whose total damage reaches or
// exceeds its power is destroyed. When one ability deals damage to several
// creatures they are damaged simultaneously and any that died are destroyed
// together, so no creature's destruction changes another's.
type DealDamage struct {
	Amount int
	Per    Count
	Target Target
	// IgnoreArmor makes the damage bypass armor (Qyxxlyx Plague Master's "this
	// damage cannot be prevented by armor").
	IgnoreArmor bool
	// AmountFrom, when set, sources the damage from a produced quantity instead of
	// Amount, rendering "deal that amount of damage" — Guardian Demon deals the
	// damage it just healed. It is distinct from Per, which multiplies Amount "for
	// each"; set one or the other, not both.
	AmountFrom Count
	// Spread, when set, deals a simultaneous batch of damage to several related
	// creatures the controller chooses (a creature and a neighbor, a flank walk),
	// each with its own amount. It replaces Amount/Per/AmountFrom/Target — the
	// spread carries its own targets and amounts and renders its own phrase.
	Spread Spread
	// PerTarget multiplies Amount by a quantity read off each target separately,
	// which Per cannot do — Word of Returning deals 1 damage to each enemy creature
	// "for each Æmber on it", a different amount per creature.
	PerTarget PerTarget
}

// PerTarget is the axis along which an amount varies from one target to the next.
// A Count is evaluated once for the whole effect; a PerTarget is evaluated again
// for every creature the effect lands on.
type PerTarget interface {
	// perTargetValue is the multiplier for one target.
	perTargetValue(ctx *EffectContext, id LocalID) int
	// perTargetText is the noun the "for each" clause repeats.
	perTargetText() string
}

// AemberOnIt scales per target by the Æmber sitting on that target.
var AemberOnIt PerTarget = aemberOnIt{}

type aemberOnIt struct{}

func (aemberOnIt) perTargetValue(ctx *EffectContext, id LocalID) int {
	return ctx.Resolver.AmberOn(id)
}

func (aemberOnIt) perTargetText() string { return "\u00c6mber on it" }

// validate requires an explicit target, or a Spread that supplies its own.
func (e DealDamage) validate() error {
	if e.Spread != nil {
		if e.Target.valid() || e.Amount != 0 || e.Per != nil || e.AmountFrom != nil {
			return fmt.Errorf(
				"DealDamage: Spread cannot combine with Target, Amount, Per, or AmountFrom",
			)
		}
		return validateSpread(e.Spread)
	}
	if !e.Target.valid() {
		return errUnsetTarget("DealDamage")
	}
	if e.AmountFrom != nil && e.Per != nil {
		return fmt.Errorf("DealDamage: set AmountFrom or Per, not both")
	}
	if e.PerTarget != nil && (e.AmountFrom != nil || e.Per != nil) {
		return fmt.Errorf("DealDamage: PerTarget cannot combine with Per or AmountFrom")
	}
	return nil
}

// Text renders the effect, e.g. "deal 2 damage to each enemy creature". A "for
// each" count leads the sentence (rule 9), e.g. "for each friendly creature in
// play, deal 1 damage to a creature". Armor-ignoring damage adds a trailing clause.
func (e DealDamage) Text() string {
	if e.Spread != nil {
		return e.Spread.spreadText()
	}
	amount := fmt.Sprintf("%d damage", e.Amount)
	if e.AmountFrom != nil {
		amount = "that amount of damage"
	}
	body := fmt.Sprintf("deal %s to %s", amount, e.Target.Text())
	if e.PerTarget != nil {
		body += " for each " + e.PerTarget.perTargetText()
	}
	if e.IgnoreArmor {
		body += ", ignoring armor"
	}
	return forEach(e.Per, body)
}

// Resolve deals the damage to every selected creature simultaneously, resolving
// destruction as part of it. A Per count multiplies the amount dealt; an AmountFrom
// count sources the amount directly; a Spread deals its own related batch of hits.
// A computed amount of zero deals nothing, so the target is never selected — a
// chosen one would otherwise be a vacuous prompt (Guardian Demon's follow-up when
// its heal removed no damage).
func (e DealDamage) Resolve(ctx *EffectContext) {
	if e.Spread != nil {
		if hits := e.Spread.hits(ctx); len(hits) > 0 {
			ctx.Resolver.DealDamage(ctx.Controller, hits)
		}
		return
	}
	if amount := e.amount(ctx); amount > 0 {
		e.dealTo(ctx, amount, e.Target.Select(ctx))
	}
}

// declinable reports that the damage lands on a single clickable creature, so a
// "you may" wrapping it (Rock-Hurling Giant) can be answered by clicking that
// creature instead of a separate Yes/No.
func (e DealDamage) declinable() bool { return e.Spread == nil && e.Target.isChosen() }

// resolveOptional is Resolve under a May: the creature is asked declinably, with
// a Done to decline.
func (e DealDamage) resolveOptional(ctx *EffectContext) bool {
	amount := e.amount(ctx)
	if amount <= 0 {
		return false
	}
	ids := e.Target.SelectOptional(ctx)
	if len(ids) == 0 {
		return false
	}
	e.dealTo(ctx, amount, ids)
	return true
}

// amount computes how much damage a non-Spread DealDamage deals, before any
// per-target multiplier.
func (e DealDamage) amount(ctx *EffectContext) int {
	amount := e.Amount
	switch {
	case e.AmountFrom != nil:
		amount = e.AmountFrom.Value(ctx)
	case e.Per != nil:
		amount *= e.Per.Value(ctx)
	}
	return amount
}

// dealTo deals amount (before any PerTarget multiplier) to an already-selected
// set of creatures simultaneously.
func (e DealDamage) dealTo(ctx *EffectContext, amount int, ids []LocalID) {
	targets := make([]DamageTarget, len(ids))
	for i, id := range ids {
		hit := amount
		if e.PerTarget != nil {
			hit *= e.PerTarget.perTargetValue(ctx, id)
		}
		targets[i] = DamageTarget{ID: id, Amount: hit, IgnoreArmor: e.IgnoreArmor}
	}
	ctx.Resolver.DealDamage(ctx.Controller, targets)
}

// DamageThenIfDestroyed deals damage to one chosen creature and, only if that damage
// destroys it, resolves a follow-up effect with the destroyed creature in context
// (ctx.It) — Seeker Needle's "deal 1 damage to a creature. If this damage destroys
// that creature, gain 1 Æmber."
type DamageThenIfDestroyed struct {
	Amount int
	Target Target
	Then   Effect
}

// validate requires a target and a well-formed follow-up effect.
func (e DamageThenIfDestroyed) validate() error {
	if !e.Target.valid() {
		return errUnsetTarget("DamageThenIfDestroyed")
	}
	return validateEffect(e.Then)
}

// Text renders the effect, e.g. "deal 2 damage to a creature. If this damage
// destroys that creature, steal 1 Æmber".
func (e DamageThenIfDestroyed) Text() string {
	return fmt.Sprintf("deal %d damage to %s. If this damage destroys that creature, %s",
		e.Amount, e.Target.Text(), e.Then.Text())
}

// Resolve deals the damage to the chosen creature, then runs Then only if the
// creature has left play. The destroyed creature is placed in context (ctx.It) so
// Then can refer to it ("purge it").
func (e DamageThenIfDestroyed) Resolve(ctx *EffectContext) {
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

// DamageThenIfSurvives deals damage to one chosen creature and, only if the creature
// is not destroyed, resolves a follow-up effect with the surviving creature in
// context (ctx.It) — Gongoozle's "deal 3 damage to a creature. If it is not
// destroyed, its owner discards a random card from their hand." This is a plain
// state branch, not a result gate.
type DamageThenIfSurvives struct {
	Amount int
	Target Target
	Then   Effect
}

// validate requires a target and a well-formed follow-up effect.
func (e DamageThenIfSurvives) validate() error {
	if !e.Target.valid() {
		return errUnsetTarget("DamageThenIfSurvives")
	}
	return validateEffect(e.Then)
}

// Text renders the effect, e.g. "deal 3 damage to a creature. If it is not
// destroyed, its owner discards a random card from their hand".
func (e DamageThenIfSurvives) Text() string {
	return fmt.Sprintf("deal %d damage to %s. If it is not destroyed, %s",
		e.Amount, e.Target.Text(), e.Then.Text())
}

// Resolve deals the damage to the chosen creature, then runs Then only if the
// creature is still in play. The surviving creature is placed in context (ctx.It)
// so Then can refer to it ("its owner").
func (e DamageThenIfSurvives) Resolve(ctx *EffectContext) {
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

// Spread turns the controller's choices into a simultaneous batch of damage hits
// and renders its own phrase — the shape a DealDamage takes when it strikes
// several related creatures at once (a chosen creature and a neighbor, an inward
// flank walk). Each renders its clause so text and behavior stay in sync.
type Spread interface {
	hits(ctx *EffectContext) []DamageTarget
	spreadText() string
}

// validateSpread returns any configuration error a spread reports (FlankWalk needs
// at least one amount); spreads that cannot be misconfigured pass.
func validateSpread(s Spread) error {
	if v, ok := s.(validator); ok {
		return v.validate()
	}
	return nil
}

// CreatureAndNeighbor deals Amount to a chosen creature and NeighborAmount to one
// of its battleline neighbors (the controller picks when it has two), all at once
// — Mighty Lance. A DealDamage Spread.
type CreatureAndNeighbor struct {
	Amount         int
	NeighborAmount int
}

// spreadText renders the clause.
func (s CreatureAndNeighbor) spreadText() string {
	return fmt.Sprintf("deal %d damage to a creature and %d damage to a neighbor of that creature",
		s.Amount, s.NeighborAmount)
}

// hits chooses a creature and one chosen neighbor.
func (s CreatureAndNeighbor) hits(ctx *EffectContext) []DamageTarget {
	chosen := (Target{Kind: TargetChosenCreature}).Select(ctx)
	if len(chosen) == 0 {
		return nil
	}
	out := []DamageTarget{{ID: chosen[0], Amount: s.Amount}}
	if ns := neighbors(ctx, chosen[0]); len(ns) > 0 {
		n := ns[0]
		if len(ns) > 1 {
			if pick, ok := ctx.ChooseCreature("Choose a neighbor", ns); ok {
				n = pick
			}
		}
		out = append(out, DamageTarget{ID: n, Amount: s.NeighborAmount})
	}
	return out
}

// CreatureAndNeighbors deals Amount to a chosen creature and Splash to each of
// that creature's neighbors, all at once — the "with N splash" wording (Lava
// Ball). NotOnFlank narrows the choice to an interior creature, which always has
// two neighbors (Booby Trap). A DealDamage Spread.
type CreatureAndNeighbors struct {
	Amount     int
	Splash     int
	NotOnFlank bool
}

// spreadText renders the clause.
func (s CreatureAndNeighbors) spreadText() string {
	creature := "a creature"
	if s.NotOnFlank {
		creature = "a creature that is not on a flank"
	}
	return fmt.Sprintf(
		"deal %d damage to %s and %d damage to each of its neighbors",
		s.Amount,
		creature,
		s.Splash,
	)
}

// hits chooses a creature, then hits it and its neighbors. NotOnFlank keeps the
// choice to an interior creature.
func (s CreatureAndNeighbors) hits(ctx *EffectContext) []DamageTarget {
	target := Target{Kind: TargetChosenCreature}
	if s.NotOnFlank {
		target = target.NotOnFlank()
	}
	chosen := target.Select(ctx)
	if len(chosen) == 0 {
		return nil
	}
	out := []DamageTarget{{ID: chosen[0], Amount: s.Amount}}
	for _, n := range neighbors(ctx, chosen[0]) {
		out = append(out, DamageTarget{ID: n, Amount: s.Splash})
	}
	return out
}

// DifferentCreatures deals First to a chosen creature and Second to a second,
// different chosen creature, all at once — Twin Bolt Emission. A DealDamage Spread.
type DifferentCreatures struct {
	First  int
	Second int
}

// spreadText renders the clause.
func (s DifferentCreatures) spreadText() string {
	return fmt.Sprintf("deal %d damage to a creature and deal %d damage to a different creature",
		s.First, s.Second)
}

// hits chooses a creature and a different creature. With only one creature in
// play, the different creature cannot be chosen and only the first is hit.
func (s DifferentCreatures) hits(ctx *EffectContext) []DamageTarget {
	chosen := (Target{Kind: TargetChosenCreature}).Select(ctx)
	if len(chosen) == 0 {
		return nil
	}
	out := []DamageTarget{{ID: chosen[0], Amount: s.First}}
	if others := creaturesExcept(ctx, chosen[0]); len(others) > 0 {
		if id, ok := ctx.ChooseCreature("Choose a different creature", others); ok {
			out = append(out, DamageTarget{ID: id, Amount: s.Second})
		}
	}
	return out
}

// FlankWalk chooses a flank creature and deals decreasing damage inward along its
// battleline: Amounts[0] to the chosen flank creature, Amounts[1] to its neighbor,
// and so on — Positron Bolt (Amounts{3, 2, 1}). The walk stops at the far flank if
// the battleline is shorter than the list of amounts. A DealDamage Spread.
type FlankWalk struct {
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

// spreadText renders the clause, e.g. "choose a flank creature. Deal 3 damage to
// it, 2 damage to its neighbor, and 1 damage to the neighbor's other neighbor."
func (s FlankWalk) spreadText() string {
	parts := make([]string, len(s.Amounts))
	for i, a := range s.Amounts {
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
func (s FlankWalk) validate() error {
	if len(s.Amounts) == 0 {
		return fmt.Errorf("FlankWalk: needs at least one amount")
	}
	return nil
}

// hits chooses a flank creature and walks the amounts inward, all at once.
func (s FlankWalk) hits(ctx *EffectContext) []DamageTarget {
	chosen := (Target{Kind: TargetChosenCreature}).OnFlank().Select(ctx)
	if len(chosen) == 0 {
		return nil
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
	var out []DamageTarget
	for k, amt := range s.Amounts {
		pos := idx + k*step
		if pos < 0 || pos >= len(bl) {
			break
		}
		out = append(out, DamageTarget{ID: bl[pos], Amount: amt})
	}
	return out
}
