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
	return because(AemberGained{Player: e.Player, Amount: e.Amount}.Text(n), e.On)
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
	base := AemberCapturedInsteadOfGain{
		Creature: e.Creature,
		Player:   e.Player,
		Amount:   e.Amount,
	}
	return because(base.Text(n), e.On)
}

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
