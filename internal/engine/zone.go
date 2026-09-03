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
	// The zones below are exported to no one: a card may not name them, but the
	// engine still has to place a move's endpoints to know what the log may say.

	// deck is a player's deck.
	deck
	// inPlay is the board — battleline, artifacts, and the upgrades on them.
	inPlay
	// purged is the pile a purged card is set aside in, out of the game.
	purged
)

// valid reports whether the zone is a real one (not the unset zero value).
func (z Zone) valid() bool { return z != zoneUnset }

// public reports whether both players can see a card sitting in the zone. The
// log names a card only where it is public (ADR 0011): a discard pile, the board,
// and the purged pile are open, while a hand, archives, and deck are not, so a
// move between two hidden zones is narrated without naming what moved.
func (z Zone) public() bool { return z == Discard || z == inPlay || z == purged }
