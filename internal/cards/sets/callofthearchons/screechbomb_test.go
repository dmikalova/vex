package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Screechbomb
//
//	House:  Brobnar
//	Type:   Artifact
//	Rarity: Uncommon
//	Traits: Weapon
//
//	Versatile.
//	Action: Destroy Screechbomb. Your opponent loses 2 Æmber.
func TestScreechbomb(t *testing.T) {
	t.Run("destroys itself and drains 2 Æmber from the opponent", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Brobnar, InPlay: ct.Cards(Screechbomb)},
			P2: ct.Side{Amber: 3},
		})

		h.P1.UseAction(Screechbomb)

		h.Expect(Screechbomb).At(ct.Discard)
		h.P2.ExpectAmber(1)
	})
}
