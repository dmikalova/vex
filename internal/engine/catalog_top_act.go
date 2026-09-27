package engine

// This file is the census of the TopAct family: one row per routing step over the
// cards a look at the top of a deck read. See catalog.go for the classification
// rule, and effect_deck.go for the steps.
//
// Each step is an ordinary zone movement or deck operation wearing a step's
// clothes, so the rows bind to the terms those rules already carry rather than
// naming the step itself.

// topActFamily is the TopAct family's census entry. A step is discovered by its
// terminal method.
func topActFamily() Family {
	return newFamily(
		"TopAct",
		"terminal",
		nil,
		TopActCatalog(),
		TopAct.clause,
	).gated()
}

// TopActCatalog returns one representative value of every TopAct, with the
// rulebook term each owes.
func TopActCatalog() []Catalogued[TopAct] {
	return []Catalogued[TopAct]{
		{
			Node:  ChooseAndMove{Cards: 1, Dest: IntoHand},
			Rules: bears("Put a Card into Another Zone"),
		},
		{
			Node:  PartitionByChosenHouse{Matching: IntoArchives, Rest: IntoDiscard},
			Rules: bears("Put a Card into Another Zone"),
		},
		{
			Node:  MayDiscardLookedAt{},
			Rules: bears("Discard"),
		},
		{
			Node:  ReorderRest{},
			Rules: bears("Reveal Top of Deck"),
		},
		{
			Node:  Shuffle{},
			Rules: bears("Shuffle"),
		},
	}
}
