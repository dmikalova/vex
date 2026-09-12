package engine

import "fmt"

// This file holds the log entries that narrate a lasting "remainder of the turn"
// effect paying out (ADR 0007). A lasting Æmber gain records the plain gain entry
// under its source-card frame; LastingDraw carries the Event that set it off,
// because a player watching cards appear wants to know which standing effect drew
// them.

// LastingDraw narrates cards a lasting reaction drew.
type LastingDraw struct {
	Player int
	Amount int
	On     Event
}

// Text renders the cards a lasting reaction drew, and its event.
func (e LastingDraw) Text(n Namer) string {
	drew := fmt.Sprintf("%s draws %s", n.PlayerName(e.Player), countNoun(e.Amount, "card"))
	return because(drew, e.On)
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
