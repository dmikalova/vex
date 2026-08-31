package card

import "github.com/dmikalova/vactrol/internal/engine"

// Discard is a player's discard pile — the pile a Purge pulls from, e.g.
// card.Purge{Zone: card.Discard}.
const Discard = engine.Discard

// Hand is a player's hand, a zone card.ShuffleIntoDeck can name.
const Hand = engine.Hand

// Archives is a player's archives, a zone card.ShuffleIntoDeck can name.
const Archives = engine.Archives

// Zone names a card pile an effect acts on (see card.Discard).
type Zone = engine.Zone
