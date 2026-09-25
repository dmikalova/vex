package card

import "github.com/dmikalova/vex/internal/engine"

// Discard is a player's discard pile — the pile a Purge pulls from, e.g.
// card.Purge{Zone: card.Discard}.
const Discard = engine.Discard

// Hand is a player's hand, a zone card.ShuffleIntoDeck can name.
const Hand = engine.Hand

// Archives is a player's archives, a zone card.ShuffleIntoDeck can name.
const Archives = engine.Archives

// Deck is a player's deck — the pile a positional archive can take the top of,
// e.g. card.ArchiveCard{Zone: card.Deck, Selection: card.Top{}}.
const Deck = engine.Deck

// Purged is the pile a purged card is set aside in, which a card may name only as
// a source — an archive out of it, e.g. card.ArchiveCard{Zone: card.Purged}.
// Setting a card aside is written as card.Purge, never as a destination.
const Purged = engine.Purged

// Battleline is the creature row, which a card may name only as a source — a
// shuffle out of it, e.g.
// card.ShuffleIntoDeck{From: []card.Zone{card.Battleline}}. Putting a card into
// play is its own verb, never a destination. A card that also reaches artifacts
// or upgrades names the wider engine.InPlay instead.
const InPlay = engine.InPlay

// Zone names a card pile an effect acts on (see card.Discard).
type Zone = engine.Zone
