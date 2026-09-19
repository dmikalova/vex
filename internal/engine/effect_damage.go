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
	// Amount is the base damage; Per multiplies it by a running count "for each";
	// Target is the creatures the damage lands on.
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
	// Then, when set, resolves on the single damaged creature once the damage has
	// landed, with After deciding which outcomes it runs after. Reach for it — not a
	// surrounding Then{} — whenever the follow-up depends on what *this* damage did:
	// only this node can tell whether this hit is what destroyed the creature, and it
	// snapshots the creature's neighbours before the hit so the follow-up can still
	// reach them. A follow-up that merely happens next, and would run the same after
	// any first half, belongs in a surrounding Then{} or Sequence instead.
	// Setting it narrows the effect to one target, so it cannot combine with Spread,
	// Per, PerTarget, or AmountFrom.
	Then Effect
	// After decides which damage outcomes Then runs after. It is required with Then
	// and meaningless without it.
	After DamageAftermath
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
	return ctx.amberOn(id)
}

func (aemberOnIt) perTargetText() string { return "\u00c6mber on it" }

// DamageOnIt scales per target by the damage already sitting on that target, so a
// single effect hits each creature for as much as it is already hurt (Cauldron
// Boil deals each creature damage equal to the damage already on it).
var DamageOnIt PerTarget = damageOnIt{}

type damageOnIt struct{}

func (damageOnIt) perTargetValue(ctx *EffectContext, id LocalID) int {
	return ctx.damageOn(id)
}

func (damageOnIt) perTargetText() string { return "point of damage on it" }

// validate requires an explicit target, or a Spread that supplies its own.
func (e DealDamage) validate() error {
	if err := e.validateAftermath(); err != nil {
		return err
	}
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

// validateAftermath holds the Then/After pair to the single-target shape it can
// resolve in: the aftermath asks what happened to one creature, so an effect that
// spreads, repeats, or varies its amount has no such creature to ask about.
func (e DealDamage) validateAftermath() error {
	if e.Then == nil {
		if e.After != aftermathUnset {
			return fmt.Errorf("DealDamage: After needs a Then to run")
		}
		return nil
	}
	switch e.After {
	case Always, IfDestroyed, IfSurvives:
	default:
		return fmt.Errorf("DealDamage: After must be Always, IfDestroyed, or IfSurvives")
	}
	if e.Spread != nil || e.Per != nil || e.PerTarget != nil || e.AmountFrom != nil {
		return fmt.Errorf(
			"DealDamage: Then cannot combine with Spread, Per, PerTarget, or AmountFrom",
		)
	}
	return validateEffect(e.Then)
}

// Text renders the effect, e.g. "deal 2 damage to each enemy creature". A "for
// each" count leads the sentence (rule 9), e.g. "for each friendly creature in
// play, deal 1 damage to a creature". Armor-ignoring damage adds a trailing clause.
func (e DealDamage) Text() string {
	if e.Spread != nil {
		return e.Spread.spreadText()
	}
	amount := damageAmount(e.Amount)
	if e.AmountFrom != nil {
		switch e.AmountFrom.(type) {
		case DamageHealed:
			amount = "that amount of damage"
		default:
			amount = "damage equal to " + e.AmountFrom.CountText()
		}
	}
	body := fmt.Sprintf("deal %s to %s", amount, e.Target.Text())
	if e.PerTarget != nil {
		body += " for each " + e.PerTarget.perTargetText()
	}
	if e.IgnoreArmor {
		body += ", ignoring armor"
	}
	if e.Then != nil {
		body += e.aftermathText()
	}
	return forEach(e.Per, body)
}

// aftermathText renders the join between the damage clause and its follow-up,
// e.g. ". If this damage destroys that creature, gain 1 Æmber".
func (e DealDamage) aftermathText() string {
	switch e.After {
	case IfDestroyed:
		return ". If this damage destroys that creature, " + e.Then.Text()
	case IfSurvives:
		return ". If it is not destroyed, " + e.Then.Text()
	default:
		return " and " + e.Then.Text()
	}
}

// Resolve deals the damage to every selected creature simultaneously, resolving
// destruction as part of it. A Per count multiplies the amount dealt; an AmountFrom
// count sources the amount directly; a Spread deals its own related batch of hits.
// A "for each" count aimed at a single chosen creature instead makes one target
// choice per instance and deals them all at once (see resolvePerInstance).
// A computed amount of zero deals nothing, so the target is never selected — a
// chosen one would otherwise be a vacuous prompt (Guardian Demon's follow-up when
// its heal removed no damage).
func (e DealDamage) Resolve(ctx *EffectContext) {
	if e.Then != nil {
		e.resolveAftermath(ctx)
		return
	}
	if e.Spread != nil {
		if hits := e.Spread.hits(ctx); len(hits) > 0 {
			owners := make([]int, len(hits))
			for i, h := range hits {
				owners[i] = ctx.Resolver.Controller(h.ID)
			}
			ctx.dealDamage(hits)
			for i, h := range hits {
				if !resolverInPlay(ctx, h.ID) {
					ctx.Produced.Destroyed[owners[i]]++
				}
			}
		}
		return
	}
	if e.Per != nil && e.Target.isChosen() {
		e.resolvePerInstance(ctx)
		return
	}
	if amount := e.amount(ctx); amount > 0 {
		e.dealTo(ctx, amount, e.Target.Select(ctx))
	}
}

// resolvePerInstance resolves a "for each" DealDamage whose target is a single
// chosen creature: it makes one choice per unit of the count — each an Amount hit
// on a creature the controller picks, and each free to name a different creature —
// then deals the whole batch at once (KeyForge's independent-target ruling; a
// creature named twice takes both hits together). Sack of Coins is the opposite
// shape — one creature chosen first takes all the damage — authored with
// ChooseCreatureThen and Target.TheChosenCreature.
func (e DealDamage) resolvePerInstance(ctx *EffectContext) {
	n := e.Per.Value(ctx)
	if n <= 0 || e.Amount <= 0 {
		return
	}
	cands := e.Target.candidates(ctx)
	if len(cands) == 0 {
		return
	}
	prompt := "Choose " + e.Target.Text()
	ctx.previewBadge(SelectionBadge{Icon: DamageIcon, Amount: e.Amount})
	defer ctx.previewBadge(SelectionBadge{})
	assigned := map[LocalID]int{}
	order := []LocalID{}
	for i := 0; i < n; i++ {
		id, ok := ctx.ChooseCreature(prompt, cands)
		if !ok {
			continue
		}
		if _, seen := assigned[id]; !seen {
			order = append(order, id)
		}
		assigned[id] += e.Amount
	}
	if len(order) == 0 {
		return
	}
	targets := make([]DamageTarget, len(order))
	for i, id := range order {
		targets[i] = DamageTarget{ID: id, Amount: assigned[id], IgnoreArmor: e.IgnoreArmor}
	}
	ctx.dealDamage(targets)
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
	if e.AmountFrom != nil {
		return e.AmountFrom.Value(ctx)
	}
	return scaled(e.Amount, e.Per, ctx)
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
		targets[i] = DamageTarget{
			ID:          id,
			Amount:      hit,
			IgnoreArmor: e.IgnoreArmor,
			Source:      ctx.Source,
		}
	}
	ctx.dealDamage(targets)
}

// DamageAftermath decides when a DealDamage's follow-up resolves and how its two
// clauses join in printed text. It has no valid zero value: an effect with a Then
// must name one so the branch is never left ambiguous.
type DamageAftermath uint8

const (
	// aftermathUnset is the invalid zero value: a Then must name its After.
	aftermathUnset DamageAftermath = iota
	// Always runs the follow-up whether or not the damage destroyed the creature,
	// joining the two clauses with "and" (Tyxl Beambuckler).
	Always
	// IfDestroyed runs the follow-up only if the damage destroyed the creature,
	// reading ". If this damage destroys that creature, …" (Seeker Needle).
	IfDestroyed
	// IfSurvives runs the follow-up only if the creature is not destroyed, reading
	// ". If it is not destroyed, …" (Gongoozle).
	IfSurvives
)

// resolveAftermath deals the damage to the chosen creature, then resolves Then
// when After's branch holds. IfDestroyed snapshots the creature's neighbors before
// the damage (a following effect may hit the destroyed creature's former
// neighbors); the creature is placed in context (ctx.It) for Then to refer to.
func (e DealDamage) resolveAftermath(ctx *EffectContext) {
	ids := e.Target.Select(ctx)
	if len(ids) == 0 {
		return
	}
	id := ids[0]
	if e.After == IfDestroyed {
		ctx.Produced.Neighbors = neighbors(ctx, id)
		captureDepartingSubject(ctx, id)
	}
	ctx.dealDamage([]DamageTarget{{ID: id, Amount: e.Amount, IgnoreArmor: e.IgnoreArmor}})
	switch e.After {
	case IfDestroyed:
		if resolverInPlay(ctx, id) {
			return
		}
	case IfSurvives:
		if !resolverInPlay(ctx, id) {
			return
		}
	}
	ctx.It, ctx.HasIt = id, true
	e.Then.Resolve(ctx)
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

// NeighborScope selects how many of a struck creature's neighbors a spread hits.
type NeighborScope uint8

const (
	// AllNeighbors hits each of the struck creature's neighbors — the "with N splash"
	// wording (Lava Ball). The zero value, so a spread need not name it.
	AllNeighbors NeighborScope = iota
	// OneNeighbor hits a single neighbor the controller chooses when the creature has
	// two (Mighty Lance).
	OneNeighbor
)

// CreatureAndNeighbors deals Amount to a chosen creature and Splash to its
// battleline neighbors, all at once — a DealDamage Spread. Scope picks whether it
// hits each neighbor (AllNeighbors, the default) or one the controller chooses
// (OneNeighbor, Mighty Lance). NotOnFlank narrows the AllNeighbors choice to an
// interior creature, which always has two neighbors (Booby Trap).
type CreatureAndNeighbors struct {
	Amount     int
	Splash     int
	Scope      NeighborScope
	NotOnFlank bool
	// Target, when set, names the primary creature to hit instead of prompting the
	// controller to choose one — Plasma Nozzle aims the spread at the creature its
	// host fights (Target.CreatureFought). Its zero value falls back to the default
	// chosen-creature behavior. It replaces the AllNeighbors chosen path; Scope
	// still selects one neighbor or each.
	Target Target
}

// spreadText renders the clause, naming one neighbor or each depending on Scope.
// The creature taking the base damage is named once and the neighbors hang off it
// as "its", so the clause never has to say "that creature" twice. A splash equal
// to the base damage is stated once rather than repeated.
func (s CreatureAndNeighbors) spreadText() string {
	neighbors := "each of its neighbors"
	if s.Scope == OneNeighbor {
		neighbors = "one of its neighbors"
	}
	prefix, subject := "choose a creature. Deal ", "the chosen creature"
	switch {
	case s.Target.Kind != targetUnset:
		prefix, subject = "deal ", s.Target.Text()
	case s.NotOnFlank:
		prefix = "choose a creature that is not on a flank. Deal "
	}
	if s.Amount == s.Splash {
		return prefix + damageAmount(s.Amount) + " to " + subject + " and " + neighbors
	}
	return prefix + damageAmount(s.Amount) + " to " + subject +
		" and " + damageAmount(s.Splash) + " to " + neighbors
}

// hits picks the primary creature — a named Target when set, otherwise one the
// controller chooses — then hits it and its neighbors: one chosen neighbor under
// OneNeighbor, every neighbor otherwise. NotOnFlank keeps the choice to an
// interior creature.
func (s CreatureAndNeighbors) hits(ctx *EffectContext) []DamageTarget {
	if s.Target.Kind != targetUnset {
		var out []DamageTarget
		for _, id := range s.Target.Select(ctx) {
			out = append(out, DamageTarget{ID: id, Amount: s.Amount})
			for _, n := range neighbors(ctx, id) {
				out = append(out, DamageTarget{ID: n, Amount: s.Splash})
			}
		}
		return out
	}
	target := Target{Kind: TargetChosenCreature}
	if s.NotOnFlank {
		target = target.NotOnFlank()
	}
	chosen := target.Select(ctx)
	if len(chosen) == 0 {
		return nil
	}
	out := []DamageTarget{{ID: chosen[0], Amount: s.Amount}}
	ns := neighbors(ctx, chosen[0])
	if s.Scope == OneNeighbor {
		if len(ns) > 0 {
			n := ns[0]
			if len(ns) > 1 {
				if pick, ok := ctx.ChooseCreature("Choose a neighbor", ns); ok {
					n = pick
				}
			}
			out = append(out, DamageTarget{ID: n, Amount: s.Splash})
		}
		return out
	}
	for _, n := range ns {
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

// spreadText renders the clause, stating the verb once and letting the second
// amount stand alone — "deal 2 damage to a creature and 2 damage to a different
// creature" — so it reads the same whether or not the amounts match.
func (s DifferentCreatures) spreadText() string {
	return dealDamageTo(s.First, "a creature") + " and " +
		damageAmount(s.Second) + " to a different creature"
}

// hits picks as many distinct creatures as the positional amounts name, each pick
// excluding the ones before it (pickCards never repeats). With only one creature
// in play, the different creature cannot be chosen and only the first is hit.
func (s DifferentCreatures) hits(ctx *EffectContext) []DamageTarget {
	amounts := []int{s.First, s.Second}
	pool := Target{Kind: TargetEachCreature}
	picked := pickCards(ctx, "Choose a creature to deal damage to", len(amounts), false,
		func() []LocalID { return pool.Select(ctx) })
	out := make([]DamageTarget, 0, len(picked))
	for i, id := range picked {
		out = append(out, DamageTarget{ID: id, Amount: amounts[i]})
	}
	return out
}

// UpToCreatures deals Amount to up to Creatures different creatures the controller
// chooses one at a time, declining any of them with Done — Throwing Stars deals 1
// damage to up to 3 creatures. Undamaged narrows the choice to creatures that have
// no damage on them (Unsuspecting Prey). WhenDamaged, when set, deals that larger
// amount instead to any chosen creature that was already damaged before this batch
// (Festering Touch deals 1, or 3 to an already-damaged creature). A DealDamage
// Spread.
type UpToCreatures struct {
	Creatures   int
	Amount      int
	Undamaged   bool
	WhenDamaged int
}

// validate requires room for at least one creature.
func (s UpToCreatures) validate() error {
	if s.Creatures < 1 {
		return fmt.Errorf("UpToCreatures: Creatures must be at least 1, got %d", s.Creatures)
	}
	if s.WhenDamaged != 0 && s.Undamaged {
		return fmt.Errorf("UpToCreatures: WhenDamaged cannot combine with Undamaged")
	}
	return nil
}

// spreadText renders the clause. The already-damaged case is phrased as its own
// subset of the chosen creatures, so the clause never shifts from "each chosen
// creature" to a singular "that creature" mid-sentence.
func (s UpToCreatures) spreadText() string {
	if s.WhenDamaged != 0 {
		return fmt.Sprintf(
			"choose up to %s. Deal %s to each chosen creature. "+
				"Deal %s instead to each chosen creature that was already damaged",
			countNoun(s.Creatures, "creature"),
			damageAmount(s.Amount),
			damageAmount(s.WhenDamaged),
		)
	}
	noun := "creature"
	if s.Undamaged {
		noun = "undamaged creature"
	}
	return fmt.Sprintf(
		"deal %s to up to %s",
		damageAmount(s.Amount),
		countNoun(s.Creatures, noun),
	)
}

// hits asks for creatures one at a time, up to Creatures, stopping when the controller
// declines or none remain. Undamaged narrows the pool through the shared
// Target.Undamaged() refinement. WhenDamaged raises the amount for a creature that
// already carries damage before this batch resolves.
func (s UpToCreatures) hits(ctx *EffectContext) []DamageTarget {
	pool := Target{Kind: TargetEachCreature}
	if s.Undamaged {
		pool = pool.Undamaged()
	}
	picked := pickCards(ctx, "Choose a creature", s.Creatures, true, func() []LocalID {
		return pool.Select(ctx)
	})
	out := make([]DamageTarget, 0, len(picked))
	for _, id := range picked {
		amount := s.Amount
		if s.WhenDamaged != 0 && ctx.Resolver.Damage(id) > 0 {
			amount = s.WhenDamaged
		}
		out = append(out, DamageTarget{ID: id, Amount: amount})
	}
	return out
}

// DivideDamage deals a pool of Amount damage — scaled by an optional Per count —
// that the controller divides one point at a time among any number of creatures,
// all landing at once. First Blood deals 2 damage for each friendly Brobnar
// creature, divided freely. A DealDamage Spread.
type DivideDamage struct {
	Amount int
	Per    Count
}

// validate requires a positive base amount.
func (s DivideDamage) validate() error {
	if s.Amount < 1 {
		return fmt.Errorf("DivideDamage: Amount must be at least 1, got %d", s.Amount)
	}
	return nil
}

// spreadText renders the clause, e.g. "deal 2 damage for each friendly Brobnar
// creature, divided among any number of creatures".
func (s DivideDamage) spreadText() string {
	amount := "deal " + damageAmount(s.Amount)
	if s.Per != nil {
		amount += " for each " + s.Per.CountText()
	}
	return amount + ", divided among any number of creatures"
}

// hits places the pool one point at a time onto creatures the controller chooses,
// accumulating the points into a simultaneous batch of hits.
func (s DivideDamage) hits(ctx *EffectContext) []DamageTarget {
	total := scaled(s.Amount, s.Per, ctx)
	assigned := map[LocalID]int{}
	order := []LocalID{}
	for i := 0; i < total; i++ {
		var cands []LocalID
		for p := 0; p < 2; p++ {
			cands = append(cands, ctx.Resolver.Battleline(p)...)
		}
		if len(cands) == 0 {
			break
		}
		id, _ := ctx.ChooseCard("Choose a creature to damage", cands)
		if _, seen := assigned[id]; !seen {
			order = append(order, id)
		}
		assigned[id]++
	}
	out := make([]DamageTarget, len(order))
	for i, id := range order {
		out[i] = DamageTarget{ID: id, Amount: assigned[id]}
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

// flankWalkStep pairs a battleline creature with the amount it receives at its
// step of an inward flank walk.
type flankWalkStep struct {
	ID     LocalID
	Amount int
}

// flankWalkSteps chooses a flank creature and walks amounts inward along its
// battleline, pairing each creature with its step's amount — Amounts[0] to the
// chosen flank creature, Amounts[1] to its neighbor, and so on, stopping at the
// far flank when the battleline is shorter than the list. Shared by the FlankWalk
// damage spread and AddPowerCounter's counter walk.
func flankWalkSteps(ctx *EffectContext, amounts []int) []flankWalkStep {
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
	var out []flankWalkStep
	for k, amt := range amounts {
		pos := idx + k*step
		if pos < 0 || pos >= len(bl) {
			break
		}
		out = append(out, flankWalkStep{ID: bl[pos], Amount: amt})
	}
	return out
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
	steps := flankWalkSteps(ctx, s.Amounts)
	if len(steps) == 0 {
		return nil
	}
	out := make([]DamageTarget, len(steps))
	for i, st := range steps {
		out[i] = DamageTarget{ID: st.ID, Amount: st.Amount}
	}
	return out
}
