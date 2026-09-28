package engine

// This file is the census of the Gather family: one row per way a BatchDestroy
// computes the set of creatures it destroys in one simultaneous batch. See
// catalog.go for the classification rule, and effect_destroy.go for the
// gatherers.
//
// A gatherer only names which creatures the batch takes; the rule a player must
// know is the simultaneous destruction itself, which the Destroy term carries.

// gatherFamily is the Gather family's census entry. A gatherer is discovered by
// its gather method.
func gatherFamily() Family {
	return newFamily(
		"Gather",
		"gather",
		[]string{"*EffectContext"},
		GatherCatalog(),
		Gather.gatherText,
	).gated()
}

// GatherCatalog returns one representative value of every Gather, with the
// rulebook term each owes.
func GatherCatalog() []Catalogued[Gather] {
	return []Catalogued[Gather]{
		{
			Node: ChosenFromEach{
				Target{Kind: TargetChosenFriendlyCreature},
				Target{Kind: TargetChosenEnemyCreature},
			},
			Rules: bears("Destroy"),
		},
		{
			Node: EachPlayerUnless{
				Spare: CardsInPlay{
					Player: Controller,
					Filter: Filter{Type: Creature, Trait: Dinosaur},
				},
				Take: MostPowerful,
			},
			Rules: bears("Destroy"),
		},
	}
}
