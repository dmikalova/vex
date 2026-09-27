package engine

// This file is the census of Destination, which — like Target — is covered as
// both halves of what it is. A Destination is a zone constant plus the modifiers
// a card writes on it, so its constants are covered as an Enum, discovered as the
// constants of the zone type they are declared with, and its modifiers as a
// Family discovered by shape, every exported method on Destination returning a
// Destination. An added zone and an added modifier are both a red build until
// they have a row.

// destinationSubject is the noun the census renders each destination's move
// around, standing in for the cards a card would name.
const destinationSubject = "a card"

// destinationEnum is the Destination enum's census, keyed by the zone half its
// constants are declared with. Each destination binds to the term for the zone
// rule it carries out rather than to a term of its own.
func destinationEnum() Enum {
	return newEnum(
		"destinationZone",
		[]string{"destUnset", "destDiscard", "destPurged"},
		len(Destinations()),
		[]Catalogued[Destination]{
			{Name: "destHand", Node: ToHand, Rules: bears("Put a Card into Another Zone")},
			{
				Name:  "destTopOfDeck",
				Node:  ToTopOfDeck,
				Rules: bears("Put a Card into Another Zone"),
			},
			{
				Name:  "destBottomOfDeck",
				Node:  ToBottomOfDeck,
				Rules: bears("Put a Card into Another Zone"),
			},
			{Name: "destDeckShuffled", Node: ToDeckShuffled, Rules: bears("Shuffle")},
			{Name: "destArchives", Node: ToArchives, Rules: bears("Archive")},
		},
		func(d Destination) string { return d.clause(destinationSubject, false) },
	)
}

// destinationFamily is the Destination modifiers' census entry. Yours is the only
// one: it redirects a move into the resolving player's own copy of the zone
// instead of the card owner's, which is the abduction rule.
func destinationFamily() Family {
	return newFamily(
		"Destination",
		"",
		nil,
		[]Catalogued[Destination]{
			{Name: "Yours", Node: ToArchives.Yours(), Rules: bears("Abduct")},
		},
		func(d Destination) string { return d.clause(destinationSubject, false) },
	).builds("Destination").gated()
}
