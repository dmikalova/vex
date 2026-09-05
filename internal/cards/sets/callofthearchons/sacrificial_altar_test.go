package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Sacrificial Altar
//
//	House:  Dis
//	Type:   Artifact
//	Rarity: Rare
//	Æmber:  1
//	Traits: Location
//
//	Action: Purge a friendly Human trait creature -> play a creature from your discard pile.
func TestSacrificialAltar(t *testing.T) {
	t.Run("trades a Human for a creature in the discard pile", func(t *testing.T) {
		var altar, human, risen ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Dis,
				InPlay: ct.Cards(
					ct.Bind(&altar, SacrificialAltar),
					ct.Bind(&human, HuntingWitch),
				),
				Discard: ct.Cards(ct.Bind(&risen, HuntingWitch)),
			},
		})

		h.P1.UseAction(altar)

		h.Expect(human).At(ct.Purge)
		h.Expect(risen).At(ct.PlayArea)
	})

	t.Run("with no Human to purge, the discard pile is left alone", func(t *testing.T) {
		var altar, risen ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:   card.House.Dis,
				InPlay:  ct.Cards(ct.Bind(&altar, SacrificialAltar)),
				Discard: ct.Cards(ct.Bind(&risen, HuntingWitch)),
			},
		})

		h.P1.UseAction(altar)

		h.Expect(risen).At(ct.Discard)
	})
}
