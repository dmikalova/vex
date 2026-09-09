package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// E'e on the Fringes
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  1
//	Traits: Imp
//
//	Elusive.
//	After you discard a Dis card, you may purge a Dis card from a discard pile -> steal 1 Æmber.
func TestEeOnTheFringes(t *testing.T) {
	t.Run("purging a Dis card steals 1 after a Dis discard", func(t *testing.T) {
		var ee, fodder, victim ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Dis,
				InPlay: ct.Cards(ct.Bind(&ee, EeOnTheFringes)),
				Hand: ct.Cards(
					ct.Bind(&fodder, ct.Creature(ct.OfHouse(card.House.Dis))),
				),
				Discard: ct.Cards(
					ct.Bind(&victim, ct.Creature(ct.OfHouse(card.House.Dis))),
				),
			},
			P2: ct.Side{Amber: 3},
		})

		h.P1.Discard(fodder)
		h.P1.ClickOption("Yes")
		h.P1.ClickCard(victim)

		h.Expect(victim).At(ct.Purge)
		h.P1.ExpectAmber(1)
		h.P2.ExpectAmber(2)
	})

	t.Run("stays quiet for an off-house discard", func(t *testing.T) {
		var ee, fodder, dis ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Dis,
				InPlay: ct.Cards(ct.Bind(&ee, EeOnTheFringes)),
				Hand: ct.Cards(
					ct.Bind(&fodder, ct.Creature(ct.OfHouse(card.House.Logos))),
				),
				Discard: ct.Cards(
					ct.Bind(&dis, ct.Creature(ct.OfHouse(card.House.Dis))),
				),
			},
			P2: ct.Side{Amber: 3},
		})

		h.P1.Discard(fodder)

		h.Expect(dis).At(ct.Discard)
		h.P2.ExpectAmber(3)
	})
}
