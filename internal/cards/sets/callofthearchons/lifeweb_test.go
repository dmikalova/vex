package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Lifeweb
//
//	House:  Untamed
//	Type:   Tactic
//	Rarity: Uncommon
//	Æmber:  1
//
//	Play: If your opponent played 3 or more creatures on their previous turn, steal 2 Æmber.
func TestLifeweb(t *testing.T) {
	t.Run("steals when the opponent played three creatures last turn", func(t *testing.T) {
		var web, a, b, c ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Untamed,
				Hand:  ct.Cards(ct.Bind(&web, Lifeweb)),
			},
			P2: ct.Side{
				Amber: 5,
				Hand: ct.Cards(
					ct.Bind(&a, ct.Creature(ct.OfHouse(card.House.Brobnar))),
					ct.Bind(&b, ct.Creature(ct.OfHouse(card.House.Brobnar))),
					ct.Bind(&c, ct.Creature(ct.OfHouse(card.House.Brobnar))),
				),
			},
		})

		h.P1.EndTurn()
		h.P2.ChooseHouse(card.House.Brobnar)
		h.P2.Play(a)
		h.P2.Play(b)
		h.P2.Play(c)
		h.P2.EndTurn()
		h.P1.ChooseHouse(card.House.Untamed)
		h.P1.Play(web)

		h.P1.ExpectAmber(3)
		h.P2.ExpectAmber(3)
	})

	t.Run("stays quiet when the opponent played too few creatures", func(t *testing.T) {
		var web, a ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Untamed,
				Hand:  ct.Cards(ct.Bind(&web, Lifeweb)),
			},
			P2: ct.Side{
				Amber: 5,
				Hand:  ct.Cards(ct.Bind(&a, ct.Creature(ct.OfHouse(card.House.Brobnar)))),
			},
		})

		h.P1.EndTurn()
		h.P2.ChooseHouse(card.House.Brobnar)
		h.P2.Play(a)
		h.P2.EndTurn()
		h.P1.ChooseHouse(card.House.Untamed)
		h.P1.Play(web)

		h.P1.ExpectAmber(1)
		h.P2.ExpectAmber(5)
	})
}
