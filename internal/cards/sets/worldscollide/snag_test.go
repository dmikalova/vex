package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Snag
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  5
//	Traits: Demon
//
//	Fight: Your opponent must choose the house of the creature Snag fights as their active house on their next turn.
func TestSnag(t *testing.T) {
	var snag, foe ct.Card
	h := ct.Play(t, ct.Setup{
		P1: ct.Side{
			House:  card.House.Dis,
			InPlay: ct.Cards(ct.Bind(&snag, Snag)),
		},
		P2: ct.Side{
			InPlay: ct.Cards(
				ct.Bind(&foe, ct.Creature(ct.OfHouse(card.House.Logos), ct.Power(3), ct.Armor(6))),
			),
		},
	})

	h.P1.Fight(snag, foe)

	if got := h.Game().State.ForcedHouseNext[1].Value; got != card.House.Logos {
		t.Errorf("forced house = %v, want Logos (the house of the creature Snag fought)", got)
	}
}
