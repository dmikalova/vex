package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Grasping Vines
//
//	House:  Untamed
//	Type:   Tactic
//	Rarity: Uncommon
//	Æmber:  1
//
//	Play: Put up to 3 artifacts into their owners' hands.
func TestGraspingVines(t *testing.T) {
	t.Run("returns artifacts from either player to their owners' hands", func(t *testing.T) {
		var a1, a2 ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Untamed,
				Hand:   ct.Cards(GraspingVines),
				InPlay: ct.Cards(ct.Bind(&a1, ct.Artifact())),
			},
			P2: ct.Side{
				InPlay: ct.Cards(ct.Bind(&a2, ct.Artifact())),
			},
		})

		h.P1.Play(GraspingVines)
		h.P1.ClickCard(a1)
		h.P1.ClickCard(a2)

		h.Expect(a1).At(ct.Hand)
		h.Expect(a2).At(ct.Hand)
	})
}
