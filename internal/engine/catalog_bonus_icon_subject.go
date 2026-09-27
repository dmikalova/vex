package engine

// This file is the census of the BonusIconSubject family: the two subjects a
// BonusIconsOf count reads its icons from. See catalog.go for the classification
// rule, and effect_count.go for the interface and the count that varies along
// it. TestFamilyTotality keeps it complete.
//
// Both rows rest on the bonus-icon rules the Bonus icon section already carries
// (ADR 0041) — which cards' icons a count totals is part of how bonus icons
// resolve, not a rule of its own — so both bind to the term that section's rules
// are named by rather than earning titles of their own.

// bonusIconSubjectFamily is the BonusIconSubject family's census entry. A
// subject is discovered by its bonusIconCards method.
func bonusIconSubjectFamily() Family {
	return newFamily(
		"BonusIconSubject",
		"bonusIconCards",
		[]string{"*EffectContext"},
		BonusIconSubjectCatalog(),
		BonusIconSubject.bonusIconNoun,
	).gated()
}

// BonusIconSubjectCatalog returns one representative value of every
// BonusIconSubject, with the rulebook term each owes.
func BonusIconSubjectCatalog() []Catalogued[BonusIconSubject] {
	return []Catalogued[BonusIconSubject]{
		{
			Node:  TheCardInContext{Noun: DiscardedCard},
			Rules: bears("Resolve Bonus Icons"),
		},
		{
			Node:  ThePurgedCards{},
			Rules: bears("Resolve Bonus Icons"),
		},
	}
}
