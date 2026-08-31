package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Smaaash
//
//	House:  Brobnar
//	Type:   Creature
//	Rarity: Common
//	Power:  5
//	Traits: Giant
//
//	Play: Stun a creature.
func TestSmaaash(t *testing.T) {
	t.Run("stuns a chosen creature when played", func(t *testing.T) {
		var foe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Brobnar, Hand: ct.Cards(Smaaash)},
			P2: ct.Side{
				InPlay: ct.Cards(
					ct.Bind(&foe, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(3))),
				),
			},
		})

		h.P1.Play(Smaaash)
		h.P1.ClickCard(foe)

		h.Expect(foe).At(ct.PlayArea).Stunned(true)
	})
}
