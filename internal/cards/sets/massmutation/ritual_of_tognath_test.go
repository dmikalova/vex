package massmutation

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Ritual of Tognath
//
//	House:  Dis
//	Type:   Tactic
//	Rarity: Rare
//	Bonus:  Æmber Æmber Æmber
//
//	Play: Destroy 2 friendly creatures.
func TestRitualOfTognath(t *testing.T) {
	t.Run("destroys 2 chosen friendly creatures", func(t *testing.T) {
		var a, b, c ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Dis,
				Hand:  ct.Cards(RitualOfTognath),
				InPlay: ct.Cards(
					ct.Bind(&a, ct.Creature(ct.OfHouse(card.House.Dis))),
					ct.Bind(&b, ct.Creature(ct.OfHouse(card.House.Dis))),
					ct.Bind(&c, ct.Creature(ct.OfHouse(card.House.Dis))),
				),
			},
		})

		h.P1.Play(RitualOfTognath)
		h.P1.ClickCard(a)
		h.P1.ClickCard(b)

		h.Expect(a).At(ct.Discard)
		h.Expect(b).At(ct.Discard)
		h.Expect(c).At(ct.PlayArea)
	})
}
