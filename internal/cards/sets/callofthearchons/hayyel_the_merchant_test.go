package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Hayyel the Merchant
//
//	House:  Sanctum
//	Type:   Creature
//	Rarity: Rare
//	Power:  3
//	Traits: Human • Merchant
//
//	After you play an artifact, gain 1 Æmber.
func TestHayyelTheMerchant(t *testing.T) {
	t.Run("gains 1 Æmber after you play an artifact", func(t *testing.T) {
		var relic ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Sanctum,
				InPlay: ct.Cards(HayyelTheMerchant),
				Hand:   ct.Cards(ct.Bind(&relic, ct.Artifact(ct.OfHouse(card.House.Sanctum)))),
			},
		})

		h.P1.Play(relic)

		h.P1.ExpectAmber(1)
	})
}
