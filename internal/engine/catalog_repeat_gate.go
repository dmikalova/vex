package engine

// This file is the census of the RepeatGate family: one row per way a Repeat
// decides whether to run its effect again. See catalog.go for the classification
// rule, and effect_repeat.go for the gates.
//
// Repeat itself is plumbing, and so are the two gates that only loop — everything
// a player reads comes from the effect being repeated and the condition already
// carrying its own term. ByExalting is not: it charges a cost to repeat, and that
// cost is a rule of its own.

// repeatGateBody is the effect the census renders each gate's clause around,
// standing in for the effect a card would repeat.
var repeatGateBody Effect = Draw{Amount: 1}

// repeatGateFamily is the RepeatGate family's census entry. A gate is discovered
// by its run method.
func repeatGateFamily() Family {
	return newFamily(
		"RepeatGate",
		"run",
		[]string{"*EffectContext", "Effect"},
		RepeatGateCatalog(),
		func(g RepeatGate) string { return g.text(repeatGateBody) },
	).gated()
}

// RepeatGateCatalog returns one representative value of every RepeatGate, with
// the rulebook term each owes.
func RepeatGateCatalog() []Catalogued[RepeatGate] {
	return []Catalogued[RepeatGate]{
		{
			Node:  While{Cond: Overwhelmed{}},
			Rules: plumbing("composition: reruns the effect while its condition holds"),
		},
		{
			Node: MayWhile{Cond: Overwhelmed{}},
			Rules: plumbing(
				"composition: reruns the effect while the controller keeps choosing to",
			),
		},
		{
			Node:  ByExalting{Creature: Target{Kind: TargetChosenFriendlyCreature}},
			Rules: bears("Exalt"),
		},
	}
}
