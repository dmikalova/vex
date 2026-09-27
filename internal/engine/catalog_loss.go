package engine

// This file is the census of the Loss family: one row per way an effect says how
// much Æmber to remove when the amount depends on the pool's current size. See
// catalog.go for the classification rule, and effect_aember.go for the interface.

// lossPossessive is the possessive the census renders each loss's object phrase
// against, standing in for the player a card would name.
const lossPossessive = "their"

// lossFamily is the Loss family's census entry. A loss is discovered by its lose
// method.
func lossFamily() Family {
	return newFamily(
		"Loss",
		"lose",
		[]string{"int"},
		LossCatalog(),
		func(l Loss) string { return l.object(lossPossessive) },
	).gated()
}

// LossCatalog returns one representative value of every Loss, with the rulebook
// term each owes. A Fraction is the odd one out: it is a Loss, but it is also the
// whole game's "portion of N" — a share of a power or of a battleline count — and
// its explicit rounding is a rule of its own rather than part of losing Æmber.
func LossCatalog() []Catalogued[Loss] {
	return []Catalogued[Loss]{
		{
			Node:  AllAember,
			Rules: bears("Lose Æmber"),
		},
		{
			Node:  AllBut(6),
			Rules: bears("Lose Æmber"),
		},
		{
			Node:  HalfRoundedDown,
			Rules: bears("Portion"),
		},
	}
}
