package engine

import "fmt"

// CardArchivedFromPurge narrates a card recovered from the purge pile into
// archives — the reverse of a purge, where the card was public in the purge pile.
type CardArchivedFromPurge struct {
	Player int
	Card   LocalID
}

// Text renders a card archived out of the purge pile.
func (e CardArchivedFromPurge) Text(n Namer) string {
	return fmt.Sprintf("%s archives %s from their purge pile",
		n.PlayerName(e.Player), nameMoved(n, e.Card, purged, Archives))
}
