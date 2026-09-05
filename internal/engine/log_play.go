package engine

import "fmt"

// This file holds the log entries that narrate a card being played or used
// (ADR 0011): entering a zone, its Æmber bonus, the toll it cost, and the reap,
// fight, or action ability it was spent on.

// CardPlayedToBattleline narrates a creature entering the battleline on a flank,
// or — when a Deploy creature lands between two others — into the battleline.
type CardPlayedToBattleline struct {
	Player    int
	Card      LocalID
	FlankLeft bool
	// Interior marks a Deploy placement that landed between two creatures rather
	// than on a flank, so it narrates as entering the battleline instead.
	Interior bool
}

// Text renders the creature played and the flank it landed on.
func (e CardPlayedToBattleline) Text(n Namer) string {
	if e.Interior {
		return fmt.Sprintf(
			"%s plays %s into their battleline",
			n.PlayerName(e.Player),
			n.Name(e.Card),
		)
	}
	side := "right"
	if e.FlankLeft {
		side = "left"
	}
	return fmt.Sprintf(
		"%s plays %s on their %s flank",
		n.PlayerName(e.Player),
		n.Name(e.Card),
		side,
	)
}

// ArtifactPlayed narrates an artifact entering the artifact row.
type ArtifactPlayed struct {
	Player int
	Card   LocalID
}

// Text renders the artifact a player played.
func (e ArtifactPlayed) Text(n Namer) string {
	return fmt.Sprintf("%s plays artifact %s", n.PlayerName(e.Player), n.Name(e.Card))
}

// ActionPlayed narrates a tactic resolving on its way to the discard pile.
type ActionPlayed struct {
	Player int
	Card   LocalID
}

// Text renders the tactic a player played.
func (e ActionPlayed) Text(n Namer) string {
	return fmt.Sprintf("%s plays action %s", n.PlayerName(e.Player), n.Name(e.Card))
}

// UpgradeAttached narrates an upgrade going onto a creature.
type UpgradeAttached struct {
	Player  int
	Upgrade LocalID
	Host    LocalID
}

// Text renders the upgrade, and the creature it went onto.
func (e UpgradeAttached) Text(n Namer) string {
	return fmt.Sprintf("%s attaches %s to %s",
		n.PlayerName(e.Player), n.Name(e.Upgrade), n.Name(e.Host))
}

// CardPutIntoPlay narrates a card entering play without being played from hand.
type CardPutIntoPlay struct {
	Player int
	Card   LocalID
}

// Text renders a card put into play under a player's control.
func (e CardPutIntoPlay) Text(n Namer) string {
	return fmt.Sprintf("%s puts %s into play under their control",
		n.PlayerName(e.Player), n.Name(e.Card))
}

// AemberBonusGained narrates the Æmber bonus printed on a card being collected.
type AemberBonusGained struct {
	Player int
	Card   LocalID
	Amount int
}

// Text renders the Æmber bonus a card paid its player.
func (e AemberBonusGained) Text(n Namer) string {
	return fmt.Sprintf("%s gains %d Æmber from %s",
		n.PlayerName(e.Player), e.Amount, n.Name(e.Card))
}

// AemberBonusCaptured narrates a printed Æmber bonus captured on its way to the
// pool.
type AemberBonusCaptured struct {
	Creature LocalID
	Card     LocalID
	Amount   int
}

// Text renders a printed Æmber bonus captured on its way to the pool.
func (e AemberBonusCaptured) Text(n Namer) string {
	return fmt.Sprintf("%s captures %d Æmber from %s's bonus",
		n.Name(e.Creature), e.Amount, n.Name(e.Card))
}

// AemberSpentToPlay narrates the Æmber a card's own play requirement cost.
type AemberSpentToPlay struct {
	Player int
	Card   LocalID
	Amount int
}

// Text renders the Æmber a card's own play requirement cost.
func (e AemberSpentToPlay) Text(n Namer) string {
	return fmt.Sprintf("%s loses %d Æmber to play %s",
		n.PlayerName(e.Player), e.Amount, n.Name(e.Card))
}

// TollPaid narrates a toll an opponent's card levied on an action.
type TollPaid struct {
	Player int
	Payee  int
	Amount int
	Action TollAction
}

// Text renders the toll paid, to whom, and the action it bought.
func (e TollPaid) Text(n Namer) string {
	return fmt.Sprintf("%s pays %d Æmber to %s to %s",
		n.PlayerName(e.Player), e.Amount, n.PlayerName(e.Payee), e.Action.phrase())
}

// Reaped narrates a reap that put its Æmber in the pool.
type Reaped struct {
	Player int
	Card   LocalID
}

// Text renders the reap and the Æmber it put in the pool.
func (e Reaped) Text(n Namer) string {
	return fmt.Sprintf("%s reaps with %s (+1 Æmber)", n.PlayerName(e.Player), n.Name(e.Card))
}

// ReapedStealing narrates a reap that a replacement turned into a steal, with
// what it actually took — zero when the opponent's pool was already empty.
type ReapedStealing struct {
	Player int
	Card   LocalID
	Amount int
}

// Text renders the reap, and what the steal actually took.
func (e ReapedStealing) Text(n Namer) string {
	if e.Amount == 0 {
		return fmt.Sprintf("%s reaps with %s (no Æmber to steal)",
			n.PlayerName(e.Player), n.Name(e.Card))
	}
	return fmt.Sprintf("%s reaps with %s, stealing %d Æmber",
		n.PlayerName(e.Player), n.Name(e.Card), e.Amount)
}

// ReapedCaptured narrates a reap whose Æmber a capturing effect intercepted.
type ReapedCaptured struct {
	Player   int
	Card     LocalID
	Creature LocalID
}

// Text renders the reap, and the creature that captured its Æmber.
func (e ReapedCaptured) Text(n Namer) string {
	return fmt.Sprintf("%s reaps with %s, but %s captures the Æmber",
		n.PlayerName(e.Player), n.Name(e.Card), n.Name(e.Creature))
}

// ActionAbilityUsed narrates a card being used for its "Action:" ability.
type ActionAbilityUsed struct {
	Player int
	Card   LocalID
}

// Text renders the card used for its action ability.
func (e ActionAbilityUsed) Text(n Namer) string {
	return fmt.Sprintf("%s uses %s's action ability", n.PlayerName(e.Player), n.Name(e.Card))
}
