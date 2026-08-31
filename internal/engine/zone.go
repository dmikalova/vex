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
)

// valid reports whether the zone is a real one (not the unset zero value).
func (z Zone) valid() bool { return z != zoneUnset }
