package engine

// This file is the census of the PerTarget family: one row per axis along which
// an amount varies from one target to the next. See catalog.go for the
// classification rule, and effect_damage.go for the interface.
//
// A per-target axis measures something already on the creature, so each row binds
// to the term for the thing it measures.

// perTargetFamily is the PerTarget family's census entry. An axis is discovered
// by its perTargetValue method.
func perTargetFamily() Family {
	return newFamily(
		"PerTarget",
		"perTargetValue",
		[]string{"*EffectContext", "LocalID"},
		PerTargetCatalog(),
		PerTarget.perTargetText,
	).gated()
}

// PerTargetCatalog returns one representative value of every PerTarget, with the
// rulebook term each owes.
func PerTargetCatalog() []Catalogued[PerTarget] {
	return []Catalogued[PerTarget]{
		{
			Node:  AemberOnIt,
			Rules: bears("Æmber"),
		},
		{
			Node:  DamageOnIt,
			Rules: bears("Damage"),
		},
		{
			Node:  ArmorLostThisWay,
			Rules: bears("Lose Armor"),
		},
		{
			Node:  UpgradesOnIt,
			Rules: bears("Upgrade"),
		},
	}
}
