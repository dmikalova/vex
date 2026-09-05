package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Commpod
//
//	House:  Mars
//	Type:   Artifact
//	Rarity: Rare
//	Traits: Item
//
//	Action: Reveal any number of Mars cards from your hand, and for each card revealed this way, ready a friendly Mars creature.
func TestCommpod(t *testing.T) {
	t.Run("readies a Mars creature for each Mars card revealed", func(t *testing.T) {
		var a, b, t1, t2 ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Mars,
				InPlay: ct.Cards(
					Commpod,
					ct.Bind(&a, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(3))),
					ct.Bind(&b, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(3))),
				),
				Hand: ct.Cards(
					ct.Bind(&t1, ct.Tactic(ct.OfHouse(card.House.Mars))),
					ct.Bind(&t2, ct.Tactic(ct.OfHouse(card.House.Mars))),
				),
			},
		})
		a.Exhaust()
		b.Exhaust()

		h.P1.UseAction(Commpod)
		h.P1.ClickCard(t1) // reveal both Mars cards from hand
		h.P1.ClickCard(t2)
		h.P1.ClickCard(a) // then the second ready has a lone candidate and auto-resolves

		h.Expect(a).Ready()
		h.Expect(b).Ready()
	})

	t.Run("readies nothing when no Mars card is revealed", func(t *testing.T) {
		var a ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Mars,
				InPlay: ct.Cards(
					Commpod,
					ct.Bind(&a, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(3))),
				),
				Hand: ct.Cards(ct.Tactic(ct.OfHouse(card.House.Untamed))),
			},
		})
		a.Exhaust()

		h.P1.UseAction(Commpod)

		h.Expect(a).Exhausted()
	})
}
