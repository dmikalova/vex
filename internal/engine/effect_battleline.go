package engine

// Swap exchanges this creature's battleline position with the creature its With
// target selects, then puts that creature in context (ctx.It) so a following
// effect can act on it — Transposition Sandals swaps with another friendly
// creature and then uses it. Only positions move: no damage, upgrades, status,
// control, or other card state travels between the two creatures, since battleline
// order is all that matters for flanks and neighbors. A With target that selects
// nothing leaves the battleline unchanged. With is a full Target, so a later card
// can swap with an enemy creature rather than a friendly one.
type Swap struct {
	With Target
}

// validate requires the creature to swap with.
func (e Swap) validate() error {
	if !e.With.valid() {
		return errUnsetTarget("Swap")
	}
	return nil
}

// Text renders the effect, e.g. "swap this creature with another friendly creature
// in your battleline".
func (e Swap) Text() string {
	return "swap this creature with " + e.With.Text() + " in your battleline"
}

// Resolve swaps this creature's position with the selected creature and puts that
// creature in context.
func (e Swap) Resolve(ctx *EffectContext) {
	for _, other := range e.With.Select(ctx) {
		ctx.Resolver.SwapBattlelinePositions(ctx.Source, other)
		ctx.It, ctx.HasIt = other, true
	}
}

// SwapChosen exchanges the battleline positions of two creatures the controller
// chooses from a single battleline — Quantum Fingertrap's "swap the positions of
// two creatures in a battleline". Only positions move; no card state travels. The
// second creature is chosen from the same battleline as the first, so the two
// always share a battleline. If no first creature, or no second in that
// battleline, is chosen, the battleline is left unchanged.
type SwapChosen struct{}

// Text renders the effect.
func (SwapChosen) Text() string {
	return "swap the positions of two creatures in a battleline"
}

// Resolve chooses two creatures in one battleline and swaps their positions.
func (SwapChosen) Resolve(ctx *EffectContext) {
	all := append(
		append([]LocalID(nil), ctx.Resolver.Battleline(ctx.Controller)...),
		ctx.Resolver.Battleline(ctx.Opponent())...,
	)
	first, ok := ctx.ChooseCreature("Choose the first creature to swap", all)
	if !ok {
		return
	}
	var others []LocalID
	for _, id := range ctx.Resolver.Battleline(ctx.Resolver.Controller(first)) {
		if id != first {
			others = append(others, id)
		}
	}
	second, ok := ctx.ChooseCreature("Choose the second creature to swap", others)
	if !ok {
		return
	}
	ctx.Resolver.SwapBattlelinePositions(first, second)
}

// MoveToFlank moves the creature its Target selects to either flank of that
// creature's controller's battleline, the controller of the effect choosing
// which flank. Only the battleline slot moves; no card state travels. A Target
// that selects nothing, or a creature no longer on a battleline (destroyed,
// purged, or now an artifact), leaves the battleline unchanged — so following a
// DealDamage that destroyed the creature is a safe no-op.
type MoveToFlank struct {
	Target Target
}

// validate requires the creature to move.
func (e MoveToFlank) validate() error {
	if !e.Target.valid() {
		return errUnsetTarget("MoveToFlank")
	}
	return nil
}

// Text renders the effect, e.g. "move it to either flank of its controller's
// battleline".
func (e MoveToFlank) Text() string {
	return "move " + e.Target.Text() + " to either flank of its controller's battleline"
}

// Resolve moves each selected creature to the flank the controller chooses.
func (e MoveToFlank) Resolve(ctx *EffectContext) {
	for _, id := range e.Target.Select(ctx) {
		// A creature the preceding damage destroyed still selects as the triggering
		// creature, but it has left the battleline: moving it is a no-op, so skip it
		// before asking a flank rather than prompt for a move that cannot happen.
		if !ctx.Resolver.InBattleline(id) {
			continue
		}
		right := ctx.ChooseOption(
			"Choose a flank", []string{"left flank", "right flank"}) == 1
		ctx.Resolver.MoveToFlank(id, right)
	}
}
