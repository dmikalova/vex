package engine

// A card that "cannot reap" is barred from one way of using it while every other
// way stays open — Tireless Crocag fights and uses its Action: normally. That is
// narrower than the timed, player-wide Restrict restrictions in
// effect_restrict.go, so it lives on the card definition rather than on state.
// UseKind names one of the three ways a card in play can be used. Reaping,
// fighting, and using an "Action:" ability are the rulebook's whole list.
type UseKind uint8

const (
	// useKindUnset is the invalid zero value; a real use kind must be named.
	useKindUnset UseKind = iota
	// ReapUse is using a creature to reap.
	ReapUse
	// FightUse is using a creature to fight.
	FightUse
	// ActionUse is using a card's "Action:" ability.
	ActionUse
	// useKindCount bounds the enum for valid checks.
	useKindCount
)

// valid reports whether the use kind is one of the three real ways to use a card.
func (k UseKind) valid() bool { return k > useKindUnset && k < useKindCount }

// verb renders the use kind as the verb a card prints after "cannot", e.g.
// "Tireless Crocag cannot reap."
func (k UseKind) verb() string {
	switch k {
	case FightUse:
		return "fight"
	case ActionUse:
		return "use its Action ability"
	default:
		return "reap"
	}
}

// cannotBeUsedTo reports whether this way of using the card is barred — by the
// card's own printed text, or by a constant ability of a card in play that grants
// the restriction to the creatures it reaches (Narp bars its neighbors from
// reaping).
func (g *Game) cannotBeUsedTo(id LocalID, kind UseKind) bool {
	if g.creaturesGloballyBarred(id, kind) {
		return true
	}
	if c := g.cat.def(id).CannotBeUsedWhile; c != nil {
		ctx := &EffectContext{Resolver: g, Source: id, Controller: g.controller(id)}
		if c.Met(ctx) {
			return true
		}
	}
	for _, k := range g.cat.def(id).CannotBeUsedTo {
		if k == kind {
			return true
		}
	}
	for _, up := range g.upgradesOf(id) {
		for _, k := range g.staticOn(id, up).CannotBeUsedTo {
			if k == kind {
				return true
			}
		}
	}
	for p := 0; p < 2; p++ {
		for _, src := range g.allInPlay(p) {
			for _, c := range g.cat.def(src).ConstantAbilities {
				for _, k := range c.CannotBeUsedTo {
					if k == kind && g.constantActive(src, c) && g.constantAffects(src, c, id) {
						return true
					}
				}
			}
		}
	}
	return false
}

// CannotBeUsedTo reports whether this way of using the card is barred by a
// restriction — the card's own text or a constant ability that reaches it (Narp
// bars its neighbors from reaping). It is the restriction-only check a client uses
// to omit an illegal verb from the buttons it offers. Unlike CanUseTo it applies no
// house, exhaustion, or ownership gate, so it stays correct for a granted off-house
// use (Universal Translator's "use as if it were yours").
func (g *Game) CannotBeUsedTo(id LocalID, kind UseKind) bool {
	return g.cannotBeUsedTo(id, kind)
}
