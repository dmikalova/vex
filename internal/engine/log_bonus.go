package engine

import "fmt"

// The records below narrate a played card's bonus icons resolving. Each names the
// card carrying the icon as the source (a Vactrol divergence: KeyForge treats the
// game itself as the source), and each says "bonus" so the log distinguishes a
// bonus icon from the same effect produced by an ability.

// BonusAemberGained narrates an Æmber bonus icon gaining its player 1 Æmber.
type BonusAemberGained struct {
	Player int
	Card   LocalID
	Amount int
}

// Text renders the bonus Æmber a card gained its player.
func (e BonusAemberGained) Text(n Namer) string {
	return fmt.Sprintf("%s gains %d bonus Æmber from %s",
		n.PlayerName(e.Player), e.Amount, n.Name(e.Card))
}

// BonusAemberCaptured narrates a bonus Æmber gain intercepted and captured on its
// way to the pool (Ether Spider).
type BonusAemberCaptured struct {
	Creature LocalID
	Card     LocalID
	Amount   int
}

// Text renders the intercepted bonus Æmber.
func (e BonusAemberCaptured) Text(n Namer) string {
	return fmt.Sprintf("%s captures %d bonus Æmber from %s",
		n.Name(e.Creature), e.Amount, n.Name(e.Card))
}

// BonusCaptured narrates a Capture bonus icon: a friendly creature captures Æmber
// from the opponent.
type BonusCaptured struct {
	Creature LocalID
	Card     LocalID
	Amount   int
}

// Text renders the bonus capture, naming the capturing creature and the source card.
func (e BonusCaptured) Text(n Namer) string {
	return fmt.Sprintf("%s captures %d bonus Æmber (%s)",
		n.Name(e.Creature), e.Amount, n.Name(e.Card))
}

// BonusDamageDealt narrates a Damage bonus icon dealing damage to a creature,
// naming the card that carries the icon as the source.
type BonusDamageDealt struct {
	Source LocalID
	Amount int
	Target LocalID
}

// Text renders the bonus damage.
func (e BonusDamageDealt) Text(n Namer) string {
	return fmt.Sprintf("%s deals %d bonus damage to %s",
		n.Name(e.Source), e.Amount, n.Name(e.Target))
}

// BonusCardDrawn narrates a Draw bonus icon drawing its player a card.
type BonusCardDrawn struct {
	Player int
	Card   LocalID
	Amount int
}

// Text renders the bonus draw.
func (e BonusCardDrawn) Text(n Namer) string {
	return fmt.Sprintf("%s draws %d bonus card from %s",
		n.PlayerName(e.Player), e.Amount, n.Name(e.Card))
}
