package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Gongoozle
//
//	House:  Dis
//	Type:   Tactic
//	Rarity: Uncommon
//	Æmber:  1
//
//	Play: Deal 3 damage to a creature. If it is not destroyed, its owner discards a random card from their hand.
func TestGongoozle(t *testing.T) {
	t.Run("a surviving creature's owner discards a random card", func(t *testing.T) {
		var foe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Dis, Hand: ct.Cards(Gongoozle)},
			P2: ct.Side{
				InPlay: ct.Cards(
					ct.Bind(&foe, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(5))),
				),
				Hand: ct.Cards(
					ct.Creature(ct.OfHouse(card.House.Mars)),
					ct.Creature(ct.OfHouse(card.House.Mars)),
				),
			},
		})

		h.P1.Play(Gongoozle)

		h.Expect(foe).Damage(3).At(ct.PlayArea)
		if got := int(h.Game().State.Hand[1].Count); got != 1 {
			t.Errorf("opponent hand = %d, want 1 (one discarded)", got)
		}
	})
}
