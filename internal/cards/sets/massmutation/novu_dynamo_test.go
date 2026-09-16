package massmutation

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Novu Dynamo
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  8
//	Armor:  2
//	Traits: Robot
//
//	At the start of your turn, discard a Logos card from your hand or archives -> gain 1 Æmber. Otherwise, destroy Novu Dynamo.
func TestNovuDynamo(t *testing.T) {
	t.Run("discards a Logos card and gains 1 Æmber", func(t *testing.T) {
		var novu, logos ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Logos,
				InPlay: ct.Cards(ct.Bind(&novu, NovuDynamo)),
				Hand:   ct.Cards(ct.Bind(&logos, ct.Creature(ct.OfHouse(card.House.Logos)))),
			},
			P2: ct.Side{},
		})

		h.P1.EndTurn()
		h.P2.ChooseHouse(card.House.Brobnar)
		h.P2.EndTurn()
		h.P1.ClickCard(logos) // discard the Logos card rather than decline

		h.Expect(logos).At(ct.Discard)
		h.Expect(novu).At(ct.PlayArea)
		h.P1.ExpectAmber(1)
	})

	t.Run("declining destroys Novu Dynamo", func(t *testing.T) {
		var novu, logos ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Logos,
				InPlay: ct.Cards(ct.Bind(&novu, NovuDynamo)),
				Hand:   ct.Cards(ct.Bind(&logos, ct.Creature(ct.OfHouse(card.House.Logos)))),
			},
			P2: ct.Side{},
		})

		h.P1.EndTurn()
		h.P2.ChooseHouse(card.House.Brobnar)
		h.P2.EndTurn()
		h.P1.ClickDone() // decline the discard

		h.Expect(novu).At(ct.Discard)
		h.Expect(logos).At(ct.Hand)
		h.P1.ExpectAmber(0)
	})

	t.Run("with no Logos card to discard, Novu Dynamo is destroyed", func(t *testing.T) {
		var novu ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Logos,
				InPlay: ct.Cards(ct.Bind(&novu, NovuDynamo)),
				Hand:   ct.Cards(ct.Creature(ct.OfHouse(card.House.Untamed))),
			},
			P2: ct.Side{},
		})

		h.P1.EndTurn()
		h.P2.ChooseHouse(card.House.Brobnar)
		h.P2.EndTurn()

		h.Expect(novu).At(ct.Discard)
		h.P1.ExpectAmber(0)
	})
}
