package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Veemos Lightbringer
//
//	House:  Sanctum
//	Type:   Creature
//	Rarity: Rare
//	Power:  6
//	Traits: Angel • Spirit
//
//	Play: Destroy each elusive creature.
func TestVeemosLightbringer(t *testing.T) {
	t.Run("destroys each elusive creature when played", func(t *testing.T) {
		var elusive, plain ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Sanctum, Hand: ct.Cards(VeemosLightbringer)},
			P2: ct.Side{
				InPlay: ct.Cards(
					ct.Bind(&elusive, ct.Creature(ct.OfHouse(card.House.Mars), ct.Keywords(card.Keyword.Elusive))),
					ct.Bind(&plain, ct.Creature(ct.OfHouse(card.House.Mars))),
				),
			},
		})

		h.P1.Play(VeemosLightbringer)

		h.Expect(elusive).At(ct.Discard)
		h.Expect(plain).At(ct.PlayArea)
	})
}
