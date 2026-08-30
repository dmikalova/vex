package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Perilous Wild
//
//	House:  Untamed
//	Type:   Tactic
//	Rarity: Rare
//	Æmber:  1
//
//	Play: Destroy each elusive creature.
func TestPerilousWild(t *testing.T) {
	t.Run("destroys each elusive creature", func(t *testing.T) {
		var elusive, plain ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Untamed, Hand: ct.Cards(PerilousWild)},
			P2: ct.Side{
				InPlay: ct.Cards(
					ct.Bind(&elusive, ct.Creature(ct.OfHouse(card.House.Mars), ct.Keywords(card.Keyword.Elusive))),
					ct.Bind(&plain, ct.Creature(ct.OfHouse(card.House.Mars))),
				),
			},
		})

		h.P1.Play(PerilousWild)

		h.Expect(elusive).At(ct.Discard)
		h.Expect(plain).At(ct.PlayArea)
	})
}
