package engine

// Destination names where an effect puts a card it moves — the KeyForge phrase a
// card uses, e.g. "into your hand" or "on top of your deck". The deck has three
// separate destinations (its top, its bottom, and shuffled in), so "the deck"
// alone is never a destination; a card always names which. An effect takes one to
// say where a card goes, e.g. PutFromDiscard{Destination: ToTopOfDeck}.
type Destination uint8

const (
	// destinationUnset is the invalid zero value: an effect must name where a card
	// goes rather than leave the destination unset.
	destinationUnset Destination = iota
	// ToHand puts the card into its owner's hand.
	ToHand
	// ToTopOfDeck puts the card on top of its owner's deck.
	ToTopOfDeck
	// ToBottomOfDeck puts the card on the bottom of its owner's deck.
	ToBottomOfDeck
	// ToDeckShuffled shuffles the card into its owner's deck.
	ToDeckShuffled
	// ToArchives puts the card into its owner's archives.
	ToArchives
)
