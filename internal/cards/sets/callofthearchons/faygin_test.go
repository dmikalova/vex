package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
)

// Faygin
//
//	House:  Shadows
//	Type:   Creature
//	Rarity: Rare
//	Power:  3
//	Traits: Human • Thief
//
//	Elusive.
//	Reap: Put an Urchin from play or your discard pile into your hand.
func TestFaygin(t *testing.T) {
	t.Run("recovers an Urchin from the discard pile", func(t *testing.T) {
		var urchin ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:   card.House.Shadows,
				InPlay:  ct.Cards(Faygin),
				Discard: ct.Cards(ct.Bind(&urchin, Urchin)),
			},
		})

		h.P1.Reap(Faygin)

		h.Expect(urchin).At(ct.Hand)
	})

	t.Run("bounces a friendly Urchin from play", func(t *testing.T) {
		var urchin ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Shadows,
				InPlay: ct.Cards(
					Faygin,
					ct.Bind(&urchin, Urchin),
				),
			},
		})

		h.P1.Reap(Faygin)

		h.Expect(urchin).At(ct.Hand)
	})
}
