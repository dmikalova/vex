package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Gambling Den
//
//	House:  Shadows
//	Type:   Artifact
//	Rarity: Rare
//	Æmber:  1
//	Traits: Location
//
//	At the start of each player's turn, you may choose a house - reveal the top card of your deck. If it is of the chosen house, gain 2 Æmber. Otherwise, lose 2 Æmber.

// TestGamblingDenGainOnMatch fires the artifact at the start of the opponent's
// turn: they name the house their revealed top card belongs to and gain 2 Æmber.
// The gain lands on the acting player (P2), not on the artifact's controller (P1).
func TestGamblingDenGainOnMatch(t *testing.T) {
	top := ct.Creature(ct.OfHouse(card.House.Mars))
	h := ct.Play(t, ct.Setup{
		P1: ct.Side{House: card.House.Shadows, InPlay: ct.Cards(GamblingDen)},
		P2: ct.Side{House: card.House.Mars, Deck: ct.Cards(top)},
	})

	h.P1.EndTurn() // P2's turn begins; Gambling Den fires with P2 as the acting player.
	h.P2.ClickCard(GamblingDen)
	h.P2.ClickOption("Mars") // the revealed top card is Mars, so the name matches.

	h.P2.ExpectAmber(2)
	h.P1.ExpectAmber(0) // the artifact's controller gains nothing.
}

// TestGamblingDenLoseOnMiss names a house the revealed top card does not match, so
// the acting player loses 2 Æmber.
func TestGamblingDenLoseOnMiss(t *testing.T) {
	top := ct.Creature(ct.OfHouse(card.House.Mars))
	h := ct.Play(t, ct.Setup{
		P1: ct.Side{House: card.House.Shadows, InPlay: ct.Cards(GamblingDen)},
		P2: ct.Side{House: card.House.Mars, Deck: ct.Cards(top), Amber: 3},
	})

	h.P1.EndTurn()
	h.P2.ClickCard(GamblingDen)
	h.P2.ClickOption("Sanctum") // the top card is Mars, not Sanctum, so the name misses.

	h.P2.ExpectAmber(1) // 3 - 2.
}

// TestGamblingDenDecline confirms the "may" is optional: declining names no house
// and changes no Æmber.
func TestGamblingDenDecline(t *testing.T) {
	top := ct.Creature(ct.OfHouse(card.House.Mars))
	h := ct.Play(t, ct.Setup{
		P1: ct.Side{House: card.House.Shadows, InPlay: ct.Cards(GamblingDen)},
		P2: ct.Side{House: card.House.Mars, Deck: ct.Cards(top), Amber: 2},
	})

	h.P1.EndTurn()
	h.P2.ClickDone()

	h.P2.ExpectAmber(2) // unchanged.
}
