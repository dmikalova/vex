package card

import "github.com/dmikalova/vactrol/internal/engine"

// To groups the destinations an effect can put a card, e.g. card.To.TopOfDeck. An
// effect that moves a card takes one — see card.PutFromDiscard's Destination.
// It mirrors the engine's destination.go.
var To = destinations{
	Hand:         engine.ToHand,
	TopOfDeck:    engine.ToTopOfDeck,
	BottomOfDeck: engine.ToBottomOfDeck,
	DeckShuffled: engine.ToDeckShuffled,
	Archives:     engine.ToArchives,
}

type destinations struct {
	// Hand puts the card into its owner's hand.
	Hand engine.Destination
	// TopOfDeck puts the card on top of its owner's deck.
	TopOfDeck engine.Destination
	// BottomOfDeck puts the card on the bottom of its owner's deck.
	BottomOfDeck engine.Destination
	// DeckShuffled puts the card into its owner's deck, then shuffles.
	DeckShuffled engine.Destination
	// Archives puts the card into its owner's archives.
	Archives engine.Destination
}

// Destination names where an effect puts a card it moves (see card.To).
type Destination = engine.Destination
