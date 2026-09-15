package ageofascension

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Earthbind
//
//	House:  Untamed
//	Type:   Upgrade
//	Rarity: Rare
//	Bonus:  Æmber
//
//	This creature cannot be used unless you have discarded a card from your hand this turn.
func TestEarthbind(t *testing.T) {
	setup := func(t *testing.T, host *ct.Card) *ct.Harness {
		t.Helper()
		return ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Untamed,
				Hand:  ct.Cards(ct.Creature()),
				InPlay: ct.Cards(
					ct.Upgraded(
						ct.Bind(host, ct.Creature(ct.OfHouse(card.House.Untamed), ct.Power(3))),
						Earthbind,
					),
				),
			},
		})
	}

	t.Run("host cannot be used before a card is discarded", func(t *testing.T) {
		var host ct.Card
		h := setup(t, &host)

		h.P1.ExpectCannotUse(host)
	})

	t.Run("host can be used once a card is discarded from hand", func(t *testing.T) {
		var host, spare ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Untamed,
				Hand:  ct.Cards(ct.Bind(&spare, ct.Creature())),
				InPlay: ct.Cards(
					ct.Upgraded(
						ct.Bind(&host, ct.Creature(ct.OfHouse(card.House.Untamed), ct.Power(3))),
						Earthbind,
					),
				),
			},
		})

		h.P1.Discard(spare)
		h.P1.Reap(host)

		h.P1.ExpectAmber(1)
	})
}
