package worldscollide

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
)

// Saury About That
//
//	House:  Saurian
//	Type:   Tactic
//	Rarity: Rare
//
//	Play: Destroy a creature -> its controller gains 1 Æmber.
func TestSauryAboutThat(t *testing.T) {
	t.Run("destroys a creature and its controller gains 1 Æmber", func(t *testing.T) {
		var target ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Saurian,
				Hand:  ct.Cards(SauryAboutThat),
			},
			P2: ct.Side{InPlay: ct.Cards(ct.Bind(&target, ct.Creature(ct.Power(3))))},
		})

		h.P1.Play(SauryAboutThat)

		h.Expect(target).At(ct.Discard)
		h.P2.ExpectAmber(1) // the destroyed creature's controller
		h.P1.ExpectAmber(0)
	})
}
