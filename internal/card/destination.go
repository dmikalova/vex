package card

import "github.com/dmikalova/vactrol/internal/engine"

// To groups the destinations an effect can put a card, e.g. card.To.TopOfDeck. An
// effect that moves a card takes one — see card.PutCard's Destination.
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
	// Archives puts the card into its owner's archives. Call Yours() on it for an
	// abduction into the archives of the player resolving the effect.
	Archives engine.Destination
}

// Destination names where an effect puts a card it moves (see card.To).
type Destination = engine.Destination

// Into groups the destinations a card.ChooseAndMove step sends the cards it takes
// off the top of a deck, e.g. card.Into.Purge. It mirrors the engine's DeckDest.
var Into = deckDests{
	Hand:         engine.IntoHand,
	Archives:     engine.IntoArchives,
	Discard:      engine.IntoDiscard,
	Purge:        engine.IntoPurge,
	BottomOfDeck: engine.IntoBottomOfDeck,
}

type deckDests struct {
	// Hand puts the chosen cards into your hand.
	Hand engine.DeckDest
	// Archives archives the chosen cards.
	Archives engine.DeckDest
	// Discard puts the chosen cards into the discard pile.
	Discard engine.DeckDest
	// Purge purges the chosen cards.
	Purge engine.DeckDest
	// BottomOfDeck puts the chosen cards on the bottom of your deck.
	BottomOfDeck engine.DeckDest
}
