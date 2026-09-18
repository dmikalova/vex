package engine

import "fmt"

// This file holds the log entries that narrate a card being played or used
// (ADR 0011): entering a zone, its Æmber bonus, the toll it cost, and the reap,
// fight, or action ability it was spent on.

// CardPlayedToBattleline narrates a creature entering the battleline on a flank,
// or — when a Deploy creature lands between two others — into the battleline.
type CardPlayedToBattleline struct {
	// Player is who played it, Card is the creature, and FlankLeft is which flank it
	// entered on.
	Player    int
	Card      LocalID
	FlankLeft bool
	// Interior marks a Deploy placement that landed between two creatures rather
	// than on a flank, so it narrates as entering the battleline instead.
	Interior bool
}

// Text renders the creature played and the flank it landed on.
func (e CardPlayedToBattleline) Text(n Namer) string {
	who, owner := actorPossessive(n, e.Player)
	if e.Interior {
		return fmt.Sprintf(
			"%s plays %s into %s battleline",
			who,
			n.Name(e.Card),
			owner,
		)
	}
	side := "right"
	if e.FlankLeft {
		side = "left"
	}
	return fmt.Sprintf(
		"%s plays %s on %s %s flank",
		who,
		n.Name(e.Card),
		owner,
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
	return fmt.Sprintf("%s plays artifact %s", subject(n, e.Player), n.Name(e.Card))
}

// TacticPlayed narrates a tactic resolving on its way to the discard pile.
type TacticPlayed struct {
	Player int
	Card   LocalID
}

// Text renders the tactic a player played.
func (e TacticPlayed) Text(n Namer) string {
	return fmt.Sprintf("%s plays tactic %s", subject(n, e.Player), n.Name(e.Card))
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
		subject(n, e.Player), n.Name(e.Upgrade), n.Name(e.Host))
}

// CardPutIntoPlay narrates a card entering play without being played from hand.
type CardPutIntoPlay struct {
	Player int
	Card   LocalID
}

// Text renders a card put into play under a player's control.
func (e CardPutIntoPlay) Text(n Namer) string {
	who, owner := actorPossessive(n, e.Player)
	return fmt.Sprintf("%s puts %s into play under %s control",
		who, n.Name(e.Card), owner)
}

// PlayedFromTopOfDeck narrates a card an ability played off the top of a deck
// (Wild Wormhole), crediting the ability's card from the record's frame and
// naming whose deck it came from, so the following placement line reads as the
// consequence of that play.
type PlayedFromTopOfDeck struct {
	Card   LocalID
	Player int
}

// Text renders the ability, the card it played, and whose deck it came off.
func (e PlayedFromTopOfDeck) Text(n Namer) string {
	who, ok := framedSource(n)
	if !ok {
		who = n.PlayerName(e.Player)
	}
	return fmt.Sprintf("%s plays %s from the top of %s's deck",
		who, n.Name(e.Card), n.PlayerName(e.Player))
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

// Reaped narrates a reap that put its Æmber in the pool.
type Reaped struct {
	Player int
	Card   LocalID
}

// Text renders the reap and the Æmber it put in the pool.
func (e Reaped) Text(n Namer) string {
	return fmt.Sprintf("%s reaps with %s (+1 Æmber)", n.PlayerName(e.Player), n.Name(e.Card))
}

// ReapedStealing narrates a reap that a replacement turned into a steal (Dimension
// Door), with what it actually took — zero when the opponent's pool was already
// empty. Cause is the card whose replacement turned the gain into a steal.
type ReapedStealing struct {
	Player int
	Card   LocalID
	Amount int
	Cause  LocalID
}

// Text renders the reap, and what the steal actually took.
func (e ReapedStealing) Text(n Namer) string {
	return replacementLine(
		n.Name(e.Cause),
		n.PlayerName(e.Player),
		"steal",
		fmt.Sprintf("%d Æmber reaping with %s", e.Amount, n.Name(e.Card)),
		"gaining it",
	)
}

// ReapedCaptured narrates a reap whose Æmber a capturing effect intercepted. The
// intercepting creature carries the replacement itself (Ether Spider), so it is
// both the cause and the actor.
type ReapedCaptured struct {
	Player   int
	Card     LocalID
	Creature LocalID
}

// Text renders the reap, and the creature that captured its Æmber.
func (e ReapedCaptured) Text(n Namer) string {
	return replacementLine(
		n.Name(e.Creature),
		"",
		"capture",
		fmt.Sprintf("1 Æmber reaping with %s", n.Name(e.Card)),
		fmt.Sprintf("%s gaining it", n.PlayerName(e.Player)),
	)
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
