package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Key Hammer
//
//	House:  Dis
//	Type:   Tactic
//	Rarity: Uncommon
//	Æmber:  1
//
//	Play: If your opponent forged a key on their previous turn, unforge one of your opponent's keys, and your opponent gains 6 Æmber.
func TestKeyHammer(t *testing.T) {
	t.Run("unforges the key the opponent just forged", func(t *testing.T) {
		var hammer, filler ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Dis,
				Hand:  ct.Cards(ct.Bind(&hammer, KeyHammer)),
			},
			P2: ct.Side{
				Amber: 6,
				Hand:  ct.Cards(ct.Bind(&filler, ct.Creature(ct.OfHouse(card.House.Brobnar)))),
			},
		})

		h.P1.EndTurn()
		h.P2.ChooseHouse(card.House.Brobnar)
		h.P2.ExpectKeys(1)
		h.P2.EndTurn()
		h.P1.ChooseHouse(card.House.Dis)
		h.P1.Play(hammer)

		h.P2.ExpectKeys(0)
		h.P2.ExpectAmber(6)
		h.P1.ExpectAmber(1)
	})

	t.Run("only pays the opponent when they forged nothing", func(t *testing.T) {
		var hammer ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Dis,
				Hand:  ct.Cards(ct.Bind(&hammer, KeyHammer)),
			},
		})

		h.P1.Play(hammer)

		h.P2.ExpectKeys(0)
		h.P2.ExpectAmber(6)
		h.P1.ExpectAmber(1)
	})
}
