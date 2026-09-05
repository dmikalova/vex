package engine

import "fmt"

// This file holds the log entries that narrate Æmber changing hands (ADR 0011):
// gaining, losing, stealing, capturing, exalting, and moving it between pools
// and cards. Each entry records the amount that actually moved, which is not
// always the amount the card asked for.

// AemberGained narrates Æmber arriving in a player's pool.
type AemberGained struct {
	Player int
	Amount int
}

// Text renders the Æmber a player gained.
func (e AemberGained) Text(n Namer) string {
	return fmt.Sprintf("%s gains %d Æmber", n.PlayerName(e.Player), e.Amount)
}

// AemberLost narrates Æmber leaving a player's pool to the common supply.
type AemberLost struct {
	Player int
	Amount int
}

// Text renders the Æmber a player lost to the common supply.
func (e AemberLost) Text(n Namer) string {
	return fmt.Sprintf("%s loses %d Æmber", n.PlayerName(e.Player), e.Amount)
}

// AemberStolen narrates Æmber moving from one pool to the other. Amount is what
// was actually taken, which is less than asked for when the victim's pool runs
// out — and zero when it was already empty. Source names the card whose ability
// stole, so the line reads from the card's perspective ("Magda the Rat steals 2
// Æmber from Player 1") rather than the controller's; HasSource is false for a
// steal with no card to credit, which falls back to the player.
type AemberStolen struct {
	// Player is who stole, From is who lost, Amount is what was actually taken, and
	// Source is the card credited (see HasSource).
	Player int
	From   int
	Amount int
	Source LocalID
	// HasSource distinguishes a card named by LocalID 0 from no source at all.
	HasSource bool
}

// Text renders the steal, and how much it actually took, crediting the source
// card when there is one and the controller otherwise.
func (e AemberStolen) Text(n Namer) string {
	if e.HasSource {
		return fmt.Sprintf("%s steals %d Æmber from %s",
			n.Name(e.Source), e.Amount, n.PlayerName(e.From))
	}
	return fmt.Sprintf("%s steals %d Æmber from %s",
		n.PlayerName(e.Player), e.Amount, n.PlayerName(e.From))
}

// AemberCaptured narrates Æmber moved onto a creature, where it stays out of
// every pool until the creature leaves play.
type AemberCaptured struct {
	Creature LocalID
	Amount   int
}

// Text renders the Æmber a creature captured.
func (e AemberCaptured) Text(n Namer) string {
	return fmt.Sprintf("%s captures %d Æmber", n.Name(e.Creature), e.Amount)
}

// AemberMovedToCommonSupply narrates Æmber removed from a creature and returned
// to the common supply, the reverse of a capture.
type AemberMovedToCommonSupply struct {
	Creature LocalID
	Amount   int
}

// Text renders the Æmber a creature moved to the common supply.
func (e AemberMovedToCommonSupply) Text(n Namer) string {
	return fmt.Sprintf("%s moves %d Æmber to the common supply", n.Name(e.Creature), e.Amount)
}

// AemberCapturedInsteadOfGain narrates a gain that a capturing effect
// intercepted, so the Æmber landed on a creature rather than in the pool.
type AemberCapturedInsteadOfGain struct {
	Creature LocalID
	Player   int
	Amount   int
}

// Text renders the gain a capturing effect intercepted.
func (e AemberCapturedInsteadOfGain) Text(n Namer) string {
	return fmt.Sprintf("%s captures %d Æmber instead of %s gaining it",
		n.Name(e.Creature), e.Amount, n.PlayerName(e.Player))
}

// AemberExalted narrates Æmber moved from a pool onto a creature as an exalt.
type AemberExalted struct {
	Creature LocalID
	Amount   int
}

// Text renders the Æmber exalted onto a creature.
func (e AemberExalted) Text(n Namer) string {
	return fmt.Sprintf("%s is exalted (%d Æmber placed)", n.Name(e.Creature), e.Amount)
}

// AemberMovedToPool narrates Æmber taken off a card and put into a pool.
type AemberMovedToPool struct {
	Player int
	From   LocalID
	To     int
	Amount int
}

// Text renders Æmber moved off a card into a player's pool.
func (e AemberMovedToPool) Text(n Namer) string {
	return fmt.Sprintf("%s moves %d Æmber from %s to %s's pool",
		n.PlayerName(e.Player), e.Amount, n.Name(e.From), n.PlayerName(e.To))
}

// AemberMovedToCard narrates Æmber moved from one card to another.
type AemberMovedToCard struct {
	Player int
	From   LocalID
	To     LocalID
	Amount int
}

// Text renders Æmber moved from one card to another.
func (e AemberMovedToCard) Text(n Namer) string {
	return fmt.Sprintf("%s moves %d Æmber from %s to %s",
		n.PlayerName(e.Player), e.Amount, n.Name(e.From), n.Name(e.To))
}

// AemberLostToCeiling narrates Æmber that never landed on a card because the
// card was already holding the most a card may hold.
type AemberLostToCeiling struct {
	Card   LocalID
	Amount int
}

// Text renders the Æmber a card was too full to hold.
func (e AemberLostToCeiling) Text(n Namer) string {
	return fmt.Sprintf("%s can hold no more Æmber; %d is lost to the ceiling",
		n.Name(e.Card), e.Amount)
}
