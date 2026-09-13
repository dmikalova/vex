package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Low Dawn
//
//	House:  Untamed
//	Type:   Tactic
//	Rarity: Uncommon
//	Æmber:  1
//
//	Play: If there are 3 or more Untamed Creatures in your discard pile, gain 2 Æmber. Shuffle each Untamed Creature from your discard pile into your deck.
func TestLowDawn(t *testing.T) {
	t.Run("three Untamed creatures gain Æmber and shuffle away", func(t *testing.T) {
		var u1, u2, u3, mars, tactic ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Untamed,
				Hand:  ct.Cards(LowDawn),
				Discard: ct.Cards(
					ct.Bind(&u1, ct.Creature(ct.OfHouse(card.House.Untamed), ct.Power(3))),
					ct.Bind(&u2, ct.Creature(ct.OfHouse(card.House.Untamed), ct.Power(3))),
					ct.Bind(&u3, ct.Creature(ct.OfHouse(card.House.Untamed), ct.Power(3))),
					ct.Bind(&mars, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(3))),
					ct.Bind(&tactic, ct.Tactic(ct.OfHouse(card.House.Untamed))),
				),
			},
		})

		h.P1.Play(LowDawn)

		// +1 from Low Dawn's own bonus, +2 from the ability.
		h.P1.ExpectAmber(3)
		h.Expect(u1).At(ct.Deck)
		h.Expect(u2).At(ct.Deck)
		h.Expect(u3).At(ct.Deck)
		// The wrong house and the non-creature stay in the discard pile.
		h.Expect(mars).At(ct.Discard)
		h.Expect(tactic).At(ct.Discard)
	})

	t.Run("two Untamed creatures gain no Æmber but still shuffle", func(t *testing.T) {
		var u1, u2 ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Untamed,
				Hand:  ct.Cards(LowDawn),
				Discard: ct.Cards(
					ct.Bind(&u1, ct.Creature(ct.OfHouse(card.House.Untamed), ct.Power(3))),
					ct.Bind(&u2, ct.Creature(ct.OfHouse(card.House.Untamed), ct.Power(3))),
				),
			},
		})

		h.P1.Play(LowDawn)

		// Only Low Dawn's own +1 bonus; the ability's 2 Æmber is not gained.
		h.P1.ExpectAmber(1)
		h.Expect(u1).At(ct.Deck)
		h.Expect(u2).At(ct.Deck)
	})
}
