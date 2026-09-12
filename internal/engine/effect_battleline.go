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
	// FromContext renders the swap as "swap it with {self}" instead of "swap this
	// creature with <With> in your battleline": the card swapped in is the contextual
	// card (ctx.It), which may rest off the battleline — Gebuk swaps in the creature
	// it just discarded — so naming a battleline position would be wrong.
	FromContext bool
}

// validate requires the creature to swap with.
func (e Swap) validate() error {
	if !e.With.valid() {
		return errUnsetTarget("Swap")
	}
	return nil
}

// Text renders the effect, e.g. "swap this creature with another friendly creature
// in your battleline", or "swap it with {self}" when the swapped-in card is the
// contextual card resting off the battleline.
func (e Swap) Text() string {
	if e.FromContext {
		return "swap it with " + SelfName
	}
	return "swap this creature with " + e.With.Text() + " in your battleline"
}

// Resolve swaps this creature with the selected card and puts that card in
// context. The two exchange places even across zones: a card selected from a
// discard pile enters play in this creature's slot while this creature leaves to
// that pile (Gebuk).
func (e Swap) Resolve(ctx *EffectContext) {
	for _, other := range e.With.Select(ctx) {
		ctx.Resolver.SwapCards(ctx.Source, other)
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
	ctx.Resolver.SwapCards(first, second)
}

// RearrangeBattleline lets the controller reorder one player's battleline by
// swapping pairs of creatures — Tactical Officer Moon. Each pass picks a creature
// and another in the same battleline and trades their positions; the controller
// keeps going until they stop, so doing nothing at all is allowed. Only positions
// move; no card state travels. The number of passes is bounded by the total
// creatures in play, enough to reach any arrangement and to keep the loop finite.
type RearrangeBattleline struct{}

// Text renders the effect.
func (RearrangeBattleline) Text() string {
	return "rearrange the creatures in a player's battleline"
}

// Resolve repeatedly swaps a chosen pair of creatures in one battleline until the
// controller declines the next pair.
func (RearrangeBattleline) Resolve(ctx *EffectContext) {
	both := append(
		append([]LocalID(nil), ctx.Resolver.Battleline(ctx.Controller)...),
		ctx.Resolver.Battleline(ctx.Opponent())...,
	)
	for range len(both) {
		all := append(
			append([]LocalID(nil), ctx.Resolver.Battleline(ctx.Controller)...),
			ctx.Resolver.Battleline(ctx.Opponent())...,
		)
		first, ok := ctx.ChooseCardOptional("Choose a creature to swap, or stop rearranging", all)
		if !ok {
			return
		}
		var others []LocalID
		for _, id := range ctx.Resolver.Battleline(ctx.Resolver.Controller(first)) {
			if id != first {
				others = append(others, id)
			}
		}
		second, ok := ctx.ChooseCreature("Choose the creature to swap it with", others)
		if !ok {
			return
		}
		ctx.Resolver.SwapCards(first, second)
	}
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
			"Choose a flank", []string{FlankLeftLabel, FlankRightLabel}) == 1
		ctx.Resolver.MoveToFlank(id, right)
	}
}

// TurnIntoCreature turns the card its Target selects into a creature and moves it
// onto a flank of its controller's battleline, the effect's controller choosing
// the flank. It is how an artifact turns itself into a creature (Auto-Legionary):
// the card keeps its exhaustion, Æmber, and power counters and reads as a creature
// until it leaves play, so power counters placed on it before the turn now count
// toward its power. A Target that selects nothing, or a card no longer in play, is
// a safe no-op.
type TurnIntoCreature struct {
	Target Target
}

// validate requires the card to convert.
func (e TurnIntoCreature) validate() error {
	if !e.Target.valid() {
		return errUnsetTarget("TurnIntoCreature")
	}
	return nil
}

// Text renders the effect, e.g. "move it to a flank of your battleline as a
// creature". The card is referred to as "it": this effect always follows an
// effect that named the source (Auto-Legionary gives itself counters first), so
// the second reference reads as a pronoun, matching KeyForge's own wording for
// Effigy of Melerukh and The Mysticeti.
func (e TurnIntoCreature) Text() string {
	return "move it to a flank of your battleline as a creature"
}

// Resolve converts each selected card and moves it to the flank the controller
// chooses.
func (e TurnIntoCreature) Resolve(ctx *EffectContext) {
	for _, id := range e.Target.Select(ctx) {
		if !ctx.Resolver.InPlay(id) {
			continue
		}
		right := ctx.ChooseOption(
			"Choose a flank", []string{FlankLeftLabel, FlankRightLabel}) == 1
		ctx.Resolver.PutIntoBattlelineAsCreature(id, right)
	}
}

// MoveWithinBattleline repositions the creature its Target selects anywhere in
// that creature's own controller's battleline, the effect's controller choosing
// the destination slot — Malison moves an enemy creature so its own controller may
// end up putting it on a flank. The moved creature is left in context (ctx.It) so
// a following effect can act on it. A Target that selects nothing, or a creature
// no longer on a battleline, leaves the battleline unchanged. resolveGate reports
// whether a creature was chosen, so a "you may" wrapper and a following Then read
// the choice.
type MoveWithinBattleline struct {
	Target Target
}

// validate requires the creature to move.
func (e MoveWithinBattleline) validate() error {
	if !e.Target.valid() {
		return errUnsetTarget("MoveWithinBattleline")
	}
	return nil
}

// Text renders the effect, e.g. "move an enemy creature anywhere in its
// controller's battleline".
func (e MoveWithinBattleline) Text() string {
	return "move " + e.Target.Text() + " anywhere in its controller's battleline"
}

// Resolve moves the selected creature within its battleline and leaves it in
// context.
func (e MoveWithinBattleline) Resolve(ctx *EffectContext) { e.resolveGate(ctx) }

// resolveGate moves the selected creature and reports whether one was chosen.
func (e MoveWithinBattleline) resolveGate(ctx *EffectContext) bool {
	ids := e.Target.Select(ctx)
	if len(ids) == 0 {
		return false
	}
	id := ids[0]
	if ctx.Resolver.InBattleline(id) {
		ctx.Resolver.MoveWithinBattleline(ctx.Controller, id)
	}
	ctx.It, ctx.HasIt = id, true
	return true
}
