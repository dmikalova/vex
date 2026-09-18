package engine

// Zone names one of a player's card zones, so an effect can say which zone it
// acts on (Purge). It has no valid zero value: an effect must name the zone it
// acts on rather than fall back to a default.
type Zone uint8

const (
	// zoneUnset is the invalid zero value: an effect must name its zone.
	zoneUnset Zone = iota
	// Discard is a player's discard pile.
	Discard
	// Hand is a player's hand.
	Hand
	// Archives is a player's archives.
	Archives
	// Deck is a player's deck.
	Deck
	// Purged is the pile a purged card is set aside in, out of the game. A card may
	// name it as a source — Universal Recycle Bin archives a card from it — but never
	// as a destination: setting a card aside is written as the Purge verb, which
	// moves through the unexported toPurged destination (ADR 0031).
	Purged
	// InPlay is every card someone controls in play — creatures, artifacts, and the
	// upgrades on them. A card under another card is not in play; it is in the
	// out-of-play zone that is "under" its host. Like Purged, a card may name InPlay
	// only as a source — putting a card into play is its own verb.
	InPlay
)

// noun names the zone as printed card text says it, so every effect that has to
// name a zone ("shuffle your discard pile into your deck", "play a creature from
// your discard pile") phrases it the same way.
func (z Zone) noun() string {
	switch z {
	case Hand:
		return "hand"
	case Archives:
		return "archives"
	case Deck:
		return "deck"
	case Purged:
		return "purge pile"
	case InPlay:
		return "play"
	default: // Discard
		return "discard pile"
	}
}

// public reports whether both players can see a card sitting in the zone. The
// log names a card only where it is public (ADR 0011): a discard pile, the board,
// and the purged pile are open, while a hand, archives, and deck are not, so a
// move between two hidden zones is narrated without naming what moved.
func (z Zone) public() bool {
	return z == Discard || z == InPlay || z == Purged
}

// ordered reports whether the zone is a stack with a top and a bottom, so a
// positional selection (Top / Bottom) can name an end of it. Only the deck and
// the discard pile are ordered; a hand, archives, and the purge pile are not.
func (z Zone) ordered() bool { return z == Deck || z == Discard }
