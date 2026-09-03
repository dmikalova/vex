package engine

import "fmt"

// This file holds the log entries that narrate a lasting "remainder of the turn"
// effect paying out (ADR 0007). They carry the Event that set them off, because
// a player watching Æmber appear needs to know which of the turn's standing
// effects produced it.

// LastingAemberGained narrates Æmber a lasting reaction produced.
type LastingAemberGained struct {
	Player int
	Amount int
	On     Event
}

// Text renders the Æmber a lasting reaction produced, and its event.
func (e LastingAemberGained) Text(n Namer) string {
	return fmt.Sprintf("%s gains %d Æmber (%s)", n.PlayerName(e.Player), e.Amount, e.On.clause())
}

// LastingAemberCaptured narrates a lasting reaction's Æmber being captured on
// the way to the pool instead.
type LastingAemberCaptured struct {
	Creature LocalID
	Player   int
	Amount   int
	On       Event
}

// Text renders a lasting reaction's Æmber captured instead of gained.
func (e LastingAemberCaptured) Text(n Namer) string {
	return fmt.Sprintf("%s captures %d Æmber instead of %s gaining it (%s)",
		n.Name(e.Creature), e.Amount, n.PlayerName(e.Player), e.On.clause())
}

// LastingDraw narrates cards a lasting reaction drew.
type LastingDraw struct {
	Player int
	Amount int
	On     Event
}

// Text renders the cards a lasting reaction drew, and its event.
func (e LastingDraw) Text(n Namer) string {
	return fmt.Sprintf("%s draws %d card(s) (%s)", n.PlayerName(e.Player), e.Amount, e.On.clause())
}

// AemberGivenAfterForging narrates the whole pool changing hands because forging
// a key triggered a standing effect that hands it over.
type AemberGivenAfterForging struct {
	Player int
	To     int
	Amount int
}

// Text renders the Æmber handed to an opponent after forging a key.
func (e AemberGivenAfterForging) Text(n Namer) string {
	return fmt.Sprintf("%s gives %d Æmber to %s after forging a key",
		n.PlayerName(e.Player), e.Amount, n.PlayerName(e.To))
}
