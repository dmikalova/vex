package massmutation

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
)

// Q-Mechs
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Common
//	Power:  1
//	Traits: Robot
//
//	Play: Draw a card.
//	Destroyed: Archive Q-Mechs.
func TestQMechs(t *testing.T) {
	t.Run("draws a card when played", func(t *testing.T) {
		var top ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Logos,
				Hand:  ct.Cards(QMechs),
				Deck:  ct.Cards(ct.Bind(&top, ct.Creature(ct.OfHouse(card.House.Logos)))),
			},
		})

		h.P1.Play(QMechs)

		h.Expect(top).At(ct.Hand)
	})

	t.Run("archives itself when destroyed", func(t *testing.T) {
		var qmechs, foe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Logos,
				InPlay: ct.Cards(ct.Bind(&qmechs, QMechs)),
			},
			P2: ct.Side{
				InPlay: ct.Cards(
					ct.Bind(&foe, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(5))),
				),
			},
		})

		h.P1.Fight(qmechs, foe)

		h.Expect(qmechs).At(ct.Archives)
	})
}
