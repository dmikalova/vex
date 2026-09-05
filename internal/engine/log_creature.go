package engine

import (
	"fmt"
	"strings"
)

// This file holds the log entries that narrate what happens to a creature or
// artifact in play (ADR 0011): readying, exhausting, stunning, damage, armor,
// destruction, and changes of control. An entry names the card by id, so a
// reader can link it to the card face without matching its name in prose.

// CreatureReadied narrates a card turning upright again.
type CreatureReadied struct{ Creature LocalID }

// Text renders the card that was readied.
func (e CreatureReadied) Text(n Namer) string {
	return fmt.Sprintf("%s is readied", n.Name(e.Creature))
}

// CreatureGainedKeyword narrates a creature gaining a keyword for the turn.
type CreatureGainedKeyword struct {
	Creature LocalID
	Keyword  Keyword
}

// Text renders the creature and the keyword it gained.
func (e CreatureGainedKeyword) Text(n Namer) string {
	return fmt.Sprintf("%s gains %s", n.Name(e.Creature), strings.ToLower(e.Keyword.String()))
}

// CreatureConsideredFlank narrates a creature being treated as a flank creature
// for the turn (Spectral Tunneler).
type CreatureConsideredFlank struct{ Creature LocalID }

// Text renders the creature now considered a flank creature.
func (e CreatureConsideredFlank) Text(n Namer) string {
	return fmt.Sprintf("%s is considered a flank creature", n.Name(e.Creature))
}

// CreatureExhausted narrates a card being turned sideways.
type CreatureExhausted struct{ Creature LocalID }

// Text renders the card that was exhausted.
func (e CreatureExhausted) Text(n Namer) string {
	return fmt.Sprintf("%s is exhausted", n.Name(e.Creature))
}

// CreatureStunned narrates a card being stunned, and by what — unless the
// stunner is the card itself (Chuff Ape enters play stunned), which reads better
// left passive. AlreadyStunned marks a stun that found its target already
// stunned: the source still had to choose it, so the choice is worth a line even
// though nothing changed (the same reasoning as a steal that finds an empty
// pool: see ReapedStealing's "no Æmber to steal").
type CreatureStunned struct {
	Creature       LocalID
	By             LocalID
	AlreadyStunned bool
}

// Text renders the card that was stunned and, when it is not self-inflicted, the
// card that stunned it.
func (e CreatureStunned) Text(n Namer) string {
	if e.AlreadyStunned {
		return fmt.Sprintf("%s is already stunned", n.Name(e.Creature))
	}
	if e.By != e.Creature {
		return fmt.Sprintf("%s stunned %s", n.Name(e.By), n.Name(e.Creature))
	}
	return fmt.Sprintf("%s is stunned", n.Name(e.Creature))
}

// NoCreatureToFight narrates a fight that found no enemy creature to attack, so
// nothing happened.
type NoCreatureToFight struct{ Creature LocalID }

// Text renders a fight that found no enemy creature to attack.
func (e NoCreatureToFight) Text(n Namer) string {
	return fmt.Sprintf("%s has no creature to fight", n.Name(e.Creature))
}

// CardsRevealedToAll narrates cards shown to both players. It is the one entry
// allowed to name a card that is otherwise hidden, because revealing it is
// exactly what made it public (ADR 0011).
type CardsRevealedToAll struct {
	Player int
	Cards  []LocalID
}

// Text renders the cards a player revealed, each by name.
func (e CardsRevealedToAll) Text(n Namer) string {
	return fmt.Sprintf("%s reveals %s", n.PlayerName(e.Player), namedCards(n, e.Cards))
}

// PositionsSwapped narrates two creatures trading places in a battleline.
type PositionsSwapped struct{ A, B LocalID }

// Text renders the two creatures that traded places.
func (e PositionsSwapped) Text(n Namer) string {
	return fmt.Sprintf("%s swaps positions with %s", n.Name(e.A), n.Name(e.B))
}

// MovedToFlank narrates a creature moving to a flank of its battleline.
type MovedToFlank struct {
	Creature LocalID
	Right    bool
}

// Text renders the creature and the flank it moved to.
func (e MovedToFlank) Text(n Namer) string {
	side := "left"
	if e.Right {
		side = "right"
	}
	return fmt.Sprintf("%s moves to the %s flank", n.Name(e.Creature), side)
}

// ControlTaken narrates a card moving into another player's rows without
// changing owner.
type ControlTaken struct {
	Player int
	Card   LocalID
}

// Text renders the card a player took control of.
func (e ControlTaken) Text(n Namer) string {
	return fmt.Sprintf("%s takes control of %s", n.PlayerName(e.Player), n.Name(e.Card))
}

// ControlReturned narrates borrowed control lapsing when its source left play.
type ControlReturned struct {
	Card  LocalID
	Owner int
}

// Text renders borrowed control lapsing back to the card's owner.
func (e ControlReturned) Text(n Namer) string {
	return fmt.Sprintf("%s returns to %s's control", n.Name(e.Card), n.PlayerName(e.Owner))
}

// CardDestroyed narrates a card in play being destroyed.
type CardDestroyed struct{ Card LocalID }

// Text renders the card that was destroyed.
func (e CardDestroyed) Text(n Namer) string {
	return fmt.Sprintf("%s is destroyed", n.Name(e.Card))
}

// CardsDestroyedBy narrates one card's effect destroying a group of cards, so the
// destruction is credited to its agent — "Strange Gizmo destroys A, B, and C" —
// rather than a passive line per creature.
type CardsDestroyedBy struct {
	Source LocalID
	Cards  []LocalID
}

// Text renders the source card and the cards it destroyed.
func (e CardsDestroyedBy) Text(n Namer) string {
	return fmt.Sprintf("%s destroys %s", n.Name(e.Source), namedCardsAnd(n, e.Cards))
}

// DestructionReplaced narrates an upgrade or ability taking a card's destruction
// on itself.
type DestructionReplaced struct {
	Card LocalID
	By   LocalID
}

// Text renders the card that took another's destruction on itself.
func (e DestructionReplaced) Text(n Namer) string {
	return fmt.Sprintf("%s would be destroyed, so %s replaces its destruction",
		n.Name(e.Card), n.Name(e.By))
}

// AemberOnCardReleased narrates the Æmber a card was holding going to a pool
// when the card left play.
type AemberOnCardReleased struct {
	Card   LocalID
	Amount int
	To     int
}

// Text renders the Æmber a departing card released to a pool.
func (e AemberOnCardReleased) Text(n Namer) string {
	return fmt.Sprintf("%d Æmber on %s goes to %s's pool",
		e.Amount, n.Name(e.Card), n.PlayerName(e.To))
}

// StunRecovered narrates a use spent removing a stun counter instead of reaping,
// fighting, or acting. It reads as the controller unstunning the creature, so the
// log speaks the same "unstun" verb as the Unstun effect (Clear Mind).
type StunRecovered struct {
	Player   int
	Creature LocalID
}

// Text renders a use spent recovering from stun.
func (e StunRecovered) Text(n Namer) string {
	return fmt.Sprintf("%s unstuns %s", n.PlayerName(e.Player), n.Name(e.Creature))
}

// CardCannotBeUsed narrates a use refused because the card was exhausted.
type CardCannotBeUsed struct{ Card LocalID }

// Text renders a use refused because the card was exhausted.
func (e CardCannotBeUsed) Text(n Namer) string {
	return fmt.Sprintf("%s is exhausted and cannot be used", n.Name(e.Card))
}
