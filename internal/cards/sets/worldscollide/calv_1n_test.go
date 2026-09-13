package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// CALV-1N
//
//	House:  Star Alliance
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  2
//	Armor:  1
//	Traits: Robot
//
//	Fight/Reap: Draw a card.
//	CALV-1N may be played as an Upgrade instead of a Creature, with the text: "This Creature gains, 'Fight/Reap: Draw a card.'"
func TestCALV1N(t *testing.T) {
	var calvin ct.Card
	h := ct.Play(t, ct.Setup{
		P1: ct.Side{
			House:  card.House.StarAlliance,
			InPlay: ct.Cards(ct.Bind(&calvin, CALV1N)),
			Deck:   ct.Cards(ct.Creature()),
		},
	})
	calvin.Ready()

	h.P1.Reap(calvin)

	h.Expect(calvin).At(ct.PlayArea)
	if got := len(h.Game().Hand(0)); got != 1 {
		t.Fatalf("hand after reap = %d, want 1", got)
	}
}
