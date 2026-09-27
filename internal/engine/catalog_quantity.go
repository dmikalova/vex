package engine

// This file is the census of the Quantity family: one row per way a zone verb
// says how many cards it takes. See catalog.go for the classification rule, and
// quantity.go for the quantities.
//
// A quantity is a sentence part, but it is not only a number: "up to 2 cards" and
// "any number of cards" both let the controller stop early, which is a rule a
// player must be taught. All three bind to the one "Quantity" umbrella term
// rather than each earning a title.

// censusQuantityNoun and censusQuantitySingle are the Selection phrases the
// census renders each quantity's object against, standing in for the selection a
// card would pair it with.
const (
	censusQuantityNoun   = "card"
	censusQuantitySingle = "a card"
)

// quantityFamily is the Quantity family's census entry. A quantity is discovered
// by its picks method.
func quantityFamily() Family {
	return newFamily(
		"Quantity",
		"picks",
		[]string{"*EffectContext"},
		QuantityCatalog(),
		func(q Quantity) string { return q.object(censusQuantityNoun, censusQuantitySingle) },
	).gated()
}

// QuantityCatalog returns one representative value of every Quantity, with the
// rulebook term each owes.
func QuantityCatalog() []Catalogued[Quantity] {
	return []Catalogued[Quantity]{
		{
			Node:  Takes{N: Fixed(2)},
			Rules: bears("Quantity"),
		},
		{
			Node:  UpTo{N: Fixed(2)},
			Rules: bears("Quantity"),
		},
		{
			Node:  AnyNumber{},
			Rules: bears("Quantity"),
		},
	}
}
