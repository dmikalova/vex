package engine

// A card that "cannot reap" is barred from one way of using it while every other
// way stays open — Tireless Crocag fights and uses its Action: normally. That is
// narrower than the timed, player-wide CannotUse/CannotFight restrictions in
// effect_restrict.go, so it lives on the card definition rather than on state.
// UseKind names one of the three ways a card in play can be used. Reaping,
// fighting, and using an "Action:" ability are the rulebook's whole list.
type UseKind uint8

const (
	useKindUnset UseKind = iota
	// ReapUse is using a creature to reap.
	ReapUse
	// FightUse is using a creature to fight.
	FightUse
	// ActionUse is using a card's "Action:" ability.
	ActionUse
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

// cannotBeUsedTo reports whether a card's printed text bars this way of using it.
func (g *Game) cannotBeUsedTo(id LocalID, kind UseKind) bool {
	for _, k := range g.cat.def(id).CannotBeUsedTo {
		if k == kind {
			return true
		}
	}
	return false
}
