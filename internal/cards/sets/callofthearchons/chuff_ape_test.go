package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Chuff Ape
//
//	House:  Mars
//	Type:   Creature
//	Rarity: Rare
//	Power:  11
//	Traits: Beast
//
//	Taunt.
//	Chuff Ape enters play stunned.
//	Fight/Reap: You may destroy another friendly creature -> fully heal Chuff Ape.
func TestChuffApe(t *testing.T) {
	t.Run("enters play stunned", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Mars, Hand: ct.Cards(ChuffApe)},
		})

		h.P1.Play(ChuffApe)

		h.Expect(ChuffApe).At(ct.PlayArea).Stunned(true)
	})

	t.Run("Fight/Reap sacrifices another friendly creature to fully heal", func(t *testing.T) {
		var chuff, friend ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Mars,
				InPlay: ct.Cards(
					ct.Bind(&chuff, ChuffApe),
					ct.Bind(&friend, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(4))),
				),
			},
		})
		chuff.Damaged(6)

		h.P1.Reap(chuff)
		h.P1.ClickCard(friend)

		h.Expect(friend).At(ct.Discard)
		h.Expect(chuff).Damage(0)
	})

	t.Run("Fight/Reap may be declined", func(t *testing.T) {
		var chuff, friend ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Mars,
				InPlay: ct.Cards(
					ct.Bind(&chuff, ChuffApe),
					ct.Bind(&friend, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(4))),
				),
			},
		})
		chuff.Damaged(6)

		h.P1.Reap(chuff)
		h.P1.ClickDone()

		h.Expect(friend).At(ct.PlayArea)
		h.Expect(chuff).Damage(6)
	})
}
