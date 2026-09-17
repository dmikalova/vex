package engine

import "fmt"

// The records below narrate a played card's bonus icons resolving. The card
// carrying the icon leads as the source (a Vactrol divergence: KeyForge treats the
// game itself as the source), and the word "bonus" marks every line — so the log
// distinguishes a bonus icon from the same effect produced by an ability. Most
// lines read "<card> bonus <verb> …"; the damage line reads "<card> deals N bonus
// damage to …", where "bonus damage" is the natural noun phrase.

// BonusAemberGained narrates an Æmber bonus icon gaining its player 1 Æmber.
type BonusAemberGained struct {
	Player int
	Card   LocalID
	Amount int
}

// Text renders the bonus Æmber a card gained its player.
func (e BonusAemberGained) Text(n Namer) string {
	return fmt.Sprintf("%s bonus gains %d Æmber for %s",
		n.Name(e.Card), e.Amount, n.PlayerName(e.Player))
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
	return fmt.Sprintf("%s bonus captures %d Æmber onto %s",
		n.Name(e.Card), e.Amount, n.Name(e.Creature))
}

// BonusCaptured narrates a Capture bonus icon: a friendly creature captures Æmber
// from the opponent.
type BonusCaptured struct {
	Creature LocalID
	Card     LocalID
	Amount   int
}

// Text renders the bonus capture, naming the source card and the capturing creature.
func (e BonusCaptured) Text(n Namer) string {
	return fmt.Sprintf("%s bonus captures %d Æmber onto %s",
		n.Name(e.Card), e.Amount, n.Name(e.Creature))
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
	return fmt.Sprintf("%s bonus draws %d card for %s",
		n.Name(e.Card), e.Amount, n.PlayerName(e.Player))
}
