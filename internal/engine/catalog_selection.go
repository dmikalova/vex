package engine

// This file is the census of the Selection family: one row per way a
// zone-movement verb picks the cards it acts on (ADR 0031). See catalog.go for
// the classification rule, and effect_selection.go for the selections.
//
// A selection is a sentence part, not a rule: "a card", "each creature", "the top
// card" and "a random card" are all the same rule about how card text names what
// an effect reaches, which the Target umbrella term carries. So every row names
// that one title rather than earning seven of its own.

// selectionFamily is the Selection family's census entry. A selection is
// discovered by its pick method.
func selectionFamily() Family {
	return newFamily(
		"Selection",
		"pick",
		[]string{"*EffectContext", "[]LocalID"},
		SelectionCatalog(),
		Selection.object,
	).gated()
}

// SelectionCatalog returns one representative value of every Selection, with the
// rulebook term each owes.
func SelectionCatalog() []Catalogued[Selection] {
	return []Catalogued[Selection]{
		{
			Node:  Chosen{Type: Creature},
			Rules: bears("Target"),
		},
		{
			Node:  Each{Type: Creature},
			Rules: bears("Target"),
		},
		{
			Node:  Random{},
			Rules: bears("Target"),
		},
		{
			Node:  Named{Name: "Velum"},
			Rules: bears("Target"),
		},
		{
			Node:  Self{},
			Rules: bears("Target"),
		},
		{
			Node:  Top{},
			Rules: bears("Target"),
		},
		{
			Node:  Bottom{},
			Rules: bears("Target"),
		},
	}
}
