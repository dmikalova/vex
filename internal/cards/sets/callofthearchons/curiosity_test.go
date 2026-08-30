package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Curiosity
//
//	House:  Untamed
//	Type:   Tactic
//	Rarity: Rare
//	Æmber:  1
//
//	Play: Destroy each Scientist trait creature.
func TestCuriosity(t *testing.T) {
	t.Run("destroys each Scientist trait creature and spares others", func(t *testing.T) {
		var sci, beast ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Untamed, Hand: ct.Cards(Curiosity)},
			P2: ct.Side{
				InPlay: ct.Cards(
					ct.Bind(&sci, ct.Creature(ct.OfHouse(card.House.Logos), ct.Power(3), ct.Traits("Scientist"))),
					ct.Bind(&beast, ct.Creature(ct.OfHouse(card.House.Untamed), ct.Power(3), ct.Traits("Beast"))),
				),
			},
		})

		h.P1.Play(Curiosity)

		h.Expect(sci).At(ct.Discard)
		h.Expect(beast).At(ct.PlayArea)
	})
}
