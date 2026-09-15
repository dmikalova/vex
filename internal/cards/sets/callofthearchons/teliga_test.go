package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Teliga
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Rare
//	Power:  3
//	Traits: Human • Witch
//
//	After your opponent plays a card, if it is a creature, gain 1 Æmber.
func TestTeliga(t *testing.T) {
	t.Run("gains Æmber when the opponent plays a creature", func(t *testing.T) {
		var teliga, beast, tactic ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Untamed,
				InPlay: ct.Cards(ct.Bind(&teliga, Teliga)),
			},
			P2: ct.Side{
				House: card.House.Brobnar,
				Hand: ct.Cards(
					ct.Bind(&beast, ct.Creature(ct.OfHouse(card.House.Brobnar), ct.Power(3))),
					ct.Bind(&tactic, ct.Tactic(ct.OfHouse(card.House.Brobnar))),
				),
			},
		})

		h.P1.EndTurn()
		h.P2.ChooseHouse(card.House.Brobnar)
		h.P2.Play(beast)

		h.P1.ExpectAmber(1)

		// An action card is not a creature, so it pays nothing.
		h.P2.Play(tactic)
		h.P1.ExpectAmber(1)
	})

	t.Run("stays quiet when its own controller plays a creature", func(t *testing.T) {
		var teliga, beast ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Untamed,
				InPlay: ct.Cards(ct.Bind(&teliga, Teliga)),
				Hand: ct.Cards(
					ct.Bind(&beast, ct.Creature(ct.OfHouse(card.House.Untamed), ct.Power(3))),
				),
			},
		})

		h.P1.Play(beast)

		h.P1.ExpectAmber(0)
	})
}
