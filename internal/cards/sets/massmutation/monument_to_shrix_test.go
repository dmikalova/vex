package massmutation

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Monument to Shrix
//
//	House:  Saurian
//	Type:   Artifact
//	Rarity: Rare
//	Traits: Location
//
//	You may spend Æmber on Monument to Shrix when forging keys.
//	Action: If Citizen Shrix is in your discard pile, move 1 Æmber from any player's pool to Monument to Shrix. Otherwise, move 1 Æmber from your pool to Monument to Shrix.
func TestMonumentToShrix(t *testing.T) {
	t.Run("banks 1 Æmber from your own pool", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Saurian,
				InPlay: ct.Cards(MonumentToShrix),
				Amber:  2,
			},
		})

		h.P1.UseAction(MonumentToShrix)

		h.Expect(MonumentToShrix).AmberOn(1)
		h.P1.ExpectAmber(1)
	})

	t.Run(
		"banks from any player's pool when Citizen Shrix is in your discard pile",
		func(t *testing.T) {
			h := ct.Play(t, ct.Setup{
				P1: ct.Side{
					House:   card.House.Saurian,
					InPlay:  ct.Cards(MonumentToShrix),
					Discard: ct.Cards(CitizenShrix),
					Amber:   2,
				},
				P2: ct.Side{Amber: 2},
			})

			h.P1.UseAction(MonumentToShrix)
			h.P1.ClickOption("your opponent's pool")

			h.Expect(MonumentToShrix).AmberOn(1)
			h.P1.ExpectAmber(2)
			h.P2.ExpectAmber(1)
		},
	)
}
