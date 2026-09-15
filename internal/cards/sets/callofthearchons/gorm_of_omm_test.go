package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Gorm of Omm
//
//	House:  Sanctum
//	Type:   Artifact
//	Rarity: Uncommon
//	Traits: Item
//
//	Versatile.
//	Action: Destroy Gorm of Omm and an artifact.
func TestGormOfOmm(t *testing.T) {
	t.Run("sacrifices itself and destroys another artifact", func(t *testing.T) {
		var cannon ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Sanctum,
				InPlay: ct.Cards(GormOfOmm),
			},
			P2: ct.Side{
				InPlay: ct.Cards(
					ct.Bind(&cannon, Cannon),
				), // the only other artifact, so it is the forced choice
			},
		})

		h.P1.UseAction(GormOfOmm)

		h.Expect(GormOfOmm).At(ct.Discard) // sacrifices itself
		h.Expect(cannon).At(ct.Discard)    // and destroys the other artifact
	})
}
