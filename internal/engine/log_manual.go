package engine

import (
	"fmt"
	"strings"
)

// This file holds the log entries for outcomes that belong to no other family:
// the turn-scoped permissions and restrictions cards hand out, chains, and the
// direct edits manual mode makes to a match.

// MayPlayOrUseGranted narrates a this-turn grant to act with cards outside the
// active house — the one log record for the whole out-of-house permission family
// (ADR 0037). Its wording narrows to the axes the grant selects: fighting,
// using, or playing a named or every house's cards.
type MayPlayOrUseGranted struct {
	Player int
	Houses HouseSelector
	Grant  HouseGrant
	Types  CardTypes
	Count  int
}

// Text renders the grant, narrowing to the houses and verbs it frees.
func (e MayPlayOrUseGranted) Text(n Namer) string {
	p := n.PlayerName(e.Player)
	if e.Houses.Controlled {
		return fmt.Sprintf("%s may play cards from other houses this turn", p)
	}
	switch e.Houses.Match.Kind {
	case MatchExceptHouse:
		return fmt.Sprintf("%s may play cards from other houses this turn", p)
	case MatchAnyHouse:
		if e.Grant&GrantFight != 0 {
			return fmt.Sprintf("%s's creatures may all fight this turn", p)
		}
		return fmt.Sprintf("%s may use friendly artifacts this turn", p)
	default: // MatchNamedHouse, MatchChosenHouse
		h := e.Houses.Match.House
		if e.Grant == GrantFight {
			return fmt.Sprintf("%s's %s creatures may fight this turn", p, h)
		}
		switch {
		case e.Grant&GrantPlay != 0 && e.Grant&GrantUse != 0:
			return fmt.Sprintf("%s may play or use %s cards this turn", p, h)
		case e.Grant&GrantPlay != 0:
			return fmt.Sprintf("%s may play %s cards this turn", p, h)
		default:
			return fmt.Sprintf("%s may use %s creatures this turn", p, h)
		}
	}
}

// MayUseTraitGranted narrates a this-turn grant to fully use friendly creatures of
// a trait even outside the active house — Mutagenic Serum.
type MayUseTraitGranted struct {
	Player int
	Trait  Trait
}

// Text renders the trait whose creatures the player may use this turn.
func (e MayUseTraitGranted) Text(n Namer) string {
	return fmt.Sprintf("%s may use friendly %s creatures this turn",
		n.PlayerName(e.Player), e.Trait)
}

// HouseForcedNextTurn narrates a card dictating next turn's active house.
type HouseForcedNextTurn struct {
	Player int
	House  House
}

// Text renders the house a player must choose next turn.
func (e HouseForcedNextTurn) Text(n Namer) string {
	return fmt.Sprintf("%s must choose house %s next turn", n.PlayerName(e.Player), e.House)
}

// HouseForbiddenNextTurn narrates a card barring a house from next turn's choice.
type HouseForbiddenNextTurn struct {
	Player int
	House  House
}

// Text renders the house a player cannot choose next turn.
func (e HouseForbiddenNextTurn) Text(n Namer) string {
	return fmt.Sprintf("%s cannot choose house %s next turn", n.PlayerName(e.Player), e.House)
}

// HouseWagerArmed narrates a card betting on a player's next active house.
type HouseWagerArmed struct {
	Predictor int
	Player    int
	House     House
	Amount    int
}

// Text renders who stands to steal if the player chooses the named house next turn.
func (e HouseWagerArmed) Text(n Namer) string {
	return fmt.Sprintf("%s steals %d Æmber if %s chooses house %s next turn",
		subject(n, e.Predictor), e.Amount, n.PlayerName(e.Player), e.House)
}

// KeywordLostByAll narrates a keyword switched off across the whole board for
// the rest of the turn.
type KeywordLostByAll struct{ Keyword Keyword }

// Text renders the keyword every creature lost for the turn.
func (e KeywordLostByAll) Text(Namer) string {
	return fmt.Sprintf("each creature loses %s for the remainder of the turn",
		strings.ToLower(e.Keyword.String()))
}

// ChainsGained narrates chains being put on a player.
type ChainsGained struct {
	Player int
	Amount int
	Total  int
}

// Text renders the chains put on a player, and their total. Under a card ability
// the card that imposed them is named too, as a pool change is (ADR 0011).
func (e ChainsGained) Text(n Namer) string {
	if s, ok := framedSource(n); ok {
		return fmt.Sprintf("%s has %s gain %d %s (%d total)",
			s, n.PlayerName(e.Player), e.Amount, chainNoun(e.Amount), e.Total)
	}
	return fmt.Sprintf("%s gains %d %s (%d total)",
		n.PlayerName(e.Player), e.Amount, chainNoun(e.Amount), e.Total)
}

// ManualCardMoved narrates manual mode putting a card in a zone directly.
type ManualCardMoved struct {
	Player int
	Card   LocalID
	To     ManualZone
}

// Text renders the zone manual mode put a card in.
func (e ManualCardMoved) Text(n Namer) string {
	return fmt.Sprintf("%s manually moves %s to %s",
		n.PlayerName(e.Player), n.Name(e.Card), e.To)
}

// ManualExhaustSet narrates manual mode turning a card sideways or upright.
type ManualExhaustSet struct {
	Card      LocalID
	Exhausted bool
}

// Text renders a card manual mode exhausted or readied.
func (e ManualExhaustSet) Text(n Namer) string {
	if e.Exhausted {
		return fmt.Sprintf("%s is manually exhausted", n.Name(e.Card))
	}
	return fmt.Sprintf("%s is manually readied", n.Name(e.Card))
}

// ManualPlacedInPlay narrates manual mode dropping a card straight into play.
type ManualPlacedInPlay struct {
	Player int
	Card   LocalID
}

// Text renders the card manual mode put into play.
func (e ManualPlacedInPlay) Text(n Namer) string {
	return fmt.Sprintf("%s manually puts %s into play", n.PlayerName(e.Player), n.Name(e.Card))
}

// ManualMatchFull narrates a card manual mode could not add, because a match's
// id space is finite.
type ManualMatchFull struct{ Player int }

// Text renders a card manual mode could not add, the match being full.
func (e ManualMatchFull) Text(n Namer) string {
	return fmt.Sprintf("%s cannot add a card: this match is full", n.PlayerName(e.Player))
}

// ManualCardAdded narrates manual mode conjuring a card into a hand.
type ManualCardAdded struct {
	Player int
	Card   LocalID
}

// Text renders the card manual mode conjured into a hand.
func (e ManualCardAdded) Text(n Namer) string {
	return fmt.Sprintf("%s manually adds %s to hand", n.PlayerName(e.Player), n.Name(e.Card))
}

// ManualAemberSet narrates manual mode dialling a pool to a number.
type ManualAemberSet struct {
	Player int
	Amount int
}

// Text renders the pool manual mode dialled a player to.
func (e ManualAemberSet) Text(n Namer) string {
	return fmt.Sprintf("%s now has %d Æmber (manual)", n.PlayerName(e.Player), e.Amount)
}

// ManualChainsSet narrates manual mode dialling a chain count to a number.
type ManualChainsSet struct {
	Player int
	Amount int
}

// Text renders the chain count manual mode dialled a player to.
func (e ManualChainsSet) Text(n Namer) string {
	return fmt.Sprintf("%s now has %d %s (manual)",
		n.PlayerName(e.Player), e.Amount, chainNoun(e.Amount))
}

// ManualHouseChosen narrates manual mode switching the active house mid-turn.
type ManualHouseChosen struct {
	Player int
	House  House
}

// Text renders the active house manual mode switched to.
func (e ManualHouseChosen) Text(n Namer) string {
	return fmt.Sprintf("%s manually chooses %s as their active house",
		n.PlayerName(e.Player), e.House)
}

// ManualKeyForged narrates manual mode adding a key.
type ManualKeyForged struct {
	Player int
	Color  KeyColor
	Keys   int
	Needed int
}

// Text renders the key manual mode forged, with its colour.
func (e ManualKeyForged) Text(n Namer) string {
	return fmt.Sprintf("%s manually forges a %s key (%d/%d)",
		n.PlayerName(e.Player), e.Color, e.Keys, e.Needed)
}

// ManualKeyUnforged narrates manual mode taking a key back.
type ManualKeyUnforged struct {
	Player int
	Keys   int
	Needed int
}

// Text renders the key manual mode took back.
func (e ManualKeyUnforged) Text(n Namer) string {
	return fmt.Sprintf("%s manually unforges a key (%d/%d)",
		n.PlayerName(e.Player), e.Keys, e.Needed)
}
