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

// AemberGiven narrates one player giving Æmber to another from their pool — the
// GiveAember effect. A toll sets Reason to the action it charged for (Customs
// Office, Tentacus), so the line names it and reads under the toll card's frame;
// a plain give leaves Reason unset.
type AemberGiven struct {
	Giver    int
	Receiver int
	Amount   int
	Reason   TollAction
}

// Text renders the Æmber one player gives another, naming the card behind it as
// the subject when a frame carries one ("Customs Office has P0 give 1 Æmber to P1
// to use an artifact") and the giver alone otherwise.
func (e AemberGiven) Text(n Namer) string {
	reason := ""
	if e.Reason != tollActionUnset {
		reason = " to " + e.Reason.phrase()
	}
	if s, ok := framedSource(n); ok {
		return fmt.Sprintf("%s has %s give %d Æmber to %s%s",
			s, n.PlayerName(e.Giver), e.Amount, n.PlayerName(e.Receiver), reason)
	}
	return fmt.Sprintf("%s gives %d Æmber to %s%s",
		n.PlayerName(e.Giver), e.Amount, n.PlayerName(e.Receiver), reason)
}
