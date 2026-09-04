package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

func TestTirelessCrocag(t *testing.T) {
	t.Run("cannot reap but fights out of house, then dies with the board", func(t *testing.T) {
		var crocag, foe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Logos,
				InPlay: ct.Cards(ct.Bind(&crocag, TirelessCrocag)),
			},
			P2: ct.Side{InPlay: ct.Cards(
				ct.Bind(&foe, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(3))),
			)},
		})

		h.P1.ExpectCannotUseTo(crocag, card.UseKind.Reap)

		// Versatile lets it fight though Brobnar is not the active house. Killing the
		// last enemy creature then makes its own destruction condition true.
		h.P1.Fight(crocag, foe)

		h.Expect(foe).At(ct.Discard)
		h.Expect(crocag).At(ct.Discard)
	})
}
