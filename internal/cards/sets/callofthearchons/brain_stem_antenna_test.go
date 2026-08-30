package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Brain Stem Antenna
//
//	House:  Mars
//	Type:   Upgrade
//	Rarity: Rare
//
//	This creature gains, "After you play a Mars creature, ready this creature, and for the remainder of the turn this creature belongs to house Mars."
func TestBrainStemAntenna(t *testing.T) {
	t.Run("readies the host and makes it Mars after you play a Mars creature", func(t *testing.T) {
		var host, martian ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Mars,
				InPlay: ct.Cards(
					ct.Upgraded(ct.Bind(&host, ct.Creature(ct.OfHouse(card.House.Brobnar), ct.Power(4))), BrainStemAntenna),
				),
				Hand: ct.Cards(ct.Bind(&martian, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(2)))),
			},
		})
		host.Exhaust()

		h.P1.Play(martian)

		h.Expect(host).Ready()
		if got := h.Game().House(host.ID()); got != card.House.Mars {
			t.Fatalf("host house = %s, want Mars", got)
		}

		h.P1.Reap(host)
		h.P1.ExpectAmber(1)

		h.P1.EndTurn()
		if got := h.Game().House(host.ID()); got != card.House.Brobnar {
			t.Fatalf("host house next turn = %s, want Brobnar", got)
		}
	})

	t.Run("does not trigger after you play a non-Mars creature", func(t *testing.T) {
		var host, logosCreature ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Logos,
				InPlay: ct.Cards(
					ct.Upgraded(ct.Bind(&host, ct.Creature(ct.OfHouse(card.House.Brobnar), ct.Power(4))), BrainStemAntenna),
				),
				Hand: ct.Cards(ct.Bind(&logosCreature, ct.Creature(ct.OfHouse(card.House.Logos), ct.Power(2)))),
			},
		})
		host.Exhaust()

		h.P1.Play(logosCreature)

		h.Expect(host).Exhausted()
		if got := h.Game().House(host.ID()); got != card.House.Brobnar {
			t.Fatalf("host house = %s, want Brobnar", got)
		}
	})

	t.Run("does not trigger after your opponent plays a Mars creature", func(t *testing.T) {
		var host, enemyMartian ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Brobnar,
				InPlay: ct.Cards(
					ct.Upgraded(ct.Bind(&host, ct.Creature(ct.OfHouse(card.House.Brobnar), ct.Power(4))), BrainStemAntenna),
				),
			},
			P2: ct.Side{
				Hand: ct.Cards(ct.Bind(&enemyMartian, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(2)))),
			},
		})

		h.P1.EndTurn()
		h.P2.ChooseHouse(card.House.Mars)
		host.Exhaust()
		h.P2.Play(enemyMartian)

		h.Expect(host).Exhausted()
		if got := h.Game().House(host.ID()); got != card.House.Brobnar {
			t.Fatalf("host house = %s, want Brobnar", got)
		}
	})
}
