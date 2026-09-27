package engine

// This file is the census of the Spread family: one row per shape a DealDamage
// takes when it strikes several related creatures at once. See catalog.go for the
// classification rule, and effect_damage.go for the spreads.
//
// A spread is how one batch of damage is distributed, so every row binds to the
// "Deal Damage" term the damage rules already carry. The shapes differ in which
// creatures are hit, not in what damage does.

// spreadFamily is the Spread family's census entry. A spread is discovered by
// its hits method.
func spreadFamily() Family {
	return newFamily(
		"Spread",
		"hits",
		[]string{"*EffectContext"},
		SpreadCatalog(),
		Spread.spreadText,
	).gated()
}

// SpreadCatalog returns one representative value of every Spread, with the
// rulebook term each owes.
func SpreadCatalog() []Catalogued[Spread] {
	return []Catalogued[Spread]{
		{
			Node:  CreatureAndNeighbors{Amount: 2, Splash: 1},
			Rules: bears("Deal Damage"),
		},
		{
			Node:  DifferentCreatures{First: 2, Second: 2},
			Rules: bears("Deal Damage"),
		},
		{
			Node:  UpToCreatures{Creatures: 3, Amount: 1},
			Rules: bears("Deal Damage"),
		},
		{
			Node:  DivideDamage{Amount: 2},
			Rules: bears("Deal Damage"),
		},
		{
			Node:  FlankWalk{Amounts: []int{3, 2, 1}},
			Rules: bears("Deal Damage"),
		},
	}
}
