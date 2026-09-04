package engine

import (
	"fmt"
)

// This file holds the log entries that narrate a turn's shape (ADR 0011): its
// boundaries, its phases, the house chosen for it, and the keys forged in it.
// The client demarcates turns and phases from these types rather than by
// matching a prefix on a line of prose.

// TurnBegan narrates the start of a player's turn.
type TurnBegan struct {
	Player int
	Turn   int
}

// Text renders the turn a player is beginning, by number.
func (e TurnBegan) Text(n Namer) string {
	return fmt.Sprintf("%s begins turn %d", n.PlayerName(e.Player), e.Turn)
}

// PhaseBegan narrates entering one of a turn's phases, so the log can be grouped
// by phase (ADR 0012). It is recorded for every phase, including one that turns
// out to do nothing, so a turn where the player plays nothing still shows a main
// phase.
type PhaseBegan struct {
	Player int
	Phase  Phase
}

// Text renders the phase that has been entered. It leaves the player unsaid: a
// phase sits inside a turn the log already opened by name, so repeating it on
// every phase says nothing the turn header did not.
func (e PhaseBegan) Text(Namer) string {
	return capitalizeFirst(e.Phase.String()) + " phase"
}

// CardsReadied narrates the ready phase, naming the cards that were turned back
// upright. A phase that readied nothing records nothing.
type CardsReadied struct {
	Player int
	Cards  []LocalID
}

// Text renders the cards a player readied, each by name.
func (e CardsReadied) Text(n Namer) string {
	return fmt.Sprintf("%s readies %s", n.PlayerName(e.Player), namedCards(n, e.Cards))
}

// CardsDrawn narrates the end-of-turn refill: how many cards the player took and
// the hand size it brought them to, so a draw that came up short against an empty
// deck reads as one.
type CardsDrawn struct {
	Player int
	Count  int
	Hand   int
}

// Text renders the refill, or the hand it stood at when nothing was drawn.
func (e CardsDrawn) Text(n Namer) string {
	if e.Count == 0 {
		return fmt.Sprintf("%s draws nothing, holding %d", n.PlayerName(e.Player), e.Hand)
	}
	return fmt.Sprintf("%s draws %s, up to %d in hand",
		n.PlayerName(e.Player), countNoun(e.Count, "card"), e.Hand)
}

// HouseChosen narrates the active house a player picked for the turn.
type HouseChosen struct {
	Player int
	House  House
}

// Text renders the house a player chose for the turn.
func (e HouseChosen) Text(n Namer) string {
	return fmt.Sprintf("%s chooses house %s", n.PlayerName(e.Player), e.House)
}

// ForgeSkipped narrates a forge phase an effect made the player sit out (Miasma).
type ForgeSkipped struct{ Player int }

// Text renders the forge phase a player had to sit out.
func (e ForgeSkipped) Text(n Namer) string {
	return fmt.Sprintf("%s skips their forge a key phase", n.PlayerName(e.Player))
}

// KeyForged narrates a forged key: its colour when the player picked one, and
// where it puts them on the way to winning.
type KeyForged struct {
	Player   int
	Color    KeyColor
	HasColor bool
	Keys     int
	Needed   int
}

// Text renders the forged key, naming its colour when it has one.
func (e KeyForged) Text(n Namer) string {
	if e.HasColor {
		return fmt.Sprintf("%s forges a %s key (%d/%d)",
			n.PlayerName(e.Player), e.Color, e.Keys, e.Needed)
	}
	return fmt.Sprintf("%s forges a key (%d/%d)", n.PlayerName(e.Player), e.Keys, e.Needed)
}

// KeyUnforged narrates a key taken back off a player (Key Charge's mirror, and
// manual mode).
type KeyUnforged struct {
	Player int
	Keys   int
	Needed int
}

// Text renders a key taken back, and the count it leaves behind.
func (e KeyUnforged) Text(n Namer) string {
	return fmt.Sprintf("%s unforges a key (%d/%d)", n.PlayerName(e.Player), e.Keys, e.Needed)
}

// ChainShed narrates a chain coming off because it actually cost the player a
// draw this turn.
type ChainShed struct {
	Player    int
	Remaining int
}

// Text renders a chain coming off, and how many are left.
func (e ChainShed) Text(n Namer) string {
	return fmt.Sprintf("%s sheds a chain (%d remaining)", n.PlayerName(e.Player), e.Remaining)
}

// GameWon narrates the third key.
type GameWon struct{ Player int }

// Text renders the player who forged their third key.
func (e GameWon) Text(n Namer) string {
	return fmt.Sprintf("%s wins the game!", n.PlayerName(e.Player))
}

// PlayerStanding narrates where a player stands as a turn ends. It states only
// facts both players can already see, so a client may show it for both.
// KeyColors holds the colour of each key forged so far, in forge order, so a
// client can draw the actual coloured keys instead of just a count.
type PlayerStanding struct {
	Player    int
	Aember    int
	KeyColors []KeyColor
}

// Text renders where a player stands.
func (e PlayerStanding) Text(n Namer) string {
	return fmt.Sprintf(
		"%s has %d Æmber and %d keys",
		n.PlayerName(e.Player),
		e.Aember,
		len(e.KeyColors),
	)
}
