package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Truebaru
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Rare
//	Power:  7
//	Traits: Demon
//
//	Taunt.
//	You must lose 3 Æmber in order to play Truebaru.
//	Destroyed: Gain 5 Æmber.
func TestTruebaru(t *testing.T) {
	t.Run("costs 3 Æmber to play", func(t *testing.T) {
		var truebaru ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Dis,
				Hand:  ct.Cards(ct.Bind(&truebaru, Truebaru)),
				Amber: 4,
			},
		})

		h.P1.Play(truebaru)

		h.P1.ExpectAmber(1)
		h.Expect(truebaru).At(ct.PlayArea)
	})

	t.Run("cannot be played without 3 Æmber", func(t *testing.T) {
		var truebaru ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Dis,
				Hand:  ct.Cards(ct.Bind(&truebaru, Truebaru)),
				Amber: 2,
			},
		})

		h.P1.ExpectCannotPlay(truebaru)
	})

	t.Run("pays back 5 Æmber when destroyed", func(t *testing.T) {
		var truebaru, killer ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Dis,
				InPlay: ct.Cards(ct.Bind(&truebaru, Truebaru)),
			},
			P2: ct.Side{
				InPlay: ct.Cards(ct.Bind(&killer, ct.Creature(ct.Power(9)))),
			},
		})

		h.P1.Fight(truebaru, killer)

		h.P1.ExpectAmber(5)
		h.Expect(truebaru).At(ct.Discard)
	})
}
