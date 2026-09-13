package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Tactical Officer Moon
//
//	House:  Star Alliance
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  4
//	Traits: Human
//
//	Assault 2.
//	Play: You may rearrange the Creatures in a player's battleline.
func TestTacticalOfficerMoon(t *testing.T) {
	var a, b ct.Card
	h := ct.Play(t, ct.Setup{
		P1: ct.Side{
			House: card.House.StarAlliance,
			Hand:  ct.Cards(TacticalOfficerMoon),
			InPlay: ct.Cards(
				ct.Bind(&a, ct.Creature(ct.OfHouse(card.House.StarAlliance))),
				ct.Bind(&b, ct.Creature(ct.OfHouse(card.House.StarAlliance))),
			),
		},
	})

	h.P1.Play(TacticalOfficerMoon)
	// Accept the optional rearrange, swap the two existing creatures, then stop.
	h.P1.ClickOption("Yes")
	h.P1.ClickCard(a)
	h.P1.ClickCard(b)
	h.P1.ClickDone()

	// Every creature stays in play; only their positions changed.
	h.Expect(a).At(ct.PlayArea)
	h.Expect(b).At(ct.PlayArea)
	h.Expect(TacticalOfficerMoon).At(ct.PlayArea)
}
