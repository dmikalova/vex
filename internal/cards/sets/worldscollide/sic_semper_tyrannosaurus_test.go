package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Sic Semper Tyrannosaurus
//
//	House:  Saurian
//	Type:   Tactic
//	Rarity: Uncommon
//
//	Play: Move all Æmber from the most powerful creature to your pool. Destroy the chosen creature.
func TestSicSemperTyrannosaurus(t *testing.T) {
	t.Run("empties the most powerful creature into your pool and destroys it", func(t *testing.T) {
		var big, small ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Saurian,
				Hand:  ct.Cards(SicSemperTyrannosaurus),
			},
			P2: ct.Side{InPlay: ct.Cards(
				ct.Bind(&big, ct.Creature(ct.Power(6))),
				ct.Bind(&small, ct.Creature(ct.Power(3))),
			)},
		})
		h.Game().State.Cards[big.ID()].Amber = 3

		h.P1.Play(SicSemperTyrannosaurus)

		h.P1.ExpectAmber(3)
		h.Expect(big).At(ct.Discard)
		h.Expect(small).At(ct.PlayArea).AmberOn(0)
	})
}
